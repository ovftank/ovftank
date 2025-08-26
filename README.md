### Windows Setup

~~~bash
powershell -nop -c "$content = [Text.Encoding]::UTF8.GetString((iwr https://raw.githubusercontent.com/ovftank/ovftank/refs/heads/master/windows-dev-setup.ps1 -UseBasicParsing).RawContentStream.ToArray()); if ($content.StartsWith([char]0xFEFF)) { $content = $content.Substring(1) }; iex $content"
~~~

</details>

## 🤝 Connect

[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/ovftank)
[![Facebook](https://img.shields.io/badge/Facebook-1877F2?style=for-the-badge&labelColor=facebook&logo=facebook)](https://www.facebook.com/ovftank/)
