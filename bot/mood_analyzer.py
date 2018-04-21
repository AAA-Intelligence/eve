from math import tanh
from typing import Tuple

from bot.data import Request
from bot.model_definitions import Mode, MoodCategory, AffectionCategory
from bot.pattern_recognizer import analyze_input

# factor 0.2 ensures steady adjustment of bots mood and affection
IMPACT_FACTOR = 0.2


def stretch_prob(probability: float) -> float:
    """
    Fits probabilities between 0.75 and 1 into the range of [-1,1] depending on the result

    :param probability: The probability with which the bot determined either mood or affection
    :return: Returns a value between 0 and 1 which is later passed to the tanh(2*x) function for
            a more realistic change in mood and affection
    """

    return 4 * probability - 3


def analyze(request: Request) -> Tuple[float, float]:
    """
    R

    :param request: The request passed by the web server to the bot instance.
                    It contains all the necessary information to determine a new bot mood/affection.
    :return: Returns the new mood/affection of the bot calculated on the text of the message and the
            previous mood and affection
    """
    # input message of the user passed by a request
    text = request.text

    # Inits bots mood. It stays unchanged if the message does not contain certain signs of specificly
    # positive or negative mood.
    mood_bot = request.mood

    # Estimate mood through the neural network
    mood_result = analyze_input(text, Mode.MOODS)
    if mood_result:
        mood_probability = mood_result.probability
        # checking for a negative or positive mood
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

    # Inits bots affection. It stays unchanged if the message does not contain certain signs of specificly
    # positive or negative affection.
    affection_bot = request.affection

    # Estimate mood through the neural network
    affection_result = analyze_input(text, Mode.AFFECTIONS)
    if affection_result:
        affection_probability = affection_result.probability

        # checking for a negative or positive affection
        if affection_result.category == AffectionCategory.A_NEG:
            sign = -1
        else:
            sign = 1

        # calculate a value from the probability which could be fed to the neural network for text processing
        affection_message = stretch_prob(affection_probability)

        # apply message to affection
        affection_bot = affection_bot + sign * affection_message * IMPACT_FACTOR
        affection_bot = tanh(2 * affection_bot)
        if affection_bot > 1:
            affection_bot = 1.0

    return (mood_bot, affection_bot)
