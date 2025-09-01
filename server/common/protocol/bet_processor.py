import logging
from typing import List, Optional

from common.network.connection_interface import ConnectionInterface
from common.utils.utils import Bet

EXPECTED_FIELDS = 6
DATA_LENGTH = 4


class BetProcessor:
    """Process bet data from clients"""

    def process_batch(self, client_connection: ConnectionInterface) -> List[Bet]:
        try:
            # TODO: HAVE FIRST BYTE AFTER HEADER BE THE DATA_LENGTH
            data_length = int.from_bytes(client_connection.receive(DATA_LENGTH), "big")
            data = client_connection.receive(data_length).decode("utf-8")

            logging.debug(f"action: receive_message | result: success | data: {data}")

            if not data:
                logging.debug("action: receive_message | result: no_data")
                raise Exception(
                    "Something went wrong when receiving bet data, stopping process . . ."
                )

            return self._parse_batch_data(data)

        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            return []

    def _parse_batch_data(self, data: str) -> List[Bet]:
        """Parse comma-separated data into list of Bet objects"""
        try:
            fields = [field.strip() for field in data.split(",")]

            bets = []

            # Process every 6 fields as one bet
            for i in range(0, len(fields), EXPECTED_FIELDS):
                # Check if we have enough fields for a complete bet
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
