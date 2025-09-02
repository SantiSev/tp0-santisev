import logging
import signal
import sys
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
        self.agencies_amount = server_config.agencies_amount
        self.connection_manager = ConnectionManager(
            port=server_config.port, listen_backlog=server_config.listen_backlog
        )
        self.lottery_service = LotteryService()
        self.clientManager: ClientManager = ClientManager(self.lottery_service)

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

                    client: ClientSession = self.clientManager.add_client(
                        client_connection
                    )

                    success, bets = client.begin()

                    if success:
                        self.processed_agencies += 1
                        self.lottery_service.store_bets(bets)
                        logging.info(
                            f"action: server_loop | result: success | agencies_processed: {self.processed_agencies} / {self.agencies_amount}"
                        )
                    else:
                        logging.info(
                            f"action: server_loop | result: fail | an error occured processing the client with {client.id} , halting client"
                        )
                        self.clientManager.remove_client(client.id)
                        continue

                except Exception as e:
                    logging.error(f"action: server_loop | result: error | error: {e}")
                    self._shutdown()
                    continue

            self.lottery_service.announce_winners()

        except Exception as e:
            logging.error(f"action: server_run | result: critical_error | error: {e}")

        finally:
            self._shutdown()

    def _running(self):
        return self.is_running and self.processed_agencies < self.agencies_amount

    def _shutdown(self, signum=None, frame=None) -> None:
        """Shutdown the server gracefully"""
        self.is_running = False
        self.clientManager.shutdown()
        self.connection_manager.shutdown()
        logging.info("action: server_shutdown | result: success")
        sys.exit(0)
