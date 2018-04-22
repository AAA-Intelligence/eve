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
    """
    Loads and returns the relevant data for executing the bot in the specified mode,
    i.e. for pattern recognition or sentiment analysis.

    Args:
        mode: The mode to load data for.

    Returns:
        A tuple contain the trained model, the feature parameters, the label
        parameters and the bag of words used for indexing words in the feature set.
    """

    elements, words = read_training_data(mode)
    # Now that all stems have been collected, we can create an array suitable
    # for training our TensorFlow model.
    # For this, we tell TensorFlow, by defining this array, which stems can lead
    # to which patterns.
    train_x, train_y = setup_traing_data(mode.category_type, elements, words)
    model = setup_nn_model(train_x, train_y)
    return model, train_x, train_y, words


def read_training_data(mode: Mode) -> Tuple[List[Tuple[Category, Set[str]]], List[str]]:
    """
    Reads the training data for the specified mode.

    Args:
        mode: The mode to load training data for.

    Returns:
        A tuple containing a list of all categories with their respective vocabulary
        and the bag of words used for indexing a string (i.e. for conversion to a
        1xN matrix as feature parameters).
    """

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
    """
    Builds a set of stems for all words used in the pattern.

    Args:
        pattern: The pattern to tokenize and stem.
        category: The category of the pattern.
        elements:
            A mutable list of (category, stem) pairs that the new stems will
            be appended to.
        total_stems:
            The set of total stems before this function was invoked.
            Will not be mutated.

    Returns:
        The union of total_stems and stems found in the pattern.
    """

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
    """
    Prepares training data for training the model.

    Args:
        CategoryType:
            The category type to be used. Can be either MoodCategory,
            AffectionCategory or PatternCategory.
            Note that these are the enum types themselves and not enum values.
        elements:
            A list containg pairs of categories of type CategoryType and their
            associated set of stems.
        words: The bag of words used for indexing.

    Returns:
        The training axes x and y, where x is the feature axis and y is the
        label (= category) axis.
    """
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
    """
    Creates a sequential neural network model that takes a 1-dimensional array
    as input and outputs a number representing a category.

    Args:
        train_x: The feature axis for the model.
        train_y: The label axis for the model.

    Returns:
        An untrained sequential model.
    """

    # Define neural network

    # Amount of words in our vocabulary / bag of words
    num_words: int = len(train_x[0])
    # Amount of defined classes
    num_classes: int = len(train_y[0])

    # The amount of neurons to work with
    # https://stackoverflow.com/a/44748370
    # Most bag-of-words training examples use 512 here
    units: int = 512

    # Probability that a neuron will be ignored while processing input
    # http://papers.nips.cc/paper/4878-understanding-dropout.pdf suggests that
    # 50% gives the best results
    dropout_rate: float = 0.5

    model = Sequential()

    model.add(
        Dense(units, input_shape=(num_words,), activation='relu'))
    model.add(Dropout(dropout_rate))
    model.add(Dense(num_classes, activation='softmax'))

    return model
