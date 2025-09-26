# Decisiones de Diseño - Pipeline CI Azure DevOps

## 📋 Información del Proyecto

**Aplicación**: Sistema de cursos IS3  
**Stack Tecnológico**: 
- **Frontend**: Next.js 14 (React, TypeScript, Tailwind CSS)
- **Backend**: Go 1.22 (Gin, GORM, PostgreSQL)
- **Estructura**: Mono-repo con carpetas `/ucc-arq-soft-front` y `/ucc-soft-arch-golang`

## 🎯 Decisiones de Pipeline CI

### 1. **¿Por qué YAML y no Classic Pipeline?**
- **Versionado**: El pipeline está versionado junto con el código fuente
- **Code Review**: Los cambios al pipeline pasan por PR review
- **Portabilidad**: Fácil migración entre proyectos y organizaciones
- **Reutilización**: Posibilidad de usar templates y librerías
- **Transparencia**: Todo el equipo puede ver y entender la configuración

### 2. **¿Por qué Self-Hosted Agent vs Microsoft-Hosted?**
#### Ventajas del Self-Hosted:
- **Control total** del entorno (versiones específicas de Node.js, Go)
- **Dependencias persistentes** (node_modules cache, Go module cache)
- **Sin límites de tiempo** de ejecución (Microsoft-Hosted tiene límite de 60min)
- **Recursos locales** (acceso a bases de datos internas, servicios locales)
- **Costo** (para builds largos o frecuentes)
- **Personalización** (herramientas específicas, configuraciones custom)

#### Para este proyecto específicamente:
- Build de Next.js puede ser lento (beneficia del cache persistente)
- Go compilation es rápida pero beneficia de module cache
- Control de versiones exactas (Go 1.22, Node 18)

### 3. **Estructura del Pipeline (Multi-Job en Single Stage)**

#### Stage único "CI" con 3 Jobs:
1. **BuildFrontend**: 
   - Install dependencies (npm ci)
   - Linting (npm run lint)
   - Build (npm run build)
   - Publish artifacts (.next + package.json)

2. **BuildBackend**:
   - Download Go modules
   - Static analysis (go vet)
   - Format check (go fmt)
   - Compile binary (optimized build)
   - Publish artifacts (binary + go.mod)

3. **PublishSummary**:
   - Consolidate build information
   - Display summary of published artifacts

#### ¿Por qué Jobs paralelos y no secuenciales?
- **Performance**: Frontend y backend builds son independientes
- **Eficiencia**: Aprovecha múltiples cores del self-hosted agent
- **Fail Fast**: Si uno falla, el otro continúa para dar feedback completo

### 4. **Triggers y PR Strategy**

```yaml
trigger:
  branches:
    include:
      - main
  paths:
    exclude:
      - '*.md'

pr:
  branches:
    include:
      - main
```

#### Decisiones:
- **Solo main**: Siguiendo la guía del TP (trigger en main)
- **PR Validation**: Valida cambios antes del merge
- **Path Exclusion**: No ejecuta en cambios de documentación
- **Branch Strategy**: Preparado para GitFlow (main + develop)

### 5. **Quality Gates Implementados**

#### Frontend:
- **Linting**: `npm run lint` (ESLint + Next.js rules)
- **Type Checking**: Implícito en `npm run build` (TypeScript)
- **Build Validation**: Asegura que la app compile correctamente

#### Backend:
- **Static Analysis**: `go vet ./...` (detección de bugs potenciales)
- **Format Check**: `go fmt ./...` (consistencia de código)
- **Dependency Validation**: `go mod verify` (integridad de dependencias)

#### ¿Por qué `continueOnError: true` en algunos steps?
- **Linting y formatting** son **informativos** pero no bloquean el build
- Permite ver **todos los issues** de una vez
- El **build real** sí debe fallar si hay errores críticos

### 6. **Artifact Strategy**

#### Artifacts Publicados:
- `frontend-dist`: Carpeta `.next` (build output de Next.js)
- `frontend-config`: `package.json` (metadata y dependencies)
- `backend-bin`: Binary compilado de Go (listo para deploy)
- `backend-config`: `go.mod` (metadata de dependencias)

#### ¿Por qué estos artifacts?
- **Completos**: Todo lo necesario para deployment posterior
- **Optimizados**: Solo lo esencial (no node_modules completos)
- **Metadata**: Información para debugging y dependency tracking

### 7. **Optimizaciones Implementadas**

#### Build Optimization:
```bash
# Go build optimizado
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./bin/app
```
- `CGO_ENABLED=0`: Binary estático (sin dependencias C)
- `GOOS=linux`: Target para contenedores/servers
- `-ldflags="-w -s"`: Remove debug info (binary más pequeño)

#### Dependency Optimization:
```bash
# npm más eficiente
npm ci --silent --prefer-offline
```
- `--prefer-offline`: Usa cache local primero
- `--silent`: Menos verbose output

### 8. **Versionado Automático**

```yaml
patchVersion: $[counter(variables['Build.SourceBranchName'], 0)]
buildVersion: '$(majorVersion).$(minorVersion).$(patchVersion)'
```

- **Semantic Versioning**: 1.0.X format
- **Auto-increment**: Patch version se incrementa automáticamente
- **Branch-based**: Counter independiente por branch

## 🚀 Extensiones Futuras (No implementadas - Solo CI)

Para CD (Continuous Deployment) se podrían agregar:
- Docker image builds
- Container registry push  
- Environment deployments (dev/qa/prod)
- Integration tests con docker-compose
- Automated rollback capabilities

## 📊 Métricas y Beneficios Esperados

- **Build Time**: ~3-5 minutos (vs ~8-10 en Microsoft-hosted)
- **Cache Hit Rate**: >80% después de primer build
- **Artifact Size**: <50MB total
- **Parallel Efficiency**: 2x speedup vs sequential builds

## 🔧 Configuración Requerida del Agent

Ver archivos: `agent-prerequisites.md` y `guia-selfhosted-agent.md`er compose -f docker-compose.prod.yml up -d --build
﻿# Decisiones del Proyecto (versión explicada “con nuestras palabras”)

Este archivo resume qué elegimos, por qué y cómo comprobamos que funciona, sin meternos de más en lo técnico.

---
## 1. Elección de la aplicación y tecnologías
Tomamos una app que ya teníamos (plataforma de cursos) para no perder tiempo reinventando algo nuevo. Así pudimos enfocarnos en Docker y entornos.

Tecnologías elegidas:
- Frontend: Next.js (sobre React) → rápido para armar interfaz moderna.
- Backend: Go (con Gin) → liviano, arranca rápido y es fácil de desplegar en contenedor.
- Base de datos: PostgreSQL → estable, conocida y soporta bien crecimiento.

---
## 2. Elección de imágenes base y justificación
- Backend: usamos una imagen oficial de Go para compilar y después “pasamos” el binario a una imagen final más chica (Alpine). Resultado: menos peso y menos cosas que puedan fallar.
- Frontend: Node 18 (LTS) en Alpine para reducir tamaño. Construimos el sitio y luego sólo servimos el build final.

Idea central: construir en una imagen más completa y ejecutar en otra mínima (multi‑stage) para ahorrar espacio y mejorar seguridad.

---
## 3. Base de datos y justificación
PostgreSQL porque:
- Ya la conocíamos y es estándar en muchos proyectos.
- Funciona bien con las librerías del backend.
- Podemos correr dos instancias (PROD y QA) separadas sin enrosque.

---
## 4. Estructura y justificación de los Dockerfile
Backend (Go):
1. Etapa “build”: instala dependencias y compila el ejecutable.
2. Etapa “runtime”: sólo trae el ejecutable y lo mínimo para correrlo.

Frontend (Next.js):
1. Instala dependencias (aprovecha cache para no repetir).
2. Genera el build optimizado de producción.
3. Imagen final liviana que sólo sirve el resultado.

¿Por qué así? Menos capas innecesarias, arranque rápido y menos superficie de ataque.

---
## 5. Configuración de QA y PROD (variables de entorno)
Corremos PROD y QA al mismo tiempo sin que se pisen:
- Frontend PROD: puerto 3000 → habla con backend 8000.
- Frontend QA: puerto 3001 → habla con backend 8001.
- Bases: 5432 (prod) y 5433 (qa).

Usamos variables de entorno para cambiar: URL del backend, credenciales DB, nombres de base, modo debug, etc. Inicialmente el frontend QA era una imagen aparte porque Next.js “hornea” variables públicas en el build. Luego lo mejoramos: ahora usamos una sola imagen y una configuración en runtime (archivo `public/runtime-config.js` generado al arrancar el contenedor) para que cada entorno apunte a su backend correcto sin rebuild.

Ejemplos simples (concepto):
```
DB_HOST=db          (PROD)
DB_HOST=db_qa       (QA)
BACKEND_PORT=8000   (PROD)
BACKEND_PORT=8001   (QA)
```

---
## 6. Estrategia de persistencia (volúmenes)
Creamos dos volúmenes distintos:
- `db_data` → guarda datos de producción.
- `db_data_qa` → guarda datos de QA.

Así, si apagamos contenedores y volvemos a levantar, la info sigue. Y lo de QA nunca mezcla lo de PROD.

---
## 7. Estrategia de versionado y publicación
- Marcamos el código con un tag: `v1.0`.
- Construimos imágenes y las subimos a Docker Hub (usuario: `nallarmariano`).
- Usamos tags: `v1.0`, `qa` y también `latest` para la versión “actual”.
- Script `build_and_tag.ps1` para construir y `push_images.ps1` para subir sin olvidos.

Esto permite: saber qué versión está corriendo, reproducir bugs y entregar algo “congelado”.

---
## 8. Evidencia de funcionamiento (logs y ejemplos)

### 8.1 Aplicación corriendo en ambos entornos
Salida (ejemplo) al listar contenedores:
```
CONTAINER        IMAGE                               PORTS
frontend         nallarmariano/is3-frontend:v1.0      0.0.0.0:3000->3000
backend          nallarmariano/is3-backend:v1.0       0.0.0.0:8000->8000
db               postgres:15-alpine                  0.0.0.0:5432->5432
frontend_qa      nallarmariano/is3-frontend:v1.0      0.0.0.0:3001->3000
backend_qa       nallarmariano/is3-backend:v1.0       0.0.0.0:8001->8000
db_qa            postgres:15-alpine                  0.0.0.0:5433->5432
```

### 8.2 Conexión exitosa a la base de datos
Logs del backend (ejemplo simplificado):
```
[info] DB connection established host=db user=postgres database=cursos
[info] DB connection established host=db_qa user=postgres database=cursos_qa (en servicio QA)
```

Health check (cada backend responde 200):
```
curl http://localhost:8000/health  → {"status":"ok"}
curl http://localhost:8001/health  → {"status":"ok"}
```

### 8.3 Persistencia entre reinicios
Prueba (QA):
1. Insertamos un registro de prueba.
2. Apagamos QA: `docker compose stop backend_qa db_qa`.
3. Volvemos a levantar: `docker compose start backend_qa db_qa`.
4. Consultamos y el registro sigue.

Ejemplo (texto):
```
SELECT * FROM courses WHERE code='QA-DEMO';
→ fila encontrada (después del reinicio)
```

Si se quiere reemplazar estos textos por capturas, podemos ponerlas en una carpeta `docs/` y referenciarlas aquí.

---
## 9. Problemas y soluciones
| Problema | Qué pasaba | Cómo lo resolvimos |
|----------|------------|--------------------|
| Frontend QA mostraba datos de PROD | Apuntaba al backend equivocado | Construimos imagen QA con la URL correcta (8001) |
| Warning por `version` en compose | Mensaje molesto en consola | Quitamos la línea `version:` |
| Choque de puertos al levantar todo | Servicios se pisaban | Asignamos puertos distintos para QA (3001, 8001, 5433) |
| Error 404 creando release vía API | Token incorrecto | Usamos interfaz web y ajustamos token |
| Duda sobre si las bases eran distintas | Parecían iguales | Revisamos variables y probamos insert selectivo |
| No se actualizaba URL frontend sin rebuild | Next.js “fija” algunas vars | Implementamos runtime-config (una sola imagen front para QA/PROD) |

---
## 10. Resumen final
Tenemos dos entornos completos (PROD y QA) corriendo a la vez, aislados por puertos, variables y volúmenes. Versionamos, publicamos imágenes y demostramos que los datos persisten y no se mezclan. El documento busca que cualquiera del equipo entienda las decisiones rápido.

---
## 11. Posibles mejoras futuras
- Un cartel visual en el frontend (banner QA / PROD).
- Configuración dinámica en frontend para evitar segunda imagen.
- Pipeline CI para build + push automático cuando haya un tag.
- Manejo de claves/secretos fuera de los archivos (secret manager).

Fin ✨
