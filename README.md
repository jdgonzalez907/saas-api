# saas-api

API REST multimodular de SaaS construida con Go, siguiendo principios de Clean Architecture y Domain-Driven Design (DDD).

## Tecnologías

| Herramienta | Propósito |
|---|---|
| Go 1.26 | Lenguaje principal |
| PostgreSQL 17 | Base de datos relacional |
| [pgx/v5](https://github.com/jackc/pgx) | Driver PostgreSQL |
| [sqlc](https://sqlc.dev/) | Generación de código Go a partir de SQL |
| [golang-migrate](https://github.com/golang-migrate/migrate) | Migraciones de base de datos |
| [go-chi/chi](https://github.com/go-chi/chi) | Router HTTP |
| [caarlos0/env](https://github.com/caarlos0/env) | Configuración desde variables de entorno |
| [mockery](https://github.com/vektra/mockery) | Generación de mocks para tests |
| [testify](https://github.com/stretchr/testify) | Assertions y mocking en tests |
| Docker + Compose | Contenerización y orquestación local |

---

## Arquitectura del Proyecto

```
saas-api/
├── cmd/
│   └── api/
│       └── main.go                  # Punto de entrada
├── db/
│   ├── migrations/                  # Migraciones SQL (golang-migrate)
│   └── queries/                     # Queries SQL (sqlc)
├── internal/
│   ├── configuration/               # Wiring de la aplicación y pools de conexión
│   │   ├── boostrap.go
│   │   └── connection.go
│   ├── postgres/                    # Código autogenerado por sqlc
│   │   ├── db.go
│   │   ├── models.go
│   │   └── users.sql.go
│   └── users/                       # Módulo de usuarios
│       ├── domain/                  # Entidades, Value Objects, interfaces de repositorio, errores de dominio
│       ├── application/             # Casos de uso (interfaces + implementaciones)
│       └── infrastructure/
│           ├── controllers/         # Handlers HTTP, router, middlewares, helpers de request/response
│           └── database/            # Implementación concreta del repositorio de usuarios
├── mocks/                           # Mocks generados por mockery (NO editar manualmente)
│   ├── application/
│   └── domain/
├── .gemini/rules.md                 # Reglas del proyecto para IAs y desarrolladores
├── .mockery.yaml                    # Configuración de mockery
├── sqlc.yaml                        # Configuración de sqlc
├── Dockerfile                       # Build multi-stage
├── docker-compose.yml               # Orquestación local
└── .env.example                     # Plantilla de variables de entorno
```

### Capas y Dependencias

```
HTTP Request
     │
     ▼
[controllers]  →  [application]  →  [domain]
                       │
                       ▼
               [infrastructure/database]  →  PostgreSQL
```

- **`domain`** (`internal/users/domain`): Núcleo del negocio. Sin dependencias externas. Define entidades (`User`), Value Objects (`Phone`, `Email`, `BirthDate`, etc.), la interfaz `UserRepository` y los errores de dominio.
- **`application`** (`internal/users/application`): Casos de uso. Solo depende de `domain`. Nunca importa paquetes de infraestructura.
- **`infrastructure`** (`internal/users/infrastructure`): Detalles de implementación. `database` implementa la interfaz del repositorio y `controllers` configura los handlers y router HTTP.
- **`configuration`** (`internal/configuration`): Hace el wiring de dependencias en `boostrap.go` e inicializa el pool de conexiones.


### Decisiones de Diseño

- **Value Objects inmutables**: Ningún campo es exportado directamente. Se crean con constructores que validan las reglas de negocio.
- **Entities con Patrón Params**: `NewUser(UserParams{...})` evita constructores con muchos argumentos posicionales.
- **Mutadores Wither**: Los métodos de actualización retornan una nueva instancia (`WithPhone`, `WithEmail`, etc.) sin mutar el estado original y actualizan `updatedAt` a `time.Now().UTC()`.
- **DTOs separados**: El dominio nunca expone sus internos directamente. La serialización pasa siempre por `.ToDTO()`.
- **UTC en todo el stack**: `time.Now().UTC()` en Go, `PGTZ=UTC` en Postgres, `TZ=UTC` en el contenedor.
- **Paginación por cursor**: `FindAll` usa keyset pagination con `id` como cursor para evitar `OFFSET` y garantizar índices deterministas.
- **Soft delete**: Los usuarios se marcan con `deleted_at`, nunca se eliminan físicamente.

---

## Configuración

Copiar `.env.example` a `.env` y completar los valores:

```bash
cp .env.example .env
```

| Variable | Requerida | Default | Descripción |
|---|---|---|---|
| `HTTP_PORT` | No | `8080` | Puerto del servidor HTTP |
| `DATABASE_URL` | **Sí** | — | DSN de conexión a Postgres |
| `POSTGRES_USER` | No | `users` | Solo para docker-compose |
| `POSTGRES_PASSWORD` | No | `secret` | Solo para docker-compose |
| `POSTGRES_DB` | No | `users_db` | Solo para docker-compose |
| `POSTGRES_PORT` | No | `5432` | Solo para docker-compose |

---

## Despliegue con Docker

### Prerrequisitos

- Docker >= 24
- Docker Compose >= 2

### Levantar todo

```bash
docker compose up --build
```

Esto hace en orden:
1. Levanta **Postgres** y espera a que el healthcheck pase.
2. Corre **migrate** — aplica todas las migraciones pendientes y termina.
3. Levanta la **API** una vez que las migraciones completaron.

### Detener y limpiar

```bash
# Detener servicios (preserva el volumen de datos)
docker compose down

# Detener y eliminar el volumen de datos
docker compose down -v
```

### Ver logs

```bash
docker compose logs -f api
docker compose logs migrate
```

---

## Desarrollo Local (sin Docker)

### Prerrequisitos

- Go 1.26+
- PostgreSQL 17 corriendo localmente
- `golang-migrate` CLI
- `sqlc` CLI
- `mockery` CLI
- `golangci-lint` CLI
- `gci` CLI


### Configurar variables de entorno

```bash
cp .env.example .env
# Editar DATABASE_URL con la conexión a tu Postgres local:
# DATABASE_URL=postgres://usuario:password@localhost:5432/users_db?sslmode=disable
```

### Aplicar migraciones

```bash
migrate -path db/migrations \
  -database "postgres://usuario:password@localhost:5432/users_db?sslmode=disable" \
  up
```

### Correr la aplicación

```bash
go run ./cmd/api
```

---

## Migraciones

Los archivos viven en `db/migrations/`. El formato es `{version}_{descripción}.{up|down}.sql`.

### Con Docker Compose (automático)

Las migraciones se aplican automáticamente al hacer `docker compose up`. Si ya están aplicadas, `migrate` termina inmediatamente sin hacer nada.

### Manualmente

```bash
# Aplicar todas las migraciones pendientes
migrate -path db/migrations -database "$DATABASE_URL" up

# Revertir la última migración
migrate -path db/migrations -database "$DATABASE_URL" down 1
```

### Agregar una nueva migración

```bash
migrate create -ext sql -dir db/migrations -seq nombre_de_la_migracion
```

Esto crea dos archivos: `.up.sql` (aplicar) y `.down.sql` (revertir).

### Qué hacer si falla una migración

**Escenario: falla antes de aplicarse** — El contenedor `migrate` termina con error y la API no arranca. Ver el error:

```bash
docker compose logs migrate
```

Corregir el SQL y volver a correr solo migrate:

```bash
docker compose up migrate
```

**Escenario: falla a mitad de ejecución (estado "dirty")** — `golang-migrate` marca la versión como sucia. En la siguiente ejecución los logs mostrarán:

```
error: Dirty database version 2. Fix and force version.
```

El número es la versión que quedó incompleta. Para identificar qué SQL falló:

1. Ver los logs de la ejecución original (disponibles hasta hacer `docker compose down`):
   ```bash
   docker compose logs migrate
   ```
2. Abrir el archivo de migración correspondiente: `db/migrations/000002_*.up.sql`

Resolución:

```bash
# 1. Forzar la versión a estado limpio (reemplazar N con el número de versión)
migrate -path db/migrations -database "$DATABASE_URL" force N

# 2. Corregir el SQL en db/migrations/000N_*.up.sql

# 3. Volver a aplicar
migrate -path db/migrations -database "$DATABASE_URL" up
```

---

## Generación de Código

### sqlc

`sqlc` genera el código Go del paquete `internal/postgres/` a partir de los archivos en `db/queries/` y el schema en `db/migrations/`.

**Cuándo correrlo**: Al agregar, modificar o eliminar cualquier query en `db/queries/users.sql`.

```bash
$(go env GOPATH)/bin/sqlc generate
```

Los archivos en `internal/postgres/` son generados — **no editarlos manualmente**.

### mockery

`mockery` genera los mocks en `mocks/` a partir de las interfaces definidas en `internal/users/domain/` e `internal/users/application/`.

**Cuándo correrlo**: Al agregar, renombrar o cambiar la firma de cualquier método en una interfaz trackeada.

```bash
$(go env GOPATH)/bin/mockery
```

Después de regenerar:
- Eliminar manualmente cualquier archivo de mock que ya no corresponda a una interfaz existente.
- Actualizar los archivos de test que referencien el nombre anterior del mock.

Los archivos en `mocks/` son generados — **no editarlos manualmente**.

---

## Calidad de Código y Linters

Para mantener un estándar de calidad alto y consistente en la industria, utilizamos:
- **`golangci-lint`**: Ejecuta múltiples analizadores estáticos esenciales (`govet`, `errcheck`, `staticcheck`, `revive`, `unused`, etc.).
- **`gci`**: Organiza los imports de forma determinista dividiéndolos en tres bloques: estándar, terceros y local (`jdgonzalez907/saas-api`).
- **`goimports`**: Aplica el formateo estándar de Go compatible con la estructuración de imports.

### Comandos de Linter

```bash
# Ejecutar el linter localmente
golangci-lint run

# Organizar imports manualmente en un archivo
gci write --section Standard --section Default --section "Prefix(jdgonzalez907/saas-api)" internal/users/domain/user.go
```

### Integración en VS Code / Cursor

Al abrir el repositorio en VS Code o Cursor, se aplicarán automáticamente el formateo y la organización de imports al guardar cualquier archivo `.go` gracias a la configuración en `.vscode/settings.json`.


---

## Tests

```bash
# Correr todos los tests
go test ./...

# Con reporte de cobertura
go test -coverprofile=coverage.out ./internal/... && go tool cover -func=coverage.out
```

La cobertura mínima requerida es **100% de sentencias** en `internal/users/domain`, `internal/users/application` e `internal/users/infrastructure/controllers`.

---

## GitFlow

### Flujo estándar para features

```bash
git checkout develop
git checkout -b feature/nombre-descriptivo

# Implementar cambios...

git add -A && git commit -m "feat: descripción del cambio funcional"
git add -A && git commit -m "test: pruebas unitarias con cobertura 100%"

# Verificar antes de merge
go build ./... && go test ./...

git checkout develop && git merge feature/nombre-descriptivo
git checkout master && git merge develop
git branch -d feature/nombre-descriptivo
```

### Hotfix (corrección urgente en producción)

```bash
git checkout master
git checkout -b hotfix/nombre-del-problema

# Corregir...

git checkout master && git merge hotfix/nombre-del-problema
git checkout develop && git merge hotfix/nombre-del-problema
git branch -d hotfix/nombre-del-problema
```

### Qué verificar antes de cada merge

1. `go build ./...` sin errores.
2. `golangci-lint run` limpio y sin advertencias/errores.
3. `go test ./...` todos en verde.
4. Cobertura 100% en paquetes `internal/`.
5. Si se cambió una interfaz: mocks regenerados con `mockery` y tests actualizados.
6. Si se modificó un query SQL: código regenerado con `sqlc generate`.
7. Si se agregó una migración: probada localmente antes del merge.


---

## Rules y Guías para IAs y Devs

Las reglas detalladas de arquitectura, patrones, tests y convenciones viven en:

```
.gemini/rules.md
```

Este archivo es la fuente de verdad del proyecto. `.cursorrules`, `.clinerules` y `.geminirules` apuntan a él.

Cualquier decisión arquitectónica nueva que se tome durante el desarrollo **debe quedar registrada en `rules.md`** antes de hacer merge a `develop`.
