from os import path
from pathlib import Path
from typing import Dict, List

from bot.data import Request
from bot.logger import logger
from bot.model_definitions import PatternCategory

# Cache for category answers
cache: Dict[str, List[str]] = {}

dir = path.dirname(__file__)


def answers_for_category(category: PatternCategory, request: Request) -> List[
    str]:
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
    if category.name == 'AFFECTION':
        direction = "POS" if request.affection > 0 else "NEG"
        p = Path(dir, 'A_%s.txt' % (direction))
    elif category.name == 'MOOD':
        direction = "POS" if request.mood > 0 else "NEG"
        logger.debug("DIRECTION:", direction)
        p = Path(dir, 'M_%s.txt' % (direction))
    else:
        p = Path(dir, category.name + '.txt')

    if not p.is_file():
        raise FileNotFoundError(
            'No answer definition file found for category {}'.format(category))

    with p.open(encoding='utf-8') as f:
        answers = f.read().splitlines()
        cache[category] = answers
        return answers
