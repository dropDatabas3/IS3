# Decisiones del Proyecto

## Aplicación elegida
Se seleccionó una aplicación web desarrollada previamente para otra materia.  
La decisión se basó en aprovechar un trabajo ya realizado que cumple con los requisitos de este proyecto, lo que permitió ahorrar tiempo y garantizar que la aplicación ya estuviera probada y funcionando.

## Motivo de la elección
La aplicación fue elegida porque ya contábamos con una base sólida de código y arquitectura.  
Esto nos permitió enfocarnos en el proceso de containerización con Docker en lugar de invertir tiempo en desarrollar una aplicación nueva desde cero.

## Configuración de Docker
- Se configuró un archivo `Dockerfile` para definir la imagen de la aplicación.  
- Se utilizó `docker-compose.yaml` para orquestar los servicios y facilitar el despliegue.  
- Se creó un archivo `.env` y un `.env.template` para manejar variables de entorno de forma segura y reproducible.  

Con esta configuración, la aplicación puede ejecutarse en contenedores, aislando dependencias y simplificando la portabilidad entre entornos.

## Estrategia QA / PROD usando la misma imagen

Decisión: construir una única imagen por servicio (frontend y backend) y ejecutar dos instancias con configuraciones distintas para QA y PROD. Esto evita diferencias entre imágenes, reduce el tiempo de construcción y asegura que lo que se prueba en QA sea exactamente la misma imagen que se despliega en producción.

Cómo se aplica:
- Variables de entorno: todas las diferencias entre QA y PROD (cadenas de conexión, modo, puertos públicos, credenciales, flags de debugging) se inyectan mediante variables de entorno. Estas variables pueden definirse en:
	- `.env` o `.env.template` para valores por defecto y ejemplos
	- variables de entorno del host antes de lanzar `docker compose`
	- `env_file` o `environment` en `docker-compose.yaml`

- Mismo artefacto, distinta configuración: en `docker-compose.yaml` se configura `backend` y `backend_qa` para usar la misma imagen (`is3-backend:latest`) o la imagen construida por el Dockerfile. Igual para `frontend` y `frontend_qa`.

- Puertos y nombres: para poder correr ambos stacks simultáneamente se exponen en puertos distintos (ej: `frontend` en 3000 / `frontend_qa` en 3001, `backend` en 8000 / `backend_qa` en 8001, `db` en 5432 / `db_qa` en 5433). Los servicios internos se resuelven por nombre dentro de la red `appnet` (ej: `backend` o `backend_qa`).

Variables clave definidas y su propósito:
- `DATABASE_URL` / `QA_DATABASE_URL`: cadena de conexión completa preferida por la aplicación.
- `PGHOST`, `PGPORT`, `PGUSER`, `PGPASSWORD`, `PGDATABASE`: variables usadas por la librería de DB si `DATABASE_URL` no está presente.
- `INTERNAL_API` / `QA_INTERNAL_API`: URL interna que los procesos server-side (SSR) usan para comunicarse entre contenedores (`http://backend:8000` o `http://backend_qa:8000`).
- `NEXT_PUBLIC_API_URL` / `QA_NEXT_PUBLIC_API_URL`: URL pública usada por el cliente en el navegador (`http://localhost:8000` o `http://localhost:3001` para QA), embebida en el bundle público.

Cómo levantar ambos entornos (ejemplos):
- Levantar producción (servicios prod actuales):
	- `docker compose up -d backend db frontend`
- Levantar QA (mismos artefactos, distinta configuración):
	- `docker compose up -d db_qa backend_qa frontend_qa`

Notas de seguridad y operación:
- No incluir secretos en el repositorio. Usa `.env` en máquina local o mecanismos de secret management en CI/CD/producción.
- Para pruebas en QA con una DB distinta, usa `QA_DATABASE_URL` apuntando a una instancia de prueba.
- Validar CORS/secretos y logs dependientes del entorno.

Conclusión: esta aproximación mantiene la paridad entre entornos y permite validar exactamente la misma imagen en QA y en producción, diferenciando únicamente la configuración a través de variables de entorno.

### Build y tagging automático

Se incluye un script PowerShell `scripts/build_and_tag.ps1` que construye las imágenes `backend` y `frontend` desde sus contextos y las etiqueta en dos tags por defecto: `prod` y `qa`.

Flujo recomendado:
1. Ejecutar el script en la máquina de desarrollo o CI:
	- `powershell -ExecutionPolicy Bypass -File .\scripts\build_and_tag.ps1`
	- Opcionalmente pasar tags personalizados: `-BackendTagProd v1.2.3 -BackendTagQA v1.2.3-qa`
2. Levantar servicios con `docker compose up -d backend frontend backend_qa frontend_qa`.

Ventajas:
- Garantiza que QA y PROD usen exactamente la misma build (mismo binario y dependencias) — solo cambia la configuración en runtime.
- Facilita reproducibilidad y debugging: si hay un bug en QA, la misma imagen se puede promover a prod.

## Comandos de ejemplo

PowerShell (construir y etiquetar imágenes):

```powershell
# Construir imágenes y etiquetarlas (prod y qa)
powershell -ExecutionPolicy Bypass -File .\scripts\build_and_tag.ps1

# Construir y pasar tags específicos
powershell -ExecutionPolicy Bypass -File .\scripts\build_and_tag.ps1 -BackendTagProd v1.2.3 -FrontendTagProd v1.2.3 -BackendTagQA v1.2.3-qa -FrontendTagQA v1.2.3-qa
```

Levantar sólo producción (usando `docker-compose.prod.yml`):

```powershell
# En el directorio del repo
docker compose -f docker-compose.prod.yml up -d --build
```

Levantar sólo QA (usando `docker-compose.qa.yml`):

```powershell
docker compose -f docker-compose.qa.yml up -d --build
```

Levantar ambos stacks simultáneamente (prod y qa) — útil para pruebas locales:

```powershell
# Levanta prod y qa (cada uno con su archivo compose)
docker compose -f docker-compose.prod.yml up -d --build
docker compose -f docker-compose.qa.yml up -d --build
```

Notas:
- Asegúrate de ajustar variables en `.env` o exportarlas en el host si necesitas conectar a bases externas (p. ej. `DATABASE_URL`).
- Si en CI quieres empujar las imágenes a un registry, añade pasos que hagan `docker push is3-backend:<tag>` y `docker push is3-frontend:<tag>`.

## Checklist para entornos reproducibles

1. Versionado de artefactos
	- Construir y etiquetar imágenes en CI (o localmente usando `scripts/build_and_tag.ps1`) y subirlas a un registry privado o público. Usar tags semánticos o commits SHA.

2. Variables de entorno
	- Mantener un `.env.template` con todas las variables necesarias y valores por defecto (sin secretos).
	- Añadir `.env` locales y `.env.qa` para pruebas locales (estas no deben subirse al repo).
	- En `docker-compose.yaml` y `docker-compose.qa.yml` usar `env_file` para inyectar variables reproducibles.

3. Volúmenes y datos persistentes
	- Definir volúmenes en `docker-compose` (ya configurado: `db_data`, `db_data_qa` o `db_data_prod`) para que la base de datos conserve el estado si es requerido.

4. Versiones de runtime
	- Documentar la versión mínima de Docker Engine y Docker Compose en el README o en CI (ej: Docker Engine >= 24.0, Compose v2).

5. Comandos de arranque reproducibles
	- En local / CI usar los mismos comandos indicados arriba para levantar los servicios. En CI preferir `docker compose -f docker-compose.prod.yml pull` antes de `up` si las imágenes están en un registry.

6. Seguridad
	- No commitear secretos. Usar secret managers o variables de entorno definidas en el entorno de ejecución (CI, servidor).

Con estos pasos, el entorno debería arrancar con la misma configuración en cualquier máquina que tenga Docker/Compose instalado.

---

## Elección de la aplicación y tecnología utilizada
Se optó por una arquitectura simple de 3 capas: frontend (Next.js/React), backend (Go + Gin + GORM) y Postgres como base de datos relacional. Razones:
- Go ofrece binarios livianos y tiempos de arranque rápidos para servicios containerizados.
- Next.js facilita SSR/ISR y empaqueta un frontend moderno en una sola imagen.
- PostgreSQL es un estándar robusto, soporta tipos avanzados y se integra bien con GORM.

## Elección de imagen base y justificación
Backend:
- `golang:1.22-alpine` en la etapa build: imagen oficial, liviana (Alpine) + herramientas necesarias.
- `alpine:latest` para runtime: minimiza superficie de ataque y tamaño final; sólo se copian el binario y certificados.
Frontend:
- `node:18-alpine` en etapas `deps`, `builder` y `runner` para mantener consistencia y optimizar tamaño. Node 18 es LTS estable.

Beneficios de multi-stage:
- Reduce a ~40MB (backend) y evita incluir toolchain Go completo.
- Separa dependencias de build (npm ci) de runtime.

## Elección de base de datos y justificación
- PostgreSQL 15 (imagen `postgres:15-alpine`): soporte a features modernas, ecosistema amplio, buena integración con ORMs. Uso de Alpine minimiza tamaño de capa.
- Separación de instancias (prod y qa) para evitar contaminación de datos y permitir pruebas paralelas.

## Estructura y justificación del Dockerfile (Backend)
Etapas:
1. build: descarga módulos, compila con `CGO_ENABLED=0` para binario estático.
2. runtime: sólo binario + certificados -> menor superficie.
Decisiones clave:
- `CGO_ENABLED=0` simplifica la portabilidad.
- `apk add build-base` sólo en build stage.
- `ENTRYPOINT` directo al binario; `EXPOSE 8000` documenta el puerto.

## Estructura y justificación del Dockerfile (Frontend)
Etapas:
1. deps: instala dependencias con `npm ci` (reproducible).
2. builder: copia código y ejecuta `npm run build` con `ARG NEXT_PUBLIC_API_URL` para incrustar endpoint.
3. runner: imagen final sólo con `.next`, `public`, `node_modules` necesarios.
Decisiones:
- `npm ci` asegura reproducción exacta según lockfile.
- Uso de build args para API pública (aunque luego se separó imagen QA para distinto endpoint).

## Configuración de QA y PROD (variables de entorno)
En `docker-compose.yaml` se definen dos conjuntos de servicios:
- `backend` / `frontend` / `db` (prod)
- `backend_qa` / `frontend_qa` / `db_qa` (qa)

Variables diferenciadas:
- Conexión DB: `DATABASE_URL` vs `QA_DATABASE_URL` + `PGHOST=db` vs `PGHOST=db_qa`.
- API pública frontend: `NEXT_PUBLIC_API_URL=http://localhost:8000` (prod) y en imagen QA `http://localhost:8001`.
- API interna SSR: `INTERNAL_API=http://backend:8000` y `INTERNAL_API=http://backend_qa:8000`.

## Estrategia de persistencia de datos (volúmenes)
Volúmenes declarados:
- `db_data` montado en `/var/lib/postgresql/data` para prod.
- `db_data_qa` para QA.
Justificación:
- Permite reiniciar contenedores sin perder datos.
- Mantiene independencia entre entornos y evita contaminación cruzada.

## Estrategia de versionado y publicación
- Tag semántico inicial: `v1.0` creado y pusheado al repositorio Git.
- Imágenes Docker etiquetadas: `is3-backend:v1.0`, `is3-frontend:v1.0`, `is3-frontend:qa` (local) y publicadas como `nallarmariano/is3-*`.
- Tags adicionales recomendados: `latest` para el tag estable actual y `qa` para pre-release / validación.
- Script `scripts/push_images.ps1` automatiza retag y push (stable/qa/latest).

## Evidencia de funcionamiento
1. Contenedores levantados (prod + qa simultáneo):
```
docker compose -f docker-compose.yaml ps
NAME               IMAGE                                   PORTS
is3-frontend-1     nallarmariano/is3-frontend:v1.0         0.0.0.0:3000->3000
is3-frontend_qa-1  nallarmariano/is3-frontend:qa           0.0.0.0:3001->3000
is3-backend-1      nallarmariano/is3-backend:v1.0          0.0.0.0:8000->8000
is3-backend_qa-1   nallarmariano/is3-backend:v1.0          0.0.0.0:8001->8000
is3-db-1           postgres:15-alpine                      0.0.0.0:5432->5432
is3-db_qa-1        postgres:15-alpine                      0.0.0.0:5433->5432
```
2. Health checks OK:
```
curl http://localhost:8000/health -> {"status":"UP"}
curl http://localhost:8001/health -> {"status":"UP"}
```
3. Diferencia de bases (tabla dummy sólo en QA):
```
docker compose exec db_qa psql -U app -d app -c "CREATE TABLE IF NOT EXISTS test_dummy(id int primary key); INSERT INTO test_dummy VALUES (1) ON CONFLICT DO NOTHING; SELECT * FROM test_dummy;"
 id
----
	1
(1 row)

docker compose exec db psql -U app -d app -c "SELECT * FROM test_dummy;" -> (0 rows)
```
4. Persistencia entre reinicios:
```
docker compose down
docker compose up -d
docker compose exec db_qa psql -U app -d app -c "SELECT * FROM test_dummy;"  # sigue presente
```
5. Publicación de imágenes (extracto):
```
docker push nallarmariano/is3-backend:v1.0
docker push nallarmariano/is3-frontend:v1.0
docker push nallarmariano/is3-frontend:qa
```

## Problemas y soluciones
| Problema | Causa | Solución |
|----------|-------|----------|
| Frontend QA apuntaba a backend PROD | Variable `NEXT_PUBLIC_API_URL` bakeada igual en la imagen | Crear imagen separada QA con build-arg distinto (`is3-frontend:qa`) |
| Warning `version` obsoleto en compose | Compose v2 ignora `version:` | Remover clave `version:` del `docker-compose.yaml` |
| Port binding conflicts (3000/8000/5432) | Contenedores previos activos | `docker compose down` y recrear, o usar puertos QA distintos (3001/8001/5433) |
| 404 al crear release vía API GitHub | Token inválido/truncado | Regenerar PAT con scopes correctos y usar UI para release |
| Confusión sobre DB compartida | Frontend llamaba a backend equivocado | Verificación env + rebuild imagen QA |
| Variables `NEXT_PUBLIC_` no cambiaban en runtime | Next.js embebe en build | Rebuild con build args o futura runtime-config |

## Próximos pasos sugeridos
- Unificar imágenes frontend (runtime config dinámico para endpoint API).
- Pipeline GitHub Actions para build/push automático en tags semánticos.
- Añadir banner visual QA/PROD usando `NEXT_PUBLIC_ENV`.
- Externalizar credenciales sensibles a un secret manager.

---
Documento elaborado para justificar decisiones de arquitectura, containerización y operaciones de despliegue multi-entorno (QA/PROD).


