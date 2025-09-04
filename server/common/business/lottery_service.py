import logging
from common.utils.utils import Bet, load_bets, store_bets


class LotteryService:

    def place_bets(self, bets: list[Bet]) -> None:
        store_bets(bets)

    def get_bets_by_agency(self, agency_id: int) -> list[Bet]:
        bets = load_bets()
        return [bet for bet in bets if bet.agency == agency_id]

