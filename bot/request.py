import json


class Context(object):
    def __init__(self, mood: float, affection: float):
        self.mood = mood
        self.affection = affection


class Request(object):
    def __init__(self, content: str, timestamp: int, context: Context, bot_id: str, user_id: str):
        self.content = content
        self.timestamp = timestamp
        self.context = context
        self.bot_id = bot_id
        self.user_id = user_id


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
