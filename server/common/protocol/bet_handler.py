import logging
import struct

from common.network.connection_interface import ConnectionInterface
from common.protocol.bet_parser import BetParser
from common.utils.utils import Bet
from common.protocol.protocol_constants import *


# TODO: DO THESE CHANGES:
# 1. REMOVE THE CLIENT MESSAGE LOOP, THE EXERCISE STATES TO SEND 1 BET, sending multiple times the same bet doesnt make sense
# 2. DONT PADD OUT THE MESSAGES, HAVE IT SEND FIRST THE LENGTH AND THEN THE DATA
# 3. ( Optional )Check to see if you can refactort the connectionINterface so you don't have to send the length and then data, have it send the length and then the data all at once

class BetHandler:
    """Handles individual client connections"""

    def __init__(self):
        self.bet_parser = BetParser()

    def get_bets(self, client_connection: ConnectionInterface) -> list[Bet]:

        header = client_connection.receive(HEADER_SIZE)

        if header != BET_HEADER:
            raise Exception(f"Unexpected header: {header}")

        bet = self.bet_parser.parse_bet(client_connection)
        return bet

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
