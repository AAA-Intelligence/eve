from typing import Optional, NamedTuple
from .static_answers import get_static_answer
from .answer_categories import Category
from .data import Request
from .logger import logger
from .train_patterns_model import load_model
from os import path
from nltk.stem.snowball import GermanStemmer
import numpy as np
import tensorflow as tf
import nltk
import pickle

dir = path.dirname(__file__)


# Load model and data
model, data = load_model()
# Create German snowball stemmer
stemmer = GermanStemmer()
# Threshold for pattern recognition
ERROR_THRESHOLD = 0.9


class PredictionResult(NamedTuple):
    """
    Data type for prediction results, used by detect_category
    """
    category: Category
    probability: float


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
    total_stems = data.total_stems
    bag = [0] * len(total_stems)
    for stem in stems:
        for i, s in enumerate(total_stems):
            if s == stem:
                bag[i] = 1

    # Convert to matrix
    input_data = np.asarray([bag])

    # Predict category
    results = model.predict(input_data)[0]
    results = [PredictionResult(Category(i), p)
               for i, p in enumerate(results) if i > 0]
    results.sort(key=lambda result: result.probability, reverse=True)

    logger.debug('Results: {}'.format(results))

    if len(results) > 0 and results[0].probability > ERROR_THRESHOLD:
        return results[0].category

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
        return get_static_answer(category)

    return None


def demo():
    """
    Demo mode for the pattern recognizer
    """

    request = Request(
        text=input('Please enter a question: '),
        mood=0.0,
        affection=0.0,
        bot_gender=0,
        bot_name='Lara',
        previous_text='Ich bin ein Baum'
    )
    answer = answer_for_pattern(request)
    if answer is None:
        print('No answer found')
    else:
        print('Answer: {}'.format(answer))
