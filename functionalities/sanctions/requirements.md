# Introduction

Gestión de sanciones (Sanctions Management) de Lumina Library. Permite a los administradores listar los estudiantes que están actualmente sancionados, agregar un estudiante a la lista de sanciones (crear sanción con motivo) y levantar una sanción (quitar al estudiante de la lista de sancionados). Los datos se almacenan en la tabla `sanctions` (student_id, reason, status, applied_at, applied_by, lifted_at, lifted_by). El status es `active` (sancionado) o `lifted` (sanción levantada). Solo se listan como “sancionados” los registros con status `active`. Stack: Go, Gin, PostgreSQL, arquitectura en capas (handlers → services → repositories). Todos los endpoints protegidos con JWT (rol admin).

# Glossary

- **Sanction (Sanción):** Registro en la tabla `sanctions`. Representa una sanción aplicada a un estudiante. Campos: student_id (FK a students), reason (motivo), status (active | lifted), applied_at, applied_by (admin que aplicó), lifted_at, lifted_by (admin que levantó).
- **Estudiante sancionado:** Estudiante que tiene al menos una sanción con status `active`. Aparece en el listado de sancionados.
- **Aplicar sanción:** Crear un nuevo registro en `sanctions` con status `active`, reason obligatorio, applied_by = admin del JWT, applied_at = NOW().
- **Levantar sanción:** Marcar una sanción activa como levantada: status = `lifted`, lifted_at = NOW(), lifted_by = admin del JWT. El estudiante deja de aparecer en el listado de sancionados (para esa sanción).
- **applied_by / lifted_by:** UUID del usuario administrador que aplicó o levantó la sanción. FK a `users(id)`.

# Requirements

### Requirement 1: Listar estudiantes sancionados

**User Story:** Como administrador, quiero ver la lista de estudiantes que están actualmente sancionados para gestionar las sanciones activas.

#### Acceptance Criteria

1. WHEN el administrador solicita GET al endpoint de sanciones con paginación (limit, offset), THE sistema SHALL devolver solo sanciones con status `active`, ordenadas de forma consistente (ej. applied_at DESC), con el total de registros. Cada elemento SHALL incluir datos de la sanción (id, reason, applied_at, applied_by si está disponible) y datos del estudiante (id, student_id_code, first_name, last_name, email según diseño) para identificar al estudiante sancionado.
2. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.
3. WHEN limit u offset son inválidos, THE sistema SHALL normalizar (limit default 20, max 100; offset >= 0). La respuesta SHALL incluir el total de registros.

### Requirement 2: Agregar estudiante a la lista de sanciones

**User Story:** Como administrador, quiero agregar un estudiante a la lista de sanciones indicando el motivo para registrar la sanción.

#### Acceptance Criteria

1. WHEN el administrador envía POST con student_id y reason válidos, THE sistema SHALL crear la sanción con status `active`, applied_at = NOW(), applied_by = admin_id (extraído del JWT) y devolver 201 con los datos de la sanción creada (incluyendo datos del estudiante si se desea consistencia con el listado).
2. WHEN student_id no existe en la tabla students, THE sistema SHALL responder 404 con mensaje claro (ej. "student not found").
3. WHEN falta reason o está vacío, THE sistema SHALL responder 400 con objeto de errores de validación estructurado.
4. WHEN no se envía token válido de administrador, THE sistema SHALL responder 401.
5. (Opcional) WHEN el estudiante ya tiene una sanción activa, THE sistema SHALL responder 409 o permitir múltiples sanciones activas según regla de negocio; el diseño SHALL documentar el criterio elegido.

### Requirement 3: Quitar estudiante de la lista de sanciones (levantar sanción)

**User Story:** Como administrador, quiero levantar una sanción para que el estudiante deje de figurar como sancionado.

#### Acceptance Criteria

1. WHEN el administrador solicita levantar la sanción (ej. PATCH o POST al recurso de la sanción) para un id de sanción existente con status `active`, THE sistema SHALL actualizar status = `lifted`, lifted_at = NOW(), lifted_by = admin_id (extraído del JWT) y devolver 200 (o 204 según diseño).
2. WHEN el id de sanción no existe o la sanción ya está lifted, THE sistema SHALL responder 404.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.

### Requirement 5: Integración con préstamos (bloqueo al prestar)

**User Story:** Como operador, quiero que un estudiante con sanción activa no pueda recibir nuevos préstamos hasta que se levanten sus sanciones, para alinear la restricción disciplinaria con el flujo de préstamos.

#### Acceptance Criteria

1. WHEN el sistema procesa **POST /api/loans** (crear préstamo) y el `student_id` tiene al menos una fila en `sanctions` con `status = active`, THE sistema SHALL rechazar la operación con **403 Forbidden** y un mensaje que indique que el estudiante tiene sanción activa (sin crear fila en `loans`).
2. WHEN no hay sanciones activas para ese estudiante (ninguna fila o todas `lifted`), THE creación de préstamo SHALL seguir las reglas del módulo loan (copias disponibles, fechas, etc.).
3. La regla SHALL documentarse en `functionalities/loan/requirements.md` y `functionalities/loan/design.md` (referencia cruzada con este módulo).

### Requirement 4: Validación y errores

**User Story:** Como desarrollador, quiero respuestas HTTP consistentes (400, 401, 404, 409, 500) y validación con mensajes por campo.

#### Acceptance Criteria

1. WHEN la entrada no cumple reglas de validación (reason vacío, student_id inválido), THE sistema SHALL responder 400 con payload estructurado (message + errors por campo).
2. WHEN ocurre un error interno no esperado, THE sistema SHALL responder 500 con mensaje genérico.
3. WHEN se accede a cualquier endpoint de sanciones sin token válido de admin, THE sistema SHALL responder 401.
