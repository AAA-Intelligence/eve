from typing import NamedTuple
import json


class Context(NamedTuple):
    mood: float
    affection: float


class Request(NamedTuple):
    content: str
    timestamp: int
    context: Context
    bot_id: str
    user_id: str


def parse_request(json_data: str):
    data = json.loads(json_data)

    return Request(
        data["content"],
        data["timestamp"],
        Context(
            data["context"]["mood"],
            data["context"]["affection"]
        ),
        data["bot_id"],
        data["user_id"]
    )
