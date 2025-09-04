import logging
from typing import List

from common.network.connection_interface import ConnectionInterface
from common.utils.utils import Bet
from common.protocol.protocol_constants import *


class BetParser:
    """Process bet data from clients"""

    def parse_batch(
        self, client_connection: ConnectionInterface, agency_id: int
    ) -> List[Bet]:
        try:
            data_length_bytes = client_connection.receive(DATA_LENGTH_SIZE)
            data_length = int.from_bytes(data_length_bytes, "big")

            if data_length == 0:
                return []

            logging.debug(
                f"action: parse_batch | result: success | data_length: {data_length}"
            )

            data = client_connection.receive(data_length).decode("utf-8")
            return self._parse_batch_data(data, agency_id)

        except Exception as e:
            logging.warning(f"action: parse_batch | result: fail | error: {e}")
            return []

    def _parse_batch_data(self, data: str, agency_id: int) -> List[Bet]:
        """Parse comma-separated data into list of Bet objects"""
        try:
            data = data.rstrip(",")
            fields = [field.strip() for field in data.split(",")]
            fields = [field for field in fields if field]

            bets = []

            for i in range(0, len(fields), EXPECTED_FIELDS):
                if i + EXPECTED_FIELDS <= len(fields):
                    first_name = fields[i + 0]
                    last_name = fields[i + 1]
                    document = fields[i + 2]
                    birthdate = fields[i + 3]
                    number = fields[i + 4]
                    bet = Bet(
                        agency=str(agency_id),
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

            logging.debug(
                f"action: parse_bet_batch | result: success | total_bets: {len(bets)}"
            )
            return bets

        except Exception as e:
            logging.error(f"action: parse_bet_batch | result: fail | error: {e}")
            return []
