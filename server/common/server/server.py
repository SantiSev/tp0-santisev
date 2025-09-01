import logging
import signal
import sys
from typing import Dict
from common.network.connection_manager import ConnectionManager
from common.protocol.bet_handler import BetHandler
from common.network.connection_interface import ConnectionInterface

# TODO: change the client bet path to be a .env variable insted of a hardcoded value


class Server:
    def __init__(self, port, listen_backlog):
        self.connection_manager = ConnectionManager(
            port=port, listen_backlog=listen_backlog
        )
        self.connectedClients: Dict[ConnectionInterface, int] = {}
        self.is_running = True
        self.bet_handler = BetHandler()
        signal.signal(signal.SIGTERM, self._shutdown)
        signal.signal(signal.SIGINT, self._shutdown)

    def run(self) -> None:
        """Start the server and handle connections"""
        try:
            self.connection_manager.start_listening()
            logging.info("action: server_start | result: success")

            while self.is_running:
                try:
                    client_connection = self._connect_client()
                    bet_counter = self.bet_handler.process_bets(client_connection)
                    self.connectedClients[client_connection] = bet_counter
                    self._disconnect_client(client_connection)

                except Exception as e:
                    logging.error(f"action: server_loop | result: error | error: {e}")
                    continue

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
            del self.connectedClients[client_connection]
        client_connection.close()
        logging.info(f"action: disconnect_client | result: success")

    def _shutdown(self, signum=None, frame=None) -> None:
        """Shutdown the server gracefully"""
        logging.info("action: server_shutdown | result: in_progress")
        self.is_running = False
        for client in self.connectedClients:
            client.close()
        logging.info("action: server_shutdown | result: complete")
        sys.exit(0)
