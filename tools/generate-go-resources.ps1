param(
  [string]$Version = "2.0-beta33-re1",
  [string]$Creator = "咸鱼阿康",
  [string]$Author = "Xyvirix"
)

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$goExe = "C:\Program Files\Go\bin\go.exe"
$gopath = & $goExe env GOPATH
$goversioninfo = Join-Path $gopath "bin\goversioninfo.exe"

if (-not (Test-Path $goversioninfo)) {
  & $goExe install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
}

function Get-VersionParts([string]$value) {
  $major = 1
  $minor = 0
  $patch = 0
  $build = 0

  if ($value -match "^(?<major>\d+)\.(?<minor>\d+)") {
    $major = [int]$Matches.major
    $minor = [int]$Matches.minor
  }

  if ($value -match "(?<build>\d+)$") {
    $build = [int]$Matches.build
  }

  return @{
    Major = $major
    Minor = $minor
    Patch = $patch
    Build = $build
  }
}

$versionParts = Get-VersionParts $Version
$iconPath = Join-Path $root "bin\icons\logo3.ico"
$versionInfoJson = Join-Path $env:TEMP "mykeymap-re-versioninfo.json"
"{}" | Set-Content -Path $versionInfoJson -Encoding ASCII

function New-VersionResource([string]$outputPath, [string]$productName, [string]$description, [string]$internalName, [string]$originalName) {
  & $goversioninfo `
    -64 `
    -icon $iconPath `
    -product-name $productName `
    -company $Author `
    -description $description `
    -internal-name $internalName `
    -original-name $originalName `
    -file-version $Version `
    -product-version $Version `
    -ver-major $versionParts.Major `
    -ver-minor $versionParts.Minor `
    -ver-patch $versionParts.Patch `
    -ver-build $versionParts.Build `
    -product-ver-major $versionParts.Major `
    -product-ver-minor $versionParts.Minor `
    -product-ver-patch $versionParts.Patch `
    -product-ver-build $versionParts.Build `
    -copyright "Created by $Creator modified by $Author" `
    -o $outputPath `
    $versionInfoJson
}

New-VersionResource `
  -outputPath (Join-Path $root "config-server\cmd\desktop\resource_windows_amd64.syso") `
  -productName "MyKeymap-re Desktop" `
  -description "MyKeymap-re desktop host for hotkey management" `
  -internalName "MyKeymapDesktop" `
  -originalName "MyKeymapDesktop.exe"

New-VersionResource `
  -outputPath (Join-Path $root "config-server\cmd\settings\resource_windows_amd64.syso") `
  -productName "MyKeymap-re Settings" `
  -description "MyKeymap-re settings service" `
  -internalName "settings" `
  -originalName "settings.exe"

Remove-Item $versionInfoJson -Force

Write-Output "Generated Windows resources for desktop and settings executables."
