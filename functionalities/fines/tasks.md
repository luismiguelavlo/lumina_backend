# Overview

Plan de implementación para el módulo de multas (Fines). La tabla `fines` ya existe en `migrations/001_initial_schema.sql`. Se construye por capas: DTOs e interfaces → FineRepository → FineService → FineHandler. Endpoints: GET listado (filtros student_id, status), POST crear, GET :id, PATCH :id/paid, PATCH :id/waived. Todos bajo `/api` con JWT de staff (admin o bibliotecario), igual que el resto de rutas API. Stack: Go, Gin, PostgreSQL, go-playground/validator/v10.

# Tasks

- [x] 1. Modelo de datos y DTOs
  - [x] 1.1 Definir entidad `Fine` (id, loan_id, student_id, amount, status, reason, created_at, paid_at)
    - _Requirements: 1, 2, 3, 4, 5_
  - [x] 1.2 Definir DTOs: CreateFineRequest (loan_id, student_id, amount, reason); FineListItem; FineResponse; ListFinesResponse; FineFilter (StudentID, Status)
    - _Requirements: 1.1, 2.1, 3.1_
  - [x] 1.3 Definir errores: ErrFineNotFound, ErrFineNotPending, ErrLoanNotFound, ErrLoanStudentMismatch
    - _Requirements: 2.2, 4.2, 5.2, 6_
  - [x] 1.4 Definir interfaces FineRepository (Create, GetByID, List, UpdateStatus) y FineService (List, Create, GetByID, MarkPaid, MarkWaived); tipo FineWithStudent para resultado del repo
    - _Requirements: 1, 2, 3, 4, 5_

- [x] 2. Checkpoint – Tipos compilando

- [x] 3. FineRepository
  - [x] 3.1 Tests: Create; GetByID (existente, no existente); List con filtro student_id, status, limit/offset y total; UpdateStatus(id, paid, paid_at); UpdateStatus(id, waived, nil)
    - _Requirements: 1.1, 2.1, 3.1, 4.1, 5.1_
  - [x] 3.2 Implementar FineRepository: Create(INSERT); GetByID(SELECT); List(WHERE student_id, status, ORDER BY created_at DESC, LIMIT/OFFSET, COUNT con JOIN students para nombre); UpdateStatus(UPDATE status, paid_at WHERE id)
    - _Requirements: 1, 2, 3, 4, 5_
  - [x] 3.3 (Opcional) Property: solo pending puede actualizarse a paid/waived
    - **Property 1**
    - **Validates: Requirements 4.2, 5.2**

- [x] 4. Checkpoint – FineRepository implementado y tests pasando

- [x] 5. FineService
  - [x] 5.1 Tests: Create con loan existente y student_id coincidente → éxito; loan no existe → ErrLoanNotFound; loan.student_id != req.StudentID → ErrLoanStudentMismatch; List con filtros; GetByID; MarkPaid con pending → éxito (paid_at set); MarkPaid con ya paid → ErrFineNotPending; MarkWaived con pending → éxito; MarkWaived con ya waived → error
    - _Requirements: 2.1, 2.2, 3.1, 4.1, 4.2, 5.1, 5.2_
  - [x] 5.2 Implementar FineService: Create valida loan (LoanRepository.GetByID) y que loan.StudentID == req.StudentID, luego FineRepository.Create; List, GetByID delegan; MarkPaid/MarkWaived obtienen fine, si status != pending → ErrFineNotPending, sino UpdateStatus
    - _Requirements: 2, 3, 4, 5_
  - [x] 5.3 (Opcional) Property: MarkPaid asigna paid_at
    - **Property 2**
    - **Validates: Requirement 4.1**

- [x] 6. Checkpoint – Servicio implementado

- [x] 7. Manejo de errores HTTP
  - [x] 7.1 Mapear ErrFineNotFound → 404; ErrFineNotPending → 400; ErrLoanNotFound → 404; ErrLoanStudentMismatch → 400; helpers 400, 401, 404, 409, 500
    - _Requirements: 4.2, 5.2, 6_
  - [x] 7.2 Tests de handlers con códigos y payloads de error
    - _Requirements: 6_

- [x] 8. FineHandler – Listar multas
  - [x] 8.1 Tests: GET /api/fines con token → 200 con data y total; query student_id y status filtran; limit/offset normalizados; sin token → 401
    - _Requirements: 1.1, 1.2, 1.3, 1.4_
  - [x] 8.2 Implementar handler: parsear student_id, status, limit, offset; FineService.List; responder ListFinesResponse
    - _Requirements: 1_

- [x] 9. FineHandler – Crear multa
  - [x] 9.1 Tests: POST con loan_id, student_id, amount, reason válidos → 201; loan inexistente o student no coincide → 404/400; amount negativo → 400; sin token → 401
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_
  - [x] 9.2 Implementar handler: binding + validator; FineService.Create; 201 con FineResponse o 400/404
    - _Requirements: 2_

- [x] 10. FineHandler – Obtener por ID
  - [x] 10.1 Tests: GET /api/fines/:id existente → 200; inexistente → 404; sin token → 401
    - _Requirements: 3.1, 3.2, 3.3_
  - [x] 10.2 Implementar handler: FineService.GetByID; 200 o 404
    - _Requirements: 3_

- [x] 11. FineHandler – Marcar como pagada / condonada
  - [x] 11.1 Tests: PATCH /api/fines/:id/paid con multa pending → 200 (paid_at set); multa ya paid/waived → 400; id inexistente → 404; sin token → 401. PATCH /api/fines/:id/waived con pending → 200; ya paid/waived → 400; sin token → 401
    - _Requirements: 4.1, 4.2, 4.3, 5.1, 5.2, 5.3_
  - [x] 11.2 Implementar handlers: MarkPaid(ctx, id) y MarkWaived(ctx, id); 200 con FineResponse o 400/404
    - _Requirements: 4, 5_

- [x] 12. Integración y protección
  - [x] 12.1 Registrar rutas GET/POST /api/fines, GET /api/fines/:id, PATCH /api/fines/:id/paid, PATCH /api/fines/:id/waived con middleware JWT admin; inyectar FineService (y LoanRepository en el servicio para validar Create)
    - _Requirements: 1.3, 2.5, 3.3, 4.3, 5.3, 6.3_
  - [x] 12.2 Cobertura en paquete fines
    - _Requirements: 6_

- [x] 13. Checkpoint final – Build verde, listado con filtros, crear multa, GET por ID, marcar pagada, marcar condonada
