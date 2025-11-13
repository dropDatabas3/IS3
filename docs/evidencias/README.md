# Evidencias requeridas (TP05)

Adjuntar capturas y/o enlaces que muestren:

1) Recursos Cloud
- Resource Groups (rg-is3-qa, rg-is3-prod, rg-is3-shared)
- ACR is3acr (repositorios is3-frontend, is3-backend con tags)
- App Service Plan asp-is3-shared (Linux, B1)
- Web Apps (is3-frontend-qa/prod, is3-backend-qa/prod)
- PostgreSQL Flexible (is3pgqa, is3pgprod)

2) Release Pipeline (CI/CD)
- Pipeline(s) creados (YAML referenciado):
  - azure-pipelines.deploy.yaml (deploy-only) y/o azure-pipelines.release.yaml (todo-en-uno)
- Ejecuciones exitosas con:
  - Despliegue a QA OK
  - Aprobación manual antes de Prod (screenshot del Environment is3-prod con “Approvals and checks” y del cuadro de aprobación con comentario)
  - Despliegue a Prod OK

3) Health checks post‑deploy
- Backend QA: 200 en https://is3-backend-qa.azurewebsites.net/health
- Frontend QA: 200 en https://is3-frontend-qa.azurewebsites.net/
- Backend Prod: 200 en https://is3-backend-prod.azurewebsites.net/health
- Frontend Prod: 200 en https://is3-frontend-prod.azurewebsites.net/

4) Variables por entorno
- Captura de App Settings en Web Apps QA y Prod (RUNTIME_PUBLIC_API_URL, INTERNAL_API, DATABASE_URL, PORT, etc.)

5) Rollback
- Evidencia de re‑deploy del último tag OK (historial de releases o captura del tag en ACR usado para volver atrás).
 - (Opcional) Captura del tag `last-known-good` apuntando al mismo digest que el despliegue estable.

6) Gates y checks (si añadidos)
- Captura de configuración de otros checks (Delay, REST, Work Items, Security scan) en Environment is3-prod.
- Evidencia de rechazo de aprobación (pipeline queda en estado Cancelled) si se simula un fallo.
