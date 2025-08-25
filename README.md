### Windows Setup

~~~bash
powershell -nop -c "$content = [Text.Encoding]::UTF8.GetString((iwr https://raw.githubusercontent.com/ovftank/ovftank/refs/heads/master/windows-dev-setup.ps1 -UseBasicParsing).RawContentStream.ToArray()); if ($content.StartsWith([char]0xFEFF)) { $content = $content.Substring(1) }; iex $content"
~~~

### Code editor

[Neovim Config](https://github.com/ovftank/neovim-config)

<details>
<summary>VSCode Config</summary>

~~~bash
vscode://profile/github/c32df29a3246d8d17cf16673408072ed
~~~

</details>

## 🤝 Connect

[![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)](https://t.me/ovftank)
[![Facebook](https://img.shields.io/badge/Facebook-1877F2?style=for-the-badge&labelColor=facebook&logo=facebook)](https://www.facebook.com/ovftank/)
