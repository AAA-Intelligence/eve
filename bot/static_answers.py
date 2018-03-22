from .answer_categories import Category
from .predefined_answers import answers_for_category
import random


def get_static_answer(category: Category) -> str:
    """
    Retrieves and formats a random predefined answer for the specified category
    from the database.

    Args:
        category: The category to retrieve an answer for.

    Returns:
        A random formatted answer for the specified category.
    """

    # TODO: Determine how to pass values for formatting, like bot name or gender

    answers = answers_for_category(category)
    answer = random.choice(answers)

    return answer.format(
        bot_name='Eve',
        bot_age=24
    )
