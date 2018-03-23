from sys import argv
from .logger import logger

target = argv[1] if len(argv) > 1 else None

if target == 'train-patterns':
    from .train_patterns_model import train_model
    logger.info('Running pattern training')
    train_model()
elif target == 'pattern-demo':
    from .pattern_recognizer import demo
    logger.info('Running pattern recognizer demo')
    demo()
else:
    from .request_handler import run_loop
    run_loop()
