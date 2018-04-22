from datetime import date
from enum import IntEnum
from os import path
from typing import Optional, NamedTuple, Generic, TypeVar, Tuple, Dict

import nltk
import numpy as np
from nltk.stem.snowball import GermanStemmer

from bot.model_definitions import Mode, Category, PatternCategory
from bot.data import Request, Gender
from bot.logger import logger
from bot.static_answers import get_static_answer
from bot.trainer import load_model

dir = path.dirname(__file__)

# Create German snowball stemmer
stemmer = GermanStemmer()


class PredictionResult(NamedTuple):
    """
    Data type for prediction results, used by analyze_input
    """
    category: Category
    probability: float  # 1 equals 100% probability


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

    # Log all results to debug the percentages
    logger.debug('Results: {}'.format(results))

    error_threshold = 0.9
    if mode == Mode.MOODS or mode == mode.AFFECTIONS:
        # Sentiment analysis should be less sensitive than pattern recognition
        # as false-positives don't have as big of an impact and real persons
        # are more likely to misinterprete the mood or affection of a sentence
        # as well.
        error_threshold = 0.75

    # Only the most-probable result is interesting for us here.
    # Check if it passes the error threshold and continue if applicable.
    if len(results) > 0 and results[0].probability > error_threshold:
        return results[0]

    return None


# Pattern transitions for context recognition as a relation of the form
# { [PatternCategory, PatternCategory]: PatternCategory }.
# Only the previously occuring pattern will be considered for transitions and
# only certain pattern combinations lead to a new pattern.
pattern_transitions: Dict[Tuple[PatternCategory, PatternCategory], PatternCategory] = {
    (PatternCategory.FATHER_AGE, PatternCategory.ANY_NAME): PatternCategory.FATHER_NAME,
    (PatternCategory.MOTHER_AGE, PatternCategory.ANY_NAME): PatternCategory.MOTHER_NAME,
    (PatternCategory.FATHER_NAME, PatternCategory.ANY_AGE): PatternCategory.FATHER_AGE,
    (PatternCategory.MOTHER_NAME, PatternCategory.ANY_AGE): PatternCategory.MOTHER_AGE,
    (PatternCategory.BOT_NAME, PatternCategory.MOTHER_ANY): PatternCategory.MOTHER_NAME,
    (PatternCategory.FATHER_NAME, PatternCategory.MOTHER_ANY): PatternCategory.MOTHER_NAME,
    (PatternCategory.MOTHER_NAME, PatternCategory.FATHER_ANY): PatternCategory.FATHER_NAME,
    (PatternCategory.BOT_AGE, PatternCategory.MOTHER_ANY): PatternCategory.MOTHER_AGE,
    (PatternCategory.FATHER_AGE, PatternCategory.MOTHER_ANY): PatternCategory.MOTHER_AGE,
    (PatternCategory.MOTHER_AGE, PatternCategory.FATHER_ANY): PatternCategory.FATHER_AGE,
    (PatternCategory.FATHER_NAME, PatternCategory.BOT_ANY): PatternCategory.BOT_NAME,
    (PatternCategory.MOTHER_NAME, PatternCategory.BOT_ANY): PatternCategory.BOT_NAME,
    (PatternCategory.FATHER_AGE, PatternCategory.BOT_ANY): PatternCategory.BOT_AGE,
    (PatternCategory.MOTHER_AGE, PatternCategory.BOT_ANY): PatternCategory.BOT_AGE,
}


def answer_for_pattern(request: Request) -> Optional[Tuple[PatternCategory, str]]:
    """
    Scans the supplied request for pre-defined patterns and returns a
    pre-defined answer if possible.

    Args:
        request: The request to scan for patterns.

    Returns:
        A tuple containing the detected pattern category and pre-defined answer
        for the scanned request, or None if no pattern was recognized.
    """
    result = analyze_input(request.text, Mode.PATTERNS)
    if result is not None:
        # Pattern found
        category = result.category
        # Check context for a possible category transition
        previous_category = request.previous_pattern
        if previous_category and (previous_category, result.category) in pattern_transitions:
            category = pattern_transitions[(
                previous_category, result.category)]
        # Retrieve pre-defined answer
        return category, get_static_answer(category, request)

    return None


def demo():
    """
    Demo mode for the pattern recognizer
    """
    request = Request(
        text=input('Please enter a question: '),
        previous_pattern=PatternCategory.BOT_NAME,
        mood=0.0,
        affection=0.0,
        bot_gender=Gender.FEMALE,
        bot_name='Lara',
        bot_birthdate=date(1995, 10, 5),
        bot_favorite_color='grün'
    )
    answer = answer_for_pattern(request)
    if answer is None:
        print('No answer found')
    else:
        print('Answer: {}'.format(answer))
