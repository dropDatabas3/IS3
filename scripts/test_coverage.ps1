param(
    [switch]$Open
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

Write-Host "Running Go tests with coverage across all packages..."

$coverageOut = "coverage.out"
$coverageHtml = "coverage.html"

# Ensure we run from the Go module root (ucc-soft-arch-golang)
$moduleDir = Join-Path $PSScriptRoot "..\ucc-soft-arch-golang"
Push-Location $moduleDir

# Run tests with atomic cover mode and aggregate profile
go test ./... -covermode=atomic -coverpkg=./... -coverprofile=$coverageOut | Write-Host

if (!(Test-Path $coverageOut)) {
    Pop-Location
    Write-Error "Coverage profile not generated."
    exit 1
}

Write-Host "\nCoverage summary (go tool cover -func):"
go tool cover -func $coverageOut | Write-Host

Write-Host "\nGenerating HTML report: $coverageHtml"
go tool cover -html $coverageOut -o $coverageHtml

if ($Open) {
    if (Test-Path $coverageHtml) { Start-Process $coverageHtml }
}

Pop-Location

Write-Host "Done."
