import logging
from typing import List

from common.network.connection_interface import ConnectionInterface
from common.utils.utils import Bet
from common.protocol.protocol_constants import *


class BetParser:
    """Process bet data from clients"""

    def parse_bet(self, client_connection: ConnectionInterface) -> List[Bet]:
        try:
            data_length_bytes = client_connection.receive(DATA_LENGTH_SIZE)
            data_length = int.from_bytes(data_length_bytes, "big")

            if data_length == 0:
                return []

            logging.debug(
                f"action: parse_batch | result: success | data_length: {data_length}"
            )

            data = client_connection.receive(data_length).decode("utf-8")
            logging.info(f"action: parse_bet | result: success | data: {data}")
            return self._parse_bet_data(data)

        except Exception as e:
            logging.warning(f"action: parse_bet | result: fail | error: {e}")
            return []

    def _parse_bet_data(self, data: str) -> List[Bet]:
        """Parse comma-separated data into list of Bet objects"""
        try:
            data = data.rstrip(",")
            fields = data.split(",")
            agency = fields[0]
            first_name = fields[1]
            last_name = fields[2]
            document = fields[3]
            birthdate = fields[4]
            number = fields[5]

            bet = Bet(
                agency=agency,
                first_name=first_name,
                last_name=last_name,
                document=document,
                birthdate=birthdate,
                number=number,
            )

            return [bet]

        except Exception as e:
            logging.error(f"action: parse_bet_batch | result: fail | error: {e}")
            return []
