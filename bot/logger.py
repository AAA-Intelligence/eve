import logging
from os import path
from sys import stderr

# setting up the logger for the entire

# determines the format of a log entry
formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')

# setting up the file to which the logs are written
file_handler = logging.FileHandler(
    path.join(path.dirname(__file__), 'bot.log'))

# only debug and more severe logs are written to the log file
file_handler.setLevel(logging.DEBUG)

# adding the formatter to the logger
file_handler.setFormatter(formatter)

# adding console output to the logger for debug purposes in the console demo
stream_handler = logging.StreamHandler(stream=stderr)
stream_handler.setLevel(logging.DEBUG)
stream_handler.setFormatter(formatter)

logger = logging.Logger('EVE')
logger.addHandler(file_handler)
logger.addHandler(stream_handler)
