param(
  [string]$SubscriptionId,
  [string]$Region = "brazilsouth",
  [string]$Owner = "<TU_USUARIO>",

  # Resource Groups
  [string]$RgQa = "rg-is3-qa",
  [string]$RgProd = "rg-is3-prod",
  [string]$RgShared = "rg-is3-shared",

  # ACR
  [string]$AcrName = "is3acr",

  # App Service Plan + Web Apps
  [string]$AppServicePlan = "asp-is3-shared",
  [string]$BackendQaApp = "is3-backend-qa",
  [string]$FrontendQaApp = "is3-frontend-qa",
  [string]$BackendProdApp = "is3-backend-prod",
  [string]$FrontendProdApp = "is3-frontend-prod",

  # PostgreSQL Flexible
  [string]$PgQaName = "is3pgqa",
  [string]$PgProdName = "is3pgprod",
  [Parameter(Mandatory=$true)][string]$PgQaPassword,
  [Parameter(Mandatory=$true)][string]$PgProdPassword
)

$ErrorActionPreference = 'Stop'

function Exec($cmd) {
  Write-Host "--> $cmd" -ForegroundColor Cyan
  Invoke-Expression $cmd
}

Write-Host "Login a Azure si es necesario..." -ForegroundColor Yellow
try { az account show 1>$null 2>$null } catch { az login | Out-Null }
if ($SubscriptionId) { Exec "az account set --subscription `$SubscriptionId" }

# 1) Resource Groups
Write-Host "[1/6] Creando Resource Groups..." -ForegroundColor Yellow
Exec "az group create --name `$RgQa --location `$Region --tags project=IS3 env=qa owner=`$Owner purpose=tp05"
Exec "az group create --name `$RgProd --location `$Region --tags project=IS3 env=prod owner=`$Owner purpose=tp05"
Exec "az group create --name `$RgShared --location `$Region --tags project=IS3 env=shared owner=`$Owner purpose=tp05"

# 2) ACR
Write-Host "[2/6] Creando ACR `$AcrName..." -ForegroundColor Yellow
Exec "az acr create --resource-group `$RgShared --name `$AcrName --sku Basic --location `$Region --tags project=IS3 env=shared owner=`$Owner purpose=tp05"
Exec "az acr update --name `$AcrName --admin-enabled true"
$acrLoginServer = (az acr show --name $AcrName --query "loginServer" -o tsv)
$acrUser = (az acr credential show --name $AcrName --query "username" -o tsv)
$acrPass = (az acr credential show --name $AcrName --query "passwords[0].value" -o tsv)

# 3) App Service Plan
Write-Host "[3/6] Creando App Service Plan Linux (B1)..." -ForegroundColor Yellow
Exec "az appservice plan create --name `$AppServicePlan --resource-group `$RgShared --sku B1 --is-linux --location `$Region --tags project=IS3 env=shared owner=`$Owner purpose=tp05"

# 4) Web Apps (contenedores)
Write-Host "[4/6] Creando Web Apps (QA/Prod)..." -ForegroundColor Yellow
Exec "az webapp create --resource-group `$RgQa   --plan `$AppServicePlan --name `$BackendQaApp   --deployment-container-image-name `$acrLoginServer/is3-backend:placeholder"
Exec "az webapp create --resource-group `$RgProd --plan `$AppServicePlan --name `$BackendProdApp --deployment-container-image-name `$acrLoginServer/is3-backend:placeholder"
Exec "az webapp create --resource-group `$RgQa   --plan `$AppServicePlan --name `$FrontendQaApp   --deployment-container-image-name `$acrLoginServer/is3-frontend:placeholder"
Exec "az webapp create --resource-group `$RgProd --plan `$AppServicePlan --name `$FrontendProdApp --deployment-container-image-name `$acrLoginServer/is3-frontend:placeholder"

# Vincular ACR credenciales (método rápido)
Write-Host "[4.1] Configurando credenciales de ACR en Web Apps..." -ForegroundColor Yellow
foreach ($pair in @(
  @{ rg=$RgQa;   app=$BackendQaApp },
  @{ rg=$RgProd; app=$BackendProdApp },
  @{ rg=$RgQa;   app=$FrontendQaApp },
  @{ rg=$RgProd; app=$FrontendProdApp }
)) {
  Exec "az webapp config appsettings set --resource-group `$($pair.rg) --name `$($pair.app) --settings DOCKER_REGISTRY_SERVER_URL=https://`$acrLoginServer DOCKER_REGISTRY_SERVER_USERNAME=`$acrUser DOCKER_REGISTRY_SERVER_PASSWORD=`"$acrPass`""
}

# 5) PostgreSQL Flexible Servers (QA/Prod)
Write-Host "[5/6] Creando PostgreSQL Flexible (QA/Prod)..." -ForegroundColor Yellow
Exec "az postgres flexible-server create --name `$PgQaName --resource-group `$RgQa --location `$Region --admin-user app --admin-password `"$PgQaPassword`" --tier Burstable --sku-name Standard_B1ms --storage-size 32 --version 15 --database-name app --tags project=IS3 env=qa owner=`$Owner purpose=tp05"
Exec "az postgres flexible-server create --name `$PgProdName --resource-group `$RgProd --location `$Region --admin-user app --admin-password `"$PgProdPassword`" --tier Burstable --sku-name Standard_B1ms --storage-size 32 --version 15 --database-name app --tags project=IS3 env=prod owner=`$Owner purpose=tp05"

# Firewall: permitir IP actual (para pruebas iniciales)
Write-Host "[5.1] Agregando regla de firewall con tu IP pública..." -ForegroundColor Yellow
$myIp = (Invoke-RestMethod -Uri 'https://ifconfig.me/ip')
Exec "az postgres flexible-server firewall-rule create --name `$PgQaName --resource-group `$RgQa --rule-name myip --start-ip-address `$myIp --end-ip-address `$myIp"
Exec "az postgres flexible-server firewall-rule create --name `$PgProdName --resource-group `$RgProd --rule-name myip --start-ip-address `$myIp --end-ip-address `$myIp"

$pgQaHost = "$PgQaName.postgres.database.azure.com"
$pgProdHost = "$PgProdName.postgres.database.azure.com"
$dbQaUrl   = "postgres://app:$PgQaPassword@$pgQaHost:5432/app?sslmode=require"
$dbProdUrl = "postgres://app:$PgProdPassword@$pgProdHost:5432/app?sslmode=require"

# 6) App Settings mínimos por entorno
Write-Host "[6/6] Configurando App Settings (QA/Prod)..." -ForegroundColor Yellow

# Backend QA/Prod
Exec "az webapp config appsettings set --resource-group `$RgQa   --name `$BackendQaApp   --settings PORT=8000 ENV=qa DATABASE_URL=`"$dbQaUrl`""
Exec "az webapp config appsettings set --resource-group `$RgProd --name `$BackendProdApp --settings PORT=8000 ENV=production DATABASE_URL=`"$dbProdUrl`""

# Frontend QA/Prod (apunta a los backends publicados de App Service)
$backendQaUrl   = "https://$BackendQaApp.azurewebsites.net"
$backendProdUrl = "https://$BackendProdApp.azurewebsites.net"
Exec "az webapp config appsettings set --resource-group `$RgQa   --name `$FrontendQaApp   --settings NODE_ENV=production RUNTIME_PUBLIC_API_URL=`"$backendQaUrl`" INTERNAL_API=`"$backendQaUrl`" NEXT_PUBLIC_API_URL=`"$backendQaUrl`""
Exec "az webapp config appsettings set --resource-group `$RgProd --name `$FrontendProdApp --settings NODE_ENV=production RUNTIME_PUBLIC_API_URL=`"$backendProdUrl`" INTERNAL_API=`"$backendProdUrl`" NEXT_PUBLIC_API_URL=`"$backendProdUrl`""

Write-Host "\nProvisioning completado. Revisa Azure Portal para validar recursos y costos." -ForegroundColor Green
Write-Host "ACR: $acrLoginServer" -ForegroundColor Green
Write-Host "Backend QA URL:  $backendQaUrl" -ForegroundColor Green
Write-Host "Backend Prod URL: $backendProdUrl" -ForegroundColor Green