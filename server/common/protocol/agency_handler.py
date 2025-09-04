import logging

from common.network.connection_interface import ConnectionInterface
from common.protocol.bet_parser import BetParser
from common.utils.utils import Bet
from common.protocol.protocol_constants import *


class AgencyHandler:
    """Handles individual client connections"""

    def __init__(self):
        self.bet_parser = BetParser()

    def get_bets(self, connection: ConnectionInterface) -> list[Bet]:

        self._process_header(connection)

        bet = self.bet_parser.parse_bet(connection)
        return bet

    def _process_header(self, client_connection: ConnectionInterface):

        headerData = client_connection.receive(HEADER_SIZE)

        header_byte = headerData[0:1]

        if header_byte != BET_HEADER:
            raise Exception(f"Unexpected header: {header_byte}, expected: {BET_HEADER}")

        logging.debug(f"action: process_header | result: success")

    def confirm_bet(
        self, bets: list[Bet], connection: ConnectionInterface, status: bool
    ) -> None:
        """
        Confirm the bet with the client
        """
        try:
            if status:
                connection.send(SUCCESS_HEADER)

                bet = bets[0]
                bet_message = f"{bet.document},{bet.number},"
                message_bytes = bet_message.encode()
                message_length = len(message_bytes)
                length_bytes = bytes([message_length])

                connection.send(length_bytes)
                connection.send(message_bytes)

                logging.debug(f"action: confirm_bet | result: success")
            else:
                connection.send(FAIL)
                logging.debug(f"action: confirm_bet | result: fail")
        except Exception as e:
            logging.error(f"action: confirm_bet | result: fail | error: {e}")
