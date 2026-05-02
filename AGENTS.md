# AGENTS.md — library_back

Guía para agentes y desarrolladores que trabajan en este repositorio.

## Stack

- **Lenguaje**: Go 1.25+
- **HTTP**: Gin
- **Validación**: go-playground/validator/v10 (struct tags + `BindAndValidate` en handlers)
- **Contraseñas**: golang.org/x/crypto/bcrypt (via `internal/pkg/crypto`)
- **Módulo**: `library_back`

## Arquitectura

Arquitectura en capas sencilla, flujo unidireccional:

```
HTTP (Gin) → handlers → services → repositories
                  ↓           ↓
              pkg/validator  pkg/crypto, models
```

- **handlers**: Reciben request, validan con `validator.Validator#BindAndValidate`, llaman al service, escriben respuesta. No contienen lógica de negocio ni acceso a datos directo.
- **services**: Lógica de negocio, orquestan repos y utilidades (crypto). Devuelven entidades o errores de dominio (ej. `ErrEmailExists`).
- **repositories**: Abstraen persistencia. Interfaces en este paquete; implementaciones (in-memory, MongoDB, etc.) también. Los handlers no conocen repos, solo services.

## Estructura del proyecto

```
library_back/
├── main.go                 # Wiring: rutas, repos, services, handlers
├── internal/
│   ├── handlers/           # Un handler por dominio o recurso (health, user, …)
│   ├── services/          # Un service por dominio
│   ├── repositories/      # Interfaces + implementaciones de persistencia
│   ├── models/            # Entidades y DTOs con tags de validación
│   └── pkg/               # Código reutilizable sin dependencias de negocio
│       ├── validator/      # Integración validator + Gin
│       └── crypto/         # Hash/compare contraseñas (bcrypt)
```

- Todo el código de aplicación vive en `internal/`. Lo que sea reutilizable y estable va en `internal/pkg/`.
- No crear paquetes vacíos ni “placeholder”. Nuevos dominios: `handlers/<dominio>.go`, `services/<dominio>_service.go`, `repositories/<dominio>_repository.go` cuando hagan falta.

## Convenciones

- **Nombres**: PascalCase para exportados, camelCase para no exportados. Nombres de archivos en snake_case cuando sea plural o compuesto (ej. `user_repository.go`).
- **Validación**: Usar tags `binding:"..."` en DTOs de request y `Validator.BindAndValidate(c, &req)` en el handler antes de llamar al service.
- **Contraseñas**: Nunca loguear ni devolver en JSON. Hashear siempre con `crypto.HashPassword`; comparar con `crypto.ComparePassword`.
- **Errores**: En services, errores de dominio (ej. “email ya existe”) como sentinel: `var ErrEmailExists = errors.New("...")`. En handlers, traducir a códigos HTTP (409, 404, etc.).
- **Inyección**: En `main.go` se instancian repos y services; los handlers reciben el service y el validator por constructor. No usar globals para servicios o DB.

## Cómo añadir un nuevo recurso (ej. “books”)

1. **models**: Definir entidad y DTOs de request/response con tags de validación si aplica.
2. **repositories**: Definir interfaz `XxxRepository` y una implementación (in-memory o DB); registrar en `main.go`.
3. **services**: Crear `XxxService` que use el repositorio; registrar en `main.go`.
4. **handlers**: Crear handler que use `Validator.BindAndValidate` y el service; registrar rutas en `main.go`.

## Comandos útiles

```bash
go build -o library_back.exe .
go run .
# Servidor por defecto: http://localhost:8080
```

**Rate limit (por IP, ventana 1 minuto UTC):** middleware global en Gin; `/health` no cuenta. Variables opcionales: `RATE_LIMIT_DISABLED=1|true`, `RATE_LIMIT_GENERAL_RPM` (default 120), `RATE_LIMIT_RANKING_RPM` (default 40). Respuesta **429** con `Retry-After: 60`. Tras proxy inverso, configurar `Gin` trusted proxies si hace falta que `ClientIP()` refleje al cliente real.

## SDD (Spec-Driven Development)

El proyecto tiene contexto SDD en Engram (`sdd-init/library_back`). Para cambios grandes se puede usar el flujo del orchestrator (explore → propose → spec → design → tasks → apply → verify → archive). Para cambios pequeños (bugs, refactors puntuales) no es obligatorio pasar por SDD.
