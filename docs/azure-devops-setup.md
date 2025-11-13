# Azure DevOps – Setup de Release (Deploy-only conectado a TP04)

Este documento lista los pasos en Azure DevOps para ejecutar los pipelines agregados en este repo sin depender todavía de que los recursos estén funcionales. Sirve como checklist para cuando el compañero con la cuenta de Azure los cree.

## 1) Service connections
- Azure Resource Manager (ARM)
  - Nombre sugerido: IS3-Azure-Subscription
  - Tipo: Service principal (automatic)
  - Alcance: Suscripción (al menos permisos de Contributor en los Resource Groups de QA/Prod/Shared).
- Docker Registry (ACR) – solo si usas el pipeline todo‑en‑uno o si tu TP04 necesita pushear
  - Nombre sugerido: is3acr-service-connection
  - Servidor: https://is3acr.azurecr.io (ajusta al nombre real de tu ACR)
  - Autenticación: usuario/clave admin del ACR (o Managed Identity con permisos) 

## 2) Environments
- Crea dos environments:
  - is3-qa
  - is3-prod (configurar “Approvals and checks” con al menos una aprobación manual).

### 2.1) Configurar aprobación manual en is3-prod
1. Azure DevOps → Pipelines → Environments → is3-prod.
2. Click en "Approvals and checks" → "Add" → "Approvals".
3. Selecciona uno o más usuarios (ej: líder técnico, docente/ayudante).
4. (Opcional) Añade un tiempo de expiración y comentarios requeridos.
5. Guarda. A partir de ahora cualquier stage/job que declare `environment: is3-prod` quedará en espera de aprobación antes de ejecutar.

### 2.2) (Opcional) Otros checks/gates en el Environment
- Check de Query Work Items (verificar no haya bugs críticos abiertos).
- Check de invocación REST (script que valide métricas o latencia en QA).
- Escaneo de seguridad (extensiones de container scanning antes de Prod).
- Tiempo mínimo en QA (gate de Delay) para observación manual.

## 3) Variables y secretos
- Variables del pipeline (o Library Variable Group):
  - Para azure-pipelines.deploy.yaml: no hay secretos obligatorios (consumes imágenes ya existentes). 
    - imageTag se pasa como parámetro al ejecutar.
  - Para azure-pipelines.release.yaml (todo‑en‑uno):
    - databaseUrlQA (secret)
    - databaseUrlProd (secret)

## 4) Crear los pipelines
- Deploy-only:
  - Pipelines → New pipeline → Existing Azure Pipelines YAML → `azure-pipelines.deploy.yaml`.
  - Ejecutar: especifica `imageTag` (el tag que publicó el build TP04 en ACR).
- All‑in‑one (opcional):
  - Pipelines → New pipeline → Existing Azure Pipelines YAML → `azure-pipelines.release.yaml`.
  - Ejecutar: construye y publica imágenes → despliega QA → aprobación → despliega Prod.

## 5) (Opcional) Trigger automático desde el TP04
- En `azure-pipelines.deploy.yaml` hay un bloque `resources.pipelines` comentado:
  - Reemplaza `source: 'IS3-CI'` por el nombre de tu pipeline de TP04.
  - Mapear el BuildId del TP04 con el `imageTag` (ejecuta el deploy pasando ese valor).

## 6) Health checks
- Los pipelines ya incluyen health checks HTTP:
  - Backend: GET /health → 200
  - Frontend: GET / → 200
- En caso de fallo, el job se marca como error y no promueve a Prod.

## 6.1) Extender health checks (opcional)
- Añadir reintentos exponenciales (Invoke-WebRequest dentro de un bucle con backoff).
- Validar contenido esperado (ej.: buscar substring "<title>" en HTML del frontend).
- Medir latencia y publicar como logging (`Write-Host`).

## 7) Rollback
- Desde Releases/Run history: re‑desplegar el mismo `imageTag` del último run OK.
- Alternativa: mantener un “alias” en ACR (tag last-known-good) y re‑apuntar la Web App.

## 8) Flujo resumido de aprobación (Stage Prod)
1. Ejecutar pipeline (release o deploy-only) → despliega QA → health checks OK.
2. Pipeline queda "Waiting" en Environment is3-prod.
3. Aprobador revisa QA: frontend operativo, `/health` 200, logs sin errores críticos.
4. Aprobador hace "Approve" con comentario (ej: "QA OK, procedemos").
5. Se despliega Prod y se ejecutan health checks Prod.
6. Capturar evidencia (pantalla de aprobación + jobs en verde) y guardarla en `docs/evidencias/`.