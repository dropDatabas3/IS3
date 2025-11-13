<#
Runs Go tests with coverage across all packages and produces a merged coverage.out and an HTML report.
Works around Go's package-by-package coverage behavior by collecting per-package profiles and merging them.
#>
param(
    [switch]$Open
)

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

# Navigate to the backend folder if script is called from repo root
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$root = (Resolve-Path (Join-Path $scriptDir '..')).Path
Set-Location $root

Write-Host "Running go tests with coverage (merged)..." -ForegroundColor Cyan

$mode = 'atomic'
$outFile = Join-Path (Get-Location) 'coverage.out'
if (Test-Path $outFile) { Remove-Item $outFile -Force }

# Write header for merged profile
"mode: $mode" | Out-File -FilePath $outFile -Encoding ascii

# List all packages in module
$packages = (& go list ./...) | Where-Object { $_ -ne '' }

foreach ($pkg in $packages) {
    $tmp = [System.IO.Path]::GetTempFileName()
    # Run tests for the package with coverage; coverpkg=./... ensures cross-package coverage accounting
    $cmd = "go test -covermode=$mode -coverpkg=./... -coverprofile=`"$tmp`" $pkg"
    Write-Host "→ $cmd" -ForegroundColor DarkGray
    $exitCode = 0
    try {
        Invoke-Expression $cmd | Out-Host
    } catch {
        $exitCode = 1
    }

    if (Test-Path $tmp) {
        # Append profile data excluding header line
        (Get-Content $tmp | Select-Object -Skip 1) | Add-Content -Path $outFile -Encoding ascii
        Remove-Item $tmp -Force
    }

    if ($exitCode -ne 0) {
        Write-Error "Tests failed for package $pkg. Aborting coverage merge."
        exit 1
    }
}

if (!(Test-Path $outFile) -or ((Get-Item $outFile).Length -eq 0)) {
    Write-Error "coverage.out not generated. Check test failures above."
    exit 1
}

# Show total coverage (robust arg passing and error checks)
$summaryOutput = & go tool cover -func $outFile 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host $summaryOutput -ForegroundColor Red
    throw "Failed to compute coverage summary from '$outFile'"
}
$total = $summaryOutput | Select-String -SimpleMatch 'total:' | ForEach-Object { $_.ToString() }
if ($total) { Write-Host $total -ForegroundColor Green }

# Generate HTML report with explicit args and validation
$report = Join-Path (Get-Location) 'coverage.html'
$htmlOutput = & go tool cover -html $outFile -o $report 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host $htmlOutput -ForegroundColor Red
    throw "Failed to generate HTML report from '$outFile'"
}
if (Test-Path $report) {
    Write-Host "HTML report generated at: $report" -ForegroundColor Green
    if ($Open) {
        try {
            Start-Process $report
        } catch {
            Write-Warning "Couldn't open coverage report automatically. Open it manually at: $report"
        }
    }
} else {
    throw "Expected report at '$report' not found after generation"
}
