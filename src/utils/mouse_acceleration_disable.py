import ctypes
import subprocess
import winreg


def mouse_acceleration_disable() -> bool:
    reg_path = r"Control Panel\Mouse"

    mouse_settings = {
        "MouseSensitivity": "10",
        "MouseSpeed": "0",
        "MouseThreshold1": "0",
        "MouseThreshold2": "0"
    }

    try:
        key = winreg.OpenKey(winreg.HKEY_CURRENT_USER, reg_path, 0, winreg.KEY_READ | winreg.KEY_WRITE)

        for setting_name, value in mouse_settings.items():
            try:
                winreg.SetValueEx(key, setting_name, 0, winreg.REG_SZ, value)
            except FileNotFoundError:
                winreg.SetValueEx(key, setting_name, 0, winreg.REG_SZ, value)

        winreg.CloseKey(key)

        return True

    except FileNotFoundError:
        return False
    except PermissionError:
        return False
    except (OSError, ValueError, TypeError):
        return False
