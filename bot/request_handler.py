import json
from datetime import date

from bot.data import Request, Response, parse_request
from bot.logger import logger
from bot.mood_analyzer import analyze
from bot.pattern_recognizer import answer_for_pattern
from bot.text_processor.generator import generate_answer


def handle_request(request: Request) -> Response:
    """
    Handles a request and generates an appropriate response.

    Args:
            request: The request to handle.

    Returns:
            The response generated for the specified request.
    """

    logger.debug('Handling request')

    """
        mood, affection : value between -1 and 1 indicating a positive or negative sentiment
       
    """
    mood_message, affection_message, mood_bot, affection_bot = analyze(request)
    # TODO is analyzed mood/affection (PredictionResult) necessary
    result = answer_for_pattern(request)
    if result:
        pattern, answer = result
    else:
        # No pattern found, fall back to generative model
        pattern = None
        answer = generate_answer(request, mood_message, affection_message)
    response = Response(text=answer,
                        pattern=pattern,
                        mood=mood_bot,
                        affection=affection_bot)
    logger.debug(response)
    return response


def run_demo():
    """
    Starts a command-line based demo request loop for debugging.
    """
    logger.info('Starting request loop')
    previous_pattern = None
    while True:
        try:
            text = input('User input: ')
            request = Request(
                text=text,
                previous_pattern=previous_pattern,
                mood=0.0,
                affection=0.0,
                bot_gender=0,
                bot_name='Lana',
                bot_birthdate=date(1995, 10, 5),
                bot_favorite_color='grün',
                father_name='Georg',
                father_age=49,
                mother_name='Agathe',
                mother_age=47
            )
            response = handle_request(request)
            print('Response: ', response.text)
            previous_pattern = response.pattern
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
