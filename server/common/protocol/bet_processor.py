import logging
from typing import Optional

from common.network.connection_interface import ConnectionInterface
from common.utils.utils import Bet

BET_DATA_SIZE = 256
EXPECTED_FIELDS = 6


class BetProcessor:
    """Process bet data from clients"""

    def process_bet(self, client_connection: ConnectionInterface) -> Optional[Bet]:
        try:
            data = client_connection.receive(BET_DATA_SIZE)

            if not data:
                logging.debug("action: receive_message | result: no_data")
                return None

            logging.debug(
                f"action: receive_message | result: success | data: {data.decode('utf-8').rstrip('\x00')}"
            )

            return self._parse_bet_data(data.decode("utf-8").rstrip("\x00"))

        except UnicodeDecodeError as e:
            logging.error(f"action: decode_message | result: fail | error: {e}")
            return None
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            return None

    def _parse_bet_data(self, data: str) -> Optional[Bet]:
        """Parse raw message string into Bet object"""
        try:
            message = data.strip()

            agency, first_name, last_name, document, birthdate, number = [
                field.strip() for field in message.split(",")
            ]

            return Bet(
                agency=agency,
                first_name=first_name,
                last_name=last_name,
                document=document,
                birthdate=birthdate,
                number=number,
            )

        except ValueError as e:
            logging.error(
                f"action: parse_bet | result: fail | reason: parsing_error | error: {e}"
            )
            return None
        except Exception as e:
            logging.error(
                f"action: parse_bet | result: fail | reason: unexpected_error | error: {e}"
            )
            return None
