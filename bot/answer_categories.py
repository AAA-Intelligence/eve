from enum import IntEnum, unique


@unique
class Category(IntEnum):
    JOKE = 0
    BOT_AGE = 1
    BOT_BIRTHDAY = 2
    BOT_NAME = 3
    BOT_GENDER = 4
    BOT_FAVORITE_COLOR = 5
