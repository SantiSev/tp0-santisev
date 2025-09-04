import logging
from common.utils.utils import Bet, has_won, load_bets, store_bets


class LotteryService:

    def place_bets(self, bets: list[Bet]) -> None:
        store_bets(bets)

    def draw_winners(self, agency_id: int) -> list[str]:
        bets = load_bets()
        winners = [
            bet.document for bet in bets if has_won(bet) and bet.agency == agency_id
        ]
        return winners

    def get_bets_by_agency(self, agency_id: int) -> list[Bet]:
        bets = load_bets()
        return [bet for bet in bets if bet.agency == agency_id]

    def announce_winners(self) -> None:
        bets = load_bets()
        winners = [bet.document for bet in bets if has_won(bet)]
        logging.info(f"action: sorteo | result: success | winners: {winners}")
