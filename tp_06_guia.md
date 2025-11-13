Optimizing tool selection...

Created 6 todos

Voy a trazar una guía paso a paso, aterrizada a tu repo (Next.js + TypeScript en front y Go en back), para montar una suite de pruebas unitarias completa con mocks, DI, “DB” en memoria y una estructura clara para que puedas defender el TP6 con fundamentos.

## Plan de alto nivel

- Frontend: configurar Jest + React Testing Library + MSW, armar utilidades de test, inyectar dependencias para desacoplar componentes y hooks, y proponer casos de prueba.
- Backend (Go): usar testing nativo + testify/mock, definir interfaces y DI en servicios/controladores, crear repos/clients in-memory y mocks para dependencias externas, y proponer casos de prueba.
- Carpetas y archivos: detallar exactamente qué agregar y dónde en tu repo.
- Comandos y dependencias: indicar instalaciones y scripts.
- Documentación: qué anotar en README y decisiones.md.
- Pipeline: un bosquejo final (lo dejamos para el cierre).

Avanzo con la arquitectura de testing y el cómo/qué/dónde por cada capa.

---

## Frontend: Next.js + TypeScript

### Objetivos de las pruebas
- Componentes y UI: render, interacción (inputs, clicks), accesibilidad básica.
- Hooks: estados y efectos, sin pegarle a la red.
- Utilidades/mappers: funciones puras con casos edge.
- Integración ligera: componentes que consumen contexto o servicios, con dependencias inyectadas y/o MSW.

### Dependencias a instalar

En ucc-arq-soft-front:

```powershell
# Jest + TS + entorno DOM
npm i -D jest @types/jest ts-jest jest-environment-jsdom

# React Testing Library (RTL)
npm i -D @testing-library/react @testing-library/jest-dom @testing-library/user-event

# MSW (Mock Service Worker) para mockear fetch/HTTP en tests
npm i -D msw

# Utilidades recomendadas
npm i -D whatwg-fetch                 # polyfill fetch en Node si hace falta
npm i -D next-router-mock @types/next-router-mock  # para mockear useRouter en pruebas de páginas
```

Notas:
- Si tu Next usa SWC (default), `ts-jest` funciona bien para TS puro. Alternativa: `babel-jest` + preset de Next, pero con TS-Jest es más directo para tu caso.
- `whatwg-fetch` asegura fetch global en Node 18/20 si alguna lib lo requiere.

### Estructura de archivos de testing (propuesta)

En ucc-arq-soft-front:

```
jest.config.ts                      # Configuración Jest TS + jsdom + mapeos
jest.setup.ts                       # Extensiones jest-dom, setup MSW si aplica
src/test/
  test-utils.tsx                    # render personalizado con Providers
  server.ts                         # instancia MSW (node)
  handlers/
    auth.handlers.ts                # handlers de rutas /auth/*
    courses.handlers.ts             # handlers de rutas /courses/*
__mocks__/
  fileMock.js
  styleMock.js
src/__tests__/
  components/
    auth/
      LoginForm.test.tsx
      RegisterForm.test.tsx
    courses/
      CourseCard.test.tsx
      CoursesList.test.tsx
  hooks/
    useNavbar.test.ts
  utils/
    courseMapper.test.ts
    commentMapper.test.ts
  pages/
    app/
      auth/
        login/page.test.tsx
```

Qué va en cada uno:
- `jest.config.ts`: entorno jsdom, transform TS via ts-jest, `setupFilesAfterEnv: ['<rootDir>/jest.setup.ts']`, `moduleNameMapper` para CSS e imágenes.
- `jest.setup.ts`: `import '@testing-library/jest-dom';` y levantar MSW en modo tests (server.listen/close/reset).
- `src/test/test-utils.tsx`: un wrapper `render` que monta `AuthProvider`, `CoursesProvider`, `UiProvider` con dependencias inyectables para reemplazar el `apiClient` real por uno fake/mocked.
- `src/test/server.ts` y `src/test/handlers/*`: definición de handlers MSW para simular respuestas HTTP en tests de integración livianos.
- `__mocks__/styleMock.js` y `__mocks__/fileMock.js`: para que Jest no intente interpretar CSS/archivos estáticos.

### Configuración Jest (qué y por qué)
- `testEnvironment: 'jsdom'`: simula browser.
- `transform: { '^.+\\.tsx?$': 'ts-jest' }`: compila TS en tests.
- `setupFilesAfterEnv`: jest-dom y arranque MSW.
- `moduleNameMapper`: mock de CSS/imagenes y alias si los usás.
- `testMatch`: por defecto `**/*.test.ts(x)`.

Ejemplo mínimo de `jest.config.ts`:
- Define environment jsdom.
- Usa ts-jest con `tsconfig` de proyecto.
- Mapea CSS/imágenes a mocks.
- Incluye `jest.setup.ts`.

### Inyección de dependencias (DI) en el Front

Objetivo: que componentes y hooks no dependan directamente de “infra” (fetch/axios, ventanas globales), y sean testeables en aislamiento.

Tu código actual usa `src/utils/api.ts`. Propongo:
- Definir una interfaz de cliente HTTP y un factory:
  - `ApiClient` con métodos `get/post/put/delete<T>(url, options)`.
  - `createApiClient(baseUrl?: string): ApiClient` que usa fetch real.
- Modificar Providers (AuthProvider, CoursesProvider) para aceptar `apiClient?: ApiClient` como prop opcional y usar el real por defecto:
  - En tests: proveer un `apiClient` falso (jest.fn()) o usar MSW para interceptar fetch.

Dónde:
- `src/utils/api.ts`: exportar `export interface ApiClient { ... }` y `export function createApiClient(...)`.
- `src/context/auth/AuthProvider.tsx`, `src/context/courses/CoursesProvider.tsx` y cualquier Provider que llame API: aceptar `apiClient?: ApiClient`.
- `src/test/test-utils.tsx`: `render(ui, { apiClientOverrides })` crea providers con el `apiClient` mockeado.

¿Por qué?
- Aísla UI de HTTP real.
- Permite tests AAAs claros: Arrange (inyecto apiClient), Act (render/click), Assert (resultado y que apiClient fue llamado con X).

### Testing con React Testing Library (patrón AAA)

Ejemplo (LoginForm):
- Arrange: render con `AuthProvider` y `apiClient` mock que responde `login` OK.
- Act: completar email/password y submit con `user-event`.
- Assert: se llamó al endpoint correcto; UI muestra “Bienvenido” o navega.

Edge cases:
- Validaciones vacías/mal formato.
- Error 401: mensaje de error.
- Loading disabled en botón mientras envía.

### MSW vs jest.fn

Cuándo usar cada uno:
- Unit puro de un hook/función → jest.fn para la dependencia.
- Componente que llama “red” vía contexto/servicio → MSW para simular rutas o inyección de `apiClient` mock si querés aislamiento total.
- Páginas Next con router → `next-router-mock` para `useRouter`.

### Hooks y utilidades

- `useNavbar.ts`: probar toggles/estados con `renderHook` de @testing-library/react (opcional) o con un componente de prueba.
- `mappers` (p. ej. `courseMapper.ts`): tests puros con entradas edge (valores nulos, arrays vacíos, campos faltantes).

### Scripts NPM

En package.json agregar:
- `"test": "jest --passWithNoTests"`
- `"test:watch": "jest --watch"`
- `"test:coverage": "jest --coverage"`

### Comandos locales

```powershell
cd ucc-arq-soft-front
npm install
npm run test
npm run test:watch
npm run test:coverage
```

---

## Backend: Go (testing nativo + testify/mock)

### Objetivos de las pruebas
- Servicios (lógica de negocio): unit puros con dependencias mockeadas.
- Controladores: probar handlers con `httptest` y servicios mock.
- Middleware: caminos success/failure (autorización, validaciones).
- Utilidades (bcrypt/jwt): tests puros, y DI si es necesario para controlar JWT.

Tu estructura relevante:
- `src/services/` (auth, categories, comments, courses, inscription, rating, user)
- `src/adapter/` y `src/clients/` (acceso externo)
- `src/controllers/` y `src/middleware/`
- `src/utils/jwt` y `src/utils/bcrypt`

### Dependencias a instalar

En ucc-soft-arch-golang:

```powershell
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
go get github.com/stretchr/testify/mock
```

### Patrones de DI y contratos

1) Define interfaces en el límite de tu dominio para cada dependencia externa:
- Si `src/services/courses_services.go` depende de un cliente HTTP, crea una interfaz `CoursesClient` con métodos usados (GetAll, GetById, Create, etc.).
- Si depende de “repositorio” (aunque hoy llame servicios), crear `CoursesRepository` interface (si preferís mental model repos).
- Sitúa las interfaces en el paquete del consumidor (services) o en `src/domain/dtos/...` si preferís evitar circularidad.

2) Crea constructores en servicios:
- `func NewCoursesService(client CoursesClient) *CoursesService { ... }`
- Servicios sólo conocen interfaces, no implementaciones concretas.

3) Crea mocks para tests con testify/mock:
- `src/mocks/courses_client_mock.go` que embeba `mock.Mock` y cumpla `CoursesClient`.

4) “DB en memoria”/Clientes en memoria (para tests de integración de servicio sin mocking):
- `src/repositories/memory/courses_memory.go` (map + mutex) o `src/clients/courses/courses_client_inmemory.go`.
- Implementan las mismas interfaces; retornan datos pre-cargados.

¿Por qué?
- Aísla y acelera tests.
- Evita hitting reales (microservicios/DB).
- Permite reproducir escenarios edge controlados (errores, latencias, vacíos).

### Estructura de archivos de testing (propuesta)

En ucc-soft-arch-golang:

```
src/mocks/
  auth_client_mock.go
  users_client_mock.go
  courses_client_mock.go
  categories_client_mock.go
  comments_client_mock.go
  rating_client_mock.go
  inscriptions_client_mock.go
src/repositories/memory/                  # o clients/*_inmemory.go si preferís
  courses_memory.go
  users_memory.go
  categories_memory.go
src/services/
  auth_service.go
  auth_service_test.go
  courses_services.go
  courses_services_test.go
  ...                                     # tests por cada servicio
src/controllers/
  courses/
    courses.controller.go
    courses_controller_test.go
  auth/
    auth_controller.go
    auth_controller_test.go
  ...
src/middleware/
  user/
    userAuth.mid.go
    userAuth.mid_test.go
  admin/
    adminAuth.mid.go
    adminAuth.mid_test.go
src/utils/jwt/
  sign_document.go
  verify_token.go
  verify_token_test.go
src/utils/bcrypt/
  hash_password.go
  compare_password.go
  bcrypt_test.go
testdata/                                  # payloads JSON de ejemplo si hace falta
```

Qué va en cada uno:
- `*_service_test.go`: unit tests de reglas de negocio con mocks (Arrange mocks → Act → Assert).
- `*_controller_test.go`: `httptest.NewRecorder`, `http.NewRequest`, router/handler y assert de status/body; el servicio es un mock inyectado.
- `*_mid_test.go`: router con middleware y handler dummy para verificar que pasa/bloquea.
- `*memory.go`: implementaciones in-memory para pruebas con más superficie sin usar mocks.
- `testdata/`: archivos JSON para inputs/outputs reproducibles.

### Ejemplo de casos a cubrir

Servicios:
- Happy path: devuelve cursos filtrados/paginados correctamente.
- Validaciones: inputs inválidos devuelven error específico.
- Errores externos: client retorna error → servicio lo mapea a error de dominio.

Controladores:
- 200 con payload esperado.
- 400/422 por inputs inválidos.
- 500 si servicio retorna error inesperado.
- Content-Type application/json y shape esperado.

Middleware:
- Token válido → deja pasar.
- Token inválido/expirado → 401.
- Sin token → 401.
- Roles insuficientes → 403 (para admin).

Utils:
- JWT: firma/verifica claims; mal secreto → falla.
- Bcrypt: hash/compare correcto; compare con password incorrecto → false.

### Cómo inyectar dependencias en Go

- Servicios reciben interfaces en constructores. Ej:
  - `type CoursesService struct { client CoursesClient }`
  - `func NewCoursesService(c CoursesClient) *CoursesService { return &CoursesService{client: c} }`
- Controladores reciben el servicio como dependencia (no crean el servicio dentro). Ej:
  - `type CoursesController struct { svc *CoursesService }`
  - `func NewCoursesController(svc *CoursesService) *CoursesController { ... }`
- Middlewares que lean JWT:
  - Extraer verificación a una interfaz `TokenVerifier` (concrete usa `utils/jwt.Verify`), de modo que en tests inyectes un stub que controla el resultado.

### Mocks (testify/mock) vs in-memory

- Unit puro → mock (control total sobre llamadas y retornos, asserts de expectativas).
- Tests “semi-integración” de servicio → in-memory client/repo (sin red).
- Controlador/middleware → preferible mock del servicio/verificador para aislar HTTP layer de lógica.

### Comandos locales

```powershell
cd ucc-soft-arch-golang
go test ./... -count=1
go test ./... -cover
```

---

## Qué, dónde y por qué (resumen con carpeta a carpeta)

### Frontend

- `jest.config.ts`: configura Jest con ts-jest, jsdom, mappers de CSS/archivos.
- `jest.setup.ts`: añade `@testing-library/jest-dom` y arranca MSW en test.
- `src/test/test-utils.tsx`: función `render(ui, { apiClient })` que envuelve con `AuthProvider`, `CoursesProvider`, `UiProvider`, inyectando `apiClient` mock si se pasa.
- `src/test/server.ts` + `src/test/handlers/*.ts`: MSW en Node para simular endpoints `/auth/*`, `/courses/*`, etc.
- `__mocks__/styleMock.js` y `fileMock.js`: mocks para imports de CSS/archivos.
- `src/__tests__/...`: tests por componente/hook/utils/página.
- `src/utils/api.ts`: define `ApiClient` + `createApiClient`.
- `src/context/*Provider.tsx`: aceptar `apiClient?: ApiClient` y usarlo.

### Backend

- `src/services/*_test.go`: unit tests con `src/mocks/*_client_mock.go`.
- `src/controllers/*/*_test.go`: `httptest` con servicios mock.
- `src/middleware/*/*_test.go`: verificación de permisos/autenticación con stubs.
- `src/repositories/memory/*_memory.go` o `src/clients/*_inmemory.go`: implementaciones en memoria para pruebas con más amplitud.
- `src/utils/jwt/*_test.go` y `src/utils/bcrypt/*_test.go`: tests de utilitarios.

---

## Patrón AAA y ejemplos breves

Frontend (component):
- Arrange: `const api = { post: jest.fn().mockResolvedValue({...}) }` y `render(<LoginForm/>, { apiClient: api })`.
- Act: completar inputs + `user.click`.
- Assert: `expect(api.post).toHaveBeenCalledWith('/auth/login', ...)` y UI muestra estado esperado.

Backend (service):
- Arrange: `client := new(mocks.CoursesClientMock); client.On("GetAll", ...).Return(expected, nil); svc := NewCoursesService(client)`
- Act: `result, err := svc.GetAll(ctx, filter)`
- Assert: `require.NoError(t, err); assert.Equal(t, expected, result); client.AssertExpectations(t)`

Backend (controller):
- Arrange: servicio mock retorna X; inyectar en controller y montar ruta con handler.
- Act: `req := httptest.NewRequest("GET", "/courses?category=...", nil)`; `rec := httptest.NewRecorder()`; `router.ServeHTTP(rec, req)`
- Assert: `assert.Equal(t, 200, rec.Code); assert.JSONEq(t, expectedJSON, rec.Body.String())`

---

## Casos edge y manejo de excepciones

Frontend:
- Inputs vacíos, formatos inválidos, submit múltiples, errores de red (500), timeouts (simulados con MSW y `ctx.delay`).
- Estados de loading/deshabilitado y retry.

Backend:
- IDs inválidos, recursos no encontrados (404 desde client → error de dominio), conflictos (409), validaciones (400/422).
- JWT expirado, firma incorrecta, roles insuficientes.

---

## Métricas y coverage

- Frontend: `npm run test:coverage` y configura `collectCoverageFrom` en `jest.config.ts` para cubrir `src/**/*.{ts,tsx}` excluyendo `__tests__` y `*.d.ts`.
- Backend: `go test ./... -cover -coverprofile=coverage.out` y `go tool cover -html=coverage.out -o coverage.html` si querés reporte HTML.

---

## Documentación que hay que sumar

En README.md (raíz o por proyecto):
- Requisitos (Node/Go).
- Instalación de dependencias.
- Cómo ejecutar tests (front/back) y ver coverage.
- Problemas frecuentes (por ej., MSW y Node ESM/TS-Jest configs).

En decisiones.md:
- Frameworks elegidos y por qué:
  - Front: Jest + RTL (+ MSW) por estándar en React/Next y foco en comportamiento del usuario.
  - Back: testing nativo + testify/mock por ergonomía y asserts/mocks expresivos.
- Estrategia de mocking:
  - Front: DI de `apiClient` y MSW para rutas.
  - Back: interfaces en servicios, mocks con testify, repos/clients en memoria.
- Patrón AAA con ejemplos.
- Principales casos probados (criticals de negocio).
- Evidencias: capturas de salida de test/coverage; breve explicación de resultados.

---

## Pipeline (bosquejo para después)

No prioritario ahora, pero para cerrar:
- Front: job que haga `npm ci` + `npm test -- --ci`.
- Back: `go test ./... -count=1 -cover`.
- Publicar reports de coverage (opcional).
- Gates: fallar pipeline si tests fallan.

---

## Cómo empezar en este repo (sugerido por pasos)

1) Frontend
- Agregar `jest.config.ts`, `jest.setup.ts`, `__mocks__`, `src/test/*`.
- Ajustar `src/utils/api.ts` para exponer `ApiClient` y `createApiClient`.
- Modificar Providers para aceptar `apiClient`.
- Escribir 3-5 tests iniciales: `LoginForm`, `CourseCard`, `courseMapper`, `useNavbar`.
- Ejecutar y estabilizar suite.

2) Backend
- Extraer interfaces para clientes en cada servicio (si no están).
- Agregar mocks en `src/mocks`.
- Escribir tests de 1 servicio (p. ej. `courses_services_test.go`) y 1 controlador y 1 middleware.
- Agregar in-memory impl para un caso (p. ej., courses) para pruebas más amplias sin mocks.
- Ampliar al resto de servicios/controladores.

3) Documentación (README y decisiones.md), luego pipeline.

---

## Comandos de referencia

Frontend:
```powershell
cd ucc-arq-soft-front
npm i
npm run test
npm run test:coverage
```

Backend:
```powershell
cd ucc-soft-arch-golang
go test ./... -count=1
go test ./... -cover
```

---

## Tips para la defensa oral

- Por qué estas herramientas: comunidad, integración con Next/React y Go, velocidad, ergonomía.
- Criterio de qué mockear: todo lo que es I/O (HTTP/DB/JWT verificador) se mockea o se reemplaza por in-memory; lo propio (lógica) se prueba real.
- Cómo validar que los tests prueban lo correcto: asserts de estado/efectos visibles, asserts de llamadas (con inputs correctos) y tests negativos (errores).
- Manejo de estado/datos externos: DI + MSW + repos en memoria → reproducible y determinista.

---

Si querés, en el próximo paso creo los archivos base de configuración (jest.config.ts, jest.setup.ts, test-utils, server MSW, mocks testify en Go) directamente en el repo para que puedas empezar a escribir pruebas encima. ¿Avanzo con eso?