⚠️ Documento consolidado

Este archivo fue reemplazado por la guía única y completa en la raíz del repositorio:

- `testing_frontend.md`

Por favor, consulta y actualiza solo ese archivo para evitar duplicaciones.


## Qué son las pruebas unitarias y por qué importan

- Las pruebas unitarias son pequeños “controles automáticos” que verifican que partes del sistema se comporten correctamente.
- Sirven para detectar errores rápido cuando alguien cambia código.
- Se ejecutan en segundos y no dependen de servicios externos (no llaman a servidores reales ni a bases de datos reales).


## Herramientas que usamos

- Jest: el motor que ejecuta los tests en Node (sin navegador real).
- React Testing Library (RTL): nos ayuda a probar componentes de interfaz simulando lo que hace un usuario (escribir en un campo, hacer clic, etc.).
- MSW (Mock Service Worker): simula las respuestas de la “API” para no depender de ningún servidor real; definimos respuestas de ejemplo.
- Mocks puntuales: reemplazamos piezas difíciles en pruebas, por ejemplo:
  - next/navigation (router de Next.js)
  - next/image (componente de imágenes de Next.js)
  - Swiper (slider) y sus CSS

Con esto, las pruebas son rápidas, repetibles y no requieren internet ni levantar servicios.


## Cómo se ejecutan las pruebas

1) Instalar dependencias (una sola vez):

```powershell
cd ucc-arq-soft-front
npm install
```

2) Ejecutar pruebas:

```powershell
npm test
```

3) Ejecutar con reporte de cobertura (opcional):

```powershell
npm run test:coverage
```

El reporte de cobertura indica qué porcentaje del código fue ejercitado por los tests. No es requisito del TP tener un mínimo, pero sirve como referencia.


## Dónde está la configuración

- `jest.config.ts`: configuración principal de Jest (entorno de pruebas, transformaciones de TypeScript, mapeos de módulos y estilos, etc.).
- `jest.setup.ts`:
  - Activa extensiones de aserciones (jest-dom).
  - Carga polyfills necesarios para Node.
  - Arranca y detiene el servidor simulado de MSW antes y después de los tests.


## Estructura de carpetas (resumen)

- `src/__tests__/`: aquí viven los archivos de prueba, organizados por tipo (components, hooks, pages, reducers, utils, etc.).
- `src/test/`: utilidades para pruebas y los “handlers” de MSW que simulan la API.
  - `src/test/server.ts`: arranca MSW con todos los handlers.
  - `src/test/handlers/*.ts`: respuestas simuladas para rutas (login, register, cursos, categorías).
- `__mocks__/`: reemplazos de librerías/archivos para que Jest no se rompa con imágenes o CSS.
  - `__mocks__/next/navigation.js`: router de Next.js simulado.
  - `__mocks__/next/image.js`: reemplaza la imagen de Next por un `<img>` simple.
  - `__mocks__/swiper-react.tsx` y `__mocks__/swiper-modules.ts`: versión simplificada del slider.
  - `__mocks__/styleMock.js` y `__mocks__/fileMock.js`: ignoran CSS/archivos estáticos.

Nota sobre la carpeta `test/` en la raíz del proyecto: hoy existen archivos vacíos en `ucc-arq-soft-front/test/server.ts` y `ucc-arq-soft-front/test/handlers/*`. No se usan. Los que realmente usa la suite están en `src/test/*`. Se pueden borrar para evitar confusión (no afecta a los tests).


## Qué se está probando hoy

- Componentes de autenticación:
  - `LoginForm`: llenado de formulario, envío y respuestas simuladas.
  - `RegisterForm`: creación de usuario con validaciones básicas.
- Componentes de cursos:
  - `CourseCard`: render de datos y formatos.
  - `CoursesList`: lista de cursos con datos de la API simulada.
  - `CourseCarrusel`: slides y selección de ítems.
- Reducers (lógica de estado):
  - `authReducer`, `coursesReducer`, `uiReducer` con sus distintas transiciones.
- Hooks y utilidades:
  - `useNavbar`: estados de navegación y scroll.
  - Mappers (transforman datos): `courseMapper`, `commentMapper`.

MSW (el simulador de API) responde a rutas como:
- `POST /auth/login`, `POST /users/register`, `POST /auth/refresh-token`
- `GET /courses`, `GET /categories`

Estas respuestas de ejemplo permiten que los componentes funcionen “como si” hablaran con un servidor, pero todo ocurre localmente y de forma controlada.


## Cómo funciona MSW en 2 pasos

1) Se definen “handlers” con rutas y respuestas de ejemplo, por ej. en `src/test/handlers/auth.handlers.ts`:
   - Si el cuerpo del request no trae email o password, devolvemos 400 (error).
   - Si viene todo bien, devolvemos un JSON con `token` y `user`.

2) `src/test/server.ts` junta todos los handlers y los registra. `jest.setup.ts` lo enciende antes de los tests y lo apaga al final.

De esta manera, las pruebas no necesitan internet y siempre reciben respuestas previsibles.


## Cómo están escritos los tests (patrón AAA)

- Arrange (Preparar): render del componente y/o datos de entrada.
- Act (Actuar): interacción del usuario (clicks, tipeo, envío de formulario).
- Assert (Afirmar): verificar lo esperado en pantalla o en el estado (mensajes, navegación, etc.).

Ejemplo simple (idea):
- Preparar un formulario de login.
- Escribir email y contraseña.
- Hacer clic en "Ingresar".
- Confirmar que aparece un mensaje de bienvenida o que se navegó.


## Cobertura (coverage)

- La cobertura es un indicador, no un objetivo en sí. Informa qué partes del código ejecutaron los tests.
- La consigna del TP no exige un mínimo de cobertura, así que hoy tomamos el reporte como guía para identificar partes no probadas.


## Preguntas frecuentes

- ¿Por qué hay polyfills (TextEncoder, Streams, BroadcastChannel) en `jest.setup.ts`?
  - Porque algunas librerías asumen APIs del navegador. Como corremos en Node (sin navegador), “simulamos” esas piezas para que todo funcione.

- ¿Por qué hay un mock de `next/image`?
  - El componente real de Next para imágenes usa optimizaciones del framework. En pruebas unitarias lo reemplazamos por un `<img>` simple para evitar advertencias.

- ¿Puedo borrar la carpeta `test/` que está en la raíz?
  - Sí, sus archivos están vacíos y no se usan. La suite utiliza `src/test/*`.


## ¿Cómo agrego un test nuevo?

1) Crear un archivo dentro de `src/__tests__/...` terminando en `.test.ts` o `.test.tsx`.
2) Importar el componente/función a probar.
3) Usar React Testing Library para render, interacción y verificaciones.
4) Si el componente llama a la “API”, usar MSW para simular la respuesta (agregar o modificar un handler en `src/test/handlers`).

Ejemplo muy resumido de estructura:

```ts
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MiComponente } from '@/components/MiComponente';

it('muestra saludo al enviar', async () => {
  render(<MiComponente />);
  await userEvent.type(screen.getByLabelText(/nombre/i), 'Ana');
  await userEvent.click(screen.getByRole('button', { name: /enviar/i }));
  expect(await screen.findByText(/hola, ana/i)).toBeInTheDocument();
});
```


## Resumen

- Las pruebas del frontend están configuradas con Jest + RTL + MSW.
- Simulamos el servidor (sin internet) para que los tests sean estables y rápidos.
- La configuración vive en `jest.config.ts` y `jest.setup.ts`.
- Los “endpoints” simulados están en `src/test/handlers/*`.
- Podés ejecutar todo con `npm test`.
- La carpeta `test/` de la raíz contiene archivos vacíos y no se utiliza; es seguro borrarla para despejar dudas.

Con esto deberías poder entender qué se prueba, cómo correrlo y cómo extender la suite si hace falta.
