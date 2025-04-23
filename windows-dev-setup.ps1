[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Warning "Vui long chay script nay voi quyen Administrator!"
    break
}
$path = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\VisualEffects"
Set-ItemProperty -Path "HKCU:\Control Panel\Mouse" -Name "MouseSpeed" -Value "0"
Set-ItemProperty -Path "HKCU:\Control Panel\Mouse" -Name "MouseThreshold1" -Value "0"
Set-ItemProperty -Path "HKCU:\Control Panel\Mouse" -Name "MouseThreshold2" -Value "0"
Set-ItemProperty -Path "HKCU:\Control Panel\Mouse" -Name "MouseSensitivity" -Value "6"
Set-ItemProperty -Path HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced -Name LaunchTo -Type DWord -Value 1
Set-ItemProperty -Path HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer -Name ShowRecent -Type DWord -Value 0
Set-ItemProperty -Path HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer -Name ShowFrequent -Type DWord -Value 0
Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced" -Name "Hidden" -Value 1
Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Explorer\Advanced" -Name "HideFileExt" -Value 0
Set-ItemProperty -Path "HKCU:\Control Panel\Accessibility\StickyKeys" -Name "Flags" -Type String -Value "506"
Set-ItemProperty -Path "HKCU:\Control Panel\Accessibility\ToggleKeys" -Name "Flags" -Type String -Value "58"
Set-ItemProperty -Path $path -Name VisualFXSetting -Value 2
Set-ItemProperty -Path "HKCU:\Control Panel\Accessibility\Keyboard Response" -Name "Flags" -Type String -Value "122"
Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\Dnscache\Parameters" -Name "MaxCacheTtl" -Type DWord -Value 86400
Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Services\Dnscache\Parameters" -Name "MaxNegativeCacheTtl" -Type DWord -Value 0
Stop-Service "WSearch" -Force
Set-Service "WSearch" -StartupType Disabled

powercfg -h off
$shell = New-Object -ComObject Shell.Application
$quickAccess = $shell.Namespace("shell:::{679F85CB-0220-4080-B29B-5540CC05AAB6}")
$items = $quickAccess.Items()

foreach ($item in $items) {
    $verb = ($item.Verbs() | Where-Object {$_.Name -eq "Unpin from Quick access"})
    if ($verb) {
        $verb.DoIt()
    }
}
Remove-Item -Path "$env:APPDATA\Microsoft\Windows\Recent\*" -Force -Recurse
Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer" -Name "ShowRecent" -Value 0
Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer" -Name "ShowFrequent" -Value 0
Set-ItemProperty -Path HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Search -Name SearchBoxTaskbarMode -Value 0 -Type DWord -Force

$gitName = Read-Host "Nhap username Git (bo qua neu khong muon cau hinh)"
$gitEmail = Read-Host "Nhap email Git (bo qua neu khong muon cau hinh)"

Write-Host "Dang cai dat..." -ForegroundColor Green

if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
    Write-Host "Dang cai dat Chocolatey..." -ForegroundColor Yellow

    $chocoDir = "C:\ProgramData\chocolatey"
    if (Test-Path $chocoDir) {
        Remove-Item -Path $chocoDir -Recurse -Force
    }

    Set-ExecutionPolicy Bypass -Scope Process -Force
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
    try {
        $ProgressPreference = 'SilentlyContinue'
        Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1')) *>&1 | Out-Null

        if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
            throw "Khong the cai dat Chocolatey"
        }
    }
    catch {
        Write-Host "Loi khi cai dat Chocolatey: $_" -ForegroundColor Red
        Write-Host "Vui long chay lai script." -ForegroundColor Red
        exit 1
    }
    finally {
        $ProgressPreference = 'Continue'
    }

    $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
}

if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Write-Host "Dang cai dat Git..." -ForegroundColor Yellow

    $gitUninstaller = "C:\Program Files\Git\unins000.exe"
    if (Test-Path $gitUninstaller) {
        Start-Process -FilePath $gitUninstaller -ArgumentList "/SILENT" -Wait
    }

    choco install git.install --params "'/GitAndUnixToolsOnPath /NoShellIntegration /NoGuiHereIntegration'" -y
}

Write-Host "Dang kiem tra pnpm..." -ForegroundColor Yellow
$pnpmPath = Join-Path $env:USERPROFILE "AppData\Local\pnpm\pnpm.exe"
if (Test-Path $pnpmPath) {
    Write-Host "pnpm da duoc cai dat, bo qua..." -ForegroundColor Green
}
else {
    Write-Host "Dang cai dat pnpm..." -ForegroundColor Yellow
    try {
        $ProgressPreference = 'SilentlyContinue'
        Invoke-WebRequest https://get.pnpm.io/install.ps1 -UseBasicParsing | Invoke-Expression *>&1 | Out-Null
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
        $ProgressPreference = 'Continue'

        if (Test-Path $pnpmPath) {
            & $pnpmPath env use --global iron *> $null
            Write-Host "Da cai dat pnpm thanh cong!" -ForegroundColor Green
        }
        else {
            Write-Host "Khong the cai dat pnpm" -ForegroundColor Red
        }
    }
    catch {
        Write-Host "Loi khi cai dat pnpm: $_" -ForegroundColor Red
    }
}

$pythonVersion = "3.10.11"
$pythonCommand = Get-Command python -ErrorAction SilentlyContinue
$currentPythonVersion = if ($pythonCommand) {
    python -c "import sys; print('.'.join(map(str, sys.version_info[:3])))"
}
else {
    $null
}

if (-not $pythonCommand -or $currentPythonVersion -ne $pythonVersion) {
    Write-Host "Dang tai xuong Python $pythonVersion..." -ForegroundColor Yellow
    $pythonUrl = "https://www.python.org/ftp/python/$pythonVersion/python-$pythonVersion-amd64.exe"
    $pythonInstaller = "$env:TEMP\python-$pythonVersion-amd64.exe"
    Invoke-WebRequest -Uri $pythonUrl -OutFile $pythonInstaller

    Write-Host "Dang cai dat Python $pythonVersion..." -ForegroundColor Yellow
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
    Remove-Item $pythonInstaller -Force

    $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
}
else {
    Write-Host "Python $pythonVersion da duoc cai dat." -ForegroundColor Green
}

if (-not [string]::IsNullOrWhiteSpace($gitName) -and -not [string]::IsNullOrWhiteSpace($gitEmail)) {
    git config --global user.name $gitName
    git config --global user.email $gitEmail
    git config --global init.defaultBranch main
    Write-Host "Da cau hinh Git thanh cong!" -ForegroundColor Green
}

$fontUrl = "https://github.com/ryanoasis/nerd-fonts/releases/download/v3.3.0/JetBrainsMono.zip"
$fontZip = "$env:TEMP\JetBrainsMono.zip"
$fontExtractPath = "$env:TEMP\JetBrainsMono"
$fontDestination = "$env:windir\Fonts"

Write-Host "Dang tai xuong JetBrains Mono Nerd Font..." -ForegroundColor Yellow
Invoke-WebRequest -Uri $fontUrl -OutFile $fontZip

if (-not (Test-Path $fontExtractPath)) {
    New-Item -ItemType Directory -Path $fontExtractPath | Out-Null
}

Write-Host "Dang giai nen font..." -ForegroundColor Yellow
Expand-Archive -Path $fontZip -DestinationPath $fontExtractPath -Force

$fonts = Get-ChildItem -Path $fontExtractPath -Include '*.ttf', '*.otf' -Recurse
foreach ($font in $fonts) {
    $destPath = Join-Path $fontDestination $font.Name
    if (-not (Test-Path $destPath)) {
        try {
            Copy-Item -Path $font.FullName -Destination $destPath -Force
        }
        catch {
            Write-Host "Khong the cai dat font $($font.Name): $_" -ForegroundColor Yellow
            continue
        }
    }
}

Remove-Item $fontZip -Force
Remove-Item $fontExtractPath -Recurse -Force

$fontRegistryPath = "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Fonts"
foreach ($font in $fonts) {
    $fontName = $font.Name
    $fontRegistryName = $fontName -replace '\.(ttf|otf)$', ' (TrueType)'
    New-ItemProperty -Path $fontRegistryPath -Name $fontRegistryName -Value $fontName -PropertyType String -Force | Out-Null
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

$clinkLatestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/chrisant996/clink/releases/latest"
$clinkExeLink = ($clinkLatestRelease.assets | Where-Object { $_.name -like "*setup.exe" }).browser_download_url
$clinkInstaller = "$env:TEMP\clink_setup.exe"

Invoke-WebRequest -Uri $clinkExeLink -OutFile $clinkInstaller

Start-Process -FilePath $clinkInstaller -ArgumentList "/S" -Wait

Start-Process -FilePath "${env:ProgramFiles(x86)}\clink\clink.bat" -ArgumentList "autorun install -- --quiet" -Wait

Remove-Item $clinkInstaller -Force

$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")

Write-Host "Dang tai xuong OpenKey..." -ForegroundColor Yellow
$openKeyUrl = "https://github.com/tuyenvm/OpenKey/releases/download/2.0.5/OpenKey64-Windows-2.0.5-RC.zip"
$openKeyZip = "$env:TEMP\OpenKey.zip"
$openKeyDestination = "$env:USERPROFILE\Documents\OpenKey"

Stop-Process -Name "OpenKey64" -Force -ErrorAction SilentlyContinue

Remove-Item -Path $openKeyDestination -Recurse -Force -ErrorAction SilentlyContinue
if (-not (Test-Path $openKeyDestination)) {
    New-Item -ItemType Directory -Path $openKeyDestination | Out-Null
}
Invoke-WebRequest -Uri $openKeyUrl -OutFile $openKeyZip
Write-Host "Dang giai nen OpenKey..." -ForegroundColor Yellow
Expand-Archive -Path $openKeyZip -DestinationPath $openKeyDestination -Force
Remove-Item $openKeyZip -Force
$WshShell = New-Object -comObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$env:USERPROFILE\Desktop\OpenKey.lnk")
$Shortcut.TargetPath = "$openKeyDestination\OpenKey64.exe"
$Shortcut.Save()
Write-Host "Dang cau hinh OpenKey..." -ForegroundColor Yellow
$regUrl = "https://raw.githubusercontent.com/ovftank/ovftank/refs/heads/master/OpenKey.reg"
$regFile = "$env:TEMP\OpenKey.reg"
Invoke-WebRequest -Uri $regUrl -OutFile $regFile
Start-Process "reg.exe" -ArgumentList "import", "`"$regFile`"" -Wait -NoNewWindow
Remove-Item $regFile -Force
Write-Host "Da cai dat va cau hinh OpenKey thanh cong!" -ForegroundColor Green

Write-Host "Dang cai dat Oh My Posh..." -ForegroundColor Yellow
choco install oh-my-posh -y *>&1 | Out-Null
$clinkPath = "${env:ProgramFiles(x86)}\clink\clink.bat"
if (Test-Path $clinkPath) {
    try {
        Start-Process "$clinkPath" -ArgumentList "config prompt use oh-my-posh"
        Start-Process "$clinkPath" -ArgumentList "set ohmyposh.theme `"${env:ProgramFiles(x86)}\oh-my-posh\themes\dracula.omp.json`""
        Write-Host "Da cau hinh Clink voi Oh My Posh thanh cong!" -ForegroundColor Green
    }
    catch {
        Write-Host "Khong the cau hinh Clink voi Oh My Posh" -ForegroundColor Yellow
    }
}
else {
    Write-Host "Khong tim thay Clink tai $clinkPath" -ForegroundColor Yellow
    Write-Host "Ban co the can cai dat lai Clink hoac cau hinh thu cong sau." -ForegroundColor Yellow
}

Write-Host "Dang cai dat Visual Studio Code..." -ForegroundColor Yellow
try {
    choco install vscode --params "/NoDesktopIcon /NoQuicklaunchIcon" -y *>&1 | Out-Null
    Write-Host "Da cai dat Visual Studio Code thanh cong!" -ForegroundColor Green
}
catch {
    Write-Host "Loi khi cai dat Visual Studio Code: $_" -ForegroundColor Red
}


$colorToolUrl = "https://raw.githubusercontent.com/waf/dracula-cmd/master/dist/ColorTool.zip"
$colorToolZip = "$env:TEMP\ColorTool.zip"
$colorToolPath = "$env:TEMP\ColorTool"

try {
    if (Test-Path $colorToolPath) {
        Remove-Item -Path $colorToolPath -Recurse -Force
    }
    New-Item -ItemType Directory -Path $colorToolPath -Force | Out-Null

    Invoke-WebRequest -Uri $colorToolUrl -OutFile $colorToolZip -UseBasicParsing
    Expand-Archive -Path $colorToolZip -DestinationPath $env:TEMP -Force

    Write-Host "Dang cai dat theme Dracula..." -ForegroundColor Yellow
    $colorToolExe = Join-Path $colorToolPath "ColorTool.exe"
    $installFolder = Join-Path $colorToolPath "install"

    if (Test-Path $colorToolExe) {
        Copy-Item -Path "$installFolder\Windows PowerShell.lnk" -Destination "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Windows PowerShell\" -Force
        Copy-Item -Path "$installFolder\Windows PowerShell (x86).lnk" -Destination "$env:APPDATA\Microsoft\Windows\Start Menu\Programs\Windows PowerShell\" -Force

        $regFile = Join-Path $installFolder "Remove Default Console Overrides.reg"
        Start-Process "reg.exe" -ArgumentList "import", "`"$regFile`"" -Wait -NoNewWindow

        Start-Process -FilePath $colorToolExe -ArgumentList "-b", "`"$installFolder\Dracula-ColorTool.itermcolors`"" -Wait -NoNewWindow

        Write-Host "Da cai dat theme Dracula thanh cong!" -ForegroundColor Green
    }
    else {
        Write-Host "Khong tim thay ColorTool.exe tai $colorToolExe" -ForegroundColor Red
    }
}
catch {
    Write-Host "Loi khi cai dat ColorTool: $_" -ForegroundColor Red
}
finally {
    if (Test-Path $colorToolZip) {
        Remove-Item $colorToolZip -Force
    }
    if (Test-Path $colorToolPath) {
        Remove-Item $colorToolPath -Recurse -Force
    }
}

Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize" -Name "EnableTransparency" -Value 0

Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize" -Name "AppsUseLightTheme" -Value 0
Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize" -Name "SystemUsesLightTheme" -Value 0
Stop-Process -Name explorer -Force
Write-Host "`nCai dat thanh cong!" -ForegroundColor Green
