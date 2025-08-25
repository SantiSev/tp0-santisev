import logging
import socket
from typing import Optional

from common.network.socket_adapter import SocketAdapter
from common.utils.utils import Bet
from common.utils.message_types import BetConfirmation


class MessageHandler:
    """Handles message parsing and validation"""

    EXPECTED_FIELDS = 6

    def process_message(self, socket: SocketAdapter) -> Optional[Bet]:
        try:
            data = socket.receive()
            if not data:
                logging.debug("action: receive_message | result: no_data")
                return None

            return self._parse_bet_data(data.decode("utf-8"))

        except UnicodeDecodeError as e:
            logging.error(f"action: decode_message | result: fail | error: {e}")
            return None
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            return None

    def confirmation_to_client(self, socket: SocketAdapter, status: bool) -> None:
        """
        Confirm the bet with the client
        """

        betConfirmation = BetConfirmation(status)

        try:
            socket.send(betConfirmation.encode())
        except Exception as e:
            logging.error(f"action: confirm_bet | result: fail | error: {e}")

    def _parse_bet_data(self, raw_message: str) -> Optional[Bet]:
        """Parse raw message string into Bet object"""
        try:
            message = raw_message.strip()
            if not message:
                logging.warning(
                    "action: parse_bet | result: fail | reason: empty_message"
                )
                raise ValueError("Empty message")

            fields = message.split(",")
            if len(fields) != self.EXPECTED_FIELDS:
                logging.warning(
                    f"action: parse_bet | result: fail | reason: invalid_field_count | expected: {self.EXPECTED_FIELDS} | got: {len(fields)}"
                )
                raise ValueError("Invalid field count")

            agency, first_name, last_name, document, birthdate, number = [
                field.strip() for field in fields
            ]

            if not all([agency, first_name, last_name, document, birthdate, number]):
                logging.warning(
                    "action: parse_bet | result: fail | reason: empty_fields"
                )
                raise ValueError("One or more required fields are empty")

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
