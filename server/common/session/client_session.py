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

    def begin(self) -> bool:
        agencyBets = []
        while True:
            try:
                betBatch, more_bets_remaining = self.protocol_handler.get_bets(
                    self.connection_interface
                )

                if not more_bets_remaining:
                    logging.info(f"action: all_bets_received | result: success")
                    break

                self.lottery_service.place_bets(betBatch)
                self.protocol_handler.confirm_batch(self.connection_interface, True)
                agencyBets.extend(betBatch)

                logging.info(f"more_bets_remaining: {more_bets_remaining}")

            except Exception as e:
                logging.error(
                    f"action: apuesta_recibida | result: fail | cantidad: {len(agencyBets)}"
                )
                self.protocol_handler.confirm_batch(self.connection_interface, False)
                return False
        logging.info(
            f"action: apuesta_recibida | result: success | cantidad: {len(agencyBets)}"
        )
        winners = self.lottery_service.draw_winners(agencyBets)
        self.protocol_handler.send_winners(self.connection_interface, winners)
        return True

    def finish(self) -> None:
        self.connection_interface.close()
