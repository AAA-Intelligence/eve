from math import tanh
from typing import Tuple

from bot.data import Request
from bot.logger import logger
from bot.model_definitions import Mode, MoodCategory, AffectionCategory
from bot.pattern_recognizer import analyze_input

# factor 0.2 ensures steady adjustment of bots mood and affection
IMPACT_FACTOR = 0.2


def stretch_prob(x: float) -> float:
    """
    Fits probabilities between 0.5 and 1 into the range of [-1,1] depending on the result

    :param x: the probability with which the bot determined either mood or affection
    :return: returns a value between -1 and 1 indicating a certain mood or affection
    """

    return 4 * x - 3


def analyze(request: Request) -> Tuple[float, float, float, float]:
    # TODO determine how which percentages influence mood and affection

    # input message of the user passed by a request
    text = request.text

    # init bots mood
    mood_bot = request.mood

    # Estimate mood through the neural network
    mood_result = analyze_input(text, Mode.MOODS)
    if mood_result:
        mood_probability = mood_result.probability

        if mood_result.category == MoodCategory.M_NEG:
            sign = -1
        else:
            sign = 1
        # calculate a value from the probability which could be fed to the neural network for text processing
        mood_message = stretch_prob(mood_probability)

        # apply message to mood
        mood_bot = mood_bot + sign * mood_message * IMPACT_FACTOR
        if mood_bot > 1:
            mood_bot = 1.0
        mood_bot = tanh(2 * mood_bot)



    else:
        # return 0 if no certain affection was found
        mood_message = 0.0

    # init bots affection
    affection_bot = request.affection

    # Estimate mood through the neural network
    affection_result = analyze_input(text, Mode.AFFECTIONS)
    if affection_result:
        affection_probability = affection_result.probability
        if affection_result.category == AffectionCategory.A_NEG:
            sign = -1
        else:
            sign = 1
        # calculate a value from the probability which could be fed to the neural network for text processing
        affection_message = stretch_prob(affection_probability)

        # apply message to affection
        affection_bot += sign * affection_message * IMPACT_FACTOR
        affection_bot = tanh(2 * affection_bot)
        if affection_bot > 1:
            affection_bot = 1.0

    else:
        # return 0 if no certain affection was found
        affection_message = 0.0

    logger.debug("MOOD: %f -> %f, AFFECTION: %f -> %f" % (
        request.mood, mood_bot, request.affection, affection_bot))

    return (mood_message, affection_message, mood_bot, affection_bot)
