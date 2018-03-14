from os import path
from sys import stderr
import logging

formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')

file_handler = logging.FileHandler(
    path.join(path.dirname(__file__), 'bot.log'))
file_handler.setLevel(logging.DEBUG)
file_handler.setFormatter(formatter)

stream_handler = logging.StreamHandler(stream=stderr)
stream_handler.setLevel(logging.DEBUG)
stream_handler.setFormatter(formatter)

logger = logging.Logger('EVE')
logger.addHandler(file_handler)
logger.addHandler(stream_handler)
