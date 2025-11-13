# IS3 - Plataforma Cursos (Backend Go + Frontend Next.js + Postgres)

[![Build Status](https://dev.azure.com/your-org/IS3/_apis/build/status%2FIS3-CI?branchName=main)](https://dev.azure.com/your-org/IS3/_build/latest?definitionId=XX&branchName=main)

## Contenido del repositorio
- **Dockerfile Backend**: `ucc-soft-arch-golang/Dockerfile`
- **Dockerfile Frontend**: `ucc-arq-soft-front/Dockerfile`
- **Compose multi-entorno**: `docker-compose.yaml` (levanta PROD y QA a la vez)
- **Compose específicos**: `docker-compose.prod.yml`, `docker-compose.qa.yml`
- **Script de build**: `scripts/build_and_tag.ps1`
- **Ejemplos env**: `.env.example`
- **Volúmenes**: `db_data` (prod), `db_data_qa` (qa)

## Imágenes Docker
Se generan localmente (no se han publicado aún en Docker Hub dentro de este repo). Ejemplos de tags locales utilizados:
- Backend: `is3-backend:v1.0` (estable) / `is3-backend:qa` (si se necesita) / `is3-backend:prod`
- Frontend: `is3-frontend:v1.0` (estable) / `is3-frontend:qa` (build con API QA) / `is3-frontend:prod`

> Si deseas publicarlas en Docker Hub:
> 1. `docker tag is3-backend:v1.0 <usuario>/is3-backend:v1.0`
> 2. `docker push <usuario>/is3-backend:v1.0`
> 3. Repetir para frontend y para tags de desarrollo (por ejemplo `:qa`)
> 4. Actualizar `docker-compose.yaml` para usar `<usuario>/is3-backend:v1.0`, etc.

## Matriz de puertos
| Servicio        | PROD | QA  | Interno contenedor |
|-----------------|------|-----|--------------------|
| Frontend        | 3000 | 3001| 3000               |
| Backend         | 8000 | 8001| 8000               |
| Postgres        | 5432 | 5433| 5432               |

## Variables de entorno clave
Backend (prod):
- `DATABASE_URL=postgres://app:password@db:5432/app?sslmode=disable`
- `PGHOST=db`, `PGPORT=5432`, `PGUSER=app`, `PGPASSWORD=password`, `PGDATABASE=app`

Backend (qa):
- `DATABASE_URL=postgres://app:password@db_qa:5432/app?sslmode=disable`
- `PGHOST=db_qa`, resto igual.

Frontend PROD runtime/build:
- `NEXT_PUBLIC_API_URL=http://localhost:8000`
- `INTERNAL_API=http://backend:8000`

Frontend QA (imagen separada):
- `NEXT_PUBLIC_API_URL=http://localhost:8001`
- `INTERNAL_API=http://backend_qa:8000`

## Construir imágenes localmente
Script unificado (build multi-stage backend y frontend, tag prod/qa):
```powershell
# En la raíz del repo
./scripts/build_and_tag.ps1
```
Build manual backend:
```powershell
docker build -t is3-backend:v1.0 -f ucc-soft-arch-golang/Dockerfile ucc-soft-arch-golang
```
Build manual frontend (PROD):
```powershell
docker build -t is3-frontend:v1.0 --build-arg NEXT_PUBLIC_API_URL=http://localhost:8000 -f ucc-arq-soft-front/Dockerfile ucc-arq-soft-front
```
Build frontend QA:
```powershell
docker build -t is3-frontend:qa --build-arg NEXT_PUBLIC_API_URL=http://localhost:8001 -f ucc-arq-soft-front/Dockerfile ucc-arq-soft-front
```

## Levantar toda la plataforma (PROD + QA simultáneo)
```powershell
docker compose -f docker-compose.yaml up -d --build
```

Ver estado:
```powershell
docker compose -f docker-compose.yaml ps
```

Logs:
```powershell
docker compose -f docker-compose.yaml logs -f backend backend_qa frontend frontend_qa db db_qa
```

Bajar servicios:
```powershell
docker compose -f docker-compose.yaml down
# (Opcional, perderás datos)
docker compose -f docker-compose.yaml down --volumes --remove-orphans
```

## Acceso a la aplicación
- Frontend PROD: http://localhost:3000
- Frontend QA:   http://localhost:3001
- Backend health PROD: http://localhost:8000/health
- Backend health QA:   http://localhost:8001/health

### URLs en Azure (cuando se despliega con CI/CD)
- Frontend QA:   https://is3-frontend-qa.azurewebsites.net
- Backend QA:    https://is3-backend-qa.azurewebsites.net/health
- Frontend PROD: https://is3-frontend-prod.azurewebsites.net
- Backend PROD:  https://is3-backend-prod.azurewebsites.net/health

## Conectarse a la base de datos
Shell Postgres PROD:
```powershell
docker compose -f docker-compose.yaml exec db psql -U app -d app
```
Shell Postgres QA:
```powershell
docker compose -f docker-compose.yaml exec db_qa psql -U app -d app
```
Ejemplo de consulta:
```powershell
docker compose -f docker-compose.yaml exec db psql -U app -d app -c "\dt;"
```

## Verificación rápida post despliegue
```powershell
# 1. Contenedores levantados
docker compose -f docker-compose.yaml ps

# 2. Health de ambos backends
Invoke-WebRequest -Uri http://localhost:8000/health -UseBasicParsing
Invoke-WebRequest -Uri http://localhost:8001/health -UseBasicParsing

# 3. Frontends responden HTML
Invoke-WebRequest -Uri http://localhost:3000 -UseBasicParsing
Invoke-WebRequest -Uri http://localhost:3001 -UseBasicParsing

# 4. Diferenciación de DB (opcional)
docker compose -f docker-compose.yaml exec db_qa psql -U app -d app -c "CREATE TABLE IF NOT EXISTS test_dummy(id int primary key); INSERT INTO test_dummy VALUES (1) ON CONFLICT DO NOTHING; SELECT * FROM test_dummy;"
docker compose -f docker-compose.yaml exec db psql -U app -d app -c "SELECT * FROM test_dummy;"  # Debe estar vacía si no lo creaste allí
```

## Testing del Frontend

Para una guía completa y sencilla sobre cómo correr los tests del frontend, leer los reportes y entender el coverage (incluyendo dónde ver `index.html` y qué significan las métricas):

- Ver `testing_frontend.md` en la raíz del repositorio.

## Publicar imágenes en Docker Hub (ejemplo)
```powershell
# Login
docker login

# Retag
docker tag is3-backend:v1.0 <usuario>/is3-backend:v1.0
docker tag is3-backend:v1.0 <usuario>/is3-backend:latest

docker tag is3-frontend:v1.0 <usuario>/is3-frontend:v1.0
docker tag is3-frontend:v1.0 <usuario>/is3-frontend:latest

# Push
docker push <usuario>/is3-backend:v1.0
docker push <usuario>/is3-backend:latest

docker push <usuario>/is3-frontend:v1.0
docker push <usuario>/is3-frontend:latest
```
Luego actualizar `docker-compose.yaml`:
```yaml
backend:
  image: <usuario>/is3-backend:v1.0
frontend:
  image: <usuario>/is3-frontend:v1.0
```
Para QA si quieres imagen separada:
```yaml
frontend_qa:
  image: <usuario>/is3-frontend:qa
```

## Notas / Próximos pasos
- (Mejora) Unificar frontend PROD/QA en una sola imagen usando runtime config.
- (Mejora) Añadir `NEXT_PUBLIC_ENV` para mostrar banner visual en la UI.
- (Mejora) Automatizar build & push con GitHub Actions.
- (Seguridad) Cambiar contraseñas por variables seguras en un `.env` no versionado.

---
## TP05/TP08 – Backend Go con ACR y Web App (Containers)

Para cubrir TP05 (CD con QA→PROD y aprobaciones) y TP08 (uso de ACR), puedes usar el pipeline `azure-pipelines.tp05-08-backend.yaml`.

Checklist de configuración:
- Azure: crear ACR y dos Web Apps Linux (containers) para backend: QA y PROD.
- Azure DevOps:
  - Service connection Azure RM: `azure-tp05-connection` (o ajusta la variable `azureSubscription`).
  - Service connection Docker Registry a tu ACR (ajusta `dockerRegistryServiceConnection`).
  - Environments: `QA` y `PROD` (agrega aprobación manual en `PROD`).
  - Variables secretas: `DATABASE_URL_QA`, `DATABASE_URL_PROD`.
- Variables del YAML: ajusta `acrLoginServer`, `backendAppQA`, `backendAppProd`, `backendUrlQA`, `backendUrlProd`.

Rollback: redeploy del tag anterior (BuildId previo) desde el historial del pipeline.

---
## CI/CD en Azure DevOps (build once → deploy many)

- Pipelines agregados:
  - `azure-pipelines.release.yaml` (build + push + deploy QA/Prod, todo en uno)
  - `azure-pipelines.deploy.yaml` (solo deploy QA/Prod consumiendo un tag de imagen existente del ACR)
  - `azure-pipelines.tp05-08-backend.yaml` (backend Go → ACR → Web App Linux Containers, QA→PROD con aprobación y health checks)
- Flujo:
  - Modo todo-en-uno (`azure-pipelines.release.yaml`):
    1) BuildAndPush: construye imágenes `is3-frontend:<tag>` e `is3-backend:<tag>` y las publica en ACR.
    2) Deploy_Testing: despliega a Web Apps QA usando el MISMO `<tag>`; valida `/health` y `/`.
    3) Aprobación manual.
    4) Deploy_Prod: despliega el MISMO `<tag>` a Prod; valida salud.

  - Modo separado (Build existente + `azure-pipelines.deploy.yaml`):
    1) Tu build (TP04) publica imágenes en ACR con un `<tag>` conocido (p. ej. BuildId).
    2) Ejecutas `azure-pipelines.deploy.yaml` pasando `imageTag=<tag>` para QA.
    3) Aprobación manual.
    4) El mismo `imageTag` se despliega a Prod.

Requisitos en Azure DevOps:
- Service Connection (Docker Registry) apuntando al ACR: `is3acr-service-connection`.
- Service Connection (Azure Resource Manager) para Web Apps: `IS3-Azure-Subscription`.
- Environments `is3-qa` e `is3-prod`; configurar aprobación en `is3-prod`.

Ejecutar el deploy-only pipeline manualmente (ejemplo):
- Variables requeridas: `imageTag` (el tag exacto existente en ACR).
- Corre el pipeline y define `imageTag` en el formulario de ejecución.

Variables por entorno (App Settings en Web Apps):
- Backend QA: `PORT=8000`, `ENV=qa`, `DATABASE_URL=postgres://app:***@.../app?sslmode=require`.
- Backend PROD: `PORT=8000`, `ENV=production`, `DATABASE_URL=postgres://app:***@.../app?sslmode=require`.
- Frontend QA: `RUNTIME_PUBLIC_API_URL=https://is3-backend-qa.azurewebsites.net`, `INTERNAL_API=...`, `NEXT_PUBLIC_API_URL=...`.
- Frontend PROD: `RUNTIME_PUBLIC_API_URL=https://is3-backend-prod.azurewebsites.net`, `INTERNAL_API=...`, `NEXT_PUBLIC_API_URL=...`.

Estrategia de rollback:
- Hacer redeploy del release anterior (mismo tag previo) desde el historial de releases.

---
Para dudas o nuevas tareas (por ejemplo agregar banner de entorno), abrir un issue o pedirlo directamente.
