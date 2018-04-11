from sys import argv

import nltk

from bot.logger import logger
from bot.model_definitions import Mode

nltk.download('punkt', quiet=True)

target = argv[1] if len(argv) > 1 else None

if target == 'train-patterns':
    from bot.trainer import train_model

    logger.info('Running pattern training')
    train_model(Mode.PATTERNS)
elif target == 'train-sentiments':
    from bot.pattern_recognizer import demo
    logger.info(
        'Running sentiments analysis through moods and affections analysis training')
    train_model(Mode.AFFECTIONS)
    train_model(Mode.MOODS)
elif target == 'train-chat':
    from bot.text_processor.train import train_and_evaluate

    logger.info('Running chat training')
    train_and_evaluate()
elif target == 'demo':
    from bot.request_handler import run_demo

    run_demo()
else:
    from bot.request_handler import run_loop

    run_loop()
