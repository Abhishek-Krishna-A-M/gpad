$ErrorActionPreference = "Stop"

$Repo   = "Abhishek-Krishna-A-M/gpad"
$Binary = "gpad.exe"
$InstallDir = "$env:LOCALAPPDATA\gpad"

# ── detect arch ──────────────────────────────────────────────────────────────

$Arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else {
  Write-Host "Only 64-bit Windows is supported."
  exit 1
}

# ── fetch latest release ──────────────────────────────────────────────────────

Write-Host "Fetching latest gpad release..."

try {
  $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
  $Tag = $Release.tag_name
} catch {
  Write-Host "Could not reach GitHub API. Check your internet connection."
  exit 1
}

Write-Host "Latest release: $Tag"

# ── download ──────────────────────────────────────────────────────────────────

$Url = "https://github.com/$Repo/releases/download/$Tag/gpad_windows_$Arch.exe"
$TmpFile = "$env:TEMP\gpad_install.exe"

Write-Host "Downloading gpad $Tag (windows/$Arch)..."
Invoke-WebRequest -Uri $Url -OutFile $TmpFile -UseBasicParsing

# ── install ───────────────────────────────────────────────────────────────────

if (-not (Test-Path $InstallDir)) {
  New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

Move-Item -Force $TmpFile "$InstallDir\$Binary"

# ── add to PATH if needed ─────────────────────────────────────────────────────

$CurrentPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
  [System.Environment]::SetEnvironmentVariable(
    "PATH",
    "$CurrentPath;$InstallDir",
    "User"
  )
  Write-Host "Added $InstallDir to your PATH."
  Write-Host "Restart your terminal for the PATH change to take effect."
}

# ── done ──────────────────────────────────────────────────────────────────────

Write-Host ""
Write-Host "  gpad installed to $InstallDir\$Binary"
Write-Host ""
Write-Host "  Get started:"
Write-Host "    gpad today               open today's daily note"
Write-Host "    gpad open my-note.md     create your first note"
Write-Host "    gpad git init <url>      connect git sync (optional)"
Write-Host ""
