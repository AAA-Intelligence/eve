import datetime
import json
import time

from numpy.matlib import random

now_time = datetime.datetime.now().time()
names = ["Simon", "Daniel", "der andere Daniel", "Leon", "Niklas"]
i = random.randint(0, len(names))
old_i = -1
# nachricht

# while True:
with open('messages/usr_msg.json', 'r') as f:
	str_file = str(f.read())
	msg_json = json.loads(str_file)
	old_msg = msg_json["content"]

message = {
	"content": "Hallo, %s!" % names[i], "timestamp": (
			"%d:%d:%d" % (
		now_time.hour, now_time.minute, now_time.second)),
	"sender": "Bot",
	"context": {"mood": -2, "affection": -3},
	"botId": 123, "userId": 110, }

while True:
	try:
		with open('messages/usr_msg.json', 'r') as f:
			text = f.read()
			d = json.loads(str(text))
			msg = d["content"]
			if old_msg != msg:
				old_msg = msg
				# print(msg)
				while i == old_i:
					i = random.randint(0, len(names))
				old_i = i
				now_time = datetime.datetime.now().time()
				with open("messages/bot_msg.json", 'w') as file:
					message["content"] = "Hallo, %s!" % names[i]
					print(names[i])
					file.write(str(message).replace("'", "\""))

	except(PermissionError, FileNotFoundError):
		print("wird gespeichert...")
		time.sleep(.5)
