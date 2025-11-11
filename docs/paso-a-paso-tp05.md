# TP05 – Paso a paso de configuración y prueba (Azure + Azure DevOps)

Este documento resume cómo dejar todo listo para la defensa del TP05. Contiene dos caminos: A) pipeline todo‑en‑uno (build + deploy) y B) deploy‑only usando imágenes publicadas por el TP04. Usá UNO solo.

---
## 0) Prerrequisitos
- Cuenta Azure con permisos (Contributor) en la suscripción.
- Azure CLI instalado y logueado:
```powershell
az login
az account show
```
- Proyecto en Azure DevOps con acceso al repo y agente self‑hosted en el pool `TP4_IS3`.
- Rama `main` con los archivos del TP05:
  - `azure-pipelines.release.yaml`, `azure-pipelines.deploy.yaml`
  - `scripts/provision_azure.ps1`
  - `docs/azure-devops-setup.md`, `docs/evidencias/README.md`

---
## 1) Provisionar recursos en Azure (Paso 1 del TP)

### Opción A – Script
Editar variables si hace falta y ejecutar:
```powershell
Set-ExecutionPolicy -Scope Process Bypass -Force
./scripts/provision_azure.ps1
```
Crea RGs, ACR, App Service Plan Linux, Web Apps QA/Prod (front/back), PostgreSQL QA/Prod y App Settings base.

### Opción B – Manual (Portal)
Crear:
- RGs: `rg-is3-shared`, `rg-is3-qa`, `rg-is3-prod` (región `brazilsouth`).
- ACR: `is3acr` (SKU Basic).
- App Service Plan Linux: `asp-is3-shared` (B1).
- Web Apps: `is3-backend-qa/prod`, `is3-frontend-qa/prod`.
- PostgreSQL Flexible: `is3pgqa`, `is3pgprod`.

App Settings clave:
- Backend QA: `PORT=8000`, `ENV=qa`, `DATABASE_URL=postgres://app:<PASS>@is3pgqa.postgres.database.azure.com:5432/app?sslmode=require`
- Backend Prod: `PORT=8000`, `ENV=production`, `DATABASE_URL=postgres://app:<PASS>@is3pgprod.postgres.database.azure.com:5432/app?sslmode=require`
- Frontend QA: `RUNTIME_PUBLIC_API_URL`, `INTERNAL_API`, `NEXT_PUBLIC_API_URL` → `https://is3-backend-qa.azurewebsites.net`
- Frontend Prod: mismas variables apuntando a backend Prod.

---
## 2) Service Connections (Azure DevOps)
- Azure Resource Manager (ARM): `IS3-Azure-Subscription` (Service principal automatic; scope Subscription).
- Docker Registry (ACR) – sólo si usás el pipeline todo‑en‑uno: `is3acr-service-connection`.

---
## 3) Environments y aprobación
- Crear `is3-qa` y `is3-prod`.
- En `is3-prod` → Approvals and checks → Approvals → agregá al menos 1 aprobador.

---
## 4) Variables / secretos
- Para `azure-pipelines.release.yaml` (all‑in‑one): definir `databaseUrlQA` y `databaseUrlProd` como secretos.
- Para `azure-pipelines.deploy.yaml` (deploy-only): no requiere secretos (usa imágenes ya publicadas + App Settings ya cargados).

---
## 5) Elegir estrategia

### A) All‑in‑one: build + deploy
1. Azure DevOps → Pipelines → New pipeline → Existing YAML → `azure-pipelines.release.yaml`.
2. Run: construye imágenes, las publica en ACR con tag `$(Build.BuildId)`, despliega a QA, espera aprobación, despliega a Prod.

### B) Deploy-only: reutilizar imágenes de TP04
Prerequisito: TP04 publica en ACR:
- `is3acr.azurecr.io/is3-backend:<tag>`
- `is3acr.azurecr.io/is3-frontend:<tag>`

Pasos:
1. Verificar tags (opcional):
```powershell
az acr repository show-tags -n is3acr --repository is3-backend
az acr repository show-tags -n is3acr --repository is3-frontend
```
2. Crear pipeline desde `azure-pipelines.deploy.yaml`.
3. Run con parámetro `imageTag=<tag>` (mismo tag de TP04).

---
## 6) Verificación post‑deploy
QA:
```powershell
Invoke-WebRequest -Uri https://is3-backend-qa.azurewebsites.net/health -UseBasicParsing
Invoke-WebRequest -Uri https://is3-frontend-qa.azurewebsites.net/ -UseBasicParsing
```
Prod (tras aprobación):
```powershell
Invoke-WebRequest -Uri https://is3-backend-prod.azurewebsites.net/health -UseBasicParsing
Invoke-WebRequest -Uri https://is3-frontend-prod.azurewebsites.net/ -UseBasicParsing
```
Debe devolver 200. Si falla, revisar App Settings, Log Stream de Web App y existencia del tag en ACR.

---
## 7) Evidencias (guardar en `docs/evidencias/`)
- Azure: RGs, ACR repos+tags, Web Apps QA/Prod (App Settings), PostgreSQL QA/Prod.
- Azure DevOps: Service connections, Environments (aprobación en `is3-prod`), pipeline con QA OK → aprobación → Prod OK.
- Health checks 200 (QA y Prod).
- (Opcional) tag `last-known-good` apuntando al digest estable.
- Rollback: redeploy de un run/tag anterior.

---
## 8) Rollback
- Desde Pipelines → histórico de runs: redeploy del run anterior exitoso.
- Alternativa con alias en ACR:
```powershell
az acr import -n is3acr --source is3acr.azurecr.io/is3-backend:<stable> --image is3-backend:last-known-good --force
az acr import -n is3acr --source is3acr.azurecr.io/is3-frontend:<stable> --image is3-frontend:last-known-good --force
```

---
## 9) Problemas comunes
| Problema | Causa | Solución |
|---------|-------|----------|
| 403 al desplegar | Service connection sin permisos | Asignar role Contributor al SP en los RG |
| Image pull error | Web App sin creds ACR | Configurar DOCKER_REGISTRY_SERVER_* o Managed Identity con ACR Pull |
| Frontend apunta mal | Vars runtime erróneas | Corregir App Settings (RUNTIME_PUBLIC_API_URL/INTERNAL_API) |
| DB SSL | Falta sslmode=require | Incluir en DATABASE_URL |
| Timeout health | Arranque frío o migraciones | Aumentar espera o añadir reintentos |

---
## 10) Smoke test consolidado
```powershell
$urls = @(
  'https://is3-backend-qa.azurewebsites.net/health',
  'https://is3-frontend-qa.azurewebsites.net/',
  'https://is3-backend-prod.azurewebsites.net/health',
  'https://is3-frontend-prod.azurewebsites.net/'
)
foreach ($u in $urls) {
  try {
    $r = Invoke-WebRequest -Uri $u -UseBasicParsing -TimeoutSec 15
    Write-Host "[OK] $u -> $($r.StatusCode)"
  } catch {
    Write-Host "[FAIL] $u -> $($_.Exception.Message)"
  }
}
```

---
## 11) Checklist final
- [ ] Recursos Azure creados (ACR, Web Apps, Postgres QA/Prod)
- [ ] App Settings configurados (backend y frontend QA/Prod)
- [ ] Service connections ARM (y ACR si aplica)
- [ ] Environments `is3-qa`/`is3-prod` con aprobación
- [ ] Pipeline elegido creado y ejecutado (release o deploy-only)
- [ ] QA OK (200), aprobación registrada, Prod OK (200)
- [ ] Evidencias en `docs/evidencias/`
- [ ] Rollback probado (redeploy tag anterior)
