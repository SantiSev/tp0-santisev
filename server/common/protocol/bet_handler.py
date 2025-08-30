import logging

from common.network.connection_interface import ConnectionInterface
from common.protocol.bet_processor import BetProcessor
from common.utils.utils import Bet, store_bets

EOF = b"\xFF"
HEADER_SIZE = 1
SUCCESS_HEADER = b"\x01"
SUCCESS_MESSAGE_SIZE = 64
FAIL = b"\xFF"


class BetHandler:
    """Handles individual client connections"""

    def __init__(self):
        self.message_handler = BetProcessor()

    def process_bets(self, client_connection: ConnectionInterface) -> None:
        """
        Process the best being sent by the client - receive and store bets
        """
        header = HEADER_SIZE

        while True:
            try:
                header = client_connection.receive(HEADER_SIZE)

                if header == EOF:
                    logging.info(f"action: end of transmission | result: success")
                    break

                bet = self.message_handler.process_bet(client_connection)
                if bet:
                    store_bets([bet])
                    self.confirmation_to_client(client_connection, bet)
                    logging.info(
                        f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}"
                    )
                else:
                    logging.info(
                        "action: an error occured during the transmission | result: fail"
                    )
                    break

            except Exception as e:
                self.confirmation_to_client(client_connection, False)
                logging.error(f"action: handle_client | result: error | error: {e}")
                break

    def confirmation_to_client(self, connection: ConnectionInterface, bet: Bet) -> None:
        """
        Confirm the bet with the client
        """
        try:
            if bet:
                connection.send(SUCCESS_HEADER)
                response_string = f"{bet.document},{bet.number}"
                response_bytes = response_string.encode("utf-8")
                padded_response = response_bytes.ljust(SUCCESS_MESSAGE_SIZE, b"\x00")
                connection.send(padded_response)
            else:
                connection.send(FAIL)
        except Exception as e:
            logging.error(f"action: confirm_bet | result: fail | error: {e}")
