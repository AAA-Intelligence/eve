from enum import Enum, IntEnum, unique, auto
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
    JOKE = auto()
    BOT_ANY = auto()
    BOT_AGE = auto()
    BOT_BIRTHDAY = auto()
    BOT_NAME = auto()
    BOT_GENDER = auto()
    BOT_FAVORITE_COLOR = auto()
    FATHER_ANY = auto()
    FATHER_AGE = auto()
    FATHER_NAME = auto()
    MOTHER_ANY = auto()
    MOTHER_AGE = auto()
    MOTHER_NAME = auto()
    ANY_AGE = auto()
    ANY_NAME = auto()
    ANY_BIRTHDAY = auto()


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
