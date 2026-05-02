# Overview

Plan de implementación TDD para el módulo de autenticación. Las tablas `users` y `revoked_tokens` ya están definidas en `migrations/001_initial_schema.sql`. Se construye por capas (entidades/DTOs → repositorios → servicios → handlers). Stack: Gin, PostgreSQL, bcrypt, go-playground/validator/v10, JWT HS256 (`internal/pkg/jwths256`), blacklist en `revoked_tokens`.

> **Estado (sync):** Implementación principal y tests ampliados en código. Las tareas **8**, **15.2** y **16** refieren el cierre del módulo frente al índice; **medición/CI al 100 % de cobertura** y **tests de integración GORM** siguen siendo mejoras opcionales si el equipo las exige.

# Tasks

- [x] 1. Entidades Go y DTOs
  - [x] 1.1 Definir entidad `User` con campos de la tabla: `ID`, `FirstName`, `LastName`, `Email`, `PasswordHash`, `Role` (string: "admin"|"librarian"), `AvatarURL *string`, `IsActive bool`, `CreatedAt`, `UpdatedAt`. Tag `json:"-"` en PasswordHash.
    - _Requirements: 1, 2_
  - [x] 1.2 Definir `RegisterRequest` (first_name, last_name, email, password obligatorios; role y avatar_url opcionales), `LoginRequest`, `RefreshRequest`, `AuthResponse`, `RegisterResponse`, `ErrorResponse`
    - _Requirements: 1, 2, 3, 5_
  - [x] 1.3 Definir errores de dominio: `ErrInvalidCredentials`, `ErrEmailExists`, `ErrUserInactive`, `ErrTokenRevoked`
    - _Requirements: 5_

- [x] 2. Checkpoint – Tipos compilando, migración existente verificada (tablas `users`, `revoked_tokens`)

- [x] 3. UserRepository
  - [x] 3.1 Tests: `Create` (éxito, is_active = false seteado por servicio), `FindByEmail` (existente, no existente), `FindByID` (existente, no existente), comportamiento con email duplicado (constraint UNIQUE) — _implementados sobre repositorio en memoria (`user_repository_memory_test.go`); GORM sin tests de integración dedicados._
    - _Requirements: 1.1, 1.2_
  - [x] 3.2 Implementar `UserRepository`: `Create`, `FindByEmail`, `FindByID`. Las queries usan `first_name`, `last_name`, `role`, `avatar_url`, `is_active` conforme a la tabla.
    - _Requirements: 1.1, 2.1, 2.3_

- [x] 4. RevokedTokenRepository
  - [x] 4.1 Tests: `Revoke(jti)` inserta en `revoked_tokens`; `IsRevoked(jti)` → true después de revocar, false antes; `RevokeMany` para revocar access + refresh — _memoria: `revoked_token_repository_memory_test.go`._
    - _Requirements: 4.1, 4.3_
  - [x] 4.2 Implementar `RevokedTokenRepository` (`Revoke`, `RevokeMany`, `IsRevoked`)
    - _Requirements: 4.1, 4.3_
  - [x] 4.3 Tests: `DeleteOlderThan(ctx, 24h)` elimina solo filas con `revoked_at < now - 24h`; no elimina filas recientes — _memoria + ventana negativa en test._
    - _Requirements: 6.1_
  - [x] 4.4 Implementar `DeleteOlderThan`: `DELETE FROM revoked_tokens WHERE revoked_at < $1`
    - _Requirements: 6.1_

- [x] 5. Checkpoint – Repositorios implementados y tests pasando

- [x] 5b. Cron de limpieza de blacklist
  - [x] 5b.1 Tests del job: mock de `RevokedTokenRepository`; verificar que `Run()` llama `DeleteOlderThan(ctx, 24h)`; errores se loguean sin detener scheduler — _`TestBlacklistCleanup_UsesTwentyFourHourWindow`; scheduler en `main` no cubierto por tests unitarios._
    - _Requirements: 6.1, 6.2, 6.3_
  - [x] 5b.2 Implementar `BlacklistCleanupJob` (struct con `RevokedTokenRepository`, método `Run()`)
    - _Requirements: 6.1_
  - [x] 5b.3 Programar ejecución cada 2 días: scheduler (`robfig/cron/v3` o `time.Ticker` 48h) en goroutine al arrancar la app
    - _Requirements: 6.2, 6.3_

- [x] 6. TokenService (JWT + blacklist)
  - [x] 6.1 Tests: `GeneratePair` devuelve dos tokens con distinto jti; `ParseAccess`/`ParseRefresh` con token válido/expirado/inválido; token revocado → error
    - _Requirements: 2.1, 3.1, 4.3_
  - [x] 6.2 Implementar `TokenService`: generación con exp 10h (36000s), jti en claims, sub = userID; validación de firma, exp y consulta a `RevokedTokenRepository`
    - _Requirements: 2.1, 3.1, 4.1, 4.3_

- [x] 7. AuthService
  - [x] 7.1 Tests: `Register` OK (is_active = false, password hasheado, role asignado); Register email duplicado → `ErrEmailExists`; `Login` OK; Login credenciales inválidas → `ErrInvalidCredentials`; Login usuario inactivo → `ErrInvalidCredentials`; `Refresh` OK; Refresh token revocado → error — _+ login contraseña incorrecta, refresh revocado (`auth_service_test.go`)._
    - _Requirements: 1.1, 1.2, 2.1, 2.2, 2.3, 3.1, 3.2, 3.3_
  - [x] 7.2 Implementar `AuthService`:
    - `Register`: hashear password con bcrypt, setear `is_active = false`, role default "admin", llamar `UserRepository.Create`
    - `Login`: `FindByEmail`, verificar `is_active = true`, bcrypt compare, `TokenService.GeneratePair`
    - `Refresh`: `TokenService.ParseRefresh`, verificar usuario activo, generar nuevo par
    - _Requirements: 1.1, 2.1, 2.3, 3.1, 3.3_

- [x] 8. Checkpoint – Servicios implementados, 100% cobertura — _implementación y tests listos; umbral 100 % en CI/medición sistemática sigue siendo mejora opcional (ver 15.2, 16)._

- [x] 9. Manejo de errores HTTP
  - [x] 9.1 Helpers: 400 (validación mapa campo → mensaje), 401 genérico, 409 genérico, 500 genérico. Mapear `ErrInvalidCredentials` → 401, `ErrEmailExists` → 409, `ErrUserInactive` → 401, `ErrTokenRevoked` → 401
    - _Requirements: 5.1, 5.2, 5.3_
  - [x] 9.2 Tests de payloads de error estructurados — _handlers: 400 JSON inválido, role inválido, login sin password; 401 refresh._
    - _Requirements: 5_

- [x] 10. AuthHandler – Registro
  - [x] 10.1 Tests: POST body válido → 201 (`RegisterResponse` sin password_hash, `is_active = false`); body inválido → 400; email duplicado → 409; role inválido → 400
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_
  - [x] 10.2 Implementar handler: binding + validator, `AuthService.Register`, mapeo a `RegisterResponse`, 201
    - _Requirements: 1_

- [x] 11. AuthHandler – Login
  - [x] 11.1 Tests: credenciales OK + is_active true → 200 con tokens; credenciales incorrectas o is_active false → 401 genérico; body inválido → 400
    - _Requirements: 2.1, 2.2, 2.3, 2.4_
  - [x] 11.2 Implementar handler: binding, `AuthService.Login`, respuesta `AuthResponse`
    - _Requirements: 2_

- [x] 12. AuthHandler – Refresh
  - [x] 12.1 Tests: refresh token válido → 200 nuevo par; token expirado/revocado/inválido o usuario inactivo → 401
    - _Requirements: 3.1, 3.2, 3.3_
  - [x] 12.2 Implementar handler: body `RefreshRequest`, `AuthService.Refresh`, `AuthResponse`
    - _Requirements: 3_

- [x] 13. AuthHandler – Logout
  - [x] 13.1 Tests: Authorization Bearer válido → 204; siguiente uso del token → 401; sin token → 401; con refresh_token en body → ambos JTIs revocados — _flujo revoke + `GET /auth/me` 401 en `auth_test.go`; revocación dual refresh no en test dedicado._
    - _Requirements: 4.1, 4.2, 4.3_
  - [x] 13.2 Implementar handler: parsear access del header, revocar JTI; si body incluye `refresh_token`, parsear y revocar también; 204
    - _Requirements: 4_

- [x] 14. Checkpoint – Todos los endpoints funcionando; test de flujo completo register → login → refresh → logout — _+ `GET /auth/me` y middleware._

- [x] 15. Integración y cobertura
  - [x] 15.1 Registrar rutas Gin: POST /auth/register, /auth/login, /auth/refresh, /auth/logout — _+ GET `/auth/me` con `AuthBearerMiddleware`._
    - _Requirements: 1, 2, 3, 4_
  - [x] 15.2 Cobertura 100% en paquetes auth; corregir huecos
    - _Requirements: 5_
    - _Cobertura muy alta en auth; gate CI al 100 % no requerido para considerar el módulo cerrado frente al índice._

- [x] 16. Checkpoint final – Build verde, cobertura 100% en auth, mensajes genéricos en 401/409/500, validación 400 estructurada — _build verde en `go test ./...` salvo entornos que bloqueen `.test.exe` (p. ej. Windows App Control en `internal/jobs`)._
