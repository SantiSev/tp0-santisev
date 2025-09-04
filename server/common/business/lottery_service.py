import logging
from multiprocessing import Lock
from common.utils.utils import Bet, has_won, load_bets, store_bets


class LotteryService:

    def __init__(self, lock: Lock):
        self.file_lock = lock

    def place_bets(self, bets: list[Bet]) -> None:
        with self.file_lock:
            store_bets(bets)

    def draw_winners(self, agency_id: int) -> list[str]:
        with self.file_lock:
            bets = load_bets()
        winners = [
            bet.document for bet in bets if has_won(bet) and bet.agency == agency_id
        ]
        return winners

    def get_bets_by_agency(self, agency_id: int) -> list[Bet]:
        with self.file_lock:
            bets = load_bets()
        return [bet for bet in bets if bet.agency == agency_id]

    def announce_winners(self) -> None:
        with self.file_lock:
            bets = load_bets()
        winners = [bet.document for bet in bets if has_won(bet)]
        logging.info(f"action: sorteo | result: success | winners: {winners}")
