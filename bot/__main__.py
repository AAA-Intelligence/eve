from sys import argv
from .request_handler import run_loop
from .pattern_recognizer import demo
from .logger import logger

if len(argv) > 1 and argv[1] == 'pattern':
    logger.info('Running pattern recognizer demo')
    demo()
else:
    run_loop()
