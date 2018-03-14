from .logger import logger
from .request import Request


def handle_request(request: Request):
    logger.debug('Handling request')
    raise NotImplementedError()
