import logging
from typing import List

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
        bets = []

        while True:
            try:
                header = client_connection.receive(HEADER_SIZE)
                logging.info(
                    f"action: receive_header | result: success | header: {header}"
                )

                if header == EOF:
                    logging.info(f"action: end of transmission | result: success")
                    break

                batchBets = self.message_handler.process_batch(client_connection)
                if batchBets:
                    bets.extend(batchBets)

                else:
                    raise Exception("An Error occured proccesing bets")

            except Exception as e:
                self.confirmation_to_client(client_connection, False)
                logging.error(f"action: process_bets | result: fail | error: {e}")
                logging.critical(
                    f"action: apuesta_recibida | result: fail | cantidad: {len(bets)}"
                )
                break
        self.store_all_bets(bets)
        self.confirmation_to_client(client_connection, True)
        logging.info(
            f"action: apuesta_recibida | result: success | cantidad: {len(bets)}"
        )
        return len(bets)

    def store_all_bets(self, bets: List[Bet]) -> None:
        """
        Store all bets in the database
        """
        store_bets(bets)
        logging.info(
            f"action: store_bets | result: success | amount of bets stored: {len(bets)}"
        )

    def confirmation_to_client(
        self, connection: ConnectionInterface, status: bool
    ) -> None:
        """
        Confirm the bet with the client
        """
        try:
            if status:
                send_status = connection.send(SUCCESS_HEADER)
                logging.info(f"action: confirm_bet | result: success")
            else:
                connection.send(FAIL)
        except Exception as e:
            logging.error(f"action: confirm_bet | result: fail | error: {e}")
