from enum import Enum, IntEnum, unique, auto
from typing import Union, Type


@unique
class AffectionCategory(IntEnum):
    """
    Categories for describing the affection data.
    """
    A_NEG = 0
    A_POS = 1


@unique
class MoodCategory(IntEnum):
    """
    Categories for describing the mood data.
    """
    M_NEG = 0
    M_POS = 1


@unique
class PatternCategory(IntEnum):
    """
    Different pattern categories which could be recognised by the bot
    The category names are equivalent to the file names for the patterns and static answers
    """
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
    BOT_RELIGION = auto()
    WEATHER = auto()
    COMPLIMENTS = auto()
    PICKUP_LINES = auto()
    BOT_HOBBIES = auto()
    DATE = auto()


Category = Union[MoodCategory, AffectionCategory, PatternCategory]


@unique
class Mode(Enum):
    """
    The modes the pattern recognizer can run in.
    Affects parameters like data source, model and error threshold.
    """
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
