import logging
import struct

from common.network.connection_interface import ConnectionInterface
from common.protocol.bet_parser import BetParser
from common.utils.utils import Bet, has_won, load_bets, store_bets
from common.protocol.protocol_constants import *


class BetHandler:
    """Handles individual client connections"""

    def __init__(self):
        self.bet_parser = BetParser()

    def process_bets(self, client_connection: ConnectionInterface) -> list[Bet]:

        bets = []

        while True:
            try:
                header = client_connection.receive(HEADER_SIZE)

                if header == EOF:
                    logging.info(f"action: end of transmission | result: success")
                    break

                if header != BET_HEADER:
                    logging.warning(
                        f"action: process_bets | result: fail | error: unexpected_header | header: {header}"
                    )
                    raise Exception("Unexpected header received")

                batchBets = self.bet_parser.parse_batch(client_connection)
                if batchBets:
                    bets.extend(batchBets)
                    self._confirmation_to_client(client_connection, True)

                else:
                    raise Exception("An Error occured proccesing bets")

            except Exception as e:
                self._confirmation_to_client(client_connection, False)
                logging.error(f"action: process_bets | result: fail | error: {e}")
                logging.critical(
                    f"action: apuesta_recibida | result: fail | cantidad: {len(bets)}"
                )
                break
        store_bets(bets)
        logging.info(
            f"action: apuesta_recibida | result: success | cantidad: {len(bets)}"
        )
        return bets

    def _confirmation_to_client(
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

    def send_winners(self, connection: ConnectionInterface, winners: list[str]) -> None:
        """
        Send the winning numbers to the client
        """
        winners_string = ",".join(winners)

        try:
            connection.send(WINNERS_HEADER)
            winners_bytes = struct.pack(">H", len(winners_string))
            connection.send(winners_bytes)
            logging.debug(
                f"action: sending_winners | result: success | winners: {winners_string}"
            )
        except Exception as e:
            logging.error(f"action: sending_winners | result: fail | error: {e}")
