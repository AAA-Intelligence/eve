from typing import Dict, List
from os import path, listdir
from ..predefined_answers import Category


def get_patterns() -> Dict[Category, List[str]]:
    patterns = {}

    dir = path.dirname(__file__)
    for file_name in listdir(dir):
        if file_name.endswith('.txt'):
            with open(file_name) as f:
                # Retrieve category name without file extension
                name = file_name.replace('.txt', '')
                if name not in Category:
                    raise Exception(
                        'Found pattern definition for invalid category {}'.format(name))
                category = Category[name]
                # Store all lines of the pattern file in the dictionary, associated to
                # the specific category by treating f as an iterator over the file's lines
                patterns[category] = list(f)
