from typing import NamedTuple, List


class TrainingData(NamedTuple):
	"""
	Data type for training data that will be saved after training,
	used by the pattern recognizer to access the list of all stems
	"""
	total_stems: List[str]
	train_x: List[int]
	train_y: List[int]
