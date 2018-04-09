from os import path
from io import StringIO

from opennmt.models import SequenceToSequence
from opennmt.config import load_model
import tensorflow as tf

from bot.trainer import dir as models_dir
from bot.data import Request

batch_size = 30

chat_model_dir = path.join(models_dir, 'chat')
model: SequenceToSequence = load_model(chat_model_dir)


def generate_answer(request: Request, mood: float, affection: float) -> str:
    """
    Generates an answer for a given request.

    Arguments:
        request: The request to generate an answer for.
        mood: The mood returned by the mood analyzer.
        affection: The affection returned by the affection analyzer.

    Returns:
        The generated answer.
    """

    input_fn = model.input_fn(
        tf.estimator.ModeKeys.PREDICT,
        batch_size,
        {},
        None,
        prefetch_buffer_size=1
    )

    stream = StringIO()
    model.print_prediction(p, stream=stream)

    # TODO: Implement

    return 'I solemnly swear that I am no bot'
