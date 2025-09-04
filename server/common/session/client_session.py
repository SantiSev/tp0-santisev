import logging
from common.business.lottery_service import LotteryService
from common.network.connection_interface import ConnectionInterface
from common.protocol.agency_handler import AgencyHandler
from common.utils.utils import Bet


class ClientSession:

    def __init__(
        self,
        agency_id: int,
        connection_interface: ConnectionInterface,
        lottery_service: LotteryService,
    ):
        self.id = agency_id
        self.connection_interface = connection_interface
        self.lottery_service = lottery_service
        self.protocol_handler = AgencyHandler()

    def begin(self) -> None:
        try:
            bets = self.protocol_handler.get_bets(
                self.connection_interface
            )

            self.lottery_service.place_bets(bets)
            self.protocol_handler.confirm_bet(bets, self.connection_interface, True)


        except Exception as e:
            logging.error(f"action: client_session | result: fail | error: {e}")
            self.protocol_handler.confirm_bet(bets, self.connection_interface, False)
            return

    def finish(self) -> None:
        self.connection_interface.close()
