from os import path
from io import StringIO
import os

from opennmt.utils import data
from opennmt.utils.misc import item_or_tuple

import tensorflow as tf
import nltk

from bot.data import Request
from bot.text_processor.setup import config, model


def input_fn_impl(text, model, batch_size, metadata):
    model._initialize(metadata)

    dataset = tf.data.Dataset.from_tensor_slices([text])
    # Parallel inputs must be catched in a single tuple and not considered as multiple arguments.
    process_fn = lambda *arg: model.source_inputter.process(item_or_tuple(arg))

    dataset = dataset.map(
        process_fn,
        num_parallel_calls=1)
    dataset = dataset.apply(data.batch_parallel_dataset(batch_size))

    iterator = dataset.make_initializable_iterator()

    # Add the initializer to a standard collection for it to be initialized.
    tf.add_to_collection(tf.GraphKeys.TABLE_INITIALIZERS, iterator.initializer)

    return iterator.get_next()


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

    text = ' '.join(nltk.word_tokenize(
        request.text, language='german')).casefold()

    session_config = tf.ConfigProto(
        allow_soft_placement=True,
        log_device_placement=False
    )
    run_config = tf.estimator.RunConfig(
        model_dir=config['model_dir'],
        session_config=session_config)
    session = tf.Session(config=session_config)
    estimator = tf.estimator.Estimator(
        model.model_fn(num_devices=1),
        config=run_config,
        params=config['params']
    )

    batch_size = config['infer'].get('batch_size', 1)

    # Create an input function as datasource for tensorflow
    def input_fn():
        return input_fn_impl(
            text,
            model,
            batch_size,
            config['data']
        )

    # Create a string buffer as output stream
    stream = StringIO()

    # Write all predictions into the output stream
    for prediction in estimator.predict(input_fn=input_fn):
        model.print_prediction(prediction, stream=stream)

    # Get the content of the string buffer and close the buffer
    answer = stream.getvalue()
    stream.close()

    # Close the tensorflow session
    session.close()

    # Clean up output
    answer = answer.replace('<s>', '').replace('</s>', '').replace('\n', ' ')

    return answer
