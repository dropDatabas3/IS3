Param(
    [string]$BackendTagProd  = "prod",
    [string]$BackendTagQA    = "qa",
    [string]$FrontendTagProd = "prod",
    [string]$FrontendTagQA   = "qa"
)

function FailIfLastExitNonZero($msg) {
    if ($LASTEXITCODE -ne 0) {
        Write-Error $msg
        exit $LASTEXITCODE
    }
}

# Script directory and repo root (one level up from scripts)
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
$repoRoot  = Split-Path -Parent $scriptDir
Write-Host "Workspace root: $repoRoot"

Write-Host "Building backend image..."
$backendContext = Join-Path $repoRoot 'ucc-soft-arch-golang'
docker build -t is3-backend:tmp $backendContext
FailIfLastExitNonZero 'Backend build failed'

Write-Host "Tagging backend image as is3-backend:$BackendTagProd and is3-backend:$BackendTagQA"
docker tag is3-backend:tmp "is3-backend:$BackendTagProd"
docker tag is3-backend:tmp "is3-backend:$BackendTagQA"
docker rmi is3-backend:tmp -f | Out-Null

Write-Host "Building frontend image..."
$frontendContext = Join-Path $repoRoot 'ucc-arq-soft-front'
docker build -t is3-frontend:tmp $frontendContext
FailIfLastExitNonZero 'Frontend build failed'

Write-Host "Tagging frontend image as is3-frontend:$FrontendTagProd and is3-frontend:$FrontendTagQA"
docker tag is3-frontend:tmp "is3-frontend:$FrontendTagProd"
docker tag is3-frontend:tmp "is3-frontend:$FrontendTagQA"
docker rmi is3-frontend:tmp -f | Out-Null

Write-Host "Done. Created tags:"
Write-Host " - is3-backend:$BackendTagProd"
Write-Host " - is3-backend:$BackendTagQA"
Write-Host " - is3-frontend:$FrontendTagProd"
Write-Host " - is3-frontend:$FrontendTagQA"

Write-Host "To run prod and qa using these images:"
Write-Host "  docker compose up -d backend frontend backend_qa frontend_qa"
