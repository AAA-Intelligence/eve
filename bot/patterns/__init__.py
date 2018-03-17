from typing import Dict, List
from os import path
from pathlib import Path
from ..predefined_answers import Category

dir = path.dirname(__file__)


def pattern_file_for_category(category: Category):
    p = Path(dir, category.name + '.txt')
    if not p.is_file():
        raise Exception(
            'No pattern definition file found for category {}'.format(category))
    return p.open()


def get_patterns() -> Dict[Category, List[str]]:
    patterns = {}
    # Iterate over all defines categories
    for category in Category:
        # Try to the open pattern definition file for the category
        with pattern_file_for_category(category) as f:
            # Store all lines of the pattern file in the dictionary, associated to
            # the specific category by treating f as an iterator over the file's lines
            patterns[category] = list(f)
