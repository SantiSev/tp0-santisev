import logging
import struct

from common.network.connection_interface import ConnectionInterface
from common.protocol.bet_processor import BetProcessor
from common.utils.utils import has_won, load_bets, store_bets

EOF = b"\xff"
HEADER_SIZE = 1
SUCCESS_HEADER = b"\x01"
WINNERS_HEADER = b"\x03"
FAIL = b"\xff"


class BetHandler:
    """Handles individual client connections"""

    def __init__(self):
        self.message_handler = BetProcessor()

    def process_bets(self, client_connection: ConnectionInterface) -> int:
        """
        Process the best being sent by the client - receive and store bets
        """
        bets = []

        while True:
            try:
                header = client_connection.receive(HEADER_SIZE)

                if header == EOF:
                    logging.info(f"action: end of transmission | result: success")
                    break

                batchBets = self.message_handler.process_batch(client_connection)
                if batchBets:
                    bets.extend(batchBets)
                    self.confirmation_to_client(client_connection, True)

                else:
                    raise Exception("An Error occured proccesing bets")

            except Exception as e:
                self.confirmation_to_client(client_connection, False)
                logging.error(f"action: process_bets | result: fail | error: {e}")
                logging.critical(
                    f"action: apuesta_recibida | result: fail | cantidad: {len(bets)}"
                )
                break
        store_bets(bets)
        logging.info(
            f"action: apuesta_recibida | result: success | cantidad: {len(bets)}"
        )
        return len(bets)

    def confirmation_to_client(
        self, connection: ConnectionInterface, status: bool
    ) -> None:
        """
        Confirm the bet with the client
        """
        try:
            if status:
                connection.send(SUCCESS_HEADER)
                logging.debug(f"action: confirm_bet | result: success")
            else:
                connection.send(FAIL)
                logging.debug(f"action: confirm_bet | result: fail")
        except Exception as e:
            logging.error(f"action: confirm_bet | result: fail | error: {e}")

    def get_winning_numbers(self) -> int:
        """
        Get the winning numbers
        """
        winner_count = 0
        bets = load_bets()
        for bet in bets:
            if has_won(bet):
                winner_count += 1

        return winner_count

    def send_winning_numbers(
        self, connection: ConnectionInterface, winners_count: int
    ) -> None:
        """
        Send the winning numbers to the client
        """
        try:
            connection.send(WINNERS_HEADER)
            winners_bytes = struct.pack(">H", winners_count)
            connection.send(winners_bytes)
            logging.info(
                f"action: send_winning_numbers | result: success | winners: {winners_count}"
            )
        except Exception as e:
            logging.error(f"action: send_winning_numbers | result: fail | error: {e}")
