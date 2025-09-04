import logging
import struct

from common.network.connection_interface import ConnectionInterface
from common.protocol.bet_parser import BetParser
from common.utils.utils import Bet
from common.protocol.protocol_constants import *


class AgencyHandler:
    """Handles individual client connections"""

    def __init__(self):
        self.bet_parser = BetParser()

    def get_bets(
        self, client_connection: ConnectionInterface
    ) -> tuple[list[Bet], bool]:

        header = client_connection.receive(HEADER_SIZE)

        if header == EOF:
            logging.info("action: end of transmission | result: success")
            return [], False

        if header != BET_HEADER:
            raise Exception(f"Unexpected header: {header}")

        batchBets = self.bet_parser.parse_batch(client_connection)
        return batchBets, True

    def confirm_batch(self, connection: ConnectionInterface, status: bool) -> None:
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
            length = len(winners_string)
            winners_bytes = length.to_bytes(DATA_LENGTH_SIZE, byteorder="big")
            logging.debug("action: sending_winners_data | result: success | length: %d", len(winners_string))
            connection.send(winners_bytes)
            connection.send(winners_string.encode())
            logging.debug(
                f"action: sending_winners | result: success | winners: [{winners_string}]"
            )
        except Exception as e:
            logging.error(f"action: sending_winners | result: fail | error: {e}")
