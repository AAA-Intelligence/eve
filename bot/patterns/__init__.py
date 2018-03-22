from typing import Dict, List, Iterator
from os import path
from pathlib import Path
from ..static_answers import Category

dir = path.dirname(__file__)


def patterns_for_category(category: Category) -> Iterator[str]:
    """
    Opens the pattern definition file for the specified category if possible
    and returns an iterator for the pattern's lines.

    Args:
        category: The category to load patterns for.

    Raises:
        FileNotFoundError:
            Raised if no pattern file is found for the specified category.

    Returns:
        A string iterator for iterating over all lines defined by the pattern
        file.
    """

    p = Path(dir, category.name + '.txt')
    if not p.is_file():
        raise FileNotFoundError(
            'No pattern definition file found for category {}'.format(category))
    with p.open(encoding='utf-8') as f:
        for line in f:
            yield line
