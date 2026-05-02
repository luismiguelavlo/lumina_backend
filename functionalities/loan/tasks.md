# Overview

Plan de implementación para Loan Control Panel. La tabla `loans` y las tablas `books`, `students`, `users` ya existen en `migrations/001_initial_schema.sql`. Se construye por capas: DTOs e interfaces → LoanRepository (y reutilización de StudentRepository/BookRepository para validación) → LoanService → LoanHandler. Endpoints: GET listado de préstamos (book title, book id, borrower, due_date, time_remaining, isbn), POST crear préstamo (student_id, book_id, due_date), PATCH devolver libro (loan_id → status = returned). Todos los endpoints protegidos con JWT staff (admin/librarian). Stack: Go, Gin, PostgreSQL, go-playground/validator/v10.

**Estado:** implementación y tests principales completos (incl. `LoanRepositoryInMemory` para tests sin PostgreSQL). El listado 401 sin token lo cubre el mismo `StaffBearerMiddleware` que el resto de `/api`.

# Tasks

- [x] 1. Modelo de datos y DTOs
  - [x] 1.1 Definir entidad `Loan` con campos id, book_id, student_id, issued_by, borrowed_at, due_date, returned_at, status, notes, created_at, updated_at
    - _Requirements: 2.1_
  - [x] 1.2 Definir DTOs: CreateLoanRequest (student_id, book_id, due_date con binding required y formato fecha); LoanListItem (loan_id, book_id, book_title, borrower, due_date, time_remaining, isbn, status); LoanDetailResponse; ListLoansResponse (data, total)
    - _Requirements: 1.1, 2.1_
  - [x] 1.3 Definir errores de dominio: ErrLoanNotFound, ErrLoanAlreadyReturned, ErrStudentNotFound, ErrBookNotFound, ErrDueDateInPast; opcional ErrNoCopiesAvailable
    - _Requirements: 2.2, 2.3, 3, 4_
  - [x] 1.4 Definir LoanFilter (Status string), interfaces LoanRepository (Create, GetByID, List, Return) y LoanService (List, Create, Return); tipo LoanWithDetails o equivalente para resultado del repo con datos de book y student
    - _Requirements: 1, 2, 3_

- [x] 2. Checkpoint – Tipos compilando

- [x] 3. LoanRepository
  - [x] 3.1 Tests: Create inserta préstamo con status active, borrowed_at y issued_by; GetByID devuelve préstamo o nil; List con status=active devuelve solo activos, ordenados por due_date ASC; List con limit/offset y total; List con status returned/overdue filtra correctamente; Return actualiza status='returned' y returned_at; Return con id inexistente devuelve error; Return con préstamo ya devuelto devuelve error — _vía `loan_repository_memory_test.go`_
    - _Requirements: 1.1, 1.2, 2.1, 3.1, 3.2, 3.3_
  - [x] 3.2 Implementar LoanRepository: Create(loan); GetByID(id); List(filter, limit, offset) con JOIN a books y students, WHERE por status y opcionalmente student_id, ORDER BY due_date ASC, LIMIT/OFFSET y COUNT; Return(id, returnedAt) UPDATE status='returned', returned_at WHERE id = $1 AND status IN ('active','overdue') — _implementación GORM: `loan_repository_gorm.go`_
    - _Requirements: 1.1, 2.1, 3.1_
  - [ ] 3.3 (Opcional) Property: listado solo incluye préstamos del status solicitado
    - **Validates: Requirement 1.2**

- [x] 4. Checkpoint – LoanRepository implementado y tests pasando

- [x] 5. LoanService
  - [x] 5.1 Tests: List delega al repo y calcula time_remaining por ítem; Create con datos válidos → éxito; Create con student_id inexistente → ErrStudentNotFound; Create con book_id inexistente → ErrBookNotFound; Create con due_date en el pasado → ErrDueDateInPast; Return con loan activo → éxito (status returned, returned_at set); Return con loan inexistente → ErrLoanNotFound; Return con loan ya devuelto → ErrLoanAlreadyReturned — _vía `loan_service_test.go` + `loan_service_extended_test.go`_
    - _Requirements: 1.1, 2.1, 2.2, 2.3, 3.1, 3.2, 3.3_
  - [x] 5.2 Implementar LoanService: List llama repo, mapea a LoanListItem calculando time_remaining; Create valida student, book, due_date >= hoy, luego LoanRepository.Create; Return: GetByID(loanID), si no existe → ErrLoanNotFound, si status == 'returned' → ErrLoanAlreadyReturned, si active/overdue → LoanRepository.Return(id, time.Now())
    - _Requirements: 1, 2, 3_
  - [x] 5.3 (Opcional) Property: due_date en pasado rechazado; issued_by asignado desde adminID — _cubierto en tests `CreateDueDatePast`, `CreateSuccessAndIssuedBy`_
    - **Property 1 y 3: due_date no pasado; issued_by desde JWT**
    - **Validates: Requirements 2.3, 2.1**

- [x] 6. Checkpoint – Servicio implementado

- [x] 7. Manejo de errores HTTP
  - [x] 7.1 Mapear ErrStudentNotFound y ErrBookNotFound → 404; ErrLoanNotFound → 404; ErrDueDateInPast → 400; ErrLoanAlreadyReturned → 409; ErrNoCopiesAvailable → 409; helpers 400 con errors por campo, 401, 404, 409, 500
    - _Requirements: 2.2, 2.3, 3.2, 3.3, 4.1, 4.4_
  - [x] 7.2 Tests de handlers que comprueben 400 (due_date pasado, validación), 401 sin token, 404 (student/book not found)
    - _Requirements: 3_

- [x] 8. LoanHandler – Listar préstamos
  - [x] 8.1 Tests: GET /api/loans con token admin → 200 con data y total; cada ítem tiene loan_id, book_id, book_title, borrower, due_date, time_remaining, isbn, status; query status=active (o default) filtra; limit/offset se normalizan; sin token → 401 — _200/errores en `loan_handler_test.go`; 401 list mismo middleware global_
    - _Requirements: 1.1, 1.2, 1.3, 1.4_
  - [x] 8.2 Implementar handler: parsear query params status (default "active"), student_id (opcional UUID), limit (default 20, max 100), offset; LoanService.List; responder ListLoansResponse
    - _Requirements: 1, 1.5_
  - [ ] 8.3 (Opcional) Property: listado incluye todos los campos requeridos
    - **Property 2 y 4: campos y time_remaining**
    - **Validates: Requirement 1.1**

- [x] 9. LoanHandler – Crear préstamo
  - [x] 9.1 Tests: POST con student_id, book_id, due_date válidos (due_date >= hoy) → 201; student_id inexistente → 404; book_id inexistente → 404; due_date en el pasado → 400; validación falla → 400; sin token → 401
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_
  - [x] 9.2 Implementar handler: binding + validator; validación de formato fecha (due_date); extraer adminID del JWT; LoanService.Create(ctx, req, adminID); 201 con LoanDetailResponse o 400/404/409
    - _Requirements: 2_

- [x] 9b. LoanHandler – Devolver libro
  - [x] 9b.1 Tests: PATCH /api/loans/:id/return con loan activo → 200; loan inexistente → 404; loan ya devuelto → 409; sin token → 401 — _401 mismo patrón middleware_
    - _Requirements: 3.1, 3.2, 3.3, 3.4_
  - [x] 9b.2 Implementar handler: extraer loan_id de la URL; `LoanService.Return(ctx, loanID, actorUserID)` con `actorUserID` desde JWT (mismo contexto que Create); 200 con LoanDetailResponse o 404/409/500 si falta usuario en contexto
    - _Requirements: 3_
  - [ ] 9b.3 (Opcional) Property: no devolver dos veces
    - **Property 6: No devolver dos veces**
    - **Validates: Requirement 3.3**

- [x] 10. Integración y protección
  - [x] 10.1 Registrar rutas GET /api/loans, POST /api/loans y PATCH /api/loans/:id/return con middleware JWT admin; inyectar LoanService (y ReputationService en el servicio para RecordReturn)
    - _Requirements: 1.3, 2.5, 3.4, 4.3_
  - [x] 10.1b Integración con reputación: tras Create préstamo exitoso, invocar ReputationService.RecordLoanCreated(ctx, student_id). Tras Return exitoso, invocar ReputationService.RecordReturn(ctx, loan). Si fallan, loguear y no fallar la operación (ver functionalities/reputation)
  - [x] 10.1c **Activity log:** tras `RecordReturn`, `ActivityLogRepository.Create` con `event_type` = `book_returned` (ver `functionalities/reputation` Req. 3); tests `loan_service_extended_test.go`, `loan_handler_test.go` (JWT en Return)
    - _Requirements: 2.1 de reputation, 5 de reputation_
  - [x] 10.2 Cobertura en paquete loan; corregir huecos
    - _Requirements: 3_

- [x] 11. Checkpoint final – Build verde, GET listado, POST crear préstamo, PATCH devolver libro (status returned, returned_at set, copia disponible de nuevo)

- [x] 12. Job sincronización `overdue`
  - [x] 12.1 `LoanRepository.MarkActiveLoansOverdue`: UPDATE `active` → `overdue` donde `due_date < CURRENT_DATE`
  - [x] 12.2 Job en `internal/jobs` (p. ej. `LoanOverdueSync`) invocado desde `main.go` con ticker (ej. 1h + run al arranque); errores solo log
  - [x] 12.3 Vista `book_availability`: subconsulta de préstamos cuenta `status IN ('active','overdue')` en `EnsureBookAvailabilityView` y en `migrations/001_initial_schema.sql`
  - [x] 12.4 Documentar en `functionalities/loan` (requirements, design) y fila del módulo en `functionalities/README.md`
    - _Requirements: 5 de loan_

---

## Pendiente opcional (producto / fase posterior)

- Multa automática al devolver tarde (módulo **fines**).
- Tests E2E contra PostgreSQL real o property-based tests (ítems 3.3, 8.3, 9b.3).
