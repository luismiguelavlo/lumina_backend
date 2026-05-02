# Introduction

Motor de reputación y gamificación (Reputation Engine) de Lumina Library. **Servicio interno** (sin endpoints HTTP propios de negocio): otros módulos (hoy principalmente **Loan**) lo invocan para registrar eventos de reputación y actualizar las estadísticas del estudiante (`student_stats`). Con esto los datos alimentan **perfil del estudiante**, y en el futuro **ranking** y **analytics** cuando esos módulos expongan consultas sobre las mismas tablas.

**Código de referencia:** `internal/services/reputation/reputation_service.go`, persistencia `internal/repositories/reputation_store.go` + `reputation_store_gorm.go`.

# Estado de implementación (resumen)

| Área | Estado | Notas |
|------|--------|--------|
| **RecordReturn** (devolución) | Implementado | Evento `book_returned_on_time` / `book_returned_late`, idempotencia por `loan_id` en `reference_id`, actualización de `student_stats`. |
| **RecordLoanCreated** (nuevo préstamo) | Implementado | `active_loans += 1` (upsert de fila `student_stats` si no existe). |
| **Puntos on-time / late** | Implementado | Constantes en servicio: **+10** a tiempo, **-5** tarde (calendario UTC día vs día). |
| **Racha** | Implementado | A tiempo: `current_streak_days = prev + 1`, `longest_streak_days` no decrece; tarde: racha a **0**; `overdue_count` incrementa en devolución tardía. |
| **Activity log** (devolución) | Implementado | `LoanService` tras `RecordReturn` inserta fila vía `ActivityLogRepository` (`event_type` = `book_returned`, `actor_id` = JWT staff). Errores solo log; no revierten la devolución. |
| **Eventos** `badge_earned`, `fine_paid`, `sanction_received`, `streak_milestone` | Parcial (enum en BD) | Tipos definidos en PostgreSQL / `models`; **no** se crean aún desde servicios de badges, fines o sanciones (ver § Extensiones futuras). |
| **global_rank** | No actualizado en flujo préstamo | Puede calcularse en SQL/vista o job; no es obligatorio en `RecordReturn`. |

# Glossary

- **Reputation event:** Registro en la tabla `reputation_events`. Tipos (enum): `book_returned_on_time`, `book_returned_late`, `badge_earned`, `streak_milestone`, `fine_paid`, `sanction_received`. Cada evento tiene `points` (positivo o negativo), `description` opcional, `reference_id` (p. ej. `loan_id` para devoluciones).
- **Student stats:** Tabla `student_stats` (una fila por estudiante). Campos usados por el motor: `total_read`, `active_loans`, `overdue_count`, `current_streak_days`, `longest_streak_days`, `total_points`, `global_rank` (opcional), `updated_at`.
- **RecordReturn:** Operación del `ReputationService` cuando un préstamo ya está `returned` con `returned_at` informado. Compara **solo día calendario** en UTC entre `due_date` y `returned_at` para clasificar on-time vs late.
- **RecordLoanCreated:** Tras crear préstamo con éxito; solo ajusta `student_stats.active_loans` (no inserta fila en `reputation_events`).
- **Idempotencia (devolución):** Antes de insertar evento, se consulta si ya existe un `reputation_events` con `reference_id = loan.id` y tipo on_time o late; si existe, **no-op** (no duplicar puntos ni lecturas).

# Requirements

### Requirement 1: Registrar devolución (RecordReturn)

**User Story:** Como sistema, cuando un préstamo se marca como devuelto, quiero registrar el evento de reputación y actualizar las estadísticas del estudiante para que el ranking y el perfil reflejen la actividad.

#### Acceptance Criteria

1. WHEN el `LoanService` invoca `ReputationService.RecordReturn(ctx, loan)` con un préstamo en status `returned` y `returned_at` definido, THE servicio SHALL crear un registro en `reputation_events` con: `student_id` del préstamo, `event_type` = `book_returned_on_time` si el día calendario (UTC) de `returned_at` **no es posterior** al día de `due_date`, o `book_returned_late` en caso contrario; `points` +10 (on time) o -5 (late); `reference_id` = `loan.id`.
2. WHEN se procesa la devolución, THE persistencia SHALL actualizar `student_stats`: `total_read` +1, `active_loans` decrementado en 1 (mínimo 0), `total_points` ajustado con los `points` del evento, racha según criterio §Requirement 4; si la devolución es tardía, `overdue_count` incrementa en 1.
3. WHEN el estudiante no tiene fila en `student_stats`, THE primera devolución SHALL crear una fila coherente (lecturas, puntos, racha, `active_loans` tras decremento).
4. La operación SHALL ser **idempotente** respecto al mismo `loan_id`: si ya existe evento de devolución (on_time o late) con `reference_id` = ese préstamo, no insertar otro ni re-aplicar stats.

### Requirement 2: Integración con LoanService

**User Story:** Como desarrollador, quiero que el flujo de devolución y alta de préstamo dispare automáticamente la reputación.

#### Acceptance Criteria

1. WHEN `LoanService.Return` persiste la devolución con éxito, THE `LoanService` SHALL invocar `RecordReturn(ctx, loan)` con el préstamo ya actualizado (`status`, `returned_at`).
2. WHEN `RecordReturn` o la persistencia fallan (error de BD), THE devolución HTTP **no** se revierte: los errores se **loguean** (`log.Printf`); el préstamo queda devuelto (criterio no bloquear al usuario).
3. WHEN `LoanService.Create` persiste el préstamo con éxito, THE `LoanService` SHALL invocar `RecordLoanCreated(ctx, studentID)`.
4. WHEN `RecordLoanCreated` falla, THE creación del préstamo **no** se revierte; se loguea el error.
5. WHEN el cliente llama a `LoanService.Return`, THE handler SHALL pasar el identificador del usuario autenticado (staff) para `activity_log.actor_id`; si falta del contexto, respuesta **500** (misma convención que Create sin auth en contexto).

### Requirement 3: Activity log en devolución

**User Story:** Como sistema, quiero que las devoluciones queden en `activity_log` para el feed "Recent activity" del dashboard.

#### Acceptance Criteria

1. WHEN `LoanService.Return` completa con éxito la persistencia del préstamo devuelto y `RecordReturn`, THE servicio SHALL insertar (si hay repositorio configurado) una fila en `activity_log` con `event_type` = `book_returned`, `title` / `description` legibles, `student_id`, `actor_id` = usuario staff del JWT, y `metadata` JSON con al menos `loan_id`, `book_id`, `returned_on_time`, `returned_at_rfc3339`.
2. WHEN la inserción en `activity_log` falla, THE devolución HTTP **no** se revierte; se **loguea** el error (misma política que reputación).
3. **Implementación:** `internal/repositories/activity_log_repository.go` + `activity_log_repository_gorm.go`; inyección en `loan.NewService(..., activityLog)`; `main.go` usa `NewActivityLogRepositoryGorm`. El handler `PATCH /api/loans/:id/return` exige contexto JWT (como Create) para obtener `actor_id`.

### Requirement 4: Reglas de puntos y racha

**User Story:** Como producto, quiero puntos y rachas coherentes para gamificación.

#### Acceptance Criteria

1. Puntos por devolución: **+10** (`book_returned_on_time`), **-5** (`book_returned_late`) — constantes en `reputation_service.go` (`pointsReturnedOnTime`, `pointsReturnedLate`).
2. Racha: devolución a tiempo → `current_streak_days = previous + 1` (con `previous` normalizado a ≥0); `longest_streak_days = max(longest, new streak)`. Devolución tardía → `current_streak_days = 0` (no incrementa `longest` en ese paso).
3. `global_rank` puede omitirse en este flujo; ranking futuro puede usar `RANK()` sobre `total_points`.

### Requirement 5: Registrar creación de préstamo (RecordLoanCreated)

**User Story:** Como sistema, al crear un préstamo quiero incrementar `active_loans` del estudiante.

#### Acceptance Criteria

1. WHEN se invoca `RecordLoanCreated(ctx, studentID)` tras un alta exitosa de préstamo, THE store SHALL incrementar `student_stats.active_loans` en 1 (crear fila con `active_loans = 1` si no existe).
2. **Idempotencia:** no hay deduplicación por `loan_id` en esta ruta; si se llamara dos veces por error de integración, se sumaría dos veces — el llamador (`LoanService`) debe invocar una sola vez por creación exitosa.
3. WHEN falla el incremento, THE alta de préstamo sigue siendo válida; solo log.

### Requirement 6: Extensiones futuras (eventos y módulos)

Objetivo: documentar qué queda **fuera** del alcance actual pero preparado en enum / schema.

1. **Multas (`fines`):** tipo `fine_paid` — invocación futura desde el flujo de marcar multa pagada (o desde un servicio de multas) para sumar puntos y opcionalmente insertar `reputation_events`.
2. **Sanciones (`sanctions`):** tipo `sanction_received` — opcional al **aplicar** sanción activa (no al levantar), si el producto asigna penalización de puntos.
3. **Badges:** tipo `badge_earned` / `streak_milestone` — hoy las insignias se otorgan por API de estudiantes sin pasar por `ReputationService`; unificar aquí sería refactor opcional.
4. **Analytics / ranking:** leerán `student_stats` y `reputation_events` sin duplicar lógica de escritura en el motor de reputación.

# Glossary (complemento técnico)

- **ReputationStore:** interfaz única de persistencia (`ExistsReturnEventForLoan`, `CreateReputationEvent`, `IncrementStudentActiveLoans`, `ApplyStudentReturn`). Implementación GORM: `reputation_store_gorm.go`. Agrupa lo que en diseños antiguos se separaba en “repositorio de eventos” + “repositorio de stats”.
- **Transacciones:** `ApplyStudentReturn` y `IncrementStudentActiveLoans` usan transacción GORM donde aplica para consistencia lectura-escritura de `student_stats`.
