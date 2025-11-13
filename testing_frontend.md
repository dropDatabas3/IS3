# Guía integral de pruebas del Frontend (Next.js + TypeScript + Jest)

Esta guía, en lenguaje simple, unifica y amplía toda la información sobre las pruebas del frontend. Explica qué se prueba, cómo se ejecuta, cómo leer el coverage, cómo funciona el simulador de API (MSW), y conceptos clave como inyección de dependencias y mocks, con ejemplos prácticos. Está pensada para que cualquier persona pueda leerla y entender cómo funcionan los tests del frontend.


## ¿Qué proyecto se testea?

- Proyecto: carpeta `ucc-arq-soft-front`
- Framework: Next.js 14 (React 18)
- Lenguaje: TypeScript
- Librerías clave de testing:
  - Jest (runner de tests, genera coverage)
  - Testing Library (pruebas de componentes de React de forma similar a como los usaría un usuario)
  - MSW (Mock Service Worker) para simular APIs en tests sin depender del backend real
  - ts-jest y babel-jest (transforman TypeScript/JavaScript para que Jest pueda ejecutarlo)


## ¿Qué son las pruebas unitarias y por qué importan?

- Son pequeños “controles automáticos” que verifican que partes del sistema se comporten correctamente.
- Detectan errores rápido cuando alguien cambia código, y ayudan a confiar en los cambios.
- Corren en segundos y no dependen de servicios externos (no llaman servidores ni bases reales).


## ¿Qué se testea? (mapa general)

Los tests viven dentro de `ucc-arq-soft-front/src/__tests__/` y cubren varias áreas:

- Pages (páginas de la app): `src/__tests__/pages/...`
  - Ejemplo: `app/courses.page.test.tsx`, `app/my-courses.page.test.tsx`, `app/auth/login/login/page.test.tsx`.
  - Qué validan: que las páginas rendericen contenido esperado, naveguen, y consuman datos simulados correctamente.

- Components (componentes UI): `src/__tests__/components/...`
  - Ejemplo: `components/courses/CourseCard.test.tsx`, `CourseDetail.test.tsx`, `NewCourseForm.test.tsx`, `CommentsList.test.tsx`, `ui/navbar/*.test.tsx`.
  - Qué validan: que los componentes muestren la información esperada, reaccionen a eventos (clicks, inputs, submit), y llamen a funciones del contexto cuando corresponde.

- Context/Providers (estado global): `src/__tests__/context/...`
  - Ejemplo: `CoursesProvider.test.tsx` (y flujos de error asociados).
  - Qué validan: que al montar los proveedores se carguen datos simulados (cursos, categorías, mis cursos), y que las acciones (inscribir, limpiar, crear, etc.) actualicen el estado.

- Reducers: `src/__tests__/reducers/...`
  - Qué validan: que las funciones reducer (puras) actualicen el estado correctamente para cada acción.

- Utils (utilidades): `src/__tests__/utils/...`
  - Ejemplo: `api.test.ts`, `courseMapper.test.ts`, `commentMapper.test.ts`, `uploadImage.test.ts`.
  - Qué validan: que las utilidades transformen datos, construyan URLs o ejecuten llamadas de forma esperada (siempre simulada/mokeada).

- Hooks: `src/__tests__/hoocks/...` y `src/__tests__/hoocks/*.test.tsx`
  - Qué validan: comportamiento de hooks personalizados (por ejemplo, la lógica del navbar).

Importante: el directorio `src/types/**` (solo definiciones de tipos TypeScript) se excluye del coverage porque no es código que “se ejecute” en runtime.


## Herramientas y configuración clave

- Jest: ejecuta tests y calcula cobertura.
- @testing-library/react + jest-dom: renderiza componentes y hace aserciones “como usuario”.
- MSW (msw/node): intercepta `fetch` y responde con datos simulados.
- ts-jest y babel-jest: permiten a Jest entender TypeScript y módulos ESM.
- jsdom: simula un navegador en Node para que React renderice en pruebas.
- Mocks puntuales para piezas de Next y otras librerías:
  - `__mocks__/next/navigation.js` (router simulado), `__mocks__/next/image.js` (imagen simple `<img>`), `__mocks__/swiper-*.{ts,tsx}` (slider simplificado), además de `styleMock.js` y `fileMock.js` para estilos/estáticos.

Archivos de configuración (en `ucc-arq-soft-front`):
- `jest.config.ts`: entorno jsdom, transforms, moduleNameMapper (alias `@/` → `src/` y mocks), cobertura y thresholds.
- `jest.setup.ts`: activa `jest-dom`, agrega polyfills (TextEncoder/Decoder, streams, requestSubmit, BroadcastChannel), y arranca/para MSW antes/después de los tests.


## ¿Cómo se testea? (patrones y utilidades)

- Renderizado con proveedores: la utilidad `src/test/test-utils.tsx` expone un `render(...)` que envuelve todo con los providers reales de la app:
  - `UiProvider` → `AuthProvider` → `CoursesProvider`
  - Esto permite que los componentes/páginas, al montarse en tests, funcionen igual que en la app real (con estado global y efectos).

- Mock de APIs con MSW: en `src/test/server.ts` se levanta un servidor MSW de Node que intercepta `fetch` y responde con handlers:
  - Handlers en `src/test/handlers/*.handlers.ts` (auth, courses, ratings, comments)
  - Simulan endpoints como `GET /courses`, `POST /courses/create`, `GET /myCourses/`, etc.
  - Ventaja: los tests no llaman al backend real. Son rápidos, predecibles y sin flakiness por red.

- Setup global de tests: `jest.setup.ts` hace varias cosas útiles:
  - Importa `@testing-library/jest-dom` para aserciones adicionales (ej. `toBeInTheDocument`).
  - Polyfills para Node/jsdom (TextEncoder/Decoder, requestSubmit en formularios, streams, BroadcastChannel) que ciertas librerías necesitan.
  - Inicia y apaga el server de MSW antes/después de los tests.

- Configuración de Jest: `jest.config.ts`
  - Ambiente: `jsdom` (simula un navegador en Node)
  - Transforms: `ts-jest` y `babel-jest` para TypeScript/ESM
  - `moduleNameMapper` para mocks (CSS, imágenes, Next, Swiper) y alias `@/` → `src/`
  - Coverage con umbrales (thresholds) globales exigentes:
    - statements: 80%
    - lines: 80%
    - functions: 75%
    - branches: 60%
  - Coverage “whitelist” (solo se miden estos directorios): `src/app`, `src/components`, `src/context`, `src/utils`, `src/providers` y subcarpetas `hooks/`
  - Exclusiones importantes: tests, d.ts, mocks y cualquier carpeta llamada `types`.


## ¿Cómo correr los tests?

1) Abrí una terminal en la carpeta del frontend:

```powershell
cd "c:\Users\Juan\OneDrive\Escritorio\IS3\IS3\ucc-arq-soft-front"
```

2) Instalá dependencias (primera vez):

```powershell
npm install
```

3) Ejecutá todos los tests:

```powershell
npm test
```

4) Ejecutá con cobertura (HTML + texto):

```powershell
npm run test:coverage
```

- “Test Suites” indica cuántos archivos de test corrieron.
- “Tests” indica cuántos casos individuales se ejecutaron.
- Si algún umbral de cobertura no se cumple, Jest falla (salida con error) para proteger la calidad.

Tip: si cambiás configuración y querés limpiar la caché de Jest:

```powershell
npx jest --clearCache
```


## Estructura de carpetas (resumen)

- `src/__tests__/`: archivos de prueba (components, hooks, pages, reducers, utils, etc.).
- `src/test/`: utilidades de pruebas y “handlers” de MSW que simulan la API.
  - `src/test/server.ts`: junta y arranca todos los handlers.
  - `src/test/handlers/*.handlers.ts`: respuestas simuladas para rutas (auth, courses, ratings, comments).
- `__mocks__/`: reemplazos de librerías/archivos para que Jest no falle con imágenes o CSS.
  - `__mocks__/next/navigation.js`, `__mocks__/next/image.js`
  - `__mocks__/swiper-react.tsx`, `__mocks__/swiper-modules.ts`
  - `__mocks__/styleMock.js`, `__mocks__/fileMock.js`

Nota: existe una carpeta `test/` en la raíz del frontend que contiene archivos vacíos y no se usa; la suite real utiliza `src/test/*`. Puede borrarse para evitar confusiones (no afecta a los tests).


## Cómo funciona MSW (en 2 pasos)

1) Definimos “handlers” con rutas y respuestas de ejemplo, por ejemplo en `src/test/handlers/auth.handlers.ts`:
   - Si falta `email` o `password` → devolver `400` (error).
   - Si todo está OK → devolver un JSON con `{ token, user }`.

2) `src/test/server.ts` junta todos los handlers y los registra. `jest.setup.ts` lo enciende antes de los tests y lo apaga al final.

Así, las pruebas no necesitan internet y siempre reciben respuestas previsibles.


## Cómo están escritos los tests (patrón AAA) con ejemplo

- Arrange (Preparar): render del componente y/o datos de entrada.
- Act (Actuar): interacción de usuario (clicks, tipeo, submit).
- Assert (Afirmar): verificar lo esperado en pantalla/estado.

Ejemplo mínimo:

```ts
import { render, screen } from '@/test/test-utils';
import userEvent from '@testing-library/user-event';
import { LoginForm } from '@/components/auth/LoginForm';

it('muestra bienvenida tras login', async () => {
  render(<LoginForm />);
  await userEvent.type(screen.getByLabelText(/email/i), 'ana@acme.dev');
  await userEvent.type(screen.getByLabelText(/contraseña/i), 'Secreta#123');
  await userEvent.click(screen.getByRole('button', { name: /ingresar/i }));
  expect(await screen.findByText(/bienvenida/i)).toBeInTheDocument();
});
```


## ¿Cómo ver el coverage (HTML) y leerlo?

Jest genera un reporte HTML en: `ucc-arq-soft-front/coverage/lcov-report/index.html`.

- Abrilo en VS Code (preview) o en tu navegador.
- Verás una tabla con carpetas y archivos, y varias columnas:
  - Statements (instrucciones): porcentaje de líneas ejecutadas vs. totales de “statement”. Similar a “lines”, pero a nivel instrucción.
  - Branches (ramas): cubre if/else, operadores ternarios, switch, etc. Es normal que sea la métrica más baja.
  - Functions (funciones): cuántas funciones se ejecutaron al menos una vez.
  - Lines (líneas): porcentaje de líneas de código ejecutadas (contabilización línea a línea).
  - “%” (porcentaje) y “Covered/Total” (cubiertas/totales) aparecen por cada métrica.
  - Colores: rojo (bajo), amarillo (medio), verde (alto). Sirven para detectar rápido qué falta cubrir.

Entrando a un archivo concreto en el reporte:
- Verás el código coloreado:
  - Verde: ejecutado en tests.
  - Rojo: nunca se ejecutó en tests.
  - Amarillo: a veces indica ramas no cubiertas (por ejemplo, un else/ternario que no se ejercitó).
- En la parte superior, suele aparecer un resumen de cuántas líneas/funciones/ramas están cubiertas.
- Podés usar esta vista para decidir qué escenarios de test faltan (ej. probar el flujo de error de una función que hoy solo probaste con respuesta exitosa).

Notas importantes del coverage en este proyecto:
- El directorio `src/types/**` no se incluye en cobertura (son solo tipos de TypeScript, no “código ejecutable”). Por eso no afectan el promedio ni deberían aparecer en la tabla.
- Los umbrales globales (80/80/75/60) están configurados. Si bajás cobertura al editar código, el comando de tests fallará y te lo va a marcar en CI.


## Conceptos clave: inyección de dependencias y mocks (con ejemplos)

En frontend, “inyectamos dependencias” de forma simple: en lugar de que un componente cree/posea todo, le “pasamos” lo que necesita desde afuera para poder sustituirlo en pruebas.

- Proveedores (Providers) como inyección: usamos `render` de `src/test/test-utils.tsx` que envuelve con `UiProvider → AuthProvider → CoursesProvider`. Así, los componentes “reciben” el contexto necesario sin tocar producción. Podemos cambiar el estado inicial o las respuestas de MSW para probar distintos escenarios.

- Mocks (dobles de prueba): reemplazan dependencias reales por versiones controladas en tests.
  - Mocks de módulos (ejemplo):

```ts
// Simular el router de Next.js en un test
jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: jest.fn(), replace: jest.fn(), prefetch: jest.fn() }),
  usePathname: () => '/login'
}));
```

  - Mocks de red con MSW (handlers):

```ts
// src/test/handlers/auth.handlers.ts (idea simplificada)
import { http, HttpResponse } from 'msw';

export const authHandlers = [
  http.post('/auth/login', async ({ request }) => {
    const body = await request.json();
    if (!body.email || !body.password) {
      return HttpResponse.json({ message: 'Faltan datos' }, { status: 400 });
    }
    return HttpResponse.json({ token: 'token-demo', user: { name: 'Ana' } });
  })
];
```

- Stubs y Spies (breve):
  - Stub: función “falsa” que devuelve algo fijo (por ejemplo, `jest.fn(() => 42)`).
  - Spy: función espía que registra llamadas/argumentos (`jest.fn()`), útil para verificar “se llamó con…”.

Beneficio: al inyectar y/o mockear, probamos el componente en aislamiento, controlando los escenarios (éxito, error, vacíos, etc.) sin depender de servicios reales.


## ¿Cómo interpretar la salida de Jest en consola?

Cuando corrés `npm test` verás logs como:

- “Test Suites: 28 passed, 28 total” → 28 archivos de test ejecutados, todos pasaron.
- “Tests: 62 passed, 62 total” → 62 casos individuales pasaron.
- Warning “act(...)” en consola → te avisa que hubo actualizaciones de estado asíncronas no envueltas explícitamente en `act`. En nuestro caso vienen de efectos dentro de providers y son benignas (los tests ya esperan correctamente los resultados). Si querés, se pueden silenciar afinando los `await`/`waitFor` o ajustando el helper de render.


## Herramientas usadas (resumen)
 (ya detalladas arriba)


## Cómo agrego un test nuevo (paso a paso)

1) Elegí qué querés probar (un componente, una página, una función de utilidades). 
2) Creá un archivo `*.test.tsx` o `*.test.ts` dentro de `src/__tests__/...` siguiendo la estructura existente.
3) Si tu caso usa red (fetch), agregá/ajustá un handler en `src/test/handlers/*.handlers.ts` para simular la respuesta.
4) Importá `render` desde `src/test/test-utils.tsx` si es un componente/página que usa providers.
5) Escribí el test: renderizá, dispará eventos (ej. `fireEvent.click`) y hacé aserciones (`expect(...).toBeInTheDocument()`).
6) Ejecutá `npm test -- --coverage` y revisá el `index.html` para ver si aumentó cobertura y si quedó alguna línea/ramas sin cubrir.


## Qué casos cubren los tests actuales (visión de alto nivel)

- Páginas principales (home, courses, course-info, my-courses, login): renderizan, consumen datos de MSW, muestran estados y listas.
- Componentes de cursos (card, carrusel, detalle, lista, formulario de creación): entradas de usuario, validaciones, submits y llamadas a contexto.
- Autenticación (Login/Register): formularios, validación básica y flujos de login/register contra MSW.
- Contexto de cursos (provider + reducer): carga inicial (cursos, categorías, mis cursos), acciones (inscripción, limpiar lista, crear/editar/eliminar cursos) y manejo de errores.
- Utilidades: construcción de URLs de API, mapeo/transformación de datos, subida de imágenes (simulada), etc.
- Hooks y Navbar: decisiones de UI según estado de autenticación/rol, visibilidad de botones/links, etc.


## Cómo ejecutar solo una parte de los tests (opcional)

- Por nombre de archivo:

```powershell
npm test -- src/__tests__/components/courses/CourseCard.test.tsx --colors
```

- Con “watch mode” (para desarrollar tests en caliente):

```powershell
npm test -- --watch
```


## Integración continua (CI) y umbrales

- Los tests del frontend están listos para ejecutarse en CI.
- Si el coverage global cae por debajo de los umbrales (80/80/75/60), los tests fallan (salida distinta de cero) y el pipeline lo va a marcar.
- Esto ayuda a mantener la calidad con el tiempo.


## Preguntas frecuentes

- ¿Por qué aparecen warnings de `act(...)`? 
  - Porque algunos efectos asíncronos actualizan estado al montar providers. Son avisos, no fallos. Se pueden reducir agregando `await`/`waitFor` en los tests donde corresponda.

- ¿Por qué no contamos `src/types/**` en cobertura?
  - Porque son definiciones de tipos de TypeScript (no se ejecutan). Incluirlos baja artificialmente el porcentaje sin aportar señal de calidad.

- ¿Dónde veo el reporte de cobertura?
  - `ucc-arq-soft-front/coverage/lcov-report/index.html`. Abrilo después de correr `npm test -- --coverage`.

- ¿Qué significan “statements”, “branches”, “functions” y “lines”?
  - Statements: instrucciones ejecutadas (piezas de código). 
  - Branches: caminos lógicos (if/else, ternarios, switch).
  - Functions: funciones invocadas.
  - Lines: líneas de código ejecutadas.


## Conclusión

Con Jest + Testing Library + MSW, podés testear el frontend de forma rápida, confiable y cercana al uso real del usuario. Con el reporte HTML de cobertura vas a ver claramente qué falta cubrir. Y con umbrales en CI, te asegurás de no perder calidad al crecer.

— Esta es la única guía de pruebas del frontend. Se consolidó el contenido previo en un solo archivo para evitar duplicados.
