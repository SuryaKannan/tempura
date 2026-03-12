import dataclasses
import hashlib
import json

import httpx
from loguru import logger

from tempura.client.models import BatterRequest, BatterResponse, Status
from tempura.constants import TEMPURA_PORT


class _TempuraClient:
    def __init__(self):
        self.base_url = f"http://localhost:{TEMPURA_PORT}"
        self._client = httpx.Client(base_url=self.base_url)

    def register_batter(
        self, func_name: str, args: tuple, kwargs: dict
    ) -> BatterResponse:

        func_inputs, func_input_hash = _prepare_inputs(args=args, kwargs=kwargs)

        req = BatterRequest(
            function_name=func_name, input_hash=func_input_hash, input=func_inputs
        )
        response = self._client.post("/batter", json=dataclasses.asdict(req))
        response.raise_for_status()

        data = response.json()
        return BatterResponse(status=Status(data["status"]), output=data.get("output"))

    def health_check(self) -> bool:

        try:
            response = self._client.get("/health")
            response.raise_for_status()
        except httpx.ConnectError:
            return False
        return True


def _prepare_inputs(
    args: tuple,
    kwargs: dict,
) -> tuple[dict, str]:

    func_inputs = {"args": list(args), "kwargs": kwargs}

    input_json = json.dumps(func_inputs, sort_keys=True)
    return func_inputs, hashlib.sha256(input_json.encode()).hexdigest()


_client: _TempuraClient | None = None
_initialized: bool = False


def get_client() -> _TempuraClient | None:
    global _client, _initialized

    if _initialized:
        return _client

    _initialized = True

    client = _TempuraClient()
    if client.health_check():
        _client = client
        return _client

    logger.warning("Tempura service unhealthy! Executing standard flow")
    return None
