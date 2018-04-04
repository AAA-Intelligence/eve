import random

from bot.data import Request
from bot.model_definitions import PatternCategory
from bot.predefined_answers import answers_for_category


def get_static_answer(category: PatternCategory, request: Request) -> str:
    """
    Retrieves and formats a random predefined answer for the specified category
    from the database.

    Args:
        category: The category to retrieve an answer for.
        request: The request the answer is directed at.

    Returns:
        A random formatted answer for the specified category.
    """

    # TODO: Determine how to pass values for formatting, like bot name or gender

    answers = answers_for_category(category)
    answer = random.choice(answers)

    return answer.format(r=request)
