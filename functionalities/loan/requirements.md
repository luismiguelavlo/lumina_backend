# Introduction

Loan Control Panel para Lumina Library. Permite a los administradores listar los préstamos activos (o filtrar por estado) mostrando nombre del libro, ID del libro, prestatario (borrower), fecha de vencimiento, tiempo restante e ISBN; crear nuevos préstamos indicando estudiante, libro y fecha de vencimiento; y registrar la devolución de un libro (marcar préstamo como returned). Los datos se almacenan en la tabla `loans` (book_id, student_id, issued_by, borrowed_at, due_date, returned_at, status). El prestatario es un estudiante (tabla `students`); el libro viene de `books`. Stack: Go, Gin, PostgreSQL, arquitectura en capas (handlers → services → repositories). Todos los endpoints protegidos con JWT (rol admin).

# Glossary

- **Loan (Préstamo):** Registro en la tabla `loans`. Relaciona un libro (`book_id`) con un estudiante (`student_id`). Campos: borrowed_at, due_date, returned_at, status (active | returned | overdue), issued_by (admin que registró el préstamo).
- **Borrower (Prestatario):** Estudiante que recibe el préstamo. Se identifica por `student_id` (FK a `students`). En listados se muestra el nombre (first_name, last_name).
- **Due date:** Fecha de vencimiento del préstamo. Columna `loans.due_date` (DATE). El estudiante debe devolver el libro en o antes de esta fecha.
- **Time remaining:** Días restantes hasta due_date. Se calcula como (due_date - fecha actual). Si due_date ya pasó, puede ser 0 o negativo; el préstamo puede considerarse overdue.
- **Status:** loan_status enum: `active` (prestado, no devuelto), `returned` (devuelto), `overdue` (vencido, no devuelto). El panel típicamente muestra préstamos activos; opcionalmente se puede filtrar por status.
- **issued_by:** UUID del admin/bibliotecario que creó el préstamo. FK a `users(id)`. Se asigna automáticamente desde el JWT.

# Requirements

### Requirement 1: Listar préstamos (Loan Control Panel)

**User Story:** Como administrador, quiero ver la lista de préstamos con nombre del libro, ID del libro, prestatario, fecha de vencimiento, tiempo restante e ISBN para gestionar los préstamos activos.

#### Acceptance Criteria

1. WHEN el administrador solicita GET al endpoint de préstamos con paginación (limit, offset), THE sistema SHALL devolver préstamos ordenados de forma consistente (ej. due_date ASC para ver primero los más urgentes). Cada elemento SHALL incluir: loan id, book id, book title (nombre del libro), borrower (nombre del estudiante: first_name + last_name o student_id_code según diseño), due_date, time_remaining (días restantes hasta due_date; 0 o negativo si ya venció), isbn del libro.
2. WHEN se envía el query param `status` (ej. active, overdue, returned), THE sistema SHALL filtrar por ese estado. Si no se envía, por defecto SHALL devolver solo préstamos con status `active` (o todos según criterio de diseño; se recomienda default `active` para el panel de control).
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.
4. WHEN limit u offset son inválidos, THE sistema SHALL normalizar (limit default 20, max 100; offset >= 0). La respuesta SHALL incluir el total de registros que cumplen el filtro.
5. WHEN se envía el query param opcional `student_id` (UUID), THE sistema SHALL filtrar los préstamos por ese estudiante, permitiendo listar "todos los préstamos de un estudiante" (p. ej. para "View All" en el perfil del estudiante).

### Requirement 2: Crear nuevo préstamo

**User Story:** Como administrador, quiero registrar un nuevo préstamo indicando el estudiante, el libro y la fecha de vencimiento para que quede registrado en el sistema.

#### Acceptance Criteria

1. WHEN el administrador envía POST con student_id, book_id y due_date válidos, THE sistema SHALL crear el préstamo con status `active`, borrowed_at = NOW(), issued_by = admin_id (extraído del JWT), y devolver 201 con los datos del préstamo creado (incluyendo al menos loan id, book, borrower, due_date, time_remaining, isbn si se desea consistencia con el listado).
2. WHEN student_id o book_id no existen en la base de datos, THE sistema SHALL responder 400 o 404 con mensaje claro (ej. "student not found", "book not found").
3. WHEN due_date es anterior a la fecha actual, THE sistema SHALL responder 400 con mensaje de validación (la fecha de vencimiento no puede ser en el pasado).
4. WHEN falta algún campo obligatorio o la validación falla, THE sistema SHALL responder 400 con objeto de errores de validación estructurado.
5. WHEN no se envía token válido de administrador, THE sistema SHALL responder 401.
6. WHEN el libro no tiene copias disponibles (todas prestadas), THE sistema SHALL responder 409 o 400 con mensaje de conflicto (opcional según regla de negocio; si se permite múltiples préstamos del mismo libro a distintos estudiantes mientras no se excedan copias, la validación se hace en servicio).
7. WHEN el estudiante (`student_id`) tiene al menos una sanción con `status = active` en la tabla `sanctions` (ver módulo `functionalities/sanctions`), THE sistema SHALL rechazar la creación del préstamo y SHALL responder **403 Forbidden** con mensaje claro (p. ej. que el estudiante tiene sanción activa y no puede recibir préstamos). WHEN todas las sanciones del estudiante están `lifted` o no hay filas de sanción, THE flujo de creación SHALL continuar con las demás validaciones.

### Requirement 3: Devolver libro (registrar devolución)

**User Story:** Como administrador, quiero registrar la devolución de un libro para que el préstamo se marque como devuelto y la copia quede disponible nuevamente.

#### Acceptance Criteria

1. WHEN el administrador envía PATCH (o POST) al endpoint de devolución para un loan_id existente con status `active` o `overdue`, THE sistema SHALL actualizar status = `returned`, returned_at = NOW() y devolver 200 con los datos del préstamo actualizado.
2. WHEN el loan_id no existe, THE sistema SHALL responder 404.
3. WHEN el préstamo ya tiene status `returned`, THE sistema SHALL responder 400 o 409 con mensaje "loan already returned".
4. WHEN no se envía token válido de administrador, THE sistema SHALL responder 401.
5. (Opcional) WHEN la devolución es de un préstamo que estaba en mora (returned_at > due_date), THE sistema SHALL crear automáticamente una multa (ver módulo fines) con loan_id, student_id, amount y reason según regla de negocio (ej. monto fijo o por día de retraso); si no se implementa, la multa se crea manualmente vía POST /api/fines.

### Requirement 4: Validación y errores

**User Story:** Como desarrollador, quiero respuestas HTTP consistentes (400, 401, 403, 404, 409, 500) y validación con mensajes por campo.

#### Acceptance Criteria

1. WHEN la entrada no cumple reglas de validación (due_date en el pasado, UUIDs inválidos, campos requeridos), THE sistema SHALL responder 400 con payload estructurado (message + errors por campo). WHEN el estudiante tiene sanción activa al crear préstamo, THE sistema SHALL responder **403** (ver criterio 7 del Requirement 2).
2. WHEN ocurre un error interno no esperado, THE sistema SHALL responder 500 con mensaje genérico.
3. WHEN se accede a cualquier endpoint de préstamos sin token válido de admin, THE sistema SHALL responder 401.
4. WHEN se intenta devolver un préstamo ya devuelto, THE sistema SHALL responder 400 o 409.

### Requirement 5: Job de sincronización `overdue` en base de datos

**User Story:** Como operador del sistema, quiero que los préstamos vigentes pasen automáticamente a estado `overdue` cuando la fecha de vencimiento (solo día, `due_date`) sea anterior a la fecha actual del calendario, para que listados, filtros y la vista de disponibilidad del catálogo reflejen la mora sin intervención manual.

#### Acceptance Criteria

1. THE sistema SHALL ejecutar periódicamente (proceso en segundo plano del servidor de aplicación, no obligatorio pg_cron) una actualización idempotente en PostgreSQL: préstamos con `status = 'active'` y `due_date < CURRENT_DATE` SHALL pasar a `status = 'overdue'`.
2. THE actualización SHALL ser segura de re-ejecutar (mismas filas no deben alternar de estado salvo que un operador las vuelva a `active` manualmente en BD, fuera de alcance API).
3. THE vista `book_availability` (o su definición equivalente en migración / `EnsureBookAvailabilityView`) SHALL contar como “prestadas” las copias en préstamos `active` **y** `overdue`, para que `checked_out` y disponibilidad sigan siendo correctos tras marcar mora.
4. WHEN el job falla (error de BD), THE proceso SHALL registrar el error en log y SHALL reintentar en el siguiente ciclo sin detener el servidor.
5. La frecuencia del job SHALL documentarse en `functionalities/loan/design.md` (valor por defecto razonable: ej. cada 1 hora, más una ejecución al arranque).
