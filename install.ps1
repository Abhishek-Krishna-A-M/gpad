# PowerShell installer for gpad
# Usage:
#   powershell -ExecutionPolicy Bypass -File install.ps1

$ErrorActionPreference = "Stop"

$repo = "Abhishek-Krishna-A-M/gpad"
$binName = "gpad.exe"
$installDir = "$env:LOCALAPPDATA\Programs\gpad"

Write-Host ">>> Detecting Architecture..."

$arch = $env:PROCESSOR_ARCHITECTURE

switch ($arch) {
    "AMD64" { $arch = "amd64" }
    "ARM64" { $arch = "arm64" }
    default {
        Write-Error "Unsupported architecture: $arch"
        exit 1
    }
}

Write-Host ">>> Architecture: $arch"

Write-Host ">>> Fetching latest version tag..."
$latest = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest"

$tag = $latest.tag_name
if (-not $tag) {
    Write-Error "Failed to retrieve latest release tag."
    exit 1
}

Write-Host ">>> Latest Version: $tag"

$file = "gpad-windows-$arch.exe"
$downloadUrl = "https://github.com/$repo/releases/download/$tag/$file"

Write-Host ">>> Downloading $downloadUrl"

$tempPath = Join-Path $env:TEMP $binName
Invoke-WebRequest -Uri $downloadUrl -OutFile $tempPath

Write-Host ">>> Creating install directory..."
New-Item -ItemType Directory -Path $installDir -Force | Out-Null

Write-Host ">>> Installing gpad to $installDir"
Copy-Item $tempPath "$installDir\$binName" -Force

Write-Host ">>> Checking PATH..."

$path = [Environment]::GetEnvironmentVariable("Path", "User")

if ($path -notlike "*$installDir*") {
    Write-Host ">>> Adding gpad to PATH..."
    $newPath = "$path;$installDir"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host ">>> PATH updated. You may need to restart your terminal."
} else {
    Write-Host ">>> PATH already contains gpad directory."
}

Write-Host ""
Write-Host ">>> Installation complete!"
Write-Host "Run 'gpad' in a new terminal window."

