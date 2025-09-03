import logging
from common.utils.utils import Bet, store_bets


class LotteryService:

    def place_bets(self, bets: list[Bet]) -> None:
        store_bets(bets)
