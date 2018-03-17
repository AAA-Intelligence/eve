from typing import Dict, List, Iterator
from os import path
from pathlib import Path
from ..predefined_answers import Category

dir = path.dirname(__file__)


def patterns_for_category(category: Category) -> Iterator[str]:
    p = Path(dir, category.name + '.txt')
    if not p.is_file():
        raise Exception(
            'No pattern definition file found for category {}'.format(category))
    with p.open() as f:
        for line in f:
            yield line


def get_patterns() -> Dict[Category, List[str]]:
    patterns = {}
    # Iterate over all defined categories
    for category in Category:
        # Store all patterns defined for the category under the respective category key
        # in the patterns dictionary as a list
        patterns[category] = list(patterns_for_category(category))

    return patterns
