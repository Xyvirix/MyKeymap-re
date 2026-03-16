param(
  [string]$Version = "2.0-beta33"
)

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$distRoot = Join-Path $root "dist"
$packageName = "MyKeymapDesktop-$Version"
$packageDir = Join-Path $distRoot $packageName
$zipPath = Join-Path $distRoot "$packageName.zip"

if (Test-Path $packageDir) {
  Remove-Item $packageDir -Recurse -Force
}
if (Test-Path $zipPath) {
  Remove-Item $zipPath -Force
}

New-Item -ItemType Directory -Path $packageDir | Out-Null
New-Item -ItemType Directory -Path (Join-Path $packageDir "shortcuts") | Out-Null

$copyTargets = @(
  "MyKeymapDesktop.exe",
  "MyKeymap.exe",
  "bin",
  "data",
  "tools",
  "readme.md",
  "readme.en.md",
  "LICENSE",
  "误报病毒时执行这个.bat"
)

foreach ($target in $copyTargets) {
  $source = Join-Path $root $target
  if (Test-Path $source) {
    Copy-Item $source $packageDir -Recurse -Force
  }
}

$packageTemplates = Join-Path $packageDir "bin\\templates"
if (-not (Test-Path $packageTemplates)) {
  $templatesSource = Join-Path $root "config-server\\templates"
  if (Test-Path $templatesSource) {
    Copy-Item $templatesSource $packageTemplates -Recurse -Force
  }
}

$shortcutsSource = Join-Path $root "shortcuts"
if (Test-Path $shortcutsSource) {
  Copy-Item (Join-Path $shortcutsSource "*") (Join-Path $packageDir "shortcuts") -Recurse -Force
}

Compress-Archive -Path (Join-Path $packageDir "*") -DestinationPath $zipPath -CompressionLevel Optimal

Write-Output "Package directory: $packageDir"
Write-Output "Package zip: $zipPath"
