docker compose -f docker-compose.prod.yml up -d --build
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

Usamos variables de entorno para cambiar: URL del backend, credenciales DB, nombres de base, modo debug, etc. En el frontend QA hicimos una imagen aparte porque en Next.js algunas variables públicas se “hornean” en el build y no cambian sólo con reiniciar.

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
frontend_qa      nallarmariano/is3-frontend:qa        0.0.0.0:3001->3000
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
| No se actualizaba URL frontend sin rebuild | Next.js “fija” algunas vars | Aceptamos segunda imagen QA y lo documentamos |

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
