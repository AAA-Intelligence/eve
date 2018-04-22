from os import path, makedirs
import yaml

from opennmt.models import SequenceToSequence
from opennmt.config import load_model


def load_config(config_path: str):
    """
    Loads an OpenNMT config file

    Arguments:
        config_path: The path to the config file

    Returns:
        A dict containing the config data
    """

    with open(config_path, encoding='utf-8') as f:
        return yaml.load(f)


__dir = path.dirname(__file__)
config_path = path.join(__dir, 'config.yml')
model_file = path.join(__dir, 'nmt_small.py')

config = load_config(config_path)

model_dir = config['model_dir']
if not path.isdir(model_dir):
    makedirs(model_dir)
model = load_model(model_dir, model_file=model_file)
