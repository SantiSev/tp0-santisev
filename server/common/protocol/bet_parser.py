import logging
from typing import List

from common.network.connection_interface import ConnectionInterface
from common.utils.utils import Bet
from common.protocol.protocol_constants import *


class BetParser:
    """Process bet data from clients"""

    def parse_batch(self, client_connection: ConnectionInterface) -> List[Bet]:
        try:
            data_length_bytes = client_connection.receive(DATA_LENGTH_SIZE)

            if data_length_bytes == 0:
                return []

            data_length = int.from_bytes(data_length_bytes, "big")
            data = client_connection.receive(data_length).decode("utf-8")

            return self._parse_batch_data(data)

        except Exception as e:
            logging.error(f"action: parse_batch | result: fail | error: {e}")
            return []

    def _parse_batch_data(self, data: str) -> List[Bet]:
        """Parse comma-separated data into list of Bet objects"""
        try:
            data = data.rstrip(",")
            fields = [field.strip() for field in data.split(",")]
            fields = [field for field in fields if field]

            bets = []

            for i in range(0, len(fields), EXPECTED_FIELDS):
                if i + EXPECTED_FIELDS <= len(fields):
                    agency = fields[i]
                    first_name = fields[i + 1]
                    last_name = fields[i + 2]
                    document = fields[i + 3]
                    birthdate = fields[i + 4]
                    number = fields[i + 5]

                    bet = Bet(
                        agency=agency,
                        first_name=first_name,
                        last_name=last_name,
                        document=document,
                        birthdate=birthdate,
                        number=number,
                    )
                    bets.append(bet)
                else:
                    remaining_fields = len(fields) - i
                    logging.warning(
                        f"action: parse_bet | result: incomplete_data | remaining_fields: {remaining_fields} | expected: {EXPECTED_FIELDS}"
                    )

            logging.info(
                f"action: parse_bet_batch | result: success | total_bets: {len(bets)}"
            )
            return bets

        except Exception as e:
            logging.error(f"action: parse_bet_batch | result: fail | error: {e}")
            return []
