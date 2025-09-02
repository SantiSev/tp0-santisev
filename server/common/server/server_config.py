class ServerConfig:
    def __init__(
        self, port: int, listen_backlog: int, agencies_amount: int, logging_level: str
    ):
        self.port = port
        self.listen_backlog = listen_backlog
        self.agencies_amount = agencies_amount
        self.logging_level = logging_level

    def __repr__(self):
        return (
            f"port: {self.port} | "
            f"listen_backlog: {self.listen_backlog} | logging_level: {self.logging_level} | agencies_amount: {self.agencies_amount}"
        )
