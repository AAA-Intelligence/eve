from typing import List, Tuple, Set, NamedTuple
from os import path, mkdir
from pathlib import Path
from nltk.stem.snowball import GermanStemmer
from .patterns import patterns_for_category
from .predefined_answers import Category
from keras.models import Sequential, model_from_json
from keras.layers import Dense, Dropout, Activation
import nltk
import random
import tensorflow as tf
import numpy as np
import pickle

dir = path.join(path.dirname(__file__), 'models')
if not path.exists(dir):
    mkdir(dir)
if not path.isdir(dir):
    raise Exception('Models path is not a directory: {}'.format(dir))

# Create a word stemmer based on the snowball stemming algorithm for the German language
stemmer = GermanStemmer()


class TrainingData(NamedTuple):
    """
    Data type for training data that will be saved after training,
    used by the pattern recognizer to access the list of all stems
    """
    total_stems: List[str]
    train_x: List[int]
    train_y: List[int]


def train_model():
    """
    Trains a neural network with the defined patterns and categories.
    Patterns will be split into words, stemmed by a German snowball stemmer and
    indexed by saving all stems in a list of total stems and assigning indices.
    The trained model will be saved in the models directory and can be loaded
    by the pattern recognizer using the load_model function.
    """

    total_stems: Set[str] = set()
    patterns: List[Tuple[Category, Set[str]]] = []

    # Iterate over all defined categories
    for category in Category:
        # Parse all patterns defined for this category
        for pattern in patterns_for_category(category):
            # Tokenize pattern into words
            words = nltk.word_tokenize(pattern)
            # Get stems for the pattern's words, as a set to avoid duplicates
            stems = {stemmer.stem(w.lower()) for w in words}
            # Add stems associated with association to the category to the
            # pattern list.
            patterns.append((category, stems))
            # Add stems to total set of stems, needed for conversion to numeric
            # TensorFlow training array
            total_stems |= stems

    words = sorted(list(total_stems))

    # Now that all stems have been collected, we can create an array suitable
    # for training our TensorFlow model.
    # For this, we tell TensorFlow, by defining this array, which stems can lead
    # to which patterns.
    train_x = []
    train_y = []
    for category, stems in patterns:
        bag_of_words = [1 if word in stems else 0 for word in words]
        output_row = [0] * len(Category)
        output_row[category] = 1
        train_x.append(bag_of_words)
        train_y.append(output_row)

    # Convert lists to numpy arrays
    train_x = np.asarray(train_x)
    train_y = np.asarray(train_y)

    # Define neural network
    model = Sequential()
    model.add(Dense(512, input_shape=(len(train_x[0]),), activation='relu'))
    model.add(Dropout(0.5))
    model.add(Dense(256, activation='sigmoid'))
    model.add(Dropout(0.5))
    model.add(Dense(len(train_y[0]), activation='softmax'))

    # Compile neural network
    model.compile(loss='categorical_crossentropy',
                  optimizer='adam', metrics=['accuracy'])
    # Train neural network
    model.fit(train_x, train_y, batch_size=32, epochs=1000,
              verbose=1, validation_split=0.1, shuffle=True)

    # Save model
    with open(path.join(dir, 'patterns-model.json'), 'w') as f:
        f.write(model.to_json())
    model.save_weights(path.join(dir, 'patterns-weights.h5'))

    # Save total_stems and training data
    with open(path.join(dir, 'patterns.dump'), 'wb') as f:
        pickle.dump(TrainingData(words, train_x, train_y), f)


def load_model() -> Tuple[Sequential, TrainingData]:
    """
    Loads a pre-trained model from disk, as well as the training data dump.

    Returns:
        A pre-trained model loaded from disk as well as an instance of
        TrainingData, containg the data used for training and the list of total
        stems used.
    """

    with open(path.join(dir, 'patterns.dump'), 'rb') as f:
        data = pickle.load(f)

    with open(path.join(dir, 'patterns-model.json')) as f:
        model = model_from_json(f.read())

    model.load_weights(path.join(dir, 'patterns-weights.h5'))

    return model, data
