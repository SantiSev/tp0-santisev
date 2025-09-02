from typing import List
from server.common.session.client_session import ClientSession
from tp0.server.common.business.lottery_service import LotteryService
from tp0.server.common.network.connection_interface import ConnectionInterface

class ClientManager:
    def __init__(self, lottery_service: LotteryService):
        self.connected_clients: List[ClientSession] = []
        self.lottery_service = lottery_service

    def add_client(self, connection: ConnectionInterface) -> ClientSession:
        client_id = len(self.connected_clients) + 1
        client = ClientSession(client_id, connection, self.lottery_service)
        self.connected_clients.append(client)
        return client

    def remove_client(self, client: ClientSession):
        client.finish()
        self.connected_clients.remove(client)

    def shutdown(self) -> None:
        for client in self.connected_clients:
            client.finish()
        self.connected_clients.clear()
