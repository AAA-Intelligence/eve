from enum import IntEnum, unique


@unique
class Category(IntEnum):
    JOKE = 0
    BOT_AGE = 1


def get_predefined_answer(category: Category) -> str:
    """
    Retrieves and formats a random predefined answer for the specified category
    from the database.

    Args:
        category: The category to retrieve an answer for.

    Returns:
        A random formatted answer for the specified category.
    """

    # TODO: Implement
    # TODO: Determine how to pass values for formatting, like bot name or gender

    return 'TODO: Implement (Category {})'.format(category.name)
