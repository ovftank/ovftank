import ctypes
from typing import cast


def is_admin() -> bool:
    try:
        result = cast(int, ctypes.windll.shell32.IsUserAnAdmin())
        return bool(result)
    except (OSError, AttributeError):
        return False
