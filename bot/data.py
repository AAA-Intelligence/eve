import json
from datetime import date
from enum import Enum
from typing import NamedTuple, Optional

from bot.model_definitions import PatternCategory
from bot.logger import logger


class Gender(Enum):
    """
    Specifies the possible bot genders and provides a method for converting
    the value a to a string in German language.
    """
    MALE = 0
    FEMALE = 1

    def __str__(self) -> str:
        if self == Gender.MALE:
            return 'männlich'
        elif self == Gender.FEMALE:
            return 'weiblich'


class Request(NamedTuple):
    """
    Specifies the structure of incoming request data.
    The JSON encoded requests will be decoded to an instance of this class.
    Also provides helper methods for formatting certain data to strings that
    can be used in the answer templates.
    """
    text: str
    previous_pattern: Optional[PatternCategory]
    mood: float
    affection: float
    bot_gender: Gender
    bot_name: str
    bot_birthdate: date
    bot_favorite_color: str
    father_name: str
    father_age: int
    mother_name: str
    mother_age: int

    @property
    def bot_birthday(self) -> str:
        return self.bot_birthdate.strftime('%d.%m.%Y')

    @property
    def bot_age(self) -> int:
        today = date.today()
        bdate = self.bot_birthdate
        return today.year - bdate.year - (
            1 if (today.month, today.day) < (bdate.month, bdate.day) else 0)


class Response(NamedTuple):
    """
    Specifies the structure of outgoing response data.
    The response will be encoded to a JSON string before being written to
    the output.
    """
    text: str
    pattern: Optional[PatternCategory]
    mood: float
    affection: float


def parse_request(json_data: str) -> Request:
    """
    Parses a JSON request string to an instance of the Request class.

    Args:
        json_data: The JSON encoded request string.

    Returns:
        The decoded data as an instance of the Request class.
    """
    logger.debug('Type: {}'.format(type(json_data)))
    data = json.loads(json_data)

    return Request(
        data["text"],
        PatternCategory(data["previous_pattern"]
                        ) if "previous_pattern" in data else None,
        data["mood"],
        data["affection"],
        Gender(data["bot_gender"]),
        data["bot_name"],
        date.fromtimestamp(data["bot_birthdate"]),
        data["bot_favorite_color"],
        data["father_name"],
        data["father_age"],
        data["mother_name"],
        data["mother_age"],
    )
