# Introduction

Gestión de multas (Fines Management) de Lumina Library. Permite a los administradores listar multas (por estudiante o por estado), crear una multa asociada a un préstamo (p. ej. por devolución tardía), consultar una multa por ID, y marcar una multa como pagada o como condonada (waived). Los datos se almacenan en la tabla `fines` (loan_id, student_id, amount, status, reason, paid_at). El status es `pending`, `paid` o `waived`. El dashboard de analytics usa la suma de multas pendientes (overdue_fines). Stack: Go, Gin, PostgreSQL, arquitectura en capas (handlers → services → repositories). Todos los endpoints protegidos con JWT (rol admin).

# Glossary

- **Fine (Multa):** Registro en la tabla `fines`. Vinculada a un préstamo (`loan_id`) y a un estudiante (`student_id`). Campos: amount (monto), status (pending | paid | waived), reason (motivo, opcional), paid_at (fecha de pago si status = paid).
- **pending:** Multa pendiente de pago. Cuenta para la métrica "overdue fines" del dashboard.
- **paid:** Multa pagada; paid_at = fecha en que se marcó como pagada.
- **waived:** Multa condonada (perdonada) por el administrador; no requiere pago.
- **loan_id:** Préstamo que originó la multa (p. ej. devolución tardía). FK a `loans(id)` ON DELETE CASCADE.

# Requirements

### Requirement 1: Listar multas

**User Story:** Como administrador, quiero listar las multas con filtros por estudiante y por estado para gestionar cobros y condonaciones.

#### Acceptance Criteria

1. WHEN el administrador solicita GET al endpoint de multas con paginación (limit, offset), THE sistema SHALL devolver multas ordenadas de forma consistente (ej. created_at DESC), con el total de registros. Cada elemento SHALL incluir: id, loan_id, student_id, amount, status, reason, created_at, paid_at (si aplica), y opcionalmente datos del estudiante (nombre) o del préstamo según diseño.
2. WHEN se envía el query param `student_id` (UUID), THE sistema SHALL filtrar por ese estudiante. WHEN se envía `status` (pending, paid, waived), THE sistema SHALL filtrar por ese estado.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.
4. WHEN limit u offset son inválidos, THE sistema SHALL normalizar (limit default 20, max 100; offset >= 0). La respuesta SHALL incluir el total de registros que cumplen el filtro.

### Requirement 2: Crear multa

**User Story:** Como administrador, quiero registrar una multa asociada a un préstamo y un estudiante (p. ej. por devolución tardía) para que quede pendiente de cobro.

#### Acceptance Criteria

1. WHEN el administrador envía POST con loan_id, student_id, amount (>= 0) y opcionalmente reason, THE sistema SHALL crear la multa con status `pending` y devolver 201 con los datos de la multa creada.
2. WHEN loan_id o student_id no existen, o el loan no pertenece al student_id, THE sistema SHALL responder 400 o 404 con mensaje claro.
3. WHEN amount es negativo, THE sistema SHALL responder 400 con objeto de errores de validación.
4. WHEN falta algún campo obligatorio (loan_id, student_id, amount), THE sistema SHALL responder 400 con payload estructurado.
5. WHEN no se envía token válido de administrador, THE sistema SHALL responder 401.
6. (Opcional) WHEN ya existe una multa pendiente para el mismo loan_id, THE sistema SHALL responder 409 o permitir según regla de negocio; el diseño SHALL documentar el criterio.

### Requirement 3: Obtener multa por ID

**User Story:** Como administrador, quiero ver el detalle de una multa.

#### Acceptance Criteria

1. WHEN el ID existe, THE sistema SHALL devolver 200 con los datos de la multa (id, loan_id, student_id, amount, status, reason, created_at, paid_at) y opcionalmente datos del estudiante y del préstamo.
2. WHEN el ID no existe, THE sistema SHALL responder 404.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.

### Requirement 4: Marcar multa como pagada

**User Story:** Como administrador, quiero marcar una multa como pagada cuando el estudiante realiza el pago.

#### Acceptance Criteria

1. WHEN el administrador envía PATCH al endpoint de marcar como pagada (ej. PATCH /api/fines/:id/paid) para un id de multa existente con status `pending`, THE sistema SHALL actualizar status = `paid`, paid_at = NOW() y devolver 200 con los datos actualizados.
2. WHEN el id no existe o la multa ya está paid o waived, THE sistema SHALL responder 404 o 400 según diseño.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.

### Requirement 5: Marcar multa como condonada (waived)

**User Story:** Como administrador, quiero condonar una multa para que el estudiante no deba pagarla.

#### Acceptance Criteria

1. WHEN el administrador envía PATCH al endpoint de condonar (ej. PATCH /api/fines/:id/waived) para un id de multa existente con status `pending`, THE sistema SHALL actualizar status = `waived` y devolver 200 (paid_at permanece null).
2. WHEN el id no existe o la multa ya está paid o waived, THE sistema SHALL responder 404 o 400.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.

### Requirement 6: Validación y errores

**User Story:** Como desarrollador, quiero respuestas HTTP consistentes (400, 401, 404, 409, 500).

#### Acceptance Criteria

1. WHEN la entrada no cumple reglas de validación (amount negativo, UUIDs inválidos), THE sistema SHALL responder 400 con payload estructurado.
2. WHEN ocurre un error interno no esperado, THE sistema SHALL responder 500 con mensaje genérico.
3. WHEN se accede a cualquier endpoint de multas sin token válido de admin, THE sistema SHALL responder 401.
