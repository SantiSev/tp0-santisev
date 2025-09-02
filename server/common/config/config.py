from configparser import ConfigParser
import logging
import os

from tp0.server.common.server.server_config import ServerConfig


def initialize_config() -> ServerConfig:
    """Parse env variables or config file to find program config params

    Function that search and parse program configuration parameters in the
    program environment variables first and the in a config file.
    If at least one of the config parameters is not found a KeyError exception
    is thrown. If a parameter could not be parsed, a ValueError is thrown.
    If parsing succeeded, the function returns a ConfigParser object
    with config parameters
    """

    config = ConfigParser(os.environ)
    # If config.ini does not exists original config object is not modified
    config.read("config.ini")

    try:
        port = int(os.getenv("SERVER_PORT", config["DEFAULT"]["SERVER_PORT"]))
        listen_backlog = int(
            os.getenv(
                "SERVER_LISTEN_BACKLOG", config["DEFAULT"]["SERVER_LISTEN_BACKLOG"]
            )
        )
        logging_level = os.getenv("LOGGING_LEVEL", config["DEFAULT"]["LOGGING_LEVEL"])
        agencies_amount = int(
            os.getenv("AGENCIES_AMOUNT", config["DEFAULT"]["AGENCIES_AMOUNT"])
        )
    except KeyError as e:
        raise KeyError("Key was not found. Error: {} .Aborting server".format(e))
    except ValueError as e:
        raise ValueError(
            "Key could not be parsed. Error: {}. Aborting server".format(e)
        )

    return ServerConfig(
        port=port,
        listen_backlog=listen_backlog,
        agencies_amount=agencies_amount,
        logging_level=logging_level,
    )


def initialize_log(logging_level):
    """
    Python custom logging initialization

    Current timestamp is added to be able to identify in docker
    compose logs the date when the log has arrived
    """
    logging.basicConfig(
        format="%(asctime)s %(levelname)-8s %(message)s",
        level=logging_level,
        datefmt="%Y-%m-%d %H:%M:%S",
    )
