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
        self.agency_id = agency_id
        self.connection_interface = connection_interface
        self.lottery_service = lottery_service
        self.protocol_handler = AgencyHandler()

    def begin(self) -> bool:
        while True:
            try:
                betBatch, more_bets_remaining = self.protocol_handler.get_bets(
                    self.connection_interface, self.agency_id
                )

                if not more_bets_remaining:
                    logging.info(f"action: all_bets_received | result: success")
                    break

                self.lottery_service.place_bets(betBatch)
                self.protocol_handler.confirm_batch(self.connection_interface, True)

                logging.info(f"more_bets_remaining: {more_bets_remaining}")

            except Exception as e:
                agencyBets = self.lottery_service.get_bets_by_agency(self.agency_id)
                logging.error(
                    f"action: apuesta_recibida | result: fail | cantidad: {len(agencyBets)} | error: {e}"
                )
                self.protocol_handler.confirm_batch(self.connection_interface, False)
                return False
        agencyBets = self.lottery_service.get_bets_by_agency(self.agency_id)
        logging.info(
            f"action: apuesta_recibida | result: success | cantidad: {len(agencyBets)}"
        )
        return True

    def tally_results(self):
        """Tally and log the results of the lottery"""
        winners = self.lottery_service.draw_winners(self.agency_id)
        self.protocol_handler.send_winners(self.connection_interface, winners)

    def finish(self) -> None:
        self.connection_interface.close()
