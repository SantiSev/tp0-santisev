import logging
from socket import socket

from common.network.socket_adapter import SocketAdapter
from common.protocol.message_handler import MessageHandler
from common.utils.utils import Bet, store_bets


class ClientHandler:
    """Handles individual client connections"""

    def __init__(self):
        self.message_handler = MessageHandler()

    def handle_client(self, client_socket: SocketAdapter) -> None:
        """
        Handle a client connection - receive and store bets
        """
        try:
            logging.info(f"action: handle_client | result: start")
            bet = self.message_handler.process_message(client_socket)
            self._process_bet(bet)
            self.message_handler.confirmation_to_client(client_socket, True)

        except Exception as e:
            logging.error(f"action: handle_client | result: error | error: {e}")

        finally:
            self._close_client_connection(client_socket)

    def _process_bet(self, bet: Bet):
        store_bets([bet])
        logging.info(
            f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}"
        )

    def _close_client_connection(self, client_socket: SocketAdapter) -> None:
        try:
            client_socket.close()
            logging.info(f"action: close_connection | result: success")
        except socket.error as e:
            logging.error(f"action: close_connection | result: error | error: {e}")
