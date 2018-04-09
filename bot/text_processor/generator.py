from os import path
from io import StringIO
import os

from opennmt.runner import Runner
import tensorflow as tf
import nltk

from bot.data import Request
from bot.text_processor.setup import config, model


def generate_answer(request: Request) -> str:
    """
    Generates an answer for a given request.

    Arguments:
        request: The request to generate an answer for.
        mood: The mood returned by the mood analyzer.
        affection: The affection returned by the affection analyzer.

    Returns:
        The generated answer.
    """

    text = ' '.join(nltk.word_tokenize(request.text, language='german'))
    with open('theinput.txt', 'w', encoding='utf-8') as f:
        f.write(text)

    runner = Runner(model, config)
    estimator = runner._estimator

    batch_size = config['infer'].get('batch_size', 1)
    input_fn = model.input_fn(
        tf.estimator.ModeKeys.PREDICT,
        batch_size,
        config['data'],
        'theinput.txt',
        prefetch_buffer_size=1
    )

    stream = StringIO()
    for prediction in estimator.predict(input_fn=input_fn):
        model.print_prediction(prediction, stream=stream)

    answer = stream.getvalue()
    stream.close()

    answer = answer.replace('<s>', '').replace('</s>', '').replace('\n', ' ')

    return answer
