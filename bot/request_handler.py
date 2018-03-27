from .logger import logger
from .data import Request, Response, parse_request
from .mood_analyzer import analyze
from .pattern_recognizer import answer_for_pattern
from .text_processor import generate_answer
from datetime import date
import json
import time


def handle_request(request: Request) -> Response:
    """
    Handles a request and generates an appropriate response.

    Args:
        request: The request to handle.

    Returns:
        The response generated for the specified request.
    """

    logger.debug('Handling request')

    mood, affection = analyze(request.text)

    answer = answer_for_pattern(request)
    if answer is None:
        # No pattern found, fall back to generative model
        answer = generate_answer(request, mood, affection)

    return Response(answer, mood, affection)


def run_demo():
    """
    Starts a command-line based demo request loop for debugging.
    """
    logger.info('Starting request loop')
    while True:
        try:
            request = Request(
                text=input('User input: '),
                previous_text='Ich bin ein Baum',
                mood=0.0,
                affection=0.0,
                bot_gender=0,
                bot_name='Lara',
                bot_birthdate=date(1995, 10, 5),
                bot_favorite_color='grün'
            )
            response = handle_request(request)
            print('Response: ', response.text)
        except KeyboardInterrupt:
            # Interrupt requested by user
            logger.info('Keyboard interrupt detected, aborting request loop')
            return
        except Exception as ex:
            logger.error('{}: {}'.format(type(ex).__name__, str(ex)))
            continue


def run_loop():
    """
    Starts a request loop that reads lines from stdin.

    Each line represents a new request in JSON format that will be parsed by the
    loop and handled by the request handler. The response returned by the
    request handler will again be formatted as a JSON string and written to
    stdout, including a newline character after every response.
    If an error is raised during parsing of the request data or the request
    handling itself, the current request will be aborted which is signaled by
    the 'error\n' string written to stdout. The loop will then wait for a new
    request.

    The loop can be interrupted by either closing the stdin pipe, resulting in
    an EOFError handled by the loop, or by sending a
    keyboard interrupt (Ctrl + C).
    """

    logger.info('Starting request loop')
    while True:
        try:
            logger.info('Waiting for request input')
            json_data = input()
            logger.info('Received request, parsing')
            request = parse_request(json_data)
            logger.info('Handling request: {}'.format(request))
            response = handle_request(request)
            print(json.dumps(response._asdict()))
        except EOFError:
            # Stdin pipe has been closed by Go
            logger.info('EOF detected, aborting request loop')
            return
        except KeyboardInterrupt:
            # Interrupt requested by developer
            logger.info('Keyboard interrupt detected, aborting request loop')
            return
        except Exception as ex:
            logger.error('{}: {}'.format(type(ex).__name__, str(ex)))
            # Pass error to Go and await next request
            print('error')
            continue
