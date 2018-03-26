import numpy as np


def determine_mood(text: str):
	# TODO: Connect to TF
	return 0.0


def affect(mood: float):
	return 0.1 * mood


def analyze(current_affection: float, text: str):
	mood: float = determine_mood(text=text)
	current_affection += affect(mood=mood)
	return [mood, current_affection]


def demo():
	wordsList = np.load('wordsList.npy')
	print('Loaded the word list!')
	wordsList = wordsList.tolist()  # Originally loaded as numpy array
	wordsList = [word.decode('UTF-8') for word in
				 wordsList]  # Encode words as UTF-8
	wordVectors = np.load('wordVectors.npy')
	print('Loaded the word vectors!')

	maxSeqLength = 10  # Maximum length of sentence
	numDimensions = 300  # Dimensions for each word vector
	firstSentence = np.zeros((maxSeqLength), dtype='int32')
	firstSentence[0] = wordsList.index("i")
	firstSentence[1] = wordsList.index("thought")
	firstSentence[2] = wordsList.index("the")
	firstSentence[3] = wordsList.index("movie")
	firstSentence[4] = wordsList.index("was")
	firstSentence[5] = wordsList.index("incredible")
	firstSentence[6] = wordsList.index("and")
	firstSentence[7] = wordsList.index("inspiring")
	# firstSentence[8] and firstSentence[9] are going to be 0
	print(firstSentence.shape)
	print(firstSentence)  # Shows the row index for each word


