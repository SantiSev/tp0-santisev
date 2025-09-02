import logging
import signal
import sys
from common.network.connection_manager import ConnectionManager
from common.protocol.bet_handler import BetHandler
from common.network.connection_interface import ConnectionInterface

class Server:
    def __init__(self, port, listen_backlog, agencies_amount):
        self.connection_manager = ConnectionManager(
            port=port, listen_backlog=listen_backlog
        )
        self.connectedClients: list[ConnectionInterface] = []
        self.is_running = True
        self.agencies_amount = agencies_amount
        self.processed_agencies = 0
        self.bet_handler = BetHandler()
        signal.signal(signal.SIGTERM, self._shutdown)
        signal.signal(signal.SIGINT, self._shutdown)

    def run(self) -> None:
        """Start the server and handle connections"""
        try:
            self.connection_manager.start_listening()
            logging.info("action: server_start | result: success")

            while self.is_running and self.processed_agencies < self.agencies_amount:
                try:
                    client_connection = self._connect_client()
                    bets = self.bet_handler.process_bets(client_connection)
                    self.connectedClients.append(client_connection)
                    self.bet_handler.send_winners(client_connection, bets)
                    self.processed_agencies += 1
                    logging.info(
                        f"action: server_loop | result: success | agencies_processed: {self.processed_agencies} / {self.agencies_amount}"
                    )

                except Exception as e:
                    logging.error(f"action: server_loop | result: error | error: {e}")
                    continue

            self._announce_winners()

        except Exception as e:
            logging.error(f"action: server_run | result: critical_error | error: {e}")

        finally:
            self._shutdown()

    def _connect_client(self) -> ConnectionInterface:
        client_connection = self.connection_manager.accept_connection()
        self.connectedClients[client_connection] = 0
        logging.info(f"action: connect_client | result: success")
        return client_connection

    def _disconnect_client(self, client_connection: ConnectionInterface) -> None:
        """Disconnect a client from the server"""
        if client_connection in self.connectedClients:
            self.connectedClients.remove(client_connection)
        client_connection.close()
        logging.info(f"action: disconnect_client | result: success")

    def _cancel_lottery(self) -> None:
        # TODO: if an error occures processing a client, then cancel the lottery and notify all waiting clients
        pass

    def _announce_winners(self) -> None:
        winners = self.bet_handler.get_all_winners()
        logging.info(
            f"action: sorteo | result: success | winners: {winners}"
        )

    def _shutdown(self, signum=None, frame=None) -> None:
        """Shutdown the server gracefully"""
        logging.info("action: server_shutdown | result: in_progress")
        self.is_running = False
        for client in self.connectedClients:
            client.close()
        logging.info("action: server_shutdown | result: success")
        sys.exit(0)
