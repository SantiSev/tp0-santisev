import socket

from common.network.connection_interface import ConnectionInterface


class ConnectionManager:
    def __init__(self, port: int, listen_backlog: int):
        self.port = port
        self.backlog = listen_backlog
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    def accept_connection(self) -> ConnectionInterface:
        client_socket, _ = self.socket.accept()
        return ConnectionInterface(client_socket)

    def start_listening(self) -> None:
        self.socket.bind(("", self.port))
        self.socket.listen(self.backlog)

    def shutdown(self) -> None:
        self.socket.close()
