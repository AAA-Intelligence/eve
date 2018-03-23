def determine_mood(text: str):
	# TODO: Connect to TF
	return 0.0


def affect(mood: float):
	return 0.1 * mood


def analyze(current_affection: float, text: str):
	mood: float = determine_mood(text=text)
	current_affection += affect(mood=mood)
	return [mood, current_affection]
