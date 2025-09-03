import logging
from common.utils.utils import Bet, has_won, load_bets, store_bets


class LotteryService:

    def place_bets(self, bets: list[Bet]) -> None:
        store_bets(bets)

    def draw_winners(self, bets: list[Bet]) -> list[str]:
        winners = [bet.document for bet in bets if has_won(bet)]
        return winners

    def announce_winners(self) -> None:
        bets = load_bets()
        winners = self.draw_winners(bets)
        logging.info(f"action: sorteo | result: success | winners: {winners} | cantidad de ganadores: {len(winners)} / {len(bets)}")
