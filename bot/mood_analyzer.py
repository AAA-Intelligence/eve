from typing import NamedTuple


class AnalyticalResult(NamedTuple):
    mood: float
    affection: float


def analyze(text: str) -> AnalyticalResult:
    # TODO: Implement
    return AnalyticalResult(1, 1)
