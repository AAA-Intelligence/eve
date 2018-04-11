from typing import List, Tuple, Type, Set

import nltk
import numpy as np
from keras import Sequential
from keras.layers import Dense, Dropout
from nltk.stem.snowball import GermanStemmer

from bot.affections import patterns_for_affection
from bot.model_definitions import Category, Mode
from bot.moods import patterns_for_mood
from bot.patterns import patterns_for_category


def setup_bot(mode: Mode) -> Tuple[
    Sequential, np.ndarray, np.ndarray, List[str]]:
    elements, words = read_training_data(mode)
    # Now that all stems have been collected, we can create an array suitable
    # for training our TensorFlow model.
    # For this, we tell TensorFlow, by defining this array, which stems can lead
    # to which patterns.
    train_x, train_y = setup_traing_data(mode.category_type, elements, words)
    model = setup_nn_model(train_x, train_y)
    return model, train_x, train_y, words


def read_training_data(mode: Mode) -> Tuple[Category, Set[str]]:
    total_stems: Set[str] = set()
    elements: List[Tuple[Category, Set[str]]] = []

    if mode == Mode.PATTERNS:
        reader_func = patterns_for_category
    elif mode == Mode.AFFECTIONS:
        reader_func = patterns_for_affection
    elif mode == Mode.MOODS:
        reader_func = patterns_for_mood
    else:
        raise ValueError('Unknown mode {}'.format(mode))

    CategoryType = mode.category_type

    # Iterate over all defined categories
    for category in CategoryType:
        # Parse all patterns defined for this category
        patterns = reader_func(category)
        for pattern in patterns:
            total_stems = build_stems(pattern, category, elements, total_stems)

    words = sorted(list(total_stems))
    return elements, words


def build_stems(
    pattern: str,
    category: Category,
    elements: List[Tuple[Category, Set[str]]],
    total_stems: Set[str]
    ) -> Set[str]:
    # Tokenize pattern into words
    words = nltk.word_tokenize(pattern)
    # Get stems for the pattern's words, as a set to avoid duplicates
    stemmer = GermanStemmer()
    stems: Set[str] = {stemmer.stem(w.lower()) for w in words}
    # Add stems associated with association to the category to the
    # pattern list.
    elements.append((category, stems))
    # Add stems to total set of stems, needed for conversion to numeric
    # TensorFlow training array
    total_stems |= stems
    return total_stems


def setup_traing_data(
    CategoryType: Type[Category],
    elements: List[Tuple[Category, Set[str]]],
    words: List[str]
    ) -> Tuple[np.ndarray, np.ndarray]:
    train_x = []
    train_y = []
    for category, stems in elements:
        bag_of_words = [1 if word in stems else 0 for word in words]
        output_row = [0] * len(CategoryType)
        output_row[category] = 1
        train_x.append(bag_of_words)
        train_y.append(output_row)
    # Convert lists to numpy arrays
    train_x = np.asarray(train_x)
    train_y = np.asarray(train_y)
    return train_x, train_y


def setup_nn_model(train_x: np.ndarray, train_y: np.ndarray) -> Sequential:
    # Define neural network

    DENSITIY: int = 512
    DROPOUT: float = 0.25

    model = Sequential()

    model.add(
        Dense(DENSITIY, input_shape=(len(train_x[0]),), activation='relu'))
    model.add(Dropout(DROPOUT))
    model.add(Dense(DENSITIY // 2, activation='sigmoid'))
    model.add(Dropout(DROPOUT))
    model.add(Dense(len(train_y[0]), activation='softmax'))

    return model
