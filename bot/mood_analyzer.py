from bot.model_definitions import Sentiment
from bot.pattern_recognizer import analyze_input


def determine_mood(text: str):
	# TODO: Connect to TF
	return 0.0


def affect(mood: float):
	return 0.1 * mood


def analyze(text):
	# TODO determine how which percentages influence mood and affection
	return analyze_input(text, Mode=Sentiment), 0
