# Tasks — Reputation Engine (as-built vs backlog)

Estado alineado con el código en `internal/services/reputation` y `internal/repositories/reputation_store*.go`. Las tareas de “repositorios separados” del diseño inicial quedan **cumplidas** vía **`ReputationStore`** unificado.

## Implementado

- [x] **1. Modelos y constantes**
  - [x] `internal/models/reputation_event.go`, `student_stats.go`, `ReputationEventType` (enum + migración SQL).
  - [x] Puntos y tipos de evento de devolución en `reputation_service.go` (`pointsReturnedOnTime` / `pointsReturnedLate`).

- [x] **2. Migración PostgreSQL**
  - [x] Tablas `reputation_events`, `student_stats`; tipo `reputation_event_type`; índices (p. ej. `reference_id`).

- [x] **3–4. Persistencia unificada (`ReputationStore`)**
  - [x] `ExistsReturnEventForLoan` (idempotencia por `loan_id` en `reference_id`).
  - [x] `CreateReputationEvent`.
  - [x] `IncrementStudentActiveLoans` (alta de préstamo).
  - [x] `ApplyStudentReturn` (evento + actualización atómica de `student_stats`, lock, racha, `overdue_count` si tarde).
  - [x] Implementación GORM: `reputation_store_gorm.go`; in-memory si aplica en tests.

- [x] **5. `ReputationService`**
  - [x] `RecordReturn(ctx, loan)` — validación `returned` + fechas; delegación al store.
  - [x] `RecordLoanCreated(ctx, studentID)`.

- [x] **6. Integración `LoanService`**
  - [x] Interfaz `reputationHook` inyectada; llamadas tras `Create` y `Return` exitosos; errores solo log.

- [x] **7. Wiring**
  - [x] `main.go`: construcción de store + servicio + inyección en `NewLoanService`.

- [x] **8. Activity log (devolución)**
  - [x] `ActivityLogRepository` + GORM; `LoanService.Return` → `logBookReturned` tras `RecordReturn`; `actor_id` desde JWT en `LoanHandler`; tests en `loan_service_extended_test.go` / `loan_handler_test.go`.

## Backlog / opcional

- [ ] **9. Eventos desde otros dominios** — `fine_paid`, `sanction_received`, `badge_earned`, `streak_milestone` desde servicios de multas/sanciones/badges si el producto lo exige.

- [ ] **10. Ranking** — job o vista SQL que actualice / exponga `global_rank` sin duplicar lógica en `RecordReturn`.

## Documentación

- [x] **11. Documentación en `functionalities/reputation`** — requirements, design y tasks alineados con `ReputationStore` y flujos Loan (esta carpeta).
