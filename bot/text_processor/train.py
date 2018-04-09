from os import path, makedirs

from opennmt.runner import Runner

from bot.text_processor.setup import config, model


def train_and_evaluate():
    runner = Runner(model, config)
    runner.train_and_evaluate()
