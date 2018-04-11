from bot.model_definitions import Mode, MoodCategory, AffectionCategory
from bot.pattern_recognizer import analyze_input


def fit_in_range(sign: int, probability: float) -> float:
    # Fits probabilites between 0.5 and 1 into the range of [-1,1] depending on the result
    return sign * (2 * probability - 1)


def analyze(text: str):
    # TODO determine how which percentages influence mood and affection
    # Estimate mood through the neural network
    analyzed_mood = analyze_input(text, Mode.MOODS)
    if analyzed_mood:
        sign = -1 if analyzed_mood.category == MoodCategory.M_NEG else 1
        # calculate a value from the probability which could be fed to the neural network for text processing
        mood = fit_in_range(sign, analyzed_mood.probability)
    else:
        # return 0 if no certain affection was found
        mood = 0.0

    # Estimate the affection through the neural network
    analyzed_affection = analyze_input(text, Mode.AFFECTIONS)
    if analyzed_affection:
        sign = -1 if analyzed_affection.category == AffectionCategory.A_NEG else 1
        # calculate a value from the probability which could be fed to the neural network for text processing
        affection = fit_in_range(sign, analyzed_affection.probability)
    else:
        # return 0 if no certain affection was found
        affection = 0.0
    return (analyzed_mood, mood, analyzed_affection, affection)
