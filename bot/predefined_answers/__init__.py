from os import path
from pathlib import Path
from typing import Dict, List

from ..model_definitions import Patterns

# Cache for category answers
cache: Dict[str, List[str]] = {}

dir = path.dirname(__file__)


def answers_for_category(category: Patterns) -> List[str]:
    """
    Returns all predefined answers for the given category if possible.
    Answers are cached per category, so the first call for a category will read
    all defined answers into memory which will be re-used by subsequent calls.

    Args:
        category: The category to retrieve answers for.

    Raises:
        FileNotFoundError:
            Raised if no answer file could be found for the given category.

    Returns:
        An array containing all answers defined for the given category.
    """

    if category in cache:
        return cache[category]

    p = Path(dir, category.name + '.txt')
    if not p.is_file():
        raise FileNotFoundError(
            'No answer definition file found for category {}'.format(category))

    with p.open(encoding='utf-8') as f:
        answers = f.read().splitlines()
        cache[category] = answers
        return answers
