# Introduction

Gestión de estudiantes (patrons) de la biblioteca Lumina. Los estudiantes NO tienen credenciales ni login; son registrados, editados y desactivados exclusivamente por administradores/bibliotecarios (tabla `users`). Los datos del estudiante se almacenan en la tabla `students` con los campos: `student_id_code` (matrícula única, ej. LUM-2024-001), `first_name`, `last_name`, `email` (contacto, opcional, único si presente), `avatar_url`, `department_id` (FK opcional a `departments`), `degree_level`, `major`, `expected_graduation_year`, `is_active` (soft delete), `member_since`, `registered_by` (FK a `users`, admin que registró). Stack: Go, Gin, PostgreSQL, arquitectura en capas (handlers → services → repositories).

# Glossary

- **Estudiante (Student):** Persona registrada en el sistema que puede solicitar libros; no tiene credenciales ni login. Se almacena en la tabla `students`.
- **student_id_code:** Código de matrícula único del estudiante (ej. LUM-2024-001). Columna `student_id_code VARCHAR(50) NOT NULL UNIQUE`.
- **Departamento (Department):** Carrera o departamento académico; tabla `departments` con `id`, `name` (unique) y `code` (unique, ej. "CS", "ENG"). FK opcional en `students.department_id`.
- **degree_level:** Nivel académico (ej. "Undergraduate Student", "Graduate Student"). Opcional.
- **major:** Especialidad o carrera del estudiante. Opcional.
- **expected_graduation_year:** Año previsto de graduación (entero). Opcional.
- **member_since:** Fecha desde que el estudiante es miembro de la biblioteca. Default: fecha de creación.
- **registered_by:** UUID del admin/bibliotecario que registró al estudiante. FK a `users(id)`.
- **is_active:** Booleano para soft delete. `true` = activo, `false` = desactivado. Los listados y búsquedas excluyen estudiantes con `is_active = false`.
- **avatar_url:** URL de la foto de perfil del estudiante. Opcional.
- **Personal Stats (estadísticas personales):** Datos agregados del estudiante en la tabla `student_stats`: total_read (libros leídos), active_loans (préstamos activos), overdue_count (préstamos vencidos), current_streak_days (racha actual en días), y opcionalmente longest_streak_days, total_points, global_rank. Se muestran en la sección "Personal Stats" del perfil.
- **Loan History (historial de préstamos):** Lista de préstamos del estudiante (tabla `loans`) con datos del libro (título, portada, autores) y fechas (borrowed_at, due_date, returned_at) y status (active | returned | overdue). Permite mostrar "Currently Reading" (activos) y "Returned" (devueltos). Opcional: valoración/rating por préstamo si en el futuro existe tabla de reseñas.
- **Badge Gallery (galería de insignias):** Insignias disponibles (tabla `badges`) y las que el estudiante ha obtenido (tabla `student_badges`). Cada badge tiene slug, name, icon_url, description, criteria. En el perfil se muestra progreso (ej. 8/12), las ganadas (earned: true, earned_at) y las no ganadas (earned: false) para mostrar como "bloqueadas" en la UI.

# Requirements

### Requirement 1: Crear estudiante

**User Story:** Como administrador, quiero registrar un estudiante con sus datos personales y académicos para gestionar quién solicita libros.

#### Acceptance Criteria

1. WHEN el administrador envía POST con los campos obligatorios (`first_name`, `last_name`, `student_id_code`) y campos opcionales válidos (`email`, `department_id`, `degree_level`, `major`, `expected_graduation_year`, `avatar_url`), THE sistema SHALL crear el estudiante con `is_active = true`, `member_since = NOW()`, `registered_by = admin_id` (extraído del JWT) y devolver 201 con los datos del estudiante. (Opcional) THE sistema MAY crear una fila en `student_stats` para ese estudiante con valores en 0 (total_read, active_loans, etc.) para que los módulos de reputación y préstamos puedan actualizarla sin tener que hacer upsert; si no se implementa, el módulo reputation hace upsert al registrar préstamos/devoluciones.
2. WHEN `student_id_code` ya existe en la base de datos, THE sistema SHALL responder 409 con mensaje de conflicto.
3. WHEN `email` ya existe en la base de datos (y no es nulo), THE sistema SHALL responder 409 con mensaje de conflicto.
4. WHEN se envía `department_id` y no existe en la tabla `departments`, THE sistema SHALL responder 400 con mensaje de validación.
5. WHEN falta algún campo obligatorio o la validación falla (ej. email con formato inválido, `student_id_code` vacío), THE sistema SHALL responder 400 con objeto de errores de validación estructurado.
6. WHEN no se envía token válido de administrador, THE sistema SHALL responder 401.

### Requirement 2: Listar estudiantes

**User Story:** Como administrador, quiero listar estudiantes activos con paginación y búsqueda por nombre o email.

#### Acceptance Criteria

1. WHEN el administrador solicita GET con paginación (`limit`, `offset`), THE sistema SHALL devolver solo estudiantes con `is_active = true`, ordenados por `created_at DESC`, con el total de registros que cumplen el filtro.
2. WHEN se envía el query param `search`, THE sistema SHALL filtrar por coincidencia parcial en nombre completo (`first_name || ' ' || last_name`) usando búsqueda trigram, o por `email` o `student_id_code`.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.
4. WHEN `limit` u `offset` son inválidos (ej. `limit > 100`, negativo), THE sistema SHALL normalizar a valores por defecto (`limit` default 20, max 100; `offset` default 0).

### Requirement 3: Obtener estudiante por ID

**User Story:** Como administrador, quiero ver el detalle completo de un estudiante.

#### Acceptance Criteria

1. WHEN el ID existe y el estudiante tiene `is_active = true`, THE sistema SHALL devolver 200 con todos los datos del estudiante (incluyendo `department` name/code si tiene `department_id`).
2. WHEN el ID no existe o el estudiante tiene `is_active = false`, THE sistema SHALL responder 404.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.

### Requirement 4: Actualizar estudiante

**User Story:** Como administrador, quiero actualizar la información de un estudiante existente.

#### Acceptance Criteria

1. WHEN el administrador envía PATCH con datos válidos para un ID existente y activo (`is_active = true`), THE sistema SHALL actualizar solo los campos enviados y devolver 200 con los datos actualizados.
2. WHEN el ID no existe o el estudiante tiene `is_active = false`, THE sistema SHALL responder 404.
3. WHEN la actualización violaría unicidad de `student_id_code` o `email` (otro estudiante ya tiene ese valor), THE sistema SHALL responder 409.
4. WHEN se envía `department_id` y no existe en `departments`, THE sistema SHALL responder 400.
5. WHEN la validación falla (ej. email mal formado), THE sistema SHALL responder 400 con errores estructurados.
6. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.

### Requirement 5: Desactivar estudiante (soft delete)

**User Story:** Como administrador, quiero desactivar un estudiante para que deje de aparecer en listados sin borrar su historial de préstamos.

#### Acceptance Criteria

1. WHEN el administrador envía DELETE para un ID existente y activo, THE sistema SHALL marcar `is_active = false` y responder 204.
2. WHEN el ID no existe o el estudiante ya está desactivado (`is_active = false`), THE sistema SHALL responder 404.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.

### Requirement 6: Listar departamentos (catálogo)

**User Story:** Como administrador, quiero obtener la lista de departamentos para usarla al crear o editar estudiantes.

#### Acceptance Criteria

1. WHEN el administrador solicita GET al endpoint de departamentos, THE sistema SHALL devolver la lista de todos los departamentos (`id`, `name`, `code`) con 200.
2. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.

### Requirement 7: Validación y errores

**User Story:** Como desarrollador, quiero respuestas HTTP consistentes y validación con mensajes por campo.

#### Acceptance Criteria

1. WHEN la entrada no cumple reglas de validación (formato email, campos requeridos, max length), THE sistema SHALL responder 400 con payload estructurado (`message` + `errors` por campo).
2. WHEN ocurre un error interno no esperado, THE sistema SHALL responder 500 con mensaje genérico.
3. WHEN se accede a cualquier endpoint de estudiantes o departamentos sin token válido de admin, THE sistema SHALL responder 401.

### Requirement 8: Endpoint de búsqueda de estudiantes

**User Story:** Como administrador, quiero buscar estudiantes por email, código de matrícula (student_id_code) o por nombre para localizar rápidamente a un estudiante.

#### Acceptance Criteria

1. WHEN el administrador envía GET al endpoint de estudiantes con el query param `search` (término de búsqueda), THE sistema SHALL devolver solo estudiantes activos que coincidan por al menos uno de: **(a)** email (coincidencia parcial, case-insensitive), **(b)** student_id_code (coincidencia parcial o exacta), **(c)** nombre completo — first_name y/o last_name (coincidencia parcial con búsqueda trigram o ILIKE).
2. WHEN se envía `search` vacío o se omite, THE sistema SHALL comportarse como listado normal (devolver lista paginada sin filtrar por búsqueda).
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.
4. Paginación (`limit`, `offset`) SHALL aplicarse igual que en el listado (limit default 20, max 100; offset >= 0). La respuesta SHALL incluir `total` con el número de registros que cumplen el criterio de búsqueda.

### Requirement 9: Perfil del estudiante (Personal Stats, Loan History, Badge Gallery)

**User Story:** Como administrador o usuario de la aplicación, quiero ver el perfil completo de un estudiante con sus estadísticas personales, historial de préstamos e insignias para tener una vista unificada del estudiante (dashboard de perfil).

#### Acceptance Criteria

1. WHEN se solicita GET al endpoint de perfil del estudiante (ej. GET `/api/students/:id/profile`) para un ID existente y activo (`is_active = true`), THE sistema SHALL devolver 200 con un objeto que incluya: **(a)** datos básicos del estudiante para el encabezado del perfil (id, first_name, last_name, avatar_url, degree_level, major, member_since, department name/code si aplica — permitiendo mostrar "Undergraduate Student - Computer Science", "Member since Sept 2021"); **(b)** **Personal Stats:** total_read, active_loans, overdue_count, current_streak_days y opcionalmente longest_streak_days, total_points, global_rank (desde `student_stats`; si el estudiante no tiene fila en student_stats, devolver valores por defecto 0 o null según diseño); **(c)** **Loan History:** lista de préstamos del estudiante (orden reciente primero), cada uno con loan id, book (title, cover_url, authors como texto o array), borrowed_at, due_date, returned_at, status (active | returned | overdue). Se MAY limitar la cantidad devuelta en el perfil (ej. últimos 10 o 20) con opción de "View All" vía endpoint de loans por estudiante; **(d)** **Badge Gallery:** total de insignias (total_badges), insignias ganadas por el estudiante (earned_count) y lista de todas las insignias con para cada una: id, slug, name, icon_url, description/criteria, earned (boolean), earned_at (si earned es true). Así la UI puede mostrar progreso (ej. 8/12) y distinguir insignias ganadas (coloreadas) vs no ganadas (bloqueadas).
2. WHEN el ID no existe o el estudiante tiene `is_active = false`, THE sistema SHALL responder 404.
3. WHEN no hay token válido de administrador (o el perfil está protegido), THE sistema SHALL responder 401.
4. Loan history SHALL incluir datos del libro (título, portada, autores) mediante JOIN o consultas a `books` y tablas de autores; los préstamos devueltos (status = returned) y activos (status = active) SHALL poder distinguirse para mostrar "Currently Reading" y "Returned" en la interfaz.

### Requirement 10: Otorgar insignia a un estudiante

**User Story:** Como administrador, quiero otorgar (asignar) una insignia a un estudiante para que aparezca en su perfil y en la galería de insignias.

#### Acceptance Criteria

1. WHEN el administrador envía POST al endpoint de otorgar insignia (ej. POST `/api/students/:id/badges`) con body `{ "badge_id": "uuid" }` para un estudiante existente y activo y una insignia existente (`badges.id`), THE sistema SHALL insertar un registro en `student_badges` (student_id, badge_id, earned_at = NOW()) y devolver 201 (o 200 con el badge otorgado). Si el estudiante ya tiene esa insignia (combinación student_id + badge_id ya existe), THE sistema SHALL responder 409 con mensaje de conflicto o 200 idempotente según criterio de negocio.
2. WHEN el `student_id` no existe o el estudiante tiene `is_active = false`, THE sistema SHALL responder 404.
3. WHEN el `badge_id` no existe en la tabla `badges`, THE sistema SHALL responder 400 o 404 con mensaje de validación.
4. WHEN no se envía token válido de administrador, THE sistema SHALL responder 401.
5. La respuesta SHALL incluir al menos el badge otorgado (id, slug, name, earned_at) para mostrar en la UI.
