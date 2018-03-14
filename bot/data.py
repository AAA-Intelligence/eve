from typing import NamedTuple
from enum import Enum
import json


class Gender(Enum):
    MALE = 0
    FEMALE = 1


class Request(NamedTuple):
    text: str
    mood: float
    affection: float
    bot_gender: Gender
    bot_name: str
    previous_text: str


class Response(NamedTuple):
    text: str
    mood: float
    affection: float


def parse_request(json_data: str):
    data = json.loads(json_data)

    return Request(
        data["text"],
        data["mood"],
        data["affection"],
        Gender(data["bot_gender"]),
        data["bot_name"],
        data["previous_text"]
    )
