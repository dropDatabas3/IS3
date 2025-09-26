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
- Variables de entorno: todas las diferencias entre QA y PROD (cadenas de conexión, modo, puertos públicos, credenciales, flags de debugging) se inyectan mediante variables de entorno.


- Puertos y nombres: para poder correr ambos stacks simultáneamente se exponen en puertos distintos (ej: `frontend` en 3000 / `frontend_qa` en 3001, `backend` en 8000 / `backend_qa` en 8001, `db` en 5432 / `db_qa` en 5433). Los servicios internos se resuelven por nombre dentro de la red `appnet` (ej: `backend` o `backend_qa`).


Cómo levantar ambos entornos (ejemplos):
- Levantar producción (servicios prod actuales):
	- `docker compose up -d backend db frontend`
- Levantar QA (mismos artefactos, distinta configuración):
	- `docker compose up -d db_qa backend_qa frontend_qa`


Conclusión: esta aproximación mantiene la paridad entre entornos y permite validar exactamente la misma imagen en QA y en producción, diferenciando únicamente la configuración a través de variables de entorno.

### Build y tagging automático

Se incluye un script PowerShell `scripts/build_and_tag.ps1` que construye las imágenes `backend` y `frontend` desde sus contextos y las etiqueta en dos tags por defecto: `prod` y `qa`.


Ventajas:
- Garantiza que QA y PROD usen exactamente la misma build (mismo binario y dependencias) — solo cambia la configuración en runtime.
- Facilita reproducibilidad y debugging: si hay un bug en QA, la misma imagen se puede promover a prod.

## Comandos de ejemplo

Levantar sólo producción (usando `docker-compose.prod.yml`):

```powershell
# En el directorio del repo
docker compose -f docker-compose.prod.yml up -d --build
```

Levantar sólo QA (usando `docker-compose.qa.yml`):
NAME               IMAGE                                   PORTS
is3-db_qa-1        postgres:15-alpine                      0.0.0.0:5433->5432
3. Diferencia de bases (tabla dummy sólo en QA):
```
# Decisiones del Proyecto (Versión resumida para exposición)

Este documento explica, en lenguaje simple, qué hicimos, por qué lo hicimos y cómo demostramos que funciona.

## 1. Aplicación y tecnologías
Elegimos una app web que ya existía (cursos) para enfocarnos en la parte de despliegue. Usamos:
- Frontend: Next.js (React) → porque acelera el desarrollo y ya trae optimizaciones.
- Backend: Go (framework Gin) → rápido, liviano y fácil de empaquetar.
- Base de datos: PostgreSQL → confiable y estándar.

## 2. Imágenes base elegidas
- Backend: partimos de una imagen oficial de Go para compilar y luego pasamos el ejecutable a una imagen mínima (Alpine). Así reducimos tamaño y riesgos.
- Frontend: usamos Node 18 (versión LTS estable) también en Alpine para que pese menos.

Idea clave: construir en una etapa “grande” y ejecutar en una etapa “chica”.

## 3. Base de datos
PostgreSQL porque:
- Es conocida, potente y soporta bien crecimiento.
- Tiene buena integración con las librerías que usamos.
- Permite correr dos instancias (PROD y QA) sin complicaciones.

## 4. Estructura de los Dockerfile (explicado sencillo)
Backend:
1. Descargar dependencias y compilar.
2. Mover el ejecutable a una imagen mínima.

Frontend:
1. Instalar dependencias.
2. Generar el build optimizado.
3. Servir sólo lo necesario (el build final) en una imagen limpia.

Objetivo: imágenes más chicas, arranque rápido y menos cosas “sobrando adentro”.

## 5. QA y PROD: cómo los separamos
Corremos los dos entornos al mismo tiempo para probar sin interferir:
- Frontend PROD: puerto 3000 → habla con backend en 8000.
- Frontend QA: puerto 3001 → habla con backend en 8001.
- Cada uno tiene su propia base: 5432 (prod), 5433 (qa).

Usamos variables de entorno para cambiar conexiones / URLs sin reconstruir todo (salvo el caso del frontend QA donde hicimos una imagen propia para apuntar correctamente al backend de QA).

## 6. Persistencia (volúmenes)
Creamos dos “cajones” separados para guardar los datos de las bases:
- `db_data` (producción)
- `db_data_qa` (qa)
Si apagamos y prendemos los contenedores, los datos siguen ahí.

## 7. Versionado y publicación
- Creamos un tag de código: `v1.0`.
- Construimos imágenes y las subimos a Docker Hub con nuestro usuario: `nallarmariano`.
- Tags usados: `v1.0` (estable) y `qa` (para pruebas); se puede usar también `latest`.
- Agregamos un script (`push_images.ps1`) para automatizar los push.

## 8. Evidencias de que funciona
1. Aplicación corriendo en ambos entornos: se ven los 6 contenedores (2 front, 2 back, 2 DB) con puertos distintos.
2. Health check: `/health` responde 200 en 8000 y 8001.
3. Bases separadas: creamos un dato sólo en QA y no aparece en PROD.
4. Persistencia: bajamos todo y al volver a levantar, el dato en QA seguía.
5. Imágenes subidas: al hacer `docker pull nallarmariano/is3-frontend:qa` (y las demás) ya están disponibles.

## 9. Problemas encontrados y cómo los resolvimos
| Problema | Qué pasaba | Solución corta |
|----------|------------|----------------|
| Frontend QA consultaba a backend PROD | Veíamos mismos datos | Crear imagen QA con URL correcta (8001) |
| Warning por `version` en compose | Mensaje molesto | Quitar la línea `version:` |
| Puertos ocupados | Conflictos al levantar todos | Asignar puertos diferentes para QA (3001, 8001, 5433) |
| 404 creando release por API | Token mal generado | Regenerar token / usar interfaz web |
| Confusión de qué base se usaba | Dudábamos si compartían datos | Ver variables dentro de cada contenedor y prueba insertando registro |
| Cambiar API del frontend sin rebuild | Variable no se actualizaba | Explicar limitación (Next.js la fija en build) y separar imagen QA |

## 10. Qué podríamos mejorar después
- Un banner que diga “QA” o “PROD” en el frontend.
- Usar una sola imagen de frontend con configuración dinámica al arrancar.
- Pipeline automático (CI) para construir y publicar imágenes al hacer un tag.
- Guardar contraseñas reales fuera de los archivos (secret manager).

## 11. Resumen final
Logramos levantar y aislar dos entornos completos (PROD y QA) con la misma base de código, manteniendo datos, configuraciones y puertos separados. Documentamos el proceso, publicamos las imágenes y demostramos que cada entorno funciona de forma independiente.

Este documento busca explicar rápidamente las decisiones sin entrar en demasiado detalle técnico, para poder presentarlo de forma clara.
is3-db-1           postgres:15-alpine                      0.0.0.0:5432->5432

is3-db_qa-1        postgres:15-alpine                      0.0.0.0:5433->5432

```
