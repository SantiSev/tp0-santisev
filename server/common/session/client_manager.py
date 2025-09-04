from typing import List
from common.session.client_session import ClientSession
from common.business.lottery_service import LotteryService
from common.network.connection_interface import ConnectionInterface

class ClientManager:
    def __init__(self, lottery_service: LotteryService):
        self.connected_clients: dict[int, ClientSession] = {}
        self.lottery_service = lottery_service

    def add_client(self, connection: ConnectionInterface) -> ClientSession:
        client_id = len(self.connected_clients) + 1
        client = ClientSession(client_id, connection, self.lottery_service)
        self.connected_clients[client_id] = client
        return client

    def remove_client(self, client_id: int):
        client = self.connected_clients.get(client_id)
        if client:
            client.finish()
        del self.connected_clients[client_id]

    def shutdown(self) -> None:
        for client in self.connected_clients.values():
            client.finish()
        self.connected_clients.clear()

    def send_results_to_all(self) -> None:
        for client in self.connected_clients.values():
            client.send_results()
