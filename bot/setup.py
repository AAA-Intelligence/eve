from collections import Set
from typing import List, Tuple

import nltk
import numpy as np
from keras import Sequential
from keras.layers import Dense, Dropout
from nltk.stem.snowball import GermanStemmer

from bot.model_definitions import Patterns, Sentiment
from bot.patterns import patterns_for_category
from bot.sentiments import patterns_for_sentiment


def setup_bot(Mode):
	enum_elements, words = read_training_data(Mode)
	# Now that all stems have been collected, we can create an array suitable
	# for training our TensorFlow model.
	# For this, we tell TensorFlow, by defining this array, which stems can lead
	# to which patterns.
	train_x, train_y = setup_traing_data(Mode, enum_elements, words)
	model = setup_nn_model(train_x, train_y)
	return model, train_x, train_y, words


def read_training_data(Mode):
	total_stems: Set[str] = set()
	enum_elements: List[Tuple[Mode, Set[str]]] = []
	# Iterate over all defined categories
	for element in Mode:
		# Parse all patterns defined for this category
		if Mode == Patterns:
			elements = patterns_for_category(element)
		elif Mode == Sentiment:
			elements = patterns_for_sentiment(element)
		else:
			raise Exception
		for e in elements:
			total_stems = build_stems(e, element, enum_elements, total_stems)

	words = sorted(list(total_stems))
	return enum_elements, words


def build_stems(e, element, enum_elements, total_stems):
	# Tokenize pattern into words
	words = nltk.word_tokenize(e)
	# Get stems for the pattern's words, as a set to avoid duplicates
	stemmer = GermanStemmer()
	stems = {stemmer.stem(w.lower()) for w in words}
	# Add stems associated with association to the category to the
	# pattern list.
	enum_elements.append((element, stems))
	# Add stems to total set of stems, needed for conversion to numeric
	# TensorFlow training array
	total_stems |= stems
	return total_stems


def setup_traing_data(Mode, enum_elements, words):
	train_x = []
	train_y = []
	for element, stems in enum_elements:
		bag_of_words = [1 if word in stems else 0 for word in words]
		output_row = [0] * len(Mode)
		output_row[element] = 1
		train_x.append(bag_of_words)
		train_y.append(output_row)
	# Convert lists to numpy arrays
	train_x = np.asarray(train_x)
	train_y = np.asarray(train_y)
	return train_x, train_y


def setup_nn_model(train_x, train_y):
	# Define neural network

	DENSITIY: int = 512
	DROPOUT: float = 0.5

	model = Sequential()
	model.add(
		Dense(DENSITIY, input_shape=(len(train_x[0]),), activation='relu'))
	model.add(Dropout(DROPOUT))
	model.add(Dense(DENSITIY // 2, activation='sigmoid'))
	model.add(Dropout(DROPOUT))
	model.add(Dense(len(train_y[0]), activation='softmax'))
	return model
