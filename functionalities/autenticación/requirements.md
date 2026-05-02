# Introduction

Sistema de autenticación para administradores y bibliotecarios de Lumina Library. Permite el registro de nuevos usuarios (endpoint abierto), el inicio de sesión con email y contraseña, y el cierre de sesión con invalidación de tokens. Los usuarios se crean con `is_active = false` (el servicio lo setea explícitamente; el default en BD es `true`) y se activan manualmente. Los usuarios se almacenan en la tabla `users` con campos `first_name`, `last_name`, `email`, `password_hash`, `role` (enum: admin, librarian), `avatar_url`, `is_active`. Backend en Go con Gin, PostgreSQL, arquitectura en capas y JWT (access + refresh) con blacklist (`revoked_tokens`) para logout.

# Glossary

- **Usuario (User):** Administrador o bibliotecario que gestiona la biblioteca. Almacenado en tabla `users`. Tiene credenciales (email + password). Roles: `admin` o `librarian`.
- **user_role:** Enum PostgreSQL con valores `admin` y `librarian`.
- **is_active:** Booleano en tabla `users`. Los usuarios registrados se crean con `is_active = false` (desde la capa de servicio) y se activan manualmente. Un usuario con `is_active = false` no puede iniciar sesión.
- **Access Token:** JWT de corta duración (10h) usado en el header `Authorization` para acceder a recursos protegidos.
- **Refresh Token:** JWT usado para obtener un nuevo par de access/refresh sin volver a enviar credenciales.
- **Blacklist (revoked_tokens):** Tabla con JTIs de tokens revocados. Columnas: `jti` (PK), `user_id` (FK), `expires_at`, `revoked_at`. Un token en blacklist no es válido aunque no haya expirado.
- **Cron de limpieza:** Job cada 2 días que elimina de `revoked_tokens` los registros con `revoked_at` > 24 horas.

# Requirements

### Requirement 1: Registro de usuario (admin/librarian)

**User Story:** Como sistema, quiero permitir el registro de un nuevo usuario con nombre, apellido, email, contraseña y rol, para que pueda ser activado manualmente y luego iniciar sesión.

#### Acceptance Criteria

1. WHEN se envía POST con `first_name`, `last_name`, `email`, `password` y opcionalmente `role` (default `admin`) y `avatar_url`, THE sistema SHALL crear el usuario con `is_active = false`, almacenar la contraseña hasheada (bcrypt) y devolver 201 con los datos del usuario (sin password_hash).
2. WHEN el email ya existe en la base de datos, THE sistema SHALL responder 409 con mensaje genérico de conflicto.
3. WHEN falta algún campo obligatorio (`first_name`, `last_name`, `email`, `password`) o la validación falla, THE sistema SHALL responder 400 con un objeto de errores de validación estructurado.
4. WHEN la petición no es JSON válido o el cuerpo está vacío, THE sistema SHALL responder 400 con mensaje genérico.
5. WHEN se envía un `role` que no es `admin` ni `librarian`, THE sistema SHALL responder 400 con error de validación.

### Requirement 2: Login de usuario

**User Story:** Como administrador/bibliotecario, quiero iniciar sesión con email y contraseña para obtener tokens y acceder al sistema.

#### Acceptance Criteria

1. WHEN se envían email y contraseña correctos y el usuario existe con `is_active = true`, THE sistema SHALL responder 200 con access_token y refresh_token (JWT) con tiempo de vida de 10 horas.
2. WHEN las credenciales son incorrectas o el usuario no existe, THE sistema SHALL responder 401 con mensaje genérico ("credenciales inválidas").
3. WHEN el usuario existe pero `is_active = false`, THE sistema SHALL responder 401 con el mismo mensaje genérico.
4. WHEN falta email o password o la validación falla, THE sistema SHALL responder 400 con objeto de errores de validación.

### Requirement 3: Refresh de tokens

**User Story:** Como usuario autenticado, quiero renovar mis tokens usando el refresh token para mantener la sesión.

#### Acceptance Criteria

1. WHEN se envía un refresh token válido (no expirado, no revocado, usuario activo), THE sistema SHALL responder 200 con un nuevo access_token y refresh_token.
2. WHEN el refresh token está expirado, revocado o es inválido, THE sistema SHALL responder 401.
3. WHEN el usuario asociado al token tiene `is_active = false`, THE sistema SHALL responder 401.

### Requirement 4: Logout (revocación de tokens)

**User Story:** Como usuario autenticado, quiero cerrar sesión para que mis tokens actuales dejen de ser válidos.

#### Acceptance Criteria

1. WHEN se envía POST con access token válido en Authorization y opcionalmente refresh_token en body, THE sistema SHALL insertar los JTI en `revoked_tokens` y responder 204.
2. WHEN no se envía token o el token es inválido/expirado, THE sistema SHALL responder 401.
3. WHEN un token está en `revoked_tokens`, THE sistema SHALL considerarlo inválido en todas las operaciones.

### Requirement 5: Validación y manejo de errores

**User Story:** Como desarrollador, quiero respuestas HTTP consistentes y validación estructurada.

#### Acceptance Criteria

1. WHEN cualquier entrada no cumple las reglas de validación (formato email, longitud password >= 8, campos requeridos, max lengths), THE sistema SHALL responder 400 con payload estructurado (message + errors por campo).
2. WHEN se produce un error interno no esperado, THE sistema SHALL responder 500 con mensaje genérico sin exponer detalles.
3. WHEN se intenta acceder a un recurso protegido sin token o con token inválido, THE sistema SHALL responder 401 con mensaje genérico.

### Requirement 6: Limpieza de la lista negra (cron)

**User Story:** Como sistema, quiero limpiar periódicamente la tabla `revoked_tokens` para no acumular registros indefinidamente.

#### Acceptance Criteria

1. WHEN han pasado 2 días desde la última ejecución, THE sistema SHALL ejecutar un job que elimine de `revoked_tokens` todos los registros con `revoked_at < NOW() - INTERVAL '24 hours'`.
2. WHEN el job se ejecuta, SHALL hacerlo sin bloquear peticiones HTTP.
3. WHEN el job falla, THE sistema SHALL registrar el error y continuar; la siguiente ejecución intentará de nuevo.
