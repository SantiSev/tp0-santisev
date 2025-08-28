import logging

from common.network.connection_interface import ConnectionInterface
from common.protocol.bet_processor import BetProcessor
from common.utils.utils import store_bets

EOF = b"\xff"
HEADER_SIZE = 1
SUCCESS = b"\x00"
FAIL = b"\xff"


class BetHandler:
    """Handles individual client connections"""

    def __init__(self):
        self.message_handler = BetProcessor()

    def process_bets(self, client_connection: ConnectionInterface) -> None:
        """
        Process the best being sent by the client - receive and store bets
        """
        header = HEADER_SIZE
        try:
            while True:
                header = client_connection.receive(HEADER_SIZE)

                if header == EOF:
                    logging.info(f"action: end of transmission | result: success")
                    break

                bet = self.message_handler.process_bet(client_connection)
                if bet:
                    store_bets([bet])
                    self.confirmation_to_client(client_connection, True)
                    logging.info(
                        f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}"
                    )

        except Exception as e:
            self.confirmation_to_client(client_connection, False)
            logging.error(f"action: handle_client | result: error | error: {e}")

    def confirmation_to_client(
        self, connection: ConnectionInterface, status: bool
    ) -> None:
        """
        Confirm the bet with the client
        """
        try:
            connection.send(SUCCESS if status else FAIL)
        except Exception as e:
            logging.error(f"action: confirm_bet | result: fail | error: {e}")
