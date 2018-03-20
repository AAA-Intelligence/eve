from typing import Optional
from .predefined_answers import Category, get_predefined_answer
from .data import Request
from .logger import logger
from os import path
from nltk.stem.snowball import GermanStemmer
import tensorflow
import tflearn
import nltk
import pickle
import numpy

dir = path.dirname(__file__)


def load_model(train_x, train_y):
    # Reset TensorFlow graph
    tensorflow.reset_default_graph()

    # Build neural network
    net = tflearn.input_data(shape=[None, len(train_x[0])])
    net = tflearn.fully_connected(net, 8)
    net = tflearn.fully_connected(net, 8)
    net = tflearn.fully_connected(net, len(train_y[0]), activation='softmax')
    net = tflearn.regression(net)

    # Define model and setup tensorboard
    model = tflearn.DNN(net, tensorboard_dir='tflearn_logs')
    # Load model
    model.load(path.join(dir, 'models', 'patterns.tflearn'))

    return model


def load_data():
    # Load data dump
    with open(path.join(dir, 'models', 'patterns.dump'), 'rb') as f:
        return pickle.load(f)


# Load data
data = load_data()
# Load model
model = load_model(data['train_x'], data['train_y'])
# Create German snowball stemmer
stemmer = GermanStemmer()

# Threshold for pattern recognition
ERROR_THRESHOLD = 0.9


def detect_category(request: Request) -> Optional[Category]:
    """
    Scans the supplied request for pre-defined patterns.

    Args:
        request: The request to scan for patterns.

    Returns:
        The category of the recognized pattern or None if none was found.
    """

    # Tokenize pattern
    words = nltk.word_tokenize(request.text)
    stems = [stemmer.stem(word.lower()) for word in words]
    total_stems = data['total_stems']
    bag = [0] * len(total_stems)
    for stem in stems:
        for i, s in enumerate(total_stems):
            if s == stem:
                bag[i] = 1

    # Predict category
    results = model.predict([numpy.array(bag)])[0]
    results = [
        (Category(category), probability)
        for category, probability in enumerate(results)
        if probability > ERROR_THRESHOLD
    ]
    results.sort(key=lambda x: x[1], reverse=True)

    logger.debug('Results: {}'.format(results))

    if len(results) > 0:
        return results[0][0]

    return None


def answer_for_pattern(request: Request) -> Optional[str]:
    """
    Scans the supplied request for pre-defined patterns and returns a
    pre-defined answer if possible.

    Args:
        request: The request to scan for patterns.

    Returns:
        A pre-defined answer for the scanned request or None if a pre-defined
        answer isn't possible.
    """

    category = detect_category(request)
    if category is not None:
        # Pattern found, retrieve pre-defined answer
        return get_predefined_answer(category)

    return None


def demo():
    request = Request(
        text='Würdest du mir sagen, wie alt du bist?',
        mood=0.0,
        affection=0.0,
        bot_gender=0,
        bot_name='Lara',
        previous_text='Ich bin ein Baum'
    )
    answer = answer_for_pattern(request)
    if answer is None:
        logger.debug('No answer found')
    else:
        logger.debug('Answer: {}'.format(answer))
