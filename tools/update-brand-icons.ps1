Add-Type -AssemblyName System.Drawing

$ErrorActionPreference = 'Stop'

function New-Color($hex, $alpha = 255) {
  $value = $hex.TrimStart('#')
  return [System.Drawing.Color]::FromArgb(
    $alpha,
    [Convert]::ToInt32($value.Substring(0, 2), 16),
    [Convert]::ToInt32($value.Substring(2, 2), 16),
    [Convert]::ToInt32($value.Substring(4, 2), 16)
  )
}

function New-RoundedRectPath([float]$x, [float]$y, [float]$width, [float]$height, [float]$radius) {
  $path = New-Object System.Drawing.Drawing2D.GraphicsPath
  $diameter = [Math]::Max(1, $radius * 2)
  $path.AddArc($x, $y, $diameter, $diameter, 180, 90)
  $path.AddArc($x + $width - $diameter, $y, $diameter, $diameter, 270, 90)
  $path.AddArc($x + $width - $diameter, $y + $height - $diameter, $diameter, $diameter, 0, 90)
  $path.AddArc($x, $y + $height - $diameter, $diameter, $diameter, 90, 90)
  $path.CloseFigure()
  return $path
}

function Write-Ico([string]$path, [byte[][]]$pngFrames, [int[]]$sizes) {
  $stream = [System.IO.File]::Open($path, [System.IO.FileMode]::Create)
  $writer = New-Object System.IO.BinaryWriter($stream)
  try {
    $writer.Write([UInt16]0)
    $writer.Write([UInt16]1)
    $writer.Write([UInt16]$pngFrames.Length)

    $offset = 6 + (16 * $pngFrames.Length)
    for ($i = 0; $i -lt $pngFrames.Length; $i++) {
      $size = $sizes[$i]
      $data = $pngFrames[$i]
      $writer.Write([byte]($(if ($size -ge 256) { 0 } else { $size })))
      $writer.Write([byte]($(if ($size -ge 256) { 0 } else { $size })))
      $writer.Write([byte]0)
      $writer.Write([byte]0)
      $writer.Write([UInt16]1)
      $writer.Write([UInt16]32)
      $writer.Write([UInt32]$data.Length)
      $writer.Write([UInt32]$offset)
      $offset += $data.Length
    }

    foreach ($data in $pngFrames) {
      $writer.Write($data)
    }
  } finally {
    $writer.Dispose()
    $stream.Dispose()
  }
}

function New-PngBytes([System.Drawing.Bitmap]$bitmap) {
  $memory = New-Object System.IO.MemoryStream
  try {
    $bitmap.Save($memory, [System.Drawing.Imaging.ImageFormat]::Png)
    return $memory.ToArray()
  } finally {
    $memory.Dispose()
  }
}

function New-BrandBitmap([int]$size, [System.Drawing.Color]$hookColor, [System.Drawing.Color]$cardColor, [System.Drawing.Color]$lineColor) {
  $bitmap = New-Object System.Drawing.Bitmap $size, $size, ([System.Drawing.Imaging.PixelFormat]::Format32bppArgb)
  $graphics = [System.Drawing.Graphics]::FromImage($bitmap)

  try {
    $graphics.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
    $graphics.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
    $graphics.PixelOffsetMode = [System.Drawing.Drawing2D.PixelOffsetMode]::HighQuality
    $graphics.Clear([System.Drawing.Color]::Transparent)

    $cardX = $size * 0.18
    $cardY = $size * 0.40
    $cardW = $size * 0.66
    $cardH = $size * 0.34
    $cardRadius = $size * 0.072

    $hookPen = New-Object System.Drawing.Pen $hookColor, ($size * 0.065)
    $hookPen.StartCap = [System.Drawing.Drawing2D.LineCap]::Round
    $hookPen.EndCap = [System.Drawing.Drawing2D.LineCap]::Round

    $cardPath = New-RoundedRectPath $cardX $cardY $cardW $cardH $cardRadius
    $cardBrush = New-Object System.Drawing.SolidBrush $cardColor
    $linePen = New-Object System.Drawing.Pen $lineColor, ($size * 0.048)
    $linePen.StartCap = [System.Drawing.Drawing2D.LineCap]::Round
    $linePen.EndCap = [System.Drawing.Drawing2D.LineCap]::Round

    $hookPath = New-Object System.Drawing.Drawing2D.GraphicsPath
    $hookPath.StartFigure()
    $hookPath.AddLine($size * 0.29, $size * 0.40, $size * 0.29, $size * 0.24)
    $hookPath.AddBezier(
      $size * 0.29, $size * 0.24,
      $size * 0.29, $size * 0.12,
      $size * 0.41, $size * 0.12,
      $size * 0.42, $size * 0.25
    )
    $hookPath.AddBezier(
      $size * 0.42, $size * 0.25,
      $size * 0.42, $size * 0.30,
      $size * 0.39, $size * 0.33,
      $size * 0.36, $size * 0.33
    )
    $graphics.DrawPath($hookPen, $hookPath)

    $graphics.FillPath($cardBrush, $cardPath)
    $graphics.DrawLine($linePen, $size * 0.29, $size * 0.50, $size * 0.69, $size * 0.50)
    $graphics.DrawLine($linePen, $size * 0.29, $size * 0.59, $size * 0.68, $size * 0.59)
    $graphics.DrawLine($linePen, $size * 0.39, $size * 0.67, $size * 0.59, $size * 0.67)

    $hookPen.Dispose()
    $cardBrush.Dispose()
    $linePen.Dispose()
    $hookPath.Dispose()
    $cardPath.Dispose()
  } catch {
    $graphics.Dispose()
    $bitmap.Dispose()
    throw
  }

  $graphics.Dispose()
  return $bitmap
}

function New-IconBitmap([int]$size, [System.Drawing.Color]$cardColor, [System.Drawing.Color]$lineColor) {
  $bitmap = New-Object System.Drawing.Bitmap $size, $size, ([System.Drawing.Imaging.PixelFormat]::Format32bppArgb)
  $graphics = [System.Drawing.Graphics]::FromImage($bitmap)

  try {
    $graphics.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
    $graphics.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
    $graphics.PixelOffsetMode = [System.Drawing.Drawing2D.PixelOffsetMode]::HighQuality
    $graphics.Clear([System.Drawing.Color]::Transparent)

    $cardX = $size * 0.11
    $cardY = $size * 0.18
    $cardW = $size * 0.78
    $cardH = $size * 0.58
    $cardRadius = $size * 0.105

    $cardPath = New-RoundedRectPath $cardX $cardY $cardW $cardH $cardRadius
    $cardBrush = New-Object System.Drawing.SolidBrush $cardColor
    $outlinePen = New-Object System.Drawing.Pen $lineColor, ($size * 0.07)
    $outlinePen.LineJoin = [System.Drawing.Drawing2D.LineJoin]::Round
    $linePen = New-Object System.Drawing.Pen $lineColor, ($size * 0.062)
    $linePen.StartCap = [System.Drawing.Drawing2D.LineCap]::Round
    $linePen.EndCap = [System.Drawing.Drawing2D.LineCap]::Round

    $graphics.FillPath($cardBrush, $cardPath)
    $graphics.DrawPath($outlinePen, $cardPath)
    $graphics.DrawLine($linePen, $size * 0.28, $size * 0.34, $size * 0.72, $size * 0.34)
    $graphics.DrawLine($linePen, $size * 0.28, $size * 0.49, $size * 0.69, $size * 0.49)
    $graphics.DrawLine($linePen, $size * 0.40, $size * 0.63, $size * 0.60, $size * 0.63)

    $cardBrush.Dispose()
    $outlinePen.Dispose()
    $linePen.Dispose()
    $cardPath.Dispose()
  } catch {
    $graphics.Dispose()
    $bitmap.Dispose()
    throw
  }

  $graphics.Dispose()
  return $bitmap
}

function Write-IconOnlyAssets([string]$icoPath, [string]$cardHex, [string]$lineHex) {
  $cardColor = New-Color $cardHex
  $lineColor = New-Color $lineHex

  $sizes = @(16, 24, 32, 48, 64, 128, 256)
  $frames = @()
  foreach ($size in $sizes) {
    $bitmap = New-IconBitmap $size $cardColor $lineColor
    try {
      $frames += ,(New-PngBytes $bitmap)
    } finally {
      $bitmap.Dispose()
    }
  }

  Write-Ico -path $icoPath -pngFrames $frames -sizes $sizes
}

function New-BrandAssets([string]$pngPath, [string]$icoPath, [string]$hookHex, [string]$cardHex, [string]$lineHex) {
  $hookColor = New-Color $hookHex
  $cardColor = New-Color $cardHex
  $lineColor = New-Color $lineHex

  $master = New-BrandBitmap 512 $hookColor $cardColor $lineColor
  try {
    $master.Save($pngPath, [System.Drawing.Imaging.ImageFormat]::Png)
  } finally {
    $master.Dispose()
  }

  $sizes = @(16, 24, 32, 48, 64, 128, 256)
  $frames = @()
  foreach ($size in $sizes) {
    $bitmap = New-IconBitmap $size $cardColor $lineColor
    try {
      $frames += ,(New-PngBytes $bitmap)
    } finally {
      $bitmap.Dispose()
    }
  }

  Write-Ico -path $icoPath -pngFrames $frames -sizes $sizes
}

$root = 'C:\Users\xlanyi\Documents\000_test\MyKeymap-re'

New-BrandAssets `
  -pngPath (Join-Path $root 'config-ui\src\assets\logo.png') `
  -icoPath (Join-Path $root 'bin\icons\logo3.ico') `
  -hookHex '#38BDF8' `
  -cardHex '#EEF4FF' `
  -lineHex '#1E3656'

Write-IconOnlyAssets `
  -icoPath (Join-Path $root 'bin\icons\logo.ico') `
  -cardHex '#F7FBFF' `
  -lineHex '#18324A'

Write-IconOnlyAssets `
  -icoPath (Join-Path $root 'bin\icons\logo2.ico') `
  -cardHex '#F3F4F6' `
  -lineHex '#586274'

Write-IconOnlyAssets `
  -icoPath (Join-Path $root 'bin\icons\logo3.ico') `
  -cardHex '#EEF4FF' `
  -lineHex '#1E3656'

Copy-Item (Join-Path $root 'bin\icons\logo.ico') (Join-Path $root 'bin\html\logo.ico') -Force
Copy-Item (Join-Path $root 'bin\icons\logo3.ico') (Join-Path $root 'config-ui\public\favicon.ico') -Force

Write-Output 'Updated brand icons.'
