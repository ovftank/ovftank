import winreg
from typing import cast


def enable_regedit() -> bool:
    reg_path = r"Software\Microsoft\Windows\CurrentVersion\Policies\System"

    try:
        key = winreg.OpenKey(winreg.HKEY_CURRENT_USER, reg_path, 0, winreg.KEY_READ | winreg.KEY_WRITE)
        try:
            disable_value, _ = cast(tuple[int, int], winreg.QueryValueEx(key, "DisableRegistryTools"))
            if disable_value == 1:
                winreg.SetValueEx(key, "DisableRegistryTools", 0, winreg.REG_DWORD, 0)
            winreg.CloseKey(key)
            return True

        except FileNotFoundError:
            winreg.CloseKey(key)
            return True

    except FileNotFoundError:
        return True
    except PermissionError:
        return False
    except (OSError, ValueError, TypeError):
        return False