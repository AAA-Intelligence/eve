from typing import Dict, List
from os import path
from pathlib import Path
from ..answer_categories import Category

# Cache for category answers
cache: Dict[str, List[str]] = {}

dir = path.dirname(__file__)


def answers_for_category(category: Category) -> List[str]:
    if category in cache:
        return cache[category]

    p = Path(dir, category.name + '.txt')
    if not p.is_file():
        raise FileNotFoundError(
            'No pattern definition file found for category {}'.format(category))

    with p.open(encoding='utf-8') as f:
        answers = [line for line in f]
        cache[category] = answers
        return answers
