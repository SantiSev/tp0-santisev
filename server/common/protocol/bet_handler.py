import logging

from common.network.connection_interface import ConnectionInterface
from common.protocol.bet_processor import BetProcessor
from common.utils.utils import Bet, store_bets

EOF = b"\xff"
HEADER_SIZE = 1
SUCCESS_HEADER = b"\x01"
FAIL = b"\xff"



class BetHandler:
    """Handles individual client connections"""

    def __init__(self):
        self.message_handler = BetProcessor()

    def process_bets(self, client_connection: ConnectionInterface) -> int:
        """
        Process the best being sent by the client - receive and store bets
        """
        header = HEADER_SIZE
        bet_counter = 0

        while True:
            try:
                header = client_connection.receive(HEADER_SIZE)

                if header == EOF:
                    logging.info(f"action: end of transmission | result: success")
                    break

                bets = self.message_handler.process_batch(client_connection)
                if bets:
                    store_bets(bets)
                    self.confirmation_to_client(client_connection)
                    bet_counter += len(bets)
                else:
                    raise Exception("An Error occured proccesing bets")

            except Exception as e:
                self.confirmation_to_client(client_connection, 0)
                logging.critical(
                    f"action: apuesta_recibida | result: fail | cantidad: ${bet_counter}"
                )
                break
        logging.info(
            f"action: apuesta_recibida | result: success | cantidad: ${bet_counter}"
        )
        return bet_counter

    # TODO: refactort client to handle this correctly
    def confirmation_to_client(self, connection: ConnectionInterface, bet_counter: int) -> None:
        """
        Confirm the bet with the client
        """
        try:
            if bet_counter > 0:
                connection.send(SUCCESS_HEADER)
                connection.send(bet_counter.to_bytes(4, "big"))
            else:
                connection.send(FAIL)
        except Exception as e:
            logging.error(f"action: confirm_bet | result: fail | error: {e}")
