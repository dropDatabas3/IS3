
Param(
  [string]$SubscriptionId,
  [string]$Region = "brazilsouth",
  [string]$ResourceGroup = "rg-tp05-ingsoft3-2025",
  [string]$PlanName = "plan-tp05-linux-b1",
  [string]$PlanSku = "B1",
  [string]$BackendQaApp = "backend-tp05-qa-palacio-nallar",
  [string]$BackendProdApp = "backend-tp05-prod-palacio-nallar",
  [Parameter(Mandatory=$true)][string]$AcrName
)

$ErrorActionPreference = 'Stop'

function Exec($cmd) {
  Write-Host "--> $cmd" -ForegroundColor Cyan
  Invoke-Expression $cmd
}

Write-Host "Login a Azure si es necesario..." -ForegroundColor Yellow
try { az account show 1>$null 2>$null } catch { az login | Out-Null }
if ($SubscriptionId) { Exec "az account set --subscription `$SubscriptionId" }

# 1) Verificar RG
Write-Host "[1/5] Asegurando Resource Group $ResourceGroup" -ForegroundColor Yellow
Exec "az group create --name `$ResourceGroup --location `$Region"

# 2) App Service Plan Linux (B1)
Write-Host "[2/5] Creando/actualizando App Service Plan Linux $PlanName ($PlanSku)" -ForegroundColor Yellow
Exec "az appservice plan create --name `$PlanName --resource-group `$ResourceGroup --is-linux --sku `$PlanSku --location `$Region"

# 3) Crear Web Apps (Linux/Containers) con imagen placeholder
$placeholderImage = "nginx:alpine"
Write-Host "[3/5] Creando Web Apps backend QA/PROD (contenedores)" -ForegroundColor Yellow
foreach ($app in @($BackendQaApp, $BackendProdApp)) {
  Exec "az webapp create --resource-group `$ResourceGroup --plan `$PlanName --name `$app --deployment-container-image-name `$placeholderImage"
}

# 4) Identidad administrada + permiso AcrPull en ACR
Write-Host "[4/5] Habilitando identidad administrada y rol AcrPull en ACR $AcrName" -ForegroundColor Yellow
$acrId = (az acr show --name $AcrName --resource-group $ResourceGroup --query id -o tsv)
if (-not $acrId) { $acrId = (az acr show --name $AcrName --query id -o tsv) }
if (-not $acrId) { throw "No se pudo resolver el ID del ACR $AcrName" }

foreach ($app in @($BackendQaApp, $BackendProdApp)) {
  $identity = (az webapp identity assign --name $app --resource-group $ResourceGroup | ConvertFrom-Json)
  $principalId = $identity.principalId
  if (-not $principalId) { throw "No se obtuvo principalId para $app" }
  Exec "az role assignment create --assignee-object-id `$principalId --assignee-principal-type ServicePrincipal --role AcrPull --scope `$acrId"
}

# 5) App Settings mínimos (puerto/entorno) y WEBSITES_PORT para contenedor custom
Write-Host "[5/5] Configurando App Settings" -ForegroundColor Yellow
Exec "az webapp config appsettings set --resource-group `$ResourceGroup --name `$BackendQaApp   --settings PORT=8000 ENV=qa WEBSITES_PORT=8000"
Exec "az webapp config appsettings set --resource-group `$ResourceGroup --name `$BackendProdApp --settings PORT=8000 ENV=production WEBSITES_PORT=8000"

Write-Host "\nListo. Web Apps (Linux/Containers) backend creadas y vinculadas con ACR (AcrPull)." -ForegroundColor Green
Write-Host "Después del primer build, el pipeline actualizará la imagen del contenedor en cada Web App." -ForegroundColor Green
Write-Host "QA URL:   https://$BackendQaApp.azurewebsites.net" -ForegroundColor Green
Write-Host "PROD URL: https://$BackendProdApp.azurewebsites.net" -ForegroundColor Green
