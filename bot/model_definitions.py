from enum import Enum, IntEnum, unique, auto
from typing import Union, Type


@unique
class AffectionCategory(IntEnum):
    # TODO separate sentiment and mood
    # A_*: Describes the affection data
    A_NEG = 0
    A_POS = 1


@unique
class MoodCategory(IntEnum):
    # M_*: Describes the mood data
    M_NEG = 0
    M_POS = 1


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
    MOOD = auto()
    AFFECTION = auto()


Category = Union[MoodCategory, AffectionCategory, PatternCategory]


@unique
class Mode(Enum):
    PATTERNS = 'patterns'
    AFFECTIONS = 'affections'
    MOODS = 'moods'

    @property
    def category_type(self) -> Type[Category]:
        if self == Mode.MOODS:
            return MoodCategory
        if self == Mode.AFFECTIONS:
            return AffectionCategory
        if self == Mode.PATTERNS:
            return PatternCategory
        raise Exception('No category defined for mode {}'.format(self))
