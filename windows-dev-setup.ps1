[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

$logFile = "$env:TEMP\windows-dev-setup.log"

# ================== UTILITY FUNCTIONS ==================
function Write-Log {
    param([string]$Message, [string]$Level = "INFO")
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $logMessage = "[$timestamp] [$Level] $Message"
    Write-Host $logMessage
    Add-Content -Path $logFile -Value $logMessage
}

function Test-Administrator {
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($currentUser)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Invoke-WithErrorHandling {
    param(
        [scriptblock]$ScriptBlock,
        [string]$ErrorMessage,
        [switch]$ContinueOnError
    )

    try {
        & $ScriptBlock
        return $true
    }
    catch {
        Write-Log "loi: $ErrorMessage - $($_.Exception.Message)" -Level "ERROR"
        if (-not $ContinueOnError) {
            throw
        }
        return $false
    }
}

function Install-ChocoPackage {
    param(
        [string]$PackageName,
        [string]$DisplayName = $PackageName,
        [string]$InstallArgs = "",
        [switch]$SkipCheck
    )

    if (-not $SkipCheck -and (Get-Command $PackageName -ErrorAction SilentlyContinue)) {
        Write-Log "$DisplayName có r" -Level "INFO"
        return
    }

    Write-Log "cài $DisplayName..." -Level "INFO"
    try {
        if ($InstallArgs) {
            choco install $PackageName --params "'$InstallArgs'" -y *>&1 | Out-Null
        } else {
            choco install $PackageName -y *>&1 | Out-Null
        }
        Write-Log "đã cài $DisplayName!" -Level "INFO"
    }
    catch {
        Write-Log "lỗi cài $DisplayName`: $($_)" -Level "ERROR"
    }
}

function Set-RegistryConfigs {
    param([array]$Configs)

    foreach ($config in $Configs) {
        try {
            if ($config.ContainsKey("Type")) {
                Set-ItemProperty -Path $config.Path -Name $config.Name -Value $config.Value -Type $config.Type -Force
            }
            else {
                Set-ItemProperty -Path $config.Path -Name $config.Name -Value $config.Value -Force
            }
        }
        catch {
            Write-Log "k config được $($config.Path)\$($config.Name): $($_.Exception.Message)" -Level "WARN"
        }
    }
}

# ================== MAIN EXECUTION ==================
if (-not (Test-Administrator)) {
    Write-Log "vui lòng chạy với quyền administrator!" -Level "ERROR"
    exit 1
}

Write-Host @"
              __ _              _
             / _| |            | |
   _____   _| |_| |_ __ _ _ __ | | __
  / _ \ \ / /  _| __/ _` | '_ \| |/ /
 | (_) \ V /| | | || (_| | | | |   <
  \___/ \_/ |_|  \__\__,_|_| |_|_|\_\


"@ -ForegroundColor Cyan

Write-Log "bắt đầu cài đặt..." -Level "INFO"

# ================== SYSTEM CONFIGURATION ==================
function Initialize-SystemConfig {
    Write-Log "bỏ chặn registry..." -Level "INFO"
    try {
        Set-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" -Name "DisableRegistryTools" -Value 0 -ErrorAction SilentlyContinue
        Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" -Name "DisableRegistryTools" -Value 0 -ErrorAction SilentlyContinue
        Write-Log "đã bỏ chặn registry!" -Level "INFO"
    }
    catch {
        Write-Log "lỗi registry: $($_.Exception.Message)" -Level "WARN"
    }

    Write-Log "config registry cơ bản..." -Level "INFO"
    $registryConfigs = @(
        @{Path = "HKCU:\Control Panel\Mouse"; Name = "MouseSpeed"; Value = "0" },
        @{Path = "HKCU:\Control Panel\Mouse"; Name = "MouseThreshold1"; Value = "0" },
        @{Path = "HKCU:\Control Panel\Mouse"; Name = "MouseThreshold2"; Value = "0" },
        @{Path = "HKCU:\Control Panel\Mouse"; Name = "MouseSensitivity"; Value = "6" },
        @{Path = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced"; Name = "LaunchTo"; Value = 1; Type = "DWord" },
        @{Path = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer"; Name = "ShowRecent"; Value = 0; Type = "DWord" },
        @{Path = "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer"; Name = "ShowFrequent"; Value = 0; Type = "DWord" },
        @{Path = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced"; Name = "Hidden"; Value = 1 },
        @{Path = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced"; Name = "HideFileExt"; Value = 0 },
        @{Path = "HKCU:\Control Panel\Accessibility\StickyKeys"; Name = "Flags"; Value = "506"; Type = "String" },
        @{Path = "HKCU:\Control Panel\Accessibility\ToggleKeys"; Name = "Flags"; Value = "58"; Type = "String" },
        @{Path = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects"; Name = "VisualFXSetting"; Value = 2 },
        @{Path = "HKCU:\Control Panel\Accessibility\Keyboard Response"; Name = "Flags"; Value = "122"; Type = "String" },
        @{Path = "HKLM:\SYSTEM\CurrentControlSet\Services\Dnscache\Parameters"; Name = "MaxCacheTtl"; Value = 86400; Type = "DWord" },
        @{Path = "HKLM:\SYSTEM\CurrentControlSet\Services\Dnscache\Parameters"; Name = "MaxNegativeCacheTtl"; Value = 0; Type = "DWord" }
    )
    Set-RegistryConfigs -Configs $registryConfigs
}

function Optimize-System {
    Write-Log "tắt windows search..." -Level "INFO"
    $searchService = Get-Service -Name "WSearch" -ErrorAction SilentlyContinue
    if ($searchService) {
        try {
            Stop-Service "WSearch" -Force -ErrorAction SilentlyContinue
            Set-Service "WSearch" -StartupType Disabled -ErrorAction SilentlyContinue
            Write-Log "đã tắt w!" -Level "INFO"
        }
        catch {
            Write-Log "k tắt được search: $($_.Exception.Message)" -Level "WARN"
        }
    }

    Write-Log "đang tắt hibernate..." -Level "INFO"
    Invoke-WithErrorHandling -ScriptBlock { powercfg -h off } -ErrorMessage "k tắt được hibernate" -ContinueOnError

    Write-Log "xoá quick access..." -Level "INFO"
    try {
        $shell = New-Object -ComObject Shell.Application
        $quickAccess = $shell.Namespace("shell:::{679F85CB-0220-4080-B29B-5540CC05AAB6}")
        $items = $quickAccess.Items()

        foreach ($item in $items) {
            $verb = ($item.Verbs() | Where-Object { $_.Name -eq "Unpin from Quick access" })
            if ($verb) {
                $verb.DoIt()
            }
        }
        Write-Log "đã xoá quick access!" -Level "INFO"
    }
    catch {
        Write-Log "k xoá được quick access: $($_.Exception.Message)" -Level "WARN"
    }
}

Initialize-SystemConfig
Optimize-System

# ================== PACKAGE MANAGERS ==================
function Install-Chocolatey {
    Write-Log "check & cài choco..." -Level "INFO"
    if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
        Write-Log "đang cài choco..." -Level "INFO"

        $chocoDir = "C:\ProgramData\chocolatey"
        if (Test-Path $chocoDir) {
            Remove-Item -Path $chocoDir -Recurse -Force -ErrorAction SilentlyContinue
        }

        try {
            Set-ExecutionPolicy Bypass -Scope Process -Force
            [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
            $ProgressPreference = 'SilentlyContinue'

            Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1')) | Out-Null

            if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
                throw "k cài được choco!"
            }

            $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
            Write-Log "đã cài choco!" -Level "INFO"
        }
        catch {
            Write-Log "lỗi cài choco: $($_)" -Level "ERROR"
            exit 1
        }
        finally {
            $ProgressPreference = 'Continue'
        }
    }
}

# ================== GIT CONFIGURATION ==================
function Install-Git {
    Write-Log "setup git..." -Level "INFO"
    $gitName = Read-Host "nhập username git (enter để skip)"
    $gitEmail = Read-Host "nhập email git (enter để skip)"

    if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
        Write-Log "đang cài git..." -Level "INFO"

        $gitUninstaller = "C:\Program Files\Git\unins000.exe"
        if (Test-Path $gitUninstaller) {
            Start-Process -FilePath $gitUninstaller -ArgumentList "/SILENT" -Wait -ErrorAction SilentlyContinue
        }

        try {
            choco install git.install --params "'/GitAndUnixToolsOnPath /NoShellIntegration /NoGuiHereIntegration'" -y --force | Out-Null
            $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
            Write-Log "đã cài git!" -Level "INFO"
        }
        catch {
            Write-Log "lỗi cài git: $($_)" -Level "ERROR"
        }
    }

    if (-not [string]::IsNullOrWhiteSpace($gitName) -and -not [string]::IsNullOrWhiteSpace($gitEmail)) {
        try {
            git config --global user.name $gitName
            git config --global user.email $gitEmail
            git config --global init.defaultBranch main
            git config --global push.autoSetupRemote true
            Write-Log "đã config git!" -Level "INFO"đã
        }
        catch {
            Write-Log "config git lỗi: $($_.Exception.Message)" -Level "WARN"
        }
    }
}

# ================== DEV TOOLS INSTALLATION ==================
function Install-DevTools {
    Write-Log "đang cài thư viện dev..." -Level "INFO"

    Install-ChocoPackage -PackageName "eza" -DisplayName "eza"
    Install-ChocoPackage -PackageName "ripgrep" -DisplayName "ripgrep (rg)"
    Install-ChocoPackage -PackageName "gzip" -DisplayName "gzip"
    Install-ChocoPackage -PackageName "mingw" -DisplayName "mingw"

    Install-ChocoPackage -PackageName "temurin" -DisplayName "Eclipse Temurin (Java)"

    Install-ChocoPackage -PackageName "neovim" -DisplayName "Neovim"
    Install-ChocoPackage -PackageName "neovide" -DisplayName "Neovide"

    Install-ChocoPackage -PackageName "lua" -DisplayName "Lua"
    Install-ChocoPackage -PackageName "luarocks" -DisplayName "LuaRocks"
}

# ================== RUNTIME INSTALLATIONS ==================
function Install-PNPM {
    Write-Log "check & cài pnpm..." -Level "INFO"
    $pnpmPath = Join-Path $env:USERPROFILE "AppData\Local\pnpm\pnpm.exe"
    if (Test-Path $pnpmPath) {
        Write-Log "pnpm có r, skip..." -Level "INFO"
    }
    else {
        Write-Log "đang cài pnpm..." -Level "INFO"

        try {
            $ProgressPreference = 'SilentlyContinue'
            Invoke-WebRequest https://get.pnpm.io/install.ps1 -UseBasicParsing | Invoke-Expression | Out-Null
            $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
            $ProgressPreference = 'Continue'

            if (Test-Path $pnpmPath) {
                Write-Log "đã cài pnpm!" -Level "INFO"

                try {
                    & $pnpmPath env use --global jod *>&1 | Out-Null
                    Write-Log "đã cài node jod!" -Level "INFO"
                }
                catch {
                    Write-Log "k cài được node: $($_.Exception.Message)" -Level "WARN"
                    Write-Log "có thể cần cài node trước" -Level "INFO"
                }
            }
            else {
                Write-Log "k cài được pnpm!" -Level "ERROR"
            }
        }
        catch {
            Write-Log "lỗi cài pnpm: $_" -Level "ERROR"
        }
    }
}

function Install-Python {
    Write-Log "check & cài python..." -Level "INFO"
    $pythonVersion = "3.10.11"
    $pythonCommand = Get-Command python -ErrorAction SilentlyContinue
    $currentPythonVersion = if ($pythonCommand) {
        try {
            python -c "import sys; print('.'.join(map(str, sys.version_info[:3])))"
        }
        catch {
            $null
        }
    }
    else {
        $null
    }

    if (-not $pythonCommand -or $currentPythonVersion -ne $pythonVersion) {
        Write-Log "đang tải python $pythonVersion..." -Level "INFO"
        try {
            $pythonUrl = "https://www.python.org/ftp/python/$pythonVersion/python-$pythonVersion-amd64.exe"
            $pythonInstaller = "$env:TEMP\python-$pythonVersion-amd64.exe"

            $ProgressPreference = 'SilentlyContinue'
            Invoke-WebRequest -Uri $pythonUrl -OutFile $pythonInstaller
            $ProgressPreference = 'Continue'

            Write-Log "đang cài python $pythonVersion..." -Level "INFO"
            $pythonArgs = @(
                "/quiet"
                "InstallAllUsers=1"
                "PrependPath=1"
                "AssociateFiles=1"
                "Include_pip=1"
                "Include_tcltk=1"
                "Include_test=0"
                "Include_doc=0"
                "Include_launcher=0"
                "InstallLauncherAllUsers=1"
                "Include_tools=1"
                "Shortcuts=0"
                "SimpleInstall=1"
            )
            Start-Process -FilePath $pythonInstaller -ArgumentList $pythonArgs -Wait
            Remove-Item $pythonInstaller -Force -ErrorAction SilentlyContinue

            $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
            Write-Log "đã cài python $pythonVersion!" -Level "INFO"
        }
        catch {
            Write-Log "lỗi cài python: $($_)" -Level "ERROR"
        }
    }
    else {
        Write-Log "python $pythonVersion có r." -Level "INFO"
    }
}

# ================== FONT INSTALLATION ==================
function Install-JetBrainsMonoFont {
    Write-Log "cài font JetBrains Mono..." -Level "INFO"
    try {
        $fontUrl = "https://github.com/ryanoasis/nerd-fonts/releases/download/v3.4.0/JetBrainsMono.zip"
        $fontZip = "$env:TEMP\JetBrainsMono.zip"
        $fontExtractPath = "$env:TEMP\JetBrainsMono"
        $fontDestination = "$env:windir\Fonts"

        $ProgressPreference = 'SilentlyContinue'
        Invoke-WebRequest -Uri $fontUrl -OutFile $fontZip -UseBasicParsing
        $ProgressPreference = 'Continue'

        if (-not (Test-Path $fontExtractPath)) {
            New-Item -ItemType Directory -Path $fontExtractPath | Out-Null
        }

        Write-Log "đang extract font..." -Level "INFO"
        Expand-Archive -Path $fontZip -DestinationPath $fontExtractPath -Force | Out-Null

        $fonts = Get-ChildItem -Path $fontExtractPath -Include '*.ttf', '*.otf' -Recurse
        $installedCount = 0

        foreach ($font in $fonts) {
            $destPath = Join-Path $fontDestination $font.Name
            if (-not (Test-Path $destPath)) {
                try {
                    Copy-Item -Path $font.FullName -Destination $destPath -Force | Out-Null
                    $installedCount++
                }
                catch {
                    Write-Log "k cài được font $($font.Name): $_" -Level "WARN"
                }
            }
        }

        $fontRegistryPath = "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts"
        foreach ($font in $fonts) {
            $fontName = $font.Name
            $fontRegistryName = $fontName -replace '\.(ttf|otf)$', ' (TrueType)'
            try {
                New-ItemProperty -Path $fontRegistryPath -Name $fontRegistryName -Value $fontName -PropertyType String -Force | Out-Null
            }
            catch {
                Write-Log "k update registry cho font $fontName" -Level "WARN"
            }
        }

        $registryPath = "HKCU:\Console"
        $fontName = "JetBrainsMono NFM"
        $fontSize = 0x140000

        if (!(Test-Path $registryPath)) {
            New-Item -Path $registryPath -Force | Out-Null
        }

        Set-ItemProperty -Path $registryPath -Name "FaceName" -Value $fontName -Type STRING
        Set-ItemProperty -Path $registryPath -Name "FontSize" -Value $fontSize -Type DWORD
        Set-ItemProperty -Path $registryPath -Name "WindowAlpha" -Value 242 -Type DWORD

        Remove-Item $fontZip -Force -ErrorAction SilentlyContinue
        Remove-Item $fontExtractPath -Recurse -Force -ErrorAction SilentlyContinue

        Write-Log "đã cài $installedCount fonts!" -Level "INFO"
    }
    catch {
        Write-Log "lỗi cài fonts: $($_)" -Level "ERROR"
    }
}

# ================== TERMINAL TOOLS ==================
function Install-Clink {
    Write-Log "cài clink..." -Level "INFO"
    try {
        $clinkLatestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/chrisant996/clink/releases/latest"
        $clinkExeLink = ($clinkLatestRelease.assets | Where-Object { $_.name -like "*setup.exe" }).browser_download_url
        $clinkInstaller = "$env:TEMP\clink_setup.exe"

        $ProgressPreference = 'SilentlyContinue'
        Invoke-WebRequest -Uri $clinkExeLink -OutFile $clinkInstaller -UseBasicParsing
        $ProgressPreference = 'Continue'

        Start-Process -FilePath $clinkInstaller -ArgumentList "/S" -Wait | Out-Null

        Write-Log "tạo aliases..." -Level "INFO"
        try {
            $aliasesDir = "$env:LOCALAPPDATA\clink"
            $aliasesFile = "$aliasesDir\aliases"

            if (-not (Test-Path $aliasesDir)) {
                New-Item -ItemType Directory -Path $aliasesDir -Force | Out-Null
            }

            "ls=eza --icons $*" | Out-File -FilePath $aliasesFile -Encoding ASCII
            Write-Log "đã tạo aliases!" -Level "INFO"
        }
        catch {
            Write-Log "lỗi tạo aliases: $($_)" -Level "ERROR"
        }

        Write-Log "config clink autorun..." -Level "INFO"
        try {
            $autorunValue = "doskey /macrofile=%LOCALAPPDATA%\clink\aliases&`"${env:ProgramFiles(x86)}\clink\clink.bat`" inject --autorun --quiet"
            Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Command Processor" -Name "Autorun" -Value $autorunValue -Type String
            Write-Log "đã config autorun!" -Level "INFO"
        }
        catch {
            Write-Log "lỗi config autorun: $($_)" -Level "ERROR"
        }

        Remove-Item $clinkInstaller -Force -ErrorAction SilentlyContinue
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")

        Write-Log "đã cài clink!" -Level "INFO"
    }
    catch {
        Write-Log "lỗi cài clink: $($_)" -Level "ERROR"
    }
}

function Install-OpenKey {
    Write-Log "cài openkey..." -Level "INFO"
    try {
        $openKeyUrl = "https://github.com/tuyenvm/OpenKey/releases/download/2.0.5/OpenKey64-Windows-2.0.5-RC.zip"
        $openKeyZip = "$env:TEMP\OpenKey.zip"
        $openKeyDestination = "$env:USERPROFILE\Documents\OpenKey"

        Stop-Process -Name "OpenKey64" -Force -ErrorAction SilentlyContinue

        Remove-Item -Path $openKeyDestination -Recurse -Force -ErrorAction SilentlyContinue
        if (-not (Test-Path $openKeyDestination)) {
            New-Item -ItemType Directory -Path $openKeyDestination | Out-Null
        }

        $ProgressPreference = 'SilentlyContinue'
        Invoke-WebRequest -Uri $openKeyUrl -OutFile $openKeyZip -UseBasicParsing
        $ProgressPreference = 'Continue'

        Write-Log "đang extract openkey..." -Level "INFO"
        Expand-Archive -Path $openKeyZip -DestinationPath $openKeyDestination -Force
        Remove-Item $openKeyZip -Force -ErrorAction SilentlyContinue

        $WshShell = New-Object -comObject WScript.Shell
        $Shortcut = $WshShell.CreateShortcut("$env:USERPROFILE\Desktop\OpenKey.lnk")
        $Shortcut.TargetPath = "$openKeyDestination\OpenKey64.exe"
        $Shortcut.Save()

        Write-Log "đang config openkey..." -Level "INFO"
        $regUrl = "https://raw.githubusercontent.com/ovftank/ovftank/refs/heads/master/OpenKey.reg"
        $regFile = "$env:TEMP\OpenKey.reg"

        $ProgressPreference = 'SilentlyContinue'
        Invoke-WebRequest -Uri $regUrl -OutFile $regFile -UseBasicParsing
        $ProgressPreference = 'Continue'

        Start-Process "reg.exe" -ArgumentList "import", "`"$regFile`"" -Wait -NoNewWindow | Out-Null
        Remove-Item $regFile -Force -ErrorAction SilentlyContinue

        Write-Log "đã cài openkey!" -Level "INFO"
    }
    catch {
        Write-Log "lỗi cài openkey: $($_)" -Level "ERROR"
    }
}

# ================== THEME & UI CUSTOMIZATION ==================
function Install-OhMyPosh {
    Write-Log "cài oh my posh..." -Level "INFO"
    try {
        choco install oh-my-posh -y --force | Out-Null

        $clinkPath = "${env:ProgramFiles(x86)}\clink\clink.bat"
        if (Test-Path $clinkPath) {
            Start-Process "$clinkPath" -ArgumentList "config prompt use oh-my-posh"

            $themePath = "${env:ProgramFiles(x86)}\oh-my-posh\themes\dracula.omp.json"
            if (Test-Path $themePath) {
                Start-Process "$clinkPath" -ArgumentList "set ohmyposh.theme `"$themePath`""
                Write-Log "đã config clink với dracula theme!" -Level "INFO"
            }
            else {
                Write-Log "k thấy theme dracula ở $themePath" -Level "WARN"
            }
        }
        else {
            Write-Log "k thấy clink ở $clinkPath" -Level "WARN"
        }
    }
    catch {
        Write-Log "lỗi cài oh my posh: $($_)" -Level "ERROR"
    }
}

function Install-DraculaTheme {
    Write-Log "cài colortool & dracula theme..." -Level "INFO"
    try {
        $colorToolUrl = "https://raw.githubusercontent.com/waf/dracula-cmd/master/dist/ColorTool.zip"
        $colorToolZip = "$env:TEMP\ColorTool.zip"
        $colorToolPath = "$env:TEMP\ColorTool"

        if (Test-Path $colorToolZip) {
            Remove-Item $colorToolZip -Force -ErrorAction SilentlyContinue
        }
        if (Test-Path $colorToolPath) {
            Remove-Item -Path $colorToolPath -Recurse -Force -ErrorAction SilentlyContinue
        }

        $ProgressPreference = 'SilentlyContinue'
        Invoke-WebRequest -Uri $colorToolUrl -OutFile $colorToolZip -UseBasicParsing
        $ProgressPreference = 'Continue'

        Expand-Archive -Path $colorToolZip -DestinationPath $colorToolPath -Force

        Write-Log "đang cài dracula theme..." -Level "INFO"
        $colorToolExe = Join-Path $colorToolPath "ColorTool.exe"
        $installFolder = Join-Path $colorToolPath "install"

        if (Test-Path $colorToolExe) {
            $shortcutPath = "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Windows PowerShell\"
            if (Test-Path "$installFolder\Windows PowerShell.lnk") {
                Copy-Item -Path "$installFolder\Windows PowerShell.lnk" -Destination $shortcutPath -Force -ErrorAction SilentlyContinue
            }
            if (Test-Path "$installFolder\Windows PowerShell (x86).lnk") {
                Copy-Item -Path "$installFolder\Windows PowerShell (x86).lnk" -Destination $shortcutPath -Force -ErrorAction SilentlyContinue
            }

            $regFile = Join-Path $installFolder "Remove Default Console Overrides.reg"
            if (Test-Path $regFile) {
                Start-Process "reg.exe" -ArgumentList "import", "`"$regFile`"" -Wait -NoNewWindow -ErrorAction SilentlyContinue
            }

            $themeFile = Join-Path $installFolder "Dracula-ColorTool.itermcolors"
            if (Test-Path $themeFile) {
                Start-Process -FilePath $colorToolExe -ArgumentList "-b", "`"$themeFile`"" -Wait -NoNewWindow -ErrorAction SilentlyContinue
            }

            Write-Log "đã cài dracula theme!" -Level "INFO"
        }
        else {
            Write-Log "k thấy colortool.exe, có thể file zip bị lỗi" -Level "WARN"
        }
    }
    catch {
        Write-Log "lỗi cài colortool: $_" -Level "ERROR"
    }
    finally {
        if (Test-Path $colorToolZip) {
            Remove-Item $colorToolZip -Force -ErrorAction SilentlyContinue
        }
        if (Test-Path $colorToolPath) {
            Remove-Item $colorToolPath -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

function Set-DarkTheme {
    Write-Log "config dark theme..." -Level "INFO"
    try {
        Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize" -Name "EnableTransparency" -Value 0
        Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize" -Name "AppsUseLightTheme" -Value 0
        Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize" -Name "SystemUsesLightTheme" -Value 0
            Write-Log "đã config dark theme!" -Level "INFO"
    }
    catch {
        Write-Log "k config được dark theme: $($_)" -Level "WARN"
    }
}

function Add-ThisPCFolders {
    Write-Log "thêm folders vào this pc..." -Level "INFO"
    $folders = @{
        "Downloads" = "{088e3905-0323-4b02-9826-5d99428e115f}"
        "Pictures"  = "{24ad3ad4-a569-4530-98e1-ab02f9417aa8}"
        "Music"     = "{3dfdf296-dbec-4fb4-81d1-6a3438bcf4de}"
        "Videos"    = "{f86fa3ab-70d2-4fc7-9c99-fcbf05467f3a}"
        "Documents" = "{d3162b92-9365-467a-956b-92703aca08af}"
        "Desktop"   = "{B4BFCC3A-DB2C-424C-B029-7FE99A87C641}"
    }

    foreach ($folder in $folders.GetEnumerator()) {
        try {
            $guid = $folder.Value
            New-Item -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\$guid" -Force | Out-Null
            New-Item -Path "HKLM:\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\$guid" -Force | Out-Null
        }
        catch {
            Write-Log "k thêm được $($folder.Key) vào this pc: $($_)" -Level "WARN"
        }
    }
}

# ================== CODE EDITOR SELECTION ==================
function Install-CodeEditor {
    Write-Log "chọn code editor..." -Level "INFO"
    Write-Host "`nChọn code editor muốn cài:" -ForegroundColor Yellow
    Write-Host "1. Cursor" -ForegroundColor Cyan
    Write-Host "2. VS Code" -ForegroundColor Cyan
    Write-Host "3. Skip" -ForegroundColor Cyan

    do {
        $choice = Read-Host "`nNhập lựa chọn (1-3)"
        switch ($choice) {
            "1" {
                Write-Log "cài cursor..." -Level "INFO"
                try {
                    choco install cursor -y --force | Out-Null
                    Write-Log "đã cài cursor!" -Level "INFO"
                }
                catch {
                    Write-Log "lỗi cài cursor: $($_)" -Level "ERROR"
                }
                break
            }
            "2" {
                Write-Log "cài vscode..." -Level "INFO"
                try {
                    choco install vscode -y --force | Out-Null
                    Write-Log "đã cài vscode!" -Level "INFO"
                }
                catch {
                    Write-Log "lỗi cài vscode: $($_)" -Level "ERROR"
                }
                break
            }
            "3" {
                Write-Log "skip" -Level "INFO"
                break
            }
            default {
                Write-Host "Chọn từ 1-3 thôi!" -ForegroundColor Red
            }
        }
    } while ($choice -notmatch "^[1-3]$")
}

# ================== EXECUTION WORKFLOW ==================
Install-Chocolatey
Install-Git
Install-DevTools

Install-PNPM
Install-Python

Install-CodeEditor
Install-JetBrainsMonoFont
Install-Clink
Install-OpenKey
Install-OhMyPosh
Install-DraculaTheme
Set-DarkTheme
Add-ThisPCFolders

Stop-Process -Name explorer -Force -ErrorAction SilentlyContinue

Write-Host "chi tiết log: $logFile" -ForegroundColor Cyan
