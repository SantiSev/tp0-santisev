#!/usr/bin/env python3

from common.server.server import Server
import logging

from common.config.config import initialize_config, initialize_log

# to run locally, cd to this directory and run: python3 -m server.main
def main():
    config_params = initialize_config()
    logging_level = config_params["logging_level"]
    port = config_params["port"]
    listen_backlog = config_params["listen_backlog"]
    amount_agencies = config_params["agencies_amount"]

    initialize_log(logging_level)

    # Log config parameters at the beginning of the program to verify the configuration
    # of the component
    logging.debug(f"action: config | result: success | port: {port} | "
                  f"listen_backlog: {listen_backlog} | logging_level: {logging_level} | agencies_amount: {amount_agencies}")

    # Initialize server and start server loop
    server = Server(port, listen_backlog, amount_agencies)
    server.run()




if __name__ == "__main__":
    main()
