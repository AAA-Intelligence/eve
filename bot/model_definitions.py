from enum import Enum, IntEnum, unique
from typing import Union, Type


@unique
class SentimentCategory(IntEnum):
    """Can either be mood or affection"""
    # TODO separate sentiment and mood
    # A_*: Describes the affection data
    A_NEG = 0
    A_POS = 1
    # M_*: Describes the mood data
    M_NEG = 2
    M_POS = 3


@unique
class PatternCategory(IntEnum):
    BLACKLIST = 0
    JOKE = 1
    BOT_AGE = 2
    BOT_BIRTHDAY = 3
    BOT_NAME = 4
    BOT_GENDER = 5
    BOT_FAVORITE_COLOR = 6


Category = Union[SentimentCategory, PatternCategory]


@unique
class Mode(Enum):
    SENTIMENTS = 'sentiments'
    PATTERNS = 'patterns'

    @property
    def category_type(self) -> Type[Category]:
        if self == Mode.SENTIMENTS:
            return SentimentCategory
        if self == Mode.PATTERNS:
            return PatternCategory
        raise Exception('No category defined for mode {}'.format(self))
