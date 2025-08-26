import socket

from common.network.socket_adapter import SocketAdapter


class ConnectionManager:
    def __init__(self, port: int, listen_backlog: int):
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(("", port))
        self._server_socket.listen(listen_backlog)
        self._is_running = True

    def accept_connection(self) -> SocketAdapter:
        client_socket, _ = self._server_socket.accept()
        return SocketAdapter(client_socket)

    def start_listening(self) -> None:
        self._server_socket.listen()

    def is_running(self) -> bool:
        return self._is_running

    # TODO: the socketAdapter should be responsible for closing, the manager just gives you a socket to connect to shut itself down
    def shutdown(self) -> None:
        self._server_socket.close()
        self._is_running = False
