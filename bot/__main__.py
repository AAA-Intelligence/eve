import datetime
import json
import logging
import sys

from setup import setup_logger

now_time = datetime.datetime.now().time()
# message as json
text_message = """{"content": "Hallo, Apache!","timestamp": "%d:%d:%d", "sender": "Bot", "context": {"mood": -2, "affection": -3},"botId": 123, "userId": 110 }""" % (
	now_time.hour, now_time.minute, now_time.second)

logger: logging.Logger = setup_logger()
a = True
while True:
	try:
		i = input()
		data = json.loads(i)
		log_message = "new Message:%s" % data["content"]
	except:
		print("Unexpected error:", sys.exc_info()[0])
		break
