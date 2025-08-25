from typing import Optional


class BetConfirmation:
    def __init__(self, success: bool, error: Optional[str] = None):
        self.success = success
        self.error = error

    def encode(self) -> bytes:
        status = "success" if self.success else "error"
        if self.error:
            return f"BET_CONFIRMATION: {status},\n{self.error}".encode("utf-8")
        return f"BET_CONFIRMATION: {status}".encode("utf-8")
