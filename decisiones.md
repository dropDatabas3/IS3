# Decisiones del Proyecto

Este documento resume qué hicimos, por qué lo hicimos y cómo comprobamos que funciona, siguiendo los puntos que pide el trabajo práctico.

---
## 1) Elegimos y preparamos la aplicación
- Selección: usamos una plataforma de cursos que ya habíamos trabajado en otra materia. Esto nos permitió enfocarnos en Docker y entornos, y no en inventar una app nueva.
- Repositorio: subimos todo a GitHub en un mono‑repo con dos carpetas claras: `ucc-arq-soft-front` (frontend) y `ucc-soft-arch-golang` (backend).
- Entorno Docker: instalamos Docker Desktop y configuramos Compose para poder levantar frontend, backend y base de datos con un solo comando. Dejamos por escrito las decisiones para que cualquiera del equipo pueda reproducirlo.

¿Por qué esta app? Porque ya la conocíamos, tenía frontend, backend, y nos servía para practicar lo que pide el TP (contenedores, imágenes, QA/PROD, volúmenes, etc.).

---
## 2) Construimos imágenes personalizadas
- Dockerfile del backend (Go): usamos un build en dos etapas (multi‑stage). Primero compilamos el binario en una imagen de Go, y después lo pasamos a una imagen final más liviana (Alpine). Resultado: imagen más chica y más segura.
- Dockerfile del frontend (Next.js): también multi‑stage. Instalamos dependencias, hacemos el build de producción y dejamos una imagen final que sólo corre la app.
- Imagen base: elegimos imágenes oficiales (Go y Node 18 en Alpine) por estabilidad y tamaño reducido.
- Etiquetado: usamos tags claros (por ejemplo, `v1.0`). La idea es que cada versión tenga su etiqueta para saber exactamente qué estamos corriendo.

Notas de estructura: menos capas, menos peso, arranque rápido y sin herramientas de desarrollo en la imagen final.

---
## 3) Publicamos en Docker Hub (versionado)
- Subimos las imágenes a Docker Hub con el usuario `nallarmariano` (se puede cambiar por variable si se usa otro usuario).
- Estrategia de tags: usamos `v1.0` como versión estable. Los entornos (QA/PROD) comparten la misma imagen y cambian por variables. Evitamos “pisar” un tag: si hay cambios, creamos `v1.1`, `v1.2`, etc.
- Beneficio: podemos “anclar” un despliegue a una versión exacta y volver atrás si hace falta.

---
- Subimos las imágenes a Docker Hub con el usuario `nallarmariano` (se puede cambiar por variable si se usa otro usuario).
- Estrategia de tags: usamos `v1.0` como versión estable. Los entornos (QA/PROD) comparten la misma imagen y cambian por variables. Evitamos “pisar” un tag: si hay cambios, creamos `v1.1`, `v1.2`, etc.
- Beneficio: podemos “anclar” un despliegue a una versión exacta y volver atrás si hace falta.

---

## 4) Base de datos en contenedor y persistencia
  - En el combinado: `db_data` (prod) y `db_data_qa` (qa)
  - En prod: `db_data_prod`
  - En qa: `db_data_qa`

## 5) QA y PROD con la misma imagen (variables de entorno)
  - PROD: frontend 3000, backend 8000, DB 5432
  - QA: frontend 3001, backend 8001, DB 5433

En resumen: una sola imagen, dos comportamientos según variables y puertos. Sin rebuild para cambiar de entorno.

## 6) Entorno reproducible con docker-compose
  - `docker-compose.yaml` (combinado: levanta QA y PROD a la vez)
  - `docker-compose.prod.yml` (sólo PROD)
  - `docker-compose.qa.yml` (sólo QA)
  - `DOCKER_USER` (usuario de Docker Hub)
  - `FRONTEND_TAG` / `BACKEND_TAG`
  - `QA_FRONTEND_TAG` / `QA_BACKEND_TAG`

Esto permite que el entorno se ejecute igual en cualquier máquina que tenga Docker y acceso a Internet.

## 7) Versión etiquetada y uso en Compose

## Evidencia de funcionamiento (resumen)

## Problemas comunes y cómo los resolvimos
 - El frontend QA mostraba datos de PROD: antes de la mejora, Next.js “horneaba” la URL. Lo solucionamos generando la configuración en runtime y usando una sola imagen.
 - Acceso a la DB desde el host: no funciona con el nombre del servicio (`db_qa`) desde afuera; se accede con `localhost:5433` (QA) o `localhost:5432` (PROD).

## Notas finales
El objetivo fue dejar algo simple, claro y que se pueda explicar en la defensa: una app conocida, imágenes pequeñas, QA y PROD corriendo al mismo tiempo con la misma imagen, datos persistentes y versiones controladas.

## TP05 – CI/CD Release (build once, deploy many)

### Arquitectura elegida (cloud + release)

### Estrategia “build once → deploy many”

### Stages y aprobaciones

### Variables y secretos por entorno
  - PORT=8000
  - DATABASE_URL (QA/PROD difieren en host/credenciales; en Azure usar sslmode=require)
  - RUNTIME_PUBLIC_API_URL: URL pública del backend del mismo entorno.
  - INTERNAL_API: en este escenario, también la URL pública del backend (front y back están en Web Apps separadas).
  - NEXT_PUBLIC_API_URL: opcional, igual a RUNTIME_PUBLIC_API_URL para consistencia.

### Evidencias previstas

### Rollback

### Archivos de pipeline agregados
  - `azure-pipelines.release.yaml`: BuildAndPush → Deploy_Testing → Deploy_Prod (todo en uno, build y deploy).
  - `azure-pipelines.deploy.yaml`: sólo despliegue QA/Prod consumiendo un `imageTag` ya existente en ACR.
  - Requisitos comunes:
    - Service connection Docker Registry (ACR): para login/push o referenciar imágenes.
    - Service connection Azure Resource Manager: para Web Apps.
    - Environments `is3-qa` (sin aprobación) e `is3-prod` (con aprobación manual configurada).

### Aprobaciones y responsables
- Environment `is3-prod` con aprobación manual requerida.
- Responsables sugeridos:
  - Líder técnico (aprobación primaria).
  - Docente/ayudante (si aplica, como observador o segundo aprobador).
- Criterios de aprobación:
  - QA en verde (health checks OK, smoke tests básicos desde el frontend).
  - No hay errores críticos en el despliegue.
 
#### Proceso concreto de aprobación (Gates)
1. El pipeline `azure-pipelines.deploy.yaml` despliega QA con el tag construido en TP04 y ejecuta health checks.
2. Si QA falla (código ≠ 200 en `/health` o `/`), el pipeline termina y no se solicita aprobación.
3. Si QA pasa, Azure DevOps queda esperando aprobación manual en el Environment `is3-prod`:
   - Los aprobadores revisan: página frontend QA cargó, endpoint `/health` en QA responde, cambios esperados OK.
   - Verifican en ACR que el tag existe y no fue modificado (inmutabilidad).
4. Aprobación manual: se hace click en "Approve" con comentario breve (ej.: "QA OK, procedemos a Prod").
5. Se ejecuta despliegue a Prod (mismo tag) y health checks.
6. Registro de evidencia: captura de pantalla del diálogo de aprobación + resultado final del job en Prod.

#### Gates adicionales (opcionales)
- Escaneo de seguridad de imagen antes de Prod (herramienta externa o extensión DevOps).
- Gate de disponibilidad: un script que verifica latencia media del backend QA < X ms.
- Gate de rollback automático: si `/health` en Prod falla 3 veces seguidas (status != 200), disparar redeploy del tag anterior manualmente.

Con esto, el binario/imágenes validadas en Testing son exactamente las que llegan a Producción tras la aprobación manual.

---
## Inventario de recursos (Paso 1)

- Región: Brazil South (brazilsouth)
- Resource Groups:
  - rg-is3-qa (QA)
  - rg-is3-prod (Prod)
  - rg-is3-shared (compartido para ACR/ASP)
- ACR: is3acr (SKU Basic) en rg-is3-shared
- App Service Plan (Linux): asp-is3-shared (SKU B1) en rg-is3-shared
- Web Apps:
  - is3-backend-qa (rg-is3-qa)
  - is3-frontend-qa (rg-is3-qa)
  - is3-backend-prod (rg-is3-prod)
  - is3-frontend-prod (rg-is3-prod)
- PostgreSQL Flexible:
  - is3pgqa (rg-is3-qa), DB app, usuario app
  - is3pgprod (rg-is3-prod), DB app, usuario app

Notas:
- DOCKER_REGISTRY_SERVER_* configurado en las Web Apps para permitir pull de imágenes privadas del ACR.
- DATABASE_URL (Flexible Server): formato `postgres://app:<PASS>@<server>.postgres.database.azure.com:5432/app?sslmode=require`.
- Deploy_Prod: despliega el mismo tag a producción y revalida salud; si falla, se recomienda redeploy del tag previo.
