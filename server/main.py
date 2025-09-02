#!/usr/bin/env python3

from common.server.server import Server
import logging

from common.config.config import initialize_config, initialize_log
from common.server.server_config import ServerConfig


# to run locally, cd to this directory and run: python3 -m server.main
def main():
    server_config: ServerConfig = initialize_config()
    initialize_log(server_config.logging_level)
    logging.debug(f"action: config | result: success | {server_config}")
    server = Server(server_config)
    server.run()


if __name__ == "__main__":
    main()
