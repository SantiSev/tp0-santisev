import logging
from multiprocessing import Lock
from multiprocessing.dummy import Process
import signal
import sys
from multiprocessing import Lock, Process
from common.network.connection_manager import ConnectionManager
from common.network.connection_interface import ConnectionInterface
from common.business.lottery_service import LotteryService
from common.server.server_config import ServerConfig
from common.session.client_manager import ClientManager
from common.session.client_session import ClientSession


class Server:
    def __init__(self, server_config: ServerConfig):
        self.is_running = True
        self.processed_agencies = 0
        self.file_lock = Lock()
        self.agencies_amount = server_config.agencies_amount
        self.connection_manager = ConnectionManager(
            port=server_config.port, listen_backlog=server_config.listen_backlog
        )

        self.lottery_service = LotteryService(self.file_lock)
        self.clientManager: ClientManager = ClientManager(self.lottery_service)

        self.processes = []

        signal.signal(signal.SIGTERM, self._shutdown)
        signal.signal(signal.SIGINT, self._shutdown)

    def run(self) -> None:
        """Start the server and handle connections"""
        try:
            self.connection_manager.start_listening()
            logging.info("action: server_start | result: success")

            while self._running():
                try:
                    client_connection: ConnectionInterface = (
                        self.connection_manager.accept_connection()
                    )

                    self.processed_agencies += 1

                    client: ClientSession = self.clientManager.add_client(
                        client_connection
                    )

                    process = Process(target=handle_client, args=(client,))
                    self.processes.append(process)
                    process.start()

                except Exception as e:
                    logging.error(f"action: server_loop | result: error | error: {e}")
                    self._shutdown()

            for process in self.processes:
                process.join()

            self._tally_results()

        except Exception as e:
            logging.error(f"action: server_run | result: critical_error | error: {e}")

        finally:
            self._shutdown()

    def _running(self):
        return self.is_running and self.processed_agencies < self.agencies_amount

    def _tally_results(self):
        """Tally and log the results of the lottery"""
        for client in self.clientManager.connected_clients:
            client.send_results()
        logging.info("action: send_results_to_all_clients | result: success")

        self.lottery_service.announce_winners()

    def _shutdown(self, signum=None, frame=None) -> None:
        """Shutdown the server gracefully"""
        self.is_running = False
        self.clientManager.shutdown()
        self.connection_manager.shutdown()
        logging.info("action: server_shutdown | result: success")
        sys.exit(0)


def handle_client(client: ClientSession) -> None:
    """Handle a client connection"""
    try:
        client.begin()
    except Exception as e:
        logging.error(
            f"action: handle_client | result: error | client_id: {client.id} | error: {e}"
        )
        client.finish()
