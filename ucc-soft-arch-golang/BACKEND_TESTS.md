# Guía completa y simple de los tests del backend (Go)

Este documento explica, en lenguaje simple, todo lo que se hizo para probar el backend en Go: qué herramientas usamos, cómo funcionan, qué se prueba en cada parte, cómo ejecutar las pruebas y cómo ver el porcentaje de cobertura con un informe HTML.

Si no tenés experiencia en programación o testing, no te preocupes: la idea es que puedas entenderlo igual.


## Objetivo general

- Contar con pruebas unitarias y de componentes para las distintas capas del backend (clientes/clients, servicios/services, controladores/controllers, middleware, rutas/routes, etc.).
- Ver el porcentaje de cobertura de código (coverage) y un informe visual (HTML) para saber qué partes están cubiertas por pruebas.
- Facilitar la ejecución con un script en Windows PowerShell y dejar documentado cómo correr paquetes individuales también.

Estado actual de cobertura (última ejecución): 80.1% del proyecto total, con informe HTML generado.


## Herramientas y librerías usadas

- Go (lenguaje y runtime) y su sistema de pruebas integrado `go test`.
- Gin (framework HTTP) en modo de pruebas ("TestMode") para simular peticiones HTTP rápidamente.
- GORM (ORM para Go) con base de datos SQLite en memoria para tests (rápida y no requiere instalar nada extra).
- `httptest` (del estándar de Go) para levantar servidores HTTP de prueba y enviar requests.
- `stretchr/testify/require` para aserciones claras en los tests (comparar valores esperados vs. reales).
- Script de PowerShell para Windows que corre todas las pruebas, calcula el coverage total y genera el informe HTML.

¿Por qué estas herramientas?
- Son rápidas, confiables y no exigen infraestructura externa durante las pruebas.
- Permiten simular escenarios reales (HTTP, base de datos) pero de forma controlada.


## Cómo correr las pruebas y ver la cobertura

Tenés dos caminos: usar el script (recomendado) o correr comandos manuales.

### Opción A: usar el script (Windows PowerShell)

- Archivo: `ucc-soft-arch-golang/scripts/test_coverage.ps1`
- Qué hace:
  - Corre todas las pruebas del backend con cobertura.
  - Calcula el porcentaje total.
  - Genera un informe HTML fácil de leer.
- Cómo usarlo (PowerShell):
  1) Abrí PowerShell.
  2) Parate en la carpeta del backend `ucc-soft-arch-golang`.
  3) Ejecutá el script.

Comando (opcional para copiar y pegar):
```
powershell -NoProfile -ExecutionPolicy Bypass -File \
  "C:\Users\Juan\OneDrive\Escritorio\IS3\IS3\ucc-soft-arch-golang\scripts\test_coverage.ps1"
```

Al terminar, vas a ver en consola el porcentaje total (por ejemplo: `total: ... 80.1%`).

El informe HTML queda en:
- `ucc-soft-arch-golang/coverage.html`

Abrilo con doble clic para ver qué archivos/funciones están cubiertos.

### Opción B: comandos manuales (alternativa)

- Correr todas las pruebas con cobertura (salida solo en consola):
```
go test ./... -cover
```

- Generar un perfil de cobertura y abrir el informe HTML (opción más técnica):
```
go test ./... -coverprofile=coverage.out

go tool cover -html=coverage.out -o coverage.html
```


## Cómo están organizadas las pruebas (visión por carpetas)

Todo el código del backend está en `ucc-soft-arch-golang/src`. A continuación te contamos qué se prueba en cada "paquete" (carpeta) de forma simple.

> Nota: Los nombres están en inglés porque así se nombró el código, pero te explicamos en castellano qué hace cada uno.

### 1) adapter/
- Archivos: `*_adapter.go` y `adapters_test.go`.
- ¿Qué es? Puente/adapter entre la capa HTTP (controladores) y la lógica/servicios. Estándar de organización.
- ¿Qué se prueba?
  - Que los adapters existan y expongan lo necesario para que los controladores y rutas funcionen.
  - Que las integraciones "simples" (sin lógica pesada) estén bien conectadas.
- Cobertura alta (100%).

### 2) clients/
Módulos que hablan con la base de datos (vía GORM + SQL) para cada entidad. Usan SQLite en memoria durante pruebas. Contienen lógica de consulta (select), creación, actualización, borrado y helpers para normalizar datos (por ejemplo, conversiones de tipos).

- categories/
  - Pruebas de listados y filtrados, obtención por id, creación/actualización/borrado.
  - Validación de errores de base de datos.
  - Cobertura alta (~83%).

- comments/
  - Pruebas de crear/listar comentarios y traerlos por relaciones.
  - Cobertura alta (~83%).

- courses/
  - CRUD de cursos, listado con filtros (compatibles con SQLite), `GetById` con caso de "no encontrado" explícito.
  - Helpers: conversiones `toBool`, `toInt`, `toFloat64`, parseo de UUIDs (para robustez entre SQLite/Postgres).
  - Cobertura muy buena (~87.8%).

- inscriptos/ (inscriptions/enrollments)
  - Funciones: `Enroll` (inscribir), `GetMyCourses` (cursos de un usuario), `GetMyStudents` (alumnos de un curso, con manejo de resultado vacío), `IsUserEnrolled` y `CourseExist`.
  - Pruebas: caminos felices, casos "no encontrado", errores de base de datos y helpers (parseos/conversiones).
  - Cobertura muy alta (~94%).

- rating/
  - Lectura/escritura de calificaciones (ratings), y relación con cursos/usuarios.
  - Cobertura buena (~78.6%).

- users/
  - Funciones: `Create`, `FindById`, `FindByEmail`, `UpdateUser`.
  - Pruebas de éxito y de error (no encontrado, errores de DB, duplicados), y que los mapeos de error sean correctos.
  - Cobertura muy buena (~88.2%).

¿Por qué tantas pruebas aquí? Porque esta capa habla con la DB y suele tener muchos caminos alternativos (datos no encontrados, entradas inválidas, errores de conexión, etc.).

### 3) config/
- Archivos: `envs.go`, `db_connection.go`.
- ¿Qué se prueba?
  - `LoadEnvs` con distintos escenarios: sin archivo `.env` (toma variables del entorno), con archivo `.env` (carga valores), y error inesperado (debe hacer panic). 
- Cobertura media (~56.2%). `db_connection.go` se ejercita indirectamente en otras capas, no tanto con pruebas unitarias puras.

### 4) controllers/
- ¿Qué son? Los "controladores" HTTP: reciben la request, validan, llaman al servicio y devuelven una respuesta (JSON + código HTTP).
- ¿Cómo se prueban?
  - Levantamos un servidor HTTP de prueba (con `gin` en modo Test) y enviamos requests con `httptest`.
  - Verificamos códigos de estado (200, 400, 401, 404, 500, etc.) y contenido básico de la respuesta.
- Subcarpetas:
  - `auth/`, `categories/`, `comments/`, `courses/`, `inscriptions/`, `rating/`, `users/` y health (salud del servicio).
- Cobertura global del paquete `controllers`: 100% (las subcarpetas rondan entre ~77% y 100%, dependiendo del caso).

### 5) middleware/
- ¿Qué es? "Puentes" que se ejecutan antes de entrar al controlador. Por ejemplo: verificar permisos, validar un token, chequear si el usuario es admin, etc.
- ¿Qué se prueba?
  - `ErrorHandler`: traduce errores de la aplicación a respuestas HTTP consistentes.
  - Subpaquetes: `admin/` (rol admin), `course/` (permisos/propiedad sobre curso), `enroll/` (reglas de inscripción), `user/` (identidad/rol).
  - Se validan casos felices y casos de bloqueo (cuando no deja pasar).
- Cobertura muy alta (muchos al 100%, otros ~86–98%).

### 6) routes/
- ¿Qué son? Registro de rutas HTTP (qué URL llama a qué controlador y con qué middleware).
- ¿Qué se prueba?
  - Tests que levantan el router y verifican que las rutas están conectadas a los controladores correctos.
  - Incluye `router_test.go` y tests de rutas específicas (`auth`, `category`, `health`). Aunque no todas las rutas tienen un test dedicado, el router completo sí está cubierto.
- Cobertura del paquete: 100%.

### 7) services/
- ¿Qué son? La lógica de negocio. Reciben datos, aplican reglas/validaciones y llaman a los `clients` para leer/escribir en la DB.
- Paquetes: `auth`, `categories`, `comments`, `courses`, `inscription`, `rating`, `user`.
- ¿Qué se prueba?
  - Casos felices (creación/lectura/actualización/borrado según corresponda), validaciones, y parte de errores.
  - En `user` se valida hashing de contraseña y que no se persista texto plano.
- Cobertura actual: ~70% (todavía hay oportunidades de sumar más casos de error y bordes).

### 8) utils/
- `utils/bcrypt`: funciones para hashear (encriptar) contraseñas y verificar.
  - Pruebas: que el hash funcione y que la verificación detecte credenciales incorrectas.
- `utils/jwt`: utilidades para generar/verificar tokens JWT.
  - Pruebas: generación y validación de tokens, y manejo básico de errores.
- Cobertura: ~80–83%.

### 9) domain/
- `domain/errors`: definiciones de errores de dominio (tienen código, mensaje y estado HTTP). Cobertura 100%.
- `domain/dtos/*`: objetos de transferencia de datos (estructuras "planas"), sin lógica; por eso no tienen tests dedicados.

### 10) model/
- Modelos de base de datos (estructuras GORM): `users`, `courses`, `comments`, `category`, `rating`, `incriptos`.
- No tienen lógica propia (son definiciones de datos), por eso aparecen con 0% de cobertura. Es normal si no hay funciones. Si se quiere, se pueden agregar tests mínimos para subir ese número, pero no aportan valor funcional (no hay computación que verificar).


## Cosas importantes que aprendimos y ajustes hechos

- Normalizaciones entre SQLite/Postgres: Para que las pruebas sean estables, agregamos helpers (`toBool`, `toInt`, `toFloat64`, parseo de UUID) y usamos filtros compatibles (por ejemplo, `LOWER(...) LIKE` en lugar de `ILIKE`).
- Manejo explícito de "no encontrado": En consultas `Raw+Scan`, cuando no hay filas, devolvemos un error de NOT_FOUND claro en lugar de dejarlo ambiguo.
- Arreglos de bugs detectados por pruebas:
  - UserService.UpdateUser: se corrigió la persistencia del hash de contraseña (antes podía quedar mal guardada).
  - Rutas/filtros en `courses client`: se ajustaron para que funcionen igual en SQLite.
- Cobertura y calidad: el informe HTML nos ayudó a detectar carpetas con 0% o bajo porcentaje y priorizar dónde sumar tests.


## Cómo interpretar el informe HTML de cobertura

- Cada archivo aparece con líneas verdes (cubiertas por tests) y rojas (no cubiertas).
- Hay números por función y por archivo.
- Sirve para decidir "dónde falta probar".
- Archivo: `ucc-soft-arch-golang/coverage.html` (se abre en el navegador).


## Problemas comunes y soluciones rápidas

- "No tengo Go instalado": instalalo desde la página oficial de Go. Luego volvé a correr el script o los comandos.
- "El script de PowerShell no corre": abrí PowerShell como administrador y probá con `-ExecutionPolicy Bypass` (ya está incluido en el comando de ejemplo de arriba).
- "No se abre el HTML": verificá que el archivo `coverage.html` exista en la carpeta del backend y abrilo con el navegador.
- "Fallan algunos tests": mirá el mensaje de error en consola. Suele indicar qué esperaba el test y qué obtuvo.


## Preguntas frecuentes (FAQ)

- ¿Necesito una base de datos real? No. En tests usamos SQLite en memoria; todo vive solo mientras dura la prueba.
- ¿Se prueban rutas/HTTP "de verdad"? Sí, pero con un servidor de prueba en memoria (`httptest`) y Gin en modo Test. No se abren puertos reales.
- ¿Puedo correr un solo paquete? Sí, por ejemplo:
```
cd ucc-soft-arch-golang/src/clients/users

go test -v -cover
```
- ¿Es normal que `model/` tenga 0%? Sí, son solo estructuras de datos sin lógica.


## Resultados actuales (resumen)

- Cobertura total: 80.1%
- Paquetes más altos: `adapter` (100%), `controllers` (100%), `routes` (100%), `clients/inscriptos` (~94%).
- Paquetes a mejorar si queremos subir más:
  - `services` (~70%): sumar más caminos de error/borde.
  - `config` (~56%): tests adicionales para la conexión y errores específicos.
  - `clients/rating` (~78.6%): pequeños casos extra pueden pasarlo de 80%.
  - `model` (0%): opcional, tests mínimos si se quiere impactar la métrica.


## Próximos pasos sugeridos (opcionales)

- Agregar casos de error adicionales en `services/*` (por ejemplo, simulando fallas de la DB).
- Añadir tests mínimos a `model/*` (aunque sea solo para quitar el 0%).
- Extender pruebas de `controllers/inscriptions` para cubrir más ramas.


## Glosario rápido

- Test: pequeño programa que verifica que otra parte del código haga lo esperado.
- Cobertura (coverage): porcentaje del código que fue ejecutado por los tests.
- Happy path: camino feliz; caso normal cuando todo sale bien.
- Caso de error: cuando simulamos que algo falla (por ejemplo, la DB), para verificar que el sistema responde bien.
- Middleware: lógica que corre antes del controlador, para validar permisos o preparar datos.
- Controller (controlador): recibe la request HTTP y devuelve la respuesta.
- Service (servicio): lógica de negocio (las reglas, validaciones, etc.).
- Client: capa que habla con la base de datos.

---

Ante cualquier duda, podés abrir el archivo `coverage.html` para ver qué partes están (o no) cubiertas por los tests, y usar los comandos de este documento para ejecutar lo que necesites.
