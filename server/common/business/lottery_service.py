import logging
from common.utils.utils import Bet, store_bets


class LotteryService:

    def place_bets(self, bets: list[Bet]) -> None:
        logging.info(
            f"action: apuesta_almacenada | result: success | dni: {bets[0].document} | numero: {bets[0].number}"
        )
        store_bets(bets)
