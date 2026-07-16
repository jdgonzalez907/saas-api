# Reglas de Desarrollo y Buenas Prácticas (AI Agent Guidelines)

Este documento sirve como la fuente de verdad (Skill / Rules) para que cualquier inteligencia artificial o desarrollador que trabaje en este repositorio (`saas-api`) mantenga la consistencia arquitectónica, patrones de diseño, y estándares de pruebas del proyecto.

---

## 1. Entities vs Value Objects (DDD)

### Value Objects
- **Definición**: Objetos sin identidad propia definidos únicamente por sus atributos.
- **Inmutabilidad Estricta**: Todos los campos internos deben ser privados (no exportados) para evitar el acoplamiento y asegurar la encapsulación. Los Value Objects nunca se mutan, solo se crean nuevas instancias. Si el Value Object es muy grande, se puede emplear el patrón "Wither" (ej. `.WithStreet(...)`) para retornar una nueva copia con el valor actualizado sin mutar la instancia original.
- **Sin Getters por Defecto**: No se deben exponer getters públicos (`Value()`, `Street()`, etc.) para sus campos de negocio, a menos que una regla de negocio específica lo requiera (YAGNI). Excepción: se permiten getters cuando la capa de infraestructura los necesita para construir queries sin pasar por un DTO.
- **Constructor**: Deben crearse mediante un constructor del tipo `New[ValueObject](...) (ValueObject, error)` que valide sus reglas de negocio.
- **Métodos**: Únicamente deben exponer el método `.ToDTO()` y los getters estrictamente necesarios. Ningún otro método o mutador, a menos que sea requerido explícitamente por reglas de negocio (YAGNI).
- **Patrón DTO**:
  - Cada VO debe tener su correspondiente estructura `[ValueObject]DTO` pública con tags `json` que define el formato de serialización.
  - Implementar el método `.ToDTO() [ValueObject]DTO` en el VO para mapear sus campos privados a la estructura plana.
- *Ejemplos en el proyecto*: `Phone`, `Email`, `Identification`, `Address`, `BirthDate`, `PaginatedUsers`.

### Entities
- **Definición**: Objetos con una identidad única (`id`) que persiste en el tiempo.
- **Inmutabilidad Estricta**: Los campos internos de la entidad deben ser privados (no exportados) para evitar mutaciones directas fuera del dominio.
- **Patrón Params**: El constructor `New[Entity](params [Entity]Params) (*[Entity], error)` recibe una estructura de parámetros pública para inicializar la entidad.
- **Patrón DTO**:
  - Una estructura `[Entity]DTO` pública con tags `json` define el formato de serialización.
  - El método `.ToDTO() [Entity]DTO` exporta el estado interno a una estructura legible.
  - La función `[Entity]FromDTO(dto *[Entity]DTO) (*[Entity], error)` reconstruye la entidad desde su representación plana.
- **Patrón Wither (Mutadores Específicos)**:
  - No se permiten actualizadores genéricos tipo `.With(...)`.
  - Se deben definir métodos "Wither" intencionales y específicos para cada caso de uso de mutación con nombres que den intención de negocio (ej. `.UpdatePersonalInformation(...)`, `.ChangePhone(...)`, `.ChangeEmail(...)`).
  - Estos métodos retornan un nuevo puntero a la entidad actualizada, manteniendo la inmutabilidad y actualizando el campo `updatedAt` con `time.Now().UTC()`.
  - Las operaciones puras en memoria que no puedan fallar no deben retornar `error`.

### Reglas de Referencias y Acoplamiento
- **Aislamiento de Dominio**: Las Entidades y Value Objects de dominio únicamente deben referenciar/llamar a otras Entidades y Value Objects. Nunca deben acoplarse ni referenciar estructuras DTO en sus campos o lógica de negocio interna.
- **Aislamiento de Serialización**: Los DTOs únicamente deben referenciar/llamar a otros DTOs (ej. `UserDTO` referencia `PhoneDTO`, `AddressDTO`, etc. en lugar de sus equivalentes VO de dominio).

---

## 2. Estructura de Casos de Uso (Application Layer)

Cada caso de uso en la capa de aplicación debe seguir una estructura estricta y limpia:
- **Interfaz**: Definida en el mismo archivo para desacoplamiento y facilidad de mocking (ej. `type UpdateUserPhoneUseCase interface`).
- **Struct de Implementación**: Estructura privada que implementa la interfaz (ej. `type updateUserPhoneUseCase struct`).
- **Único Método Público**: Solo debe exponer el método ejecutor principal: `Execute(...) error` (o retornar entidad/value object y error para consultas).
- **Propagación de Contexto**: Todos los métodos `Execute` reciben `context.Context` como primer argumento y lo propagan al repositorio.
- **Tipos de Datos de Entrada/Salida**: Los casos de uso deben recibir únicamente entidades/value objects como entrada y retornar únicamente entidades/value objects y/o errores como salida. Nunca deben recibir ni retornar estructuras DTO en sus firmas públicas.
- **Métodos Auxiliares Privados**: Cualquier lógica de validación, mapeo o cálculo adicional debe encapsularse en métodos privados del struct de implementación.
- **Manejo de Errores (Error Wrapping)**: Cuando se capture un error proveniente de la capa de infraestructura (ej. del repositorio), este debe envolverse con el error de dominio correspondiente utilizando `fmt.Errorf("%v: %w", domain.ErrXxx, err)`. Esto permite que el error original sea inspeccionable mediante `errors.Unwrap` manteniendo a la vez la semántica de negocio en las pruebas unitarias y capas externas.

---

## 3. Guía de Pruebas Unitarias y Cobertura

- **Table-Driven Tests**: Obligatorio utilizar pruebas basadas en tablas (table-driven) cuando existan más de 2 caminos de ejecución / flujos lógicos en la función a probar.
- **Nomenclatura de Casos**: Cada escenario dentro del slice de casos de prueba (`testCases`) debe nombrarse explícitamente bajo los formatos:
  - `"success - [descripción corta del éxito]"`
  - `"fail - [descripción corta del error/falla]"`
- **Mocking**: Utilizar `mockery` para la generación automática de mocks. Los mocks deben estar ubicados bajo la carpeta `mocks/`.
- **Cobertura de Código**: Se requiere obligatoriamente una cobertura del **100.0% de sentencias** (statement coverage) para todos los paquetes de negocio y de infraestructura bajo el directorio `internal/` (incluyendo `internal/<modulo>/domain`, `internal/<modulo>/application` e `internal/<modulo>/infrastructure/controllers`).
- **Prohibición de Reflexión (No Reflection)**: Está estrictamente prohibido usar el paquete `reflect` o funciones como `reflect.DeepEqual` en las pruebas unitarias. Las validaciones de aserción e igualdad de entidades y value objects deben realizarse de forma explícita y manual, campo por campo. Ante campos opcionales o punteros, se debe validar si son diferentes de `nil` y, de ser así, verificar sus campos internos individualmente. Esto asegura claridad, legibilidad y facilidad de depuración.
- **Nombres de Variables con Intención**: Queda prohibido el uso de nombres de variables genéricos, numerados o confusos (`x1`, `x2`, `user1`, `user2`, etc.) en los archivos de pruebas. Las variables deben ser nombradas con nombres que demuestren claramente su intención (ej. `firstUser`, `secondUser`, `personalInfo`, `badEmail`).
- **Verificación**: Ejecutar la validación con:
  ```bash
  go test -coverprofile=coverage.out ./internal/... && go tool cover -func=coverage.out
  ```

---

## 4. Flujo de Trabajo Git (GitFlow)

Cada cambio y nueva funcionalidad debe apegarse estrictamente al siguiente flujo:
1. **Rama de Origen**: Partir siempre de `develop`.
2. **Creación de Rama**: Crear una rama de funcionalidad con el prefijo `feature/` (ej. `feature/nombre-de-la-tarea`). Para correcciones urgentes en producción usar `hotfix/`.
3. **Mínimo 2 Commits**:
   - **Commit Funcional**: Implementa la lógica de negocio (`domain` y `application`).
   - **Commit de Tests**: Implementa las pruebas unitarias y mocks necesarios para lograr el 100% de cobertura.
4. **Verificación antes de merge**: `go build ./...` y `go test ./...` deben pasar sin errores.
5. **Merge a Develop**: Una vez validados los tests y cobertura localmente, realizar merge a `develop`.
6. **Merge a Master**: Integrar cambios a `master` para releases estables.
7. **Limpieza**: Eliminar la rama de funcionalidad después del merge a `master`.
8. **Aprobación**: Todos los merges deben ser revisados y aprobados por el desarrollador líder.

---

## 5. Capa de Controllers HTTP (`internal/<modulo>/infrastructure/controllers`)

Cada controller HTTP debe seguir una estructura estricta y limpia:
- **Router**: Configurado en `router.go` usando `go-chi/chi/v5`. Middlewares globales registrados en orden: `RequestID` → `Logger` → `ErrorLoggerMiddleware` → `Recoverer` → `JSONContentTypeMiddleware`. No usar `RealIP` (deprecated).
- **Struct de Controller**: Estructura pública que recibe sus dependencias (casos de uso) por constructor (ej. `NewUserController(...)`).
- **Responsabilidad del Handler**: El método handler parsea la request y construye las entidades/VOs de dominio necesarias, llama al caso de uso (que recibe y devuelve solo entidades/VOs), y convierte el resultado a DTO vía `.ToDTO()` para responder. Nunca pasar DTOs directamente a los casos de uso.
- **Mapeo de DTOs Planos y Embebidos**:
  - Para campos que usan DTOs de tipo primitivo (ej. `EmailDTO` definido como `string`), el handler debe decodificar la request directamente a este tipo y luego pasar el valor subyacente para construir el Value Object del dominio (`Email`) antes de llamar al caso de uso.
  - Para DTOs que usan composición por struct embedding (ej. `UserDTO` que embebe anónimamente `PersonalInformationDTO`), el handler debe extraer los campos del DTO embebido para construir primero el Value Object de dominio (`PersonalInformation`) y luego pasarlo al constructor de la entidad (`User`).
- **Parseo de Parámetros de Ruta**: Usar siempre el helper centralizado `ParseRouteIntParam(r, "paramName")` para parámetros enteros. Agregar nuevos helpers en `request.go` para otros tipos de parseo. Cualquier fallo de parseo responde con `400 Bad Request`.
- **Respuestas de Éxito**: Usar `RespondWithJSON(w, statusCode, entity.ToDTO())` sin envoltura. Código `200` para consultas y actualizaciones, `201` para creaciones, `204` para eliminaciones (pasar `nil` como data).
- **Respuestas de Error de Dominio**: Usar siempre `RespondWithDomainError(w, err)`. Esta función inspecciona el error mediante `errors.Is` contra el mapa privado `domainErrorStatus` que centraliza el mapeo de cada error de dominio a su código HTTP. Al crear un nuevo error de dominio, añadir su entrada al mapa (una línea). Errores 5xx responden con el mensaje genérico `"internal server error"` para no filtrar detalles internos; errores 4xx responden con el mensaje del centinela de dominio.
- **Pruebas Unitarias**: Usar `httptest.NewRecorder()` y `httptest.NewRequest(...)`. Levantar el router real con `controllers.NewRouter(controller)` para que los middlewares se apliquen en pruebas. Cobertura mínima del 100%. Para `RespondWithDomainError` basta con 3 casos (uno por rama): error conocido no-500, error conocido 500, y error desconocido no presente en el mapa.

---

## 6. Manejo de Tiempo (UTC)

- **Todo `time.Now()` debe llamarse como `time.Now().UTC()`** en cualquier capa del proyecto. Esto incluye creación de entidades, métodos Wither y cualquier lógica de negocio.
- **Al leer timestamps de Postgres** (`pgtype.Timestamptz`), normalizar explícitamente con `.Time.UTC()` al construir entidades de dominio. El driver puede retornar el location del servidor si no está configurado como UTC.
- **`_ "time/tzdata"`** se importa en `cmd/api/main.go` para embeber la base de datos de zonas horarias dentro del binario. Necesario porque la imagen de runtime (`distroless/static`) no tiene archivos de timezone en disco.
- A nivel de infraestructura, los servicios Docker tienen `TZ=UTC` / `PGTZ=UTC` configurados como variables de entorno para garantizar consistencia en todo el stack.

---

## 7. Convenciones de Nomenclatura Go

- **Acrónimos en mayúsculas completas**: `ID` no `Id`, `URL` no `Url`, `HTTP` no `Http`. Aplica a interfaces, structs, constructores, métodos y nombres de archivo cuando sean parte del nombre exportado.
- Esta convención es estándar de Go y debe respetarse en toda adición al proyecto.

---

## 8. Mocks (`mockery`)

- **Cuándo regenerar**: Al renombrar, agregar o eliminar métodos en cualquier interfaz trackeada por mockery (ver `.mockery.yaml`).
- **Comando**:
  ```bash
  $(go env GOPATH)/bin/mockery
  ```
- **Limpieza manual**: Si una interfaz es renombrada, el archivo de mock anterior no se elimina automáticamente. Eliminarlo manualmente antes de commitear.
- **Actualizar tests**: Después de regenerar mocks, actualizar todos los archivos de test que referencien el tipo o constructor renombrado.

---

## 9. Infraestructura Docker y Configuración

- **Dockerfile**: Multi-stage. Stage `builder` usa `golang:alpine`; stage `runtime` usa `distroless/static-debian12:nonroot`. La imagen final no tiene shell.
- **docker-compose.yml**: Tres servicios en orden de dependencia: `postgres` → `migrate` → `api`.
  - `migrate` tiene `restart: "no"` porque es una tarea de una sola ejecución. Si falla, se investiga y se vuelve a correr manualmente.
  - `api` solo arranca si `migrate` completó exitosamente (`service_completed_successfully`).
- **Variables de entorno**: Definidas en `.env` (ignorado por git). Usar `.env.example` como plantilla. La aplicación las lee con `github.com/caarlos0/env/v11`.
- **Config de la aplicación**: Vive en `internal/configuration/boostrap.go`. `DATABASE_URL` es requerida; `HTTP_PORT` tiene default `8080`.

---

## 10. Comentarios en Código

Solo se permiten comentarios que expliquen comportamiento genuinamente no obvio para un desarrollador con contexto del proyecto. Están prohibidos:
- Comentarios de agrupación que replican lo que dice el nombre de la variable (`// Repositories`, `// Controllers`).
- Comentarios que describen lo que hace la siguiente línea de código (`// Parse the ID`, `// Return error`).

---

## 11. Calidad de Código, Formateo y Linters (`golangci-lint`)

- **Linters Obligatorios**:
  - `golangci-lint` es la herramienta de análisis estático obligatoria del proyecto. Ningún cambio puede ser merged a `develop` si reporta advertencias o errores.
  - La configuración vive en `.golangci.yml` y excluye explícitamente el código generado (`internal/postgres/` y `mocks/`).
- **Formateo e Imports**:
  - Se debe utilizar `goimports` para el formateo estándar del código y `gci` para la organización y ordenación determinista de imports.
  - Los imports deben estar agrupados exactamente en 3 bloques separados por una línea en blanco:
    1. Biblioteca estándar de Go.
    2. Librerías y dependencias externas de terceros.
    3. Importaciones locales del módulo (`jdgonzalez907/saas-api`).
- **Integración con Editor**:
  - El editor (VS Code / Cursor) debe estar configurado con `.vscode/settings.json` para formatear y organizar los imports de forma automática al guardar.
