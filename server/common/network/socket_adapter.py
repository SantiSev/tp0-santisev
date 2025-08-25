import logging
import socket

from typing import Optional

BUFFER_SIZE = 512

class SocketAdapter:

    def __init__(self, socket: socket.socket, buffer_size: int = BUFFER_SIZE):
        self.socket = socket
        self.buffer_size = buffer_size

    def send(self, data: bytes) -> bool:
        try:
            self.socket.sendall(data)
            return True
        except Exception as e:
            logging.error(f"action: send | result: fail | error: {e}")
            return False

    def receive(self) -> Optional[bytes]:
        try:
            data = self.socket.recv(self.buffer_size)
            if not data:
                logging.debug("action: receive | result: no_data")
                return None
            return data
        except Exception as e:
            logging.error(f"action: receive | result: fail | section: SocketAdapter | error: {e}")
            return None

    def close(self) -> None:
        try:
            self.socket.close()
            logging.info("action: close | result: success")
        except Exception as e:
            logging.error(f"action: close | result: fail | error: {e}")
