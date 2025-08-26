from typing import Optional


class BetConfirmation:
    def __init__(self, success: bool, error: Optional[str] = None):
        self.success = success
        self.error = error


    def encode(self) -> bytes:
        status = "success" if self.success else "error"
        if self.error:
            message = f"BET_CONFIRMATION: {status},\n{self.error}"
        else:
            message = f"BET_CONFIRMATION: {status}"

        # Pad or truncate to exactly 256 bytes
        encoded = message.encode("utf-8")

        if len(encoded) > 256:
            # Truncate if too long
            return encoded[:256]
        else:
            # Pad with null bytes if too short
            return encoded + b"\x00" * (256 - len(encoded))
