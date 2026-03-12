import functools
import inspect

from loguru import logger

from tempura.client import get_client
from tempura.client.models import Status


def batter(func):

    @functools.wraps(func)
    def wrapper(*args, **kwargs):

        client = get_client()

        if not client:
            return func(*args, **kwargs)

        # get function defaults if they exist and apply defaults if no args passed
        sig = inspect.signature(func)
        bound = sig.bind(*args, **kwargs)
        bound.apply_defaults()

        execution = client.register_batter(func.__name__, bound.args, bound.kwargs)
        logger.info(execution)
        if execution.status == Status.COMPLETED:
            return execution.output

        return func(*args, **kwargs)

    return wrapper
