from typing import List, Tuple, Set
from os import path, mkdir
from pathlib import Path
from nltk.stem.snowball import GermanStemmer
from .patterns import patterns_for_category
from .predefined_answers import Category
import nltk
import random
import tensorflow as tf
import tflearn
import pickle

dir = path.join(path.dirname(__file__), 'models')
if not path.exists(dir):
    mkdir(dir)
if not path.isdir(dir):
    raise Exception('Models path is not a directory: {}'.format(dir))

# Define all punctuation we want to ignore in texts
punctuation = ['.', ',', ';', '?', '!', '-', '(', ')', '{', '}', '/', '\\']
# Create a word stemmer based on the snowball stemming algorithm for the German language
stemmer = GermanStemmer()


def remove_punctuation(text: str) -> str:
    return ''.join(c for c in text if c not in punctuation)


def train_model():
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
    training_data = []
    for category, stems in patterns:
        bag_of_words = [1 if word in stems else 0 for word in words]
        output_row = [0] * len(Category)
        output_row[category] = 1
        training_data.append([bag_of_words, output_row])

    random.shuffle(training_data)

    train_x = [row[0] for row in training_data]
    train_y = [row[1] for row in training_data]

    # Reset TensorFlow graph
    tf.reset_default_graph()

    # Build neural network
    net = tflearn.input_data(shape=[None, len(train_x[0])])
    net = tflearn.fully_connected(net, 8)
    net = tflearn.fully_connected(net, 8)
    net = tflearn.fully_connected(net, len(train_y[0]), activation='softmax')
    net = tflearn.regression(net)

    # Define model and setup tensorboard
    model = tflearn.DNN(net, tensorboard_dir='tflearn_logs')

    # Start training (apply gradient descent algorithm)
    model.fit(train_x, train_y, n_epoch=1000, batch_size=8, show_metric=True)
    model.save(path.join(dir, 'patterns.tflearn'))

    # Save total_stems and training data
    with open(path.join(dir, 'patterns.dump'), 'wb') as f:
        pickle.dump({
            'total_stems': words,
            'train_x': train_x,
            'train_y': train_y
        }, f)
