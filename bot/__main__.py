from sys import argv

import nltk

from bot.logger import logger
from bot.trainer import train_model

nltk.download('punkt', quiet=True)

target = argv[1] if len(argv) > 1 else None

if target == 'train-patterns':

	logger.info('Running pattern training')
	train_model("patterns")
elif target == 'train-sentiments':
	from bot import trainer

	logger.info('Running sentiment analysis training')
	trainer.train_model("sentiment")
elif target == 'console-demo':
	from bot.pattern_recognizer import demo

	logger.info('Running pattern recognizer demo')
	demo("%s" % argv[2])
elif target == 'demo':
	from bot.request_handler import run_demo

	run_demo()
else:
	from .request_handler import run_loop

	run_loop()
