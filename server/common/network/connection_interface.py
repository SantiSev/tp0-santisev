import logging
import socket

from typing import Optional


class ConnectionInterface:

    def __init__(self, socket: socket.socket):
        self.socket = socket

    def send(self, data: bytes) -> bool:
        try:
            self.socket.sendall(data)
            return True
        except Exception as e:
            logging.error(f"action: send | result: fail | error: {e}")
            return False

    def receive(self, buffer_size: int) -> Optional[bytes]:
        try:
            if buffer_size < 0:
                logging.error(
                    f"action: receive | result: fail | error: invalid_buffer_size"
                )
                return None
            data = self._receive_all(buffer_size)
            if not data:
                return None
            return data

        except Exception as e:
            logging.error(
                f"action: receive | result: fail | section: ConnectionInterface | error: {e}"
            )
            return None

    def _receive_all(self, buffer_size: int) -> Optional[bytes]:
        """Ensure we receive exactly 'buffer_size' bytes, handling short reads"""
        data = b""
        while len(data) < buffer_size:
            try:
                chunk = self.socket.recv(buffer_size - len(data))
                if not chunk:
                    return None
                data += chunk
            except socket.timeout:
                logging.warning("action: receive_all | result: timeout")
                return None
            except Exception as e:
                logging.error(f"action: receive_all | result: fail | error: {e}")
                return None
        return data

    def close(self) -> None:
        try:
            self.socket.close()
            logging.info("action: close_connection | result: success")
        except Exception as e:
            logging.error(f"action: close_connection | result: fail | error: {e}")
