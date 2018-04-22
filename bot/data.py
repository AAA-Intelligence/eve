import json
from datetime import date
from enum import Enum
from typing import NamedTuple, Optional

from bot.model_definitions import PatternCategory


class Gender(Enum):
    MALE = 0
    FEMALE = 1

    def __str__(self) -> str:
        if self == Gender.MALE:
            return 'männlich'
        elif self == Gender.FEMALE:
            return 'weiblich'


class Request(NamedTuple):
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
    text: str
    pattern: Optional[PatternCategory]
    mood: float
    affection: float


def parse_request(json_data: str) -> Request:
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
