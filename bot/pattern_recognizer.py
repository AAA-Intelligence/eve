from datetime import date
from enum import IntEnum
from os import path
from typing import Optional, NamedTuple, Generic, TypeVar

import nltk
import numpy as np
from nltk.stem.snowball import GermanStemmer

from bot.model_definitions import Mode, Category
from bot.data import Request, Gender
from bot.logger import logger
from bot.static_answers import get_static_answer
from bot.trainer import load_model

dir = path.dirname(__file__)

# Create German snowball stemmer
stemmer = GermanStemmer()
# Threshold for pattern recognition
ERROR_THRESHOLD = 0.9


class PredictionResult(NamedTuple):
    """
    Data type for prediction results, used by analyze_input
    """
    category: Category
    probability: float


def analyze_input(text: str, mode: Mode) -> Optional[PredictionResult]:
    """
    Scans the supplied request for pre-defined patterns.

    Args:
        request: The request to scan for patterns.

    Returns:
        The category of the recognized pattern or None if none was found.
    """

    # Load model and data
    model, data = load_model(mode)
    # Tokenize pattern
    words = nltk.word_tokenize(text)
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
    lower_bound = 0 if mode == Mode.PATTERNS else -1

    CategoryType = mode.category_type

    results = [PredictionResult(CategoryType(i), p) for i, p in enumerate(results)
               if i > lower_bound]
    results.sort(key=lambda result: result.probability, reverse=True)

    logger.debug('Results: {}'.format(results))

    if len(results) > 0 and results[0].probability > ERROR_THRESHOLD:
        return results[0]

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
    result = analyze_input(request.text, Mode.PATTERNS)
    if result is not None:
        # Pattern found, retrieve pre-defined answer
        return get_static_answer(result.category, request)

    return None


def demo(mode: str):
    """
    Demo mode for the pattern recognizer
    """

    request = Request(
        text=input('Please enter a question: '),
        previous_text='Ich bin ein Baum',
        mood=0.0,
        affection=0.0,
        bot_gender=Gender.APACHE,
        bot_name='Lara',
        bot_birthdate=date(1995, 10, 5),
        bot_favorite_color='grün'
    )
    answer = answer_for_pattern(request)
    if answer is None:
        print('No answer found')
    else:
        print('Answer: {}'.format(answer))
