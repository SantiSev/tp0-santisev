import logging
import signal
import sys
from common.network.connection_manager import ConnectionManager
from common.protocol.client_handler import ClientHandler


class Server:
    def __init__(self, port, listen_backlog):
        self.connection_manager = ConnectionManager(
            port=port, listen_backlog=listen_backlog
        )
        self.client_handler = ClientHandler()

        signal.signal(signal.SIGTERM, self._shutdown)
        signal.signal(signal.SIGINT, self._shutdown) # for local testing

    def run(self) -> None:
        """Start the server and handle connections"""
        try:
            self.connection_manager.start_listening()
            logging.info("action: server_start | result: success")

            while self.connection_manager.is_running:
                try:
                    client_socket = (
                        self.connection_manager.accept_connection()
                    )

                    # Handle the client (receive and log message)
                    self.client_handler.handle_client(client_socket)

                except Exception as e:
                    logging.error(f"action: server_loop | result: error | error: {e}")
                    continue

        except Exception as e:
            logging.error(f"action: server_run | result: critical_error | error: {e}")

        finally:
            self._shutdown()

    def _shutdown(self, signum=None, frame=None) -> None:
        """Shutdown the server gracefully"""
        logging.info("action: server_shutdown | result: in_progress")
        self.connection_manager.shutdown()
        logging.info("action: server_shutdown | result: complete")
        sys.exit(0)
