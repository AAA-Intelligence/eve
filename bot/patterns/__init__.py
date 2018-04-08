from enum import IntEnum
from os import path
from pathlib import Path
from typing import Iterator

from bot.model_definitions import PatternCategory

dir = path.dirname(__file__)


def patterns_for_category(category: PatternCategory) -> Iterator[str]:
    """
    Opens the pattern definition file for the specified category if possible
    and returns an iterator for the pattern's lines.

    Args:
            mode: The category to load patterns for.

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
            yield line.rstrip('\n')
