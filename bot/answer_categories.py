from enum import IntEnum, unique


@unique
class Category(IntEnum):
    BLACKLIST = 0
    JOKE = 1
    BOT_AGE = 2
    BOT_BIRTHDAY = 3
    BOT_NAME = 4
    BOT_GENDER = 5
    BOT_FAVORITE_COLOR = 6
