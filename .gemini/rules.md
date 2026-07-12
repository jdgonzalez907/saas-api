# Reglas de Desarrollo y Buenas Prácticas (AI Agent Guidelines)

Este documento sirve como la fuente de verdad (Skill / Rules) para que cualquier inteligencia artificial o desarrollador que trabaje en este repositorio (`users-api`) mantenga la consistencia arquitectónica, patrones de diseño, y estándares de pruebas del proyecto.

---

## 1. Entities vs Value Objects (DDD)

### Value Objects
- **Definición**: Objetos sin identidad propia definidos únicamente por sus atributos.
- **Inmutabilidad Estricta**: Todos los campos internos deben ser privados (no exportados) para evitar el acoplamiento y asegurar la encapsulación. Los Value Objects nunca se mutan, solo se crean nuevas instancias. Si el Value Object es muy grande, se puede emplear el patrón "Wither" (ej. `.WithStreet(...)`) para retornar una nueva copia con el valor actualizado sin mutar la instancia original.
- **Sin Getters por Defecto**: No se deben exponer getters públicos (`Value()`, `Street()`, etc.) para sus campos de negocio, a menos que una regla de negocio específica lo requiera (YAGNI).
- **Constructor**: Deben crearse mediante un constructor del tipo `New[ValueObject](...) (ValueObject, error)` que valide sus reglas de negocio.
- **Métodos**: Únicamente deben exponer el método `.ToDTO()` y ningún otro método o mutador, a menos que sea requerido explícitamente por reglas de negocio (YAGNI).
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
  - Se deben definir métodos "Wither" intencionales y específicos para cada caso de uso de mutación (ej. `.WithPersonalInformation(...)`, `.WithPhone(...)`, `.WithEmail(...)`).
  - Estos métodos retornan un nuevo puntero a la entidad actualizada, manteniendo la inmutabilidad y actualizando el campo `updatedAt`.
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
- **Cobertura de Código**: Se requiere obligatoriamente una cobertura del **100.0% de sentencias** (statement coverage) para los paquetes `internal/domain` e `internal/application`.
- **Verificación**: Ejecutar la validación con:
  ```bash
  go test -coverprofile=coverage.out ./internal/... && go tool cover -func=coverage.out
  ```

---

## 4. Flujo de Trabajo Git (GitFlow)

Cada cambio y nueva funcionalidad debe apegarse estrictamente al siguiente flujo:
1. **Rama de Origen**: Partir siempre de `develop`.
2. **Creación de Rama**: Crear una rama de funcionalidad con el prefijo `feature/` (ej. `feature/nombre-de-la-tarea`).
3. **Mínimo 2 Commits**:
   - **Commit Funcional**: Implementa la lógica de negocio (`domain` y `application`).
   - **Commit de Tests**: Implementa las pruebas unitarias y mocks necesarios para lograr el 100% de cobertura.
4. **Merge a Develop**: Una vez validados los tests y cobertura localmente, realizar merge a `develop`.
5. **Merge a Master**: Integrar cambios a `master` para releases estables.
6. **Aprobación**: Todos los merges deben ser revisados y aprobados por el desarrollador líder.
