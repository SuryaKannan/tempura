import functools

from tempura.client import get_client
from tempura.client.models import Status


def batter(func):

    @functools.wraps(func)
    def wrapper(*args, **kwargs):

        client = get_client()

        if not client:
            return func(*args, **kwargs)

        execution = client.register_batter(func.__name__, args, kwargs)

        if execution.status == Status.COMPLETED:
            return execution.output

        return func(*args, **kwargs)

    return wrapper
