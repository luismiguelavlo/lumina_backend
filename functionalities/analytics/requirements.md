# Introduction

Módulo de analíticas (dashboard) para administradores de Lumina Library. Expone métricas agregadas de solo lectura: **total books**, **active students**, **overdue fines** (suma de multas pendientes), **most borrowed books** (top N títulos más prestados) y **recent activity** (últimos eventos del `activity_log`). Los datos se derivan de las tablas existentes `books`, `students`, `fines`, `loans` y `activity_log`. Solo usuarios autenticados con rol admin pueden acceder. Stack: Go, Gin, PostgreSQL, arquitectura en capas (handlers → services → repositories). Sin escritura en BD desde este módulo; solo consultas de lectura.

# Glossary

- **Total books:** Número total de registros en la tabla `books` (títulos/copias catalogados; `total_copies` por libro se usa para inventario, pero la métrica "total books" se refiere al conteo de filas en `books` o, según criterio de negocio, a la suma de `books.total_copies`).
- **Active students:** Estudiantes con `is_active = true` en la tabla `students`. Conteo de filas que cumplen esta condición.
- **Overdue fines:** Suma monetaria de multas pendientes. Suma de `fines.amount` donde `fines.status = 'pending'`. Representa el total adeudado por mora u otras razones no pagadas.
- **Most borrowed books:** Libros ordenados por número de préstamos (cuenta de filas en `loans` por `book_id`). Se devuelve un top N (ej. 5 o 10) con título, autor(es) y cantidad de préstamos. Se consideran todos los préstamos (activos y devueltos) para la cuenta.
- **Recent activity:** Últimas entradas de la tabla `activity_log`, ordenadas por `created_at DESC`, con límite configurable (ej. 10 o 20). Cada entrada incluye tipo de evento, título, descripción, metadatos opcionales y fecha; opcionalmente nombres de actor (admin) y estudiante relacionado.
- **Dashboard:** Vista consolidada que muestra las cinco métricas anteriores en una sola respuesta o en endpoints separados, según diseño.

# Requirements

### Requirement 1: Total books

**User Story:** Como administrador, quiero ver el número total de libros en el catálogo para tener una visión del tamaño de la colección.

#### Acceptance Criteria

1. WHEN el administrador solicita la métrica de total books (como parte del dashboard o en un endpoint dedicado), THE sistema SHALL devolver el conteo de filas en la tabla `books` (o la suma de `total_copies` si se define como criterio de negocio; por defecto conteo de filas).
2. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.
3. WHEN la consulta se ejecuta correctamente, THE respuesta SHALL incluir un valor numérico entero (total_books).

### Requirement 2: Active students

**User Story:** Como administrador, quiero ver cuántos estudiantes activos hay registrados en el sistema.

#### Acceptance Criteria

1. WHEN el administrador solicita la métrica de active students, THE sistema SHALL devolver el conteo de filas en `students` donde `is_active = true`.
2. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.
3. WHEN la consulta se ejecuta correctamente, THE respuesta SHALL incluir un valor numérico entero (active_students).

### Requirement 3: Overdue fines

**User Story:** Como administrador, quiero ver el monto total de multas pendientes de pago (overdue fines).

#### Acceptance Criteria

1. WHEN el administrador solicita la métrica de overdue fines, THE sistema SHALL devolver la suma de `fines.amount` para todas las filas en `fines` donde `status = 'pending'`. Si no hay multas pendientes, THE suma SHALL ser 0.
2. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.
3. WHEN la consulta se ejecuta correctamente, THE respuesta SHALL incluir un valor numérico (decimal o entero según tipo en BD) con dos decimales para representar moneda (overdue_fines).

### Requirement 4: Most borrowed books

**User Story:** Como administrador, quiero ver los libros más prestados (top N) para conocer la demanda del catálogo.

#### Acceptance Criteria

1. WHEN el administrador solicita most borrowed books con un límite N (ej. 5 o 10), THE sistema SHALL devolver una lista ordenada por número de préstamos (count de `loans` por `book_id`) descendente, con hasta N elementos. Cada elemento SHALL incluir al menos: identificador del libro, título, cantidad de préstamos; y opcionalmente autor(es) o catalog_code.
2. WHEN el límite N no se envía o es inválido, THE sistema SHALL usar un valor por defecto (ej. 5) o un máximo (ej. 20) según diseño.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.
4. WHEN no hay préstamos, THE sistema SHALL devolver una lista vacía o los libros con 0 préstamos según criterio de diseño (recomendado: solo libros con al menos un préstamo, ordenados por count descendente).

### Requirement 5: Recent activity

**User Story:** Como administrador, quiero ver las actividades recientes del sistema (devoluciones, altas de estudiantes, avisos de mora, etc.) para tener un resumen de lo que ha ocurrido.

#### Acceptance Criteria

1. WHEN el administrador solicita recent activity con un límite N (ej. 10 o 20), THE sistema SHALL devolver las últimas N filas de `activity_log` ordenadas por `created_at DESC`, incluyendo al menos: id, event_type, title, description, created_at; y opcionalmente metadata, actor_id/actor_name, student_id/student_name.
2. WHEN el límite N no se envía o es inválido, THE sistema SHALL usar un valor por defecto (ej. 10) o un máximo (ej. 50) según diseño.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.
4. WHEN no hay registros en activity_log, THE sistema SHALL devolver una lista vacía.

### Requirement 6: Endpoint(s) de dashboard y errores

**User Story:** Como desarrollador, quiero un contrato API claro y respuestas HTTP consistentes para las analíticas.

#### Acceptance Criteria

1. WHEN el diseño define un único endpoint de dashboard (ej. GET /api/analytics/dashboard), THE respuesta SHALL incluir al menos total_books, active_students, overdue_fines, most_borrowed_books (array) y recent_activity (array), con códigos 200 o 401/500 según corresponda.
2. WHEN ocurre un error interno (ej. fallo de BD), THE sistema SHALL responder 500 con mensaje genérico y no exponer detalles internos.
3. WHEN se accede sin token válido de administrador, THE sistema SHALL responder 401 en cualquier endpoint de analíticas.
