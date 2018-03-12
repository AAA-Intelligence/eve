import logging


def setup_logger():
	logger = logging.getLogger(__name__)
	# default level: NOTSET (lowest level)
	fh = logging.FileHandler('bot/bot.log')

	ch = logging.StreamHandler()
	ch.setLevel(logging.ERROR)

	formatter = logging.Formatter(
		'%(asctime)s - %(levelname)s - %(message)s')
	ch.setFormatter(formatter)
	fh.setFormatter(formatter)

	logger.addHandler(ch)
	logger.addHandler(fh)

	# levels below warning wont work
	logger.warning(
		"logging is setup. error's will be logged to console as well")
	return logger


if __name__ == '__main__':
	# test1.py executed as script
	# do something
	setup_logger()
