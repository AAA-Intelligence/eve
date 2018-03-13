from .logger import logger
from .request import parse_request
from .request_handler import handle_request
import json

def run_loop():
    logger.debug('Starting loop')
    while True:
        try:
            json_data = input()
            request = parse_request(json_data)
            response = handle_request(request)
            print(json.dumps(response))
        except EOFError:
            # Stdin pipe has been closed by Go
            return
        except KeyboardInterrupt:
            # Interrupt requested by developer
            return
        except:
            # Pass error to Go and await next request
            print('error')
            continue
