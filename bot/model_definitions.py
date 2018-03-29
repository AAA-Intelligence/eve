from enum import IntEnum, unique


# Can either be mood or affection
@unique
class Sentiment(IntEnum):
	# TODO separate sentiment and mood
	# A_*: Describes the affection data
	A_NEG = 0
	A_POS = 1
	# M_*: Describes the mood data
	M_NEG = 2
	M_POS = 3


@unique
class Patterns(IntEnum):
	BLACKLIST = 0
	JOKE = 1
	BOT_AGE = 2
	BOT_BIRTHDAY = 3
	BOT_NAME = 4
	BOT_GENDER = 5
	BOT_FAVORITE_COLOR = 6
