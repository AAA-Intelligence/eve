from typing import Tuple, Optional

from bot.model_definitions import Mode, SentimentCategory
from bot.pattern_recognizer import analyze_input


def determine_mood(text: str) -> float:
    # TODO: Connect to TF
    return 0.0


def affect(mood: float) -> float:
    return 0.1 * mood


def analyze(text: str) -> Tuple[SentimentCategory, float]:
    # TODO determine how which percentages influence mood and affection
    return analyze_input(text, Mode.SENTIMENTS), 0.0
