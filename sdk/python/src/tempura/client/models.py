from dataclasses import dataclass
from enum import StrEnum, auto


class Status(StrEnum):
    NEW = auto()
    COMPLETED = auto()
    INCOMPLETE = auto()


@dataclass
class BatterRequest:
    function_name: str
    input_hash: str
    input: dict


@dataclass
class BatterResponse:
    status: Status
    output: dict | None = None
