import logging
from common.business.lottery_service import LotteryService
from common.network.connection_interface import ConnectionInterface
from common.protocol.bet_handler import BetHandler
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
        self.protocol_handler = BetHandler()

    def begin(self) -> tuple[bool, list[Bet]]:
        agencyBets = []
        while True:
            try:
                betBatch, more_bets_remaining = self.protocol_handler.get_bets(
                    self.connection_interface
                )
                self.lottery_service.place_bets(betBatch)
                self.protocol_handler.confirm_batch(self.connection_interface, True)
                agencyBets.extend(betBatch)

                if more_bets_remaining is False:
                    logging.info(f"action: all_bets_received | result: success")
                    break

            except Exception as e:
                logging.error(f"action: client_session | result: fail | error: {e}")
                self.protocol_handler.confirm_batch(self.connection_interface, False)
                return False, []

        winners = self.lottery_service.draw_winners(agencyBets)
        self.protocol_handler.send_winners(self.connection_interface, winners)

    def finish(self) -> None:
        self.connection_interface.close()
