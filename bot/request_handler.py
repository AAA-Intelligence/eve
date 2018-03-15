from .logger import logger
from .data import Request, Response, parse_request
from .mood_analyzer import analyze
from .pattern_recognizer import recognize_pattern
from .text_processor import generate_answer
import json


def handle_request(request: Request) -> Response:
    """
    Handles a request and generates an appropriate response.

    Args:
        request: The request to handle.

    Returns:
        The response generated for the specified request.
    """

    logger.debug('Handling request')

    (mood, affection) = analyze(request.text)

    answer = recognize_pattern(request)
    if answer is None:
        # No pattern found, fall back to generative model
        answer = generate_answer(request, mood, affection)

    return Response(answer, mood, affection)


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

    logger.debug('Starting loop')
    while True:
        try:
            json_data = input()
            request = parse_request(json_data)
            response = handle_request(request)
            print(json.dumps(response))
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
