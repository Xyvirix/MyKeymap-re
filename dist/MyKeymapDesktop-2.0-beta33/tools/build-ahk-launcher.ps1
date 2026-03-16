param(
  [string]$Version = "2.0-beta33"
)

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$tempRoot = Join-Path $env:TEMP "MyKeymap-re-build"
$compilerDir = Join-Path $tempRoot "Ahk2Exe"
$compilerExe = Join-Path $compilerDir "Ahk2Exe.exe"
$releaseApi = "https://api.github.com/repos/AutoHotkey/Ahk2Exe/releases/latest"

function Ensure-Ahk2Exe() {
  if (Test-Path $compilerExe) {
    return
  }

  New-Item -ItemType Directory -Path $compilerDir -Force | Out-Null
  $release = Invoke-RestMethod -Headers @{ "User-Agent" = "Codex" } $releaseApi
  $asset = $release.assets | Where-Object { $_.name -like "*.zip" } | Select-Object -First 1
  if (-not $asset) {
    throw "Unable to find Ahk2Exe zip asset in latest release."
  }

  $zipPath = Join-Path $tempRoot $asset.name
  Invoke-WebRequest -UseBasicParsing $asset.browser_download_url -OutFile $zipPath
  Expand-Archive -Path $zipPath -DestinationPath $compilerDir -Force
}

function Build-Launcher() {
  $launcherScript = Join-Path $root "bin\Launcher.ahk"
  $iconFile = Join-Path $root "bin\icons\logo3.ico"
  $baseFile = Join-Path $root "bin\AutoHotkey64.exe"
  $outputFile = Join-Path $root "MyKeymap.exe"
  $tempOutput = Join-Path $tempRoot "MyKeymap.exe"

  if (-not (Test-Path $launcherScript)) {
    throw "Launcher script not found: $launcherScript"
  }
  if (-not (Test-Path $baseFile)) {
    throw "AutoHotkey base executable not found: $baseFile"
  }
  if (-not (Test-Path $iconFile)) {
    throw "Launcher icon not found: $iconFile"
  }

  if (Test-Path $tempOutput) {
    Remove-Item $tempOutput -Force
  }

  & $compilerExe `
    /in $launcherScript `
    /out $tempOutput `
    /icon $iconFile `
    /base $baseFile `
    /silent verbose

  for ($i = 0; $i -lt 20 -and -not (Test-Path $tempOutput); $i++) {
    Start-Sleep -Milliseconds 250
  }

  if (-not (Test-Path $tempOutput)) {
    throw "Ahk2Exe did not produce $tempOutput"
  }

  Copy-Item $tempOutput $outputFile -Force
}

New-Item -ItemType Directory -Path $tempRoot -Force | Out-Null
Ensure-Ahk2Exe
Build-Launcher

Write-Output "Built MyKeymap.exe with Ahk2Exe."
