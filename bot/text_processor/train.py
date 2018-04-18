from os import path, makedirs

from opennmt.runner import Runner
import tensorflow as tf

from bot.text_processor.setup import config, model


def train_and_evaluate():
    tf.logging.set_verbosity(tf.logging.INFO)
    runner = Runner(model, config)
    runner.train()
