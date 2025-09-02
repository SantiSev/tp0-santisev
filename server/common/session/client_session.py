import logging
from tp0.server.common.business.lottery_service import LotteryService
from tp0.server.common.network.connection_interface import ConnectionInterface
from tp0.server.common.protocol.bet_handler import BetHandler
from tp0.server.common.utils.utils import Bet


class ClientSession:

    def __init__(self, client_id: int, connection_interface: ConnectionInterface, lottery_service: LotteryService):
        self.id = client_id
        self.connection_interface = connection_interface
        self.lottery_service = lottery_service
        self.protocol_handler = BetHandler()

    def begin(self) -> tuple[bool, list[Bet]]:

        try:
            bets: list[Bet] = self.protocol_handler.process_bets(
                self.connection_interface
            )
            winners: list[str] = self.lottery_service.get_winners(bets)
            self.protocol_handler.send_winners(self.connection_interface, winners)
            return True, bets
        except Exception as e:
            logging.error(f"action: client_session | result: fail | error: {e}")
            return False, []

    def finish(self) -> None:
        self.connection_interface.close()
