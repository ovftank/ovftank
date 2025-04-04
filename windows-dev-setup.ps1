[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Warning "Vui long chay script nay voi quyen Administrator!"
    break
}

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
        Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))

        if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
            throw "Khong the cai dat Chocolatey"
        }
    }
    catch {
        Write-Host "Loi khi cai dat Chocolatey: $_" -ForegroundColor Red
        Write-Host "Vui long chay lai script." -ForegroundColor Red
        exit 1
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

Write-Host "Dang cai dat pnpm..." -ForegroundColor Yellow
$pnpmOutput = Invoke-WebRequest https://get.pnpm.io/install.ps1 -UseBasicParsing | Invoke-Expression *>&1
$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")

Write-Host "Dang cai dat NodeJS Iron..." -ForegroundColor Yellow
$nodeOutput = pnpm env use --global iron *>&1

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

$gitName = Read-Host "Nhap username Git (bo qua neu khong muon cau hinh)"
$gitEmail = Read-Host "Nhap email Git (bo qua neu khong muon cau hinh)"

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
$webClient = New-Object System.Net.WebClient
$webClient.DownloadProgressChanged = {
    $downloadedMB = [math]::Round($EventArgs.BytesReceived / 1MB, 2)
    $totalMB = [math]::Round($EventArgs.TotalBytesToReceive / 1MB, 2)
    $percentComplete = $EventArgs.ProgressPercentage

    Write-Progress -Activity "Dang tai xuong JetBrains Mono Nerd Font" `
        -Status "$downloadedMB MB / $totalMB MB" `
        -PercentComplete $percentComplete
}
$webClient.DownloadFileCompleted = {
    Write-Progress -Activity "Dang tai xuong JetBrains Mono Nerd Font" -Completed
}

$webClient.DownloadFileAsync($fontUrl, $fontZip)
do { Start-Sleep -Milliseconds 100 }
while ($webClient.IsBusy)

if (-not (Test-Path $fontExtractPath)) {
    New-Item -ItemType Directory -Path $fontExtractPath | Out-Null
}

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
Set-ItemProperty -Path $registryPath -Name "WindowAlpha" -Value 230 -Type DWORD

$clinkLatestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/chrisant996/clink/releases/latest"
$clinkExeLink = ($clinkLatestRelease.assets | Where-Object { $_.name -like "*setup.exe" }).browser_download_url
$clinkInstaller = "$env:TEMP\clink_setup.exe"

Invoke-WebRequest -Uri $clinkExeLink -OutFile $clinkInstaller

Start-Process -FilePath $clinkInstaller -ArgumentList "/S" -Wait

Remove-Item $clinkInstaller -Force

$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")

python -m pip install -q cursor-manager

Write-Host "Dang cai dat Cursor..." -ForegroundColor Yellow
cursor-manager tat-update
cursor-manager downgrade


Write-Host "Dang tai xuong EVKey..." -ForegroundColor Yellow
$evkeyUrl = "https://github.com/lamquangminh/EVKey/releases/download/Release/EVKey.zip"
$evkeyZip = "$env:TEMP\EVKey.zip"
$evkeyDestination = "$env:USERPROFILE\Documents\EVKey"

Stop-Process -Name "EVKey64" -Force -ErrorAction SilentlyContinue

Remove-Item -Path $evkeyDestination -Recurse -Force -ErrorAction SilentlyContinue

Invoke-WebRequest -Uri $evkeyUrl -OutFile $evkeyZip
$evkeySetting = "$evkeyDestination\setting.ini"

if (-not (Test-Path $evkeyDestination)) {
    New-Item -ItemType Directory -Path $evkeyDestination | Out-Null
}

Expand-Archive -Path $evkeyZip -DestinationPath $evkeyDestination -Force

Invoke-WebRequest -Uri "https://raw.githubusercontent.com/ovftank/ovftank/refs/heads/master/setting.ini" -OutFile $evkeySetting

Remove-Item $evkeyZip -Force

$WshShell = New-Object -comObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$env:USERPROFILE\Desktop\EVKey.lnk")
$Shortcut.TargetPath = "$evkeyDestination\EVKey64.exe"
$Shortcut.Save()

Write-Host "Dang cai dat Oh My Posh..." -ForegroundColor Yellow
choco install oh-my-posh -y --force

$clinkPath = "${env:ProgramFiles(x86)}\clink\clink.bat"
if (Test-Path $clinkPath) {
    try {
        & "$clinkPath" config prompt use oh-my-posh
        & "$clinkPath" set ohmyposh.theme "${env:ProgramFiles(x86)}\oh-my-posh\themes\dracula.omp.json"
        Write-Host "Da cau hinh Clink voi Oh My Posh thanh cong!" -ForegroundColor Green
    }
    catch {
        Write-Host "Khong the cau hinh Clink voi Oh My Posh: $_" -ForegroundColor Yellow
        Write-Host "Ban co the can cau hinh thu cong sau." -ForegroundColor Yellow
    }
}
else {
    Write-Host "Khong tim thay Clink tai $clinkPath" -ForegroundColor Yellow
    Write-Host "Ban co the can cai dat lai Clink hoac cau hinh thu cong sau." -ForegroundColor Yellow
}

$colorToolUrl = "https://raw.githubusercontent.com/waf/dracula-cmd/master/dist/ColorTool.zip"
$colorToolZip = "$env:TEMP\ColorTool.zip"
$colorToolPath = "$env:TEMP\ColorTool"

Invoke-WebRequest -Uri $colorToolUrl -OutFile $colorToolZip

if (-not (Test-Path $colorToolPath)) {
    New-Item -ItemType Directory -Path $colorToolPath | Out-Null
}

Expand-Archive -Path $colorToolZip -DestinationPath $colorToolPath -Force

Start-Process -FilePath "$colorToolPath\install.cmd" -Wait -NoNewWindow

Remove-Item $colorToolZip -Force
Remove-Item $colorToolPath -Recurse -Force

Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize" -Name "EnableTransparency" -Value 0

Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize" -Name "AppsUseLightTheme" -Value 0
Set-ItemProperty -Path "HKCU:\SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize" -Name "SystemUsesLightTheme" -Value 0

Write-Host "`nCai dat thanh cong!" -ForegroundColor Green