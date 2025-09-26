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
## 4) Base de datos en contenedor y persistencia
- Elegimos PostgreSQL porque es conocida, estable y se integra bien con Go.
- Persistencia: definimos volúmenes nombrados para que los datos no se pierdan al reiniciar contenedores.
  - En el combinado: `db_data` (prod) y `db_data_qa` (qa)
  - En prod: `db_data_prod`
  - En qa: `db_data_qa`
- Conexión: el backend se conecta por variables (host del servicio `db` o `db_qa` dentro de la red de Docker). Desde el host, para herramientas externas, usamos los puertos publicados (5432 para PROD y 5433 para QA).

---
## 5) QA y PROD con la misma imagen (variables de entorno)
- Levantamos QA y PROD al mismo tiempo usando la misma imagen para frontend y backend. Lo que cambia son las variables de entorno y los puertos.
- Puertos del host:
  - PROD: frontend 3000, backend 8000, DB 5432
  - QA: frontend 3001, backend 8001, DB 5433
- Frontend (punto importante): Next.js suele “fijar” variables en build. Para evitar tener dos imágenes, agregamos una configuración en tiempo de arranque (un archivo `public/runtime-config.js` que se genera cuando inicia el contenedor). Así, una sola imagen sirve tanto para PROD como para QA, y el navegador apunta al backend correcto según `RUNTIME_PUBLIC_API_URL`.
- Backend: toma la conexión a la base por variables (`DATABASE_URL` o `PGHOST`, `PGPORT`, etc.). La misma imagen sirve para ambos entornos.

En resumen: una sola imagen, dos comportamientos según variables y puertos. Sin rebuild para cambiar de entorno.

---
## 6) Entorno reproducible con docker-compose
- Archivos de Compose:
  - `docker-compose.yaml` (combinado: levanta QA y PROD a la vez)
  - `docker-compose.prod.yml` (sólo PROD)
  - `docker-compose.qa.yml` (sólo QA)
- Qué levantan: frontend, backend y su base de datos correspondiente, con sus puertos y volúmenes.
- Reproducibilidad: las imágenes se bajan desde Docker Hub por defecto (`nallarmariano/is3-frontend` y `nallarmariano/is3-backend` con `v1.0`). Si otro usuario quiere usar su propio namespace o tags, puede cambiarlo con variables de entorno sin tocar archivos:
  - `DOCKER_USER` (usuario de Docker Hub)
  - `FRONTEND_TAG` / `BACKEND_TAG`
  - `QA_FRONTEND_TAG` / `QA_BACKEND_TAG`

Esto permite que el entorno se ejecute igual en cualquier máquina que tenga Docker y acceso a Internet.

---
## 7) Versión etiquetada y uso en Compose
- Etiquetamos la versión estable como `v1.0` y actualizamos los Compose para usar esa etiqueta por defecto.
- Convención de versionado: simple y clara (vX.Y). Cuando cambia algo relevante, subimos la versión (v1.1, v1.2…). Evitamos usar `latest` en producción para no llevarnos sorpresas.

---
## Evidencia de funcionamiento (resumen)
- Ambos entornos corren a la vez sin pisarse: PROD (8000/3000/5432) y QA (8001/3001/5433).
- El frontend de QA apunta a su backend (8001) gracias al archivo de configuración en runtime.
- Los datos permanecen tras reiniciar contenedores porque usamos volúmenes.

---
## Problemas comunes y cómo los resolvimos
- El frontend QA mostraba datos de PROD: antes de la mejora, Next.js “horneaba” la URL. Lo solucionamos generando la configuración en runtime y usando una sola imagen.
- Acceso a la DB desde el host: no funciona con el nombre del servicio (`db_qa`) desde afuera; se accede con `localhost:5433` (QA) o `localhost:5432` (PROD).

