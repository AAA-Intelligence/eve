from typing import Tuple

from bot.data import Request
from bot.logger import logger
from bot.model_definitions import Mode, MoodCategory, AffectionCategory
from bot.pattern_recognizer import analyze_input


def fit_in_range(sign: int, probability: float) -> float:
    """
    Fits probabilities between 0.5 and 1 into the range of [-1,1] depending on the result

    :param sign: sign is equivalent to mood or affection either positive or negative
    :param probability: the probability with which the bot determined either mood or affection
    :return: returns a value between -1 and 1 indicating a certain mood or affection
    """

    return sign * (2 * probability - 1)


def analyze(request: Request) -> Tuple[float, float, float, float]:
    # TODO determine how which percentages influence mood and affection

    # input message of the user passed by a request
    text = request.text

    # Estimate mood through the neural network
    analyzed_mood = analyze_input(text, Mode.MOODS)
    if analyzed_mood:
        sign = -1 if analyzed_mood.category == MoodCategory.M_NEG else 1
        # calculate a value from the probability which could be fed to the neural network for text processing
        mood = fit_in_range(sign, analyzed_mood.probability)
    else:
        # return 0 if no certain affection was found
        mood = 0.0

    mood_bot = request.mood

    # TODO rethink e.g. function hyperbel like form towards 1
    mood_bot += 0.1 * mood
    if mood_bot > 1:
        mood_bot = 1.0
    # Estimate the affection through the neural network
    analyzed_affection = analyze_input(text, Mode.AFFECTIONS)
    if analyzed_affection:
        sign = -1 if analyzed_affection.category == AffectionCategory.A_NEG else 1
        # calculate a value from the probability which could be fed to the neural network for text processing
        affection = fit_in_range(sign, analyzed_affection.probability)
    else:
        # return 0 if no certain affection was found
        affection = 0.0
    affection_bot = request.affection
    affection_bot += 0.1 * affection
    if affection_bot > 1:
        affection_bot = 1.0

    logger.debug("MOOD: %f -> %f, AFFECTION: %f -> %f" % (
        request.mood, mood_bot, request.affection, affection_bot))

    return (mood, affection, mood_bot, affection_bot)
