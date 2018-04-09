from os import path, makedirs

from opennmt.models import SequenceToSequence
from opennmt.config import load_config, load_model


__dir = path.dirname(__file__)
config_path = path.join(__dir, 'config.yml')
model_file = path.join(__dir, 'models', 'transformer.py')

config = load_config([config_path])

model_dir = config['model_dir']
if not path.isdir(model_dir):
    makedirs(model_dir)
model = load_model(model_dir, model_file=model_file)
