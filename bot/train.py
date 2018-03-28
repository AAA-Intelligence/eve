import pickle
from enum import IntEnum
from os import path, mkdir
from typing import Tuple

from keras.models import Sequential, model_from_json

from bot.setup import setup_bot
from bot.training_data import TrainingData
from .model_definitions import Sentiment, Patterns


def setup_models_dir():
	global dir
	dir = path.join(path.dirname(__file__), 'models')
	if not path.exists(dir):
		mkdir(dir)
	if not path.isdir(dir):
		raise Exception('Models path is not a directory: {}'.format(dir))


def train_model(mode: str):
	# creates a directory where the trained models are stored
	setup_models_dir()

	"""
	Trains a neural network with the defined patterns and categories.
	Patterns will be split into words, stemmed by a German snowball stemmer and
	indexed by saving all stems in a list of total stems and assigning indices.
	The trained model will be saved in the models directory and can be loaded
	by the pattern recognizer using the load_model function.
	"""

	Mode = set_mode(mode)
	model, train_x, train_y, words = setup_bot(Mode)

	# Compile neural network
	model.compile(loss='categorical_crossentropy',
				  optimizer='adam', metrics=['accuracy'])
	# Train neural network
	model.fit(train_x, train_y, batch_size=32, epochs=1000,
			  verbose=1, validation_split=0.1, shuffle=True)

	file_name = get_file_name_by_mode(mode)

	save_training(file_name, model, train_x, train_y, words)


def save_training(file_name, model, train_x, train_y, words):
	# Save model
	with open(path.join(dir, '%s-model.json' % file_name), 'w') as f:
		f.write(model.to_json())
	model.save_weights(path.join(dir, '%s-weights.h5' % file_name))
	# Save total_stems and training data
	with open(path.join(dir, '%s.dump' % file_name), 'wb') as f:
		pickle.dump(TrainingData(words, train_x, train_y), f)


def set_mode(mode):
	if mode == "patterns":
		Mode: IntEnum = Patterns
	elif mode == "sentiment":
		Mode: IntEnum = Sentiment
	else:
		raise Exception
	return Mode


def load_model(mode: str) -> Tuple[Sequential, TrainingData]:
	"""
	Loads a pre-trained model from disk, as well as the training data dump.

	Returns:
		A pre-trained model loaded from disk as well as an instance of
		TrainingData, containg the data used for training and the list of total
		stems used.
	"""
	file_name = get_file_name_by_mode(mode)

	with open(path.join(dir, '%s.dump' % file_name), 'rb') as f:
		data = pickle.load(f)

	with open(path.join(dir, '%s-model.json' % file_name)) as f:
		model = model_from_json(f.read())

	model.load_weights(path.join(dir, '%s-weights.h5' % file_name))

	return model, data


def get_file_name_by_mode(mode):
	return "patterns" if mode == "Category" else "sentiments"
