# Overview

Plan de implementación para el módulo de sanciones. La tabla `sanctions` y la tabla `students` ya existen en `migrations/001_initial_schema.sql`. Se construye por capas: DTOs e interfaces → SanctionRepository → SanctionService → SanctionHandler. Endpoints: GET listado de estudiantes sancionados (sanciones activas con datos del estudiante), POST crear sanción (student_id, reason), PATCH levantar sanción (por id de sanción). Todos bajo `/api` con JWT de staff (admin o bibliotecario). Stack: Go, Gin, PostgreSQL, go-playground/validator/v10.

# Tasks

- [x] 1. Modelo de datos y DTOs
  - [x] 1.1 Definir entidad `Sanction` con campos id, student_id, reason, status, applied_at, applied_by, lifted_at, lifted_by, created_at
    - _Requirements: 1, 2, 3_
  - [x] 1.2 Definir DTOs: CreateSanctionRequest (student_id, reason con binding required); SanctionListItem (sanction_id, student_id, student_id_code, first_name, last_name, email, reason, applied_at, applied_by_id); SanctionResponse; ListSanctionsResponse (data, total); LiftSanctionRequest opcional (status) para PATCH
    - _Requirements: 1.1, 2.1, 3.1_
  - [x] 1.3 Definir errores de dominio: ErrSanctionNotFound, ErrStudentNotFound, ErrSanctionAlreadyLifted
    - _Requirements: 2.2, 3.2, 4_
  - [x] 1.4 Definir interfaces SanctionRepository (Create, GetByID, ListActive, Lift) y SanctionService (ListActive, Create, Lift); tipo SanctionWithStudent para resultado del repo con datos del estudiante
    - _Requirements: 1, 2, 3_

- [x] 2. Checkpoint – Tipos compilando

- [x] 3. SanctionRepository
  - [x] 3.1 Tests: Create inserta con status active, applied_at y applied_by; GetByID (existente activo, existente lifted, no existente); ListActive devuelve solo status active, con JOIN students, orden applied_at DESC, limit/offset y total; Lift actualiza status, lifted_at, lifted_by (éxito; id inexistente o ya lifted → error o 0 rows)
    - _Requirements: 1.1, 2.1, 3.1, 3.2_
  - [x] 3.2 Implementar SanctionRepository: Create(INSERT); GetByID(SELECT por id); ListActive(SELECT con JOIN students WHERE status = 'active', ORDER BY applied_at DESC, LIMIT/OFFSET, COUNT); Lift(UPDATE status = 'lifted', lifted_at = NOW(), lifted_by WHERE id = $1 AND status = 'active')
    - _Requirements: 1, 2, 3_
  - [x] 3.3 (Opcional) Property: listado solo incluye sanciones con status active
    - **Property 1: Listado solo sanciones activas**
    - **Validates: Requirement 1.1**

- [x] 4. Checkpoint – SanctionRepository implementado y tests pasando

- [x] 5. SanctionService
  - [x] 5.1 Tests: ListActive delega al repo y mapea a SanctionListItem; Create con student_id y reason válidos → éxito, applied_by asignado; Create con student_id inexistente → ErrStudentNotFound; Lift con id existente activo → éxito, lifted_by asignado; Lift con id inexistente o ya lifted → ErrSanctionNotFound
    - _Requirements: 1.1, 2.1, 2.2, 3.1, 3.2_
  - [x] 5.2 Implementar SanctionService: ListActive llama repo ListActive, mapea a SanctionListItem; Create valida que student exista (StudentRepository.GetByID), luego SanctionRepository.Create con applied_at = NOW(), status = 'active', applied_by = adminID; Lift obtiene sanción por ID, si no existe o status != 'active' devolver ErrSanctionNotFound, sino SanctionRepository.Lift(id, time.Now(), adminID)
    - _Requirements: 1, 2, 3_
  - [x] 5.3 (Opcional) Property: applied_by y lifted_by asignados desde adminID; no levantar dos veces
    - **Property 2, 3, 4**
    - **Validates: Requirements 2.1, 3.1, 3.2**

- [x] 6. Checkpoint – Servicio implementado

- [x] 7. Manejo de errores HTTP
  - [x] 7.1 Mapear ErrSanctionNotFound y ErrStudentNotFound → 404; ErrSanctionAlreadyLifted → 404; helpers 400 con errors por campo, 401, 404, 500
    - _Requirements: 2.2, 3.2, 4.1, 4.3_
  - [x] 7.2 Tests de handlers que comprueben 400 (validación), 401 sin token, 404 (student not found, sanction not found / already lifted)
    - _Requirements: 4_

- [x] 8. SanctionHandler – Listar sancionados
  - [x] 8.1 Tests: GET /api/sanctions con token admin → 200 con data y total; cada ítem tiene sanction_id, student_id, student_id_code, first_name, last_name, reason, applied_at; limit/offset normalizados; sin token → 401
    - _Requirements: 1.1, 1.2, 1.3_
  - [x] 8.2 Implementar handler: parsear limit (default 20, max 100), offset (default 0); SanctionService.ListActive; responder ListSanctionsResponse
    - _Requirements: 1_

- [x] 9. SanctionHandler – Crear sanción (agregar a lista)
  - [x] 9.1 Tests: POST con student_id y reason válidos → 201; student_id inexistente → 404; reason vacío → 400; sin token → 401
    - _Requirements: 2.1, 2.2, 2.3, 2.4_
  - [x] 9.2 Implementar handler: binding + validator; extraer adminID del JWT; SanctionService.Create(ctx, req, adminID); 201 con SanctionResponse o 400/404
    - _Requirements: 2_

- [x] 10. SanctionHandler – Levantar sanción (quitar de lista)
  - [x] 10.1 Tests: PATCH /api/sanctions/:id con body status=lifted (o sin body si el diseño es solo levantar) → 200; id inexistente o sanción ya lifted → 404; sin token → 401
    - _Requirements: 3.1, 3.2, 3.3_
  - [x] 10.2 Implementar handler: extraer sanction id de la URL; extraer adminID del JWT; SanctionService.Lift(ctx, id, adminID); 200 o 404
    - _Requirements: 3_

- [x] 11. Integración y protección
  - [x] 11.1 Registrar rutas GET /api/sanctions, POST /api/sanctions, PATCH /api/sanctions/:id con middleware JWT admin; inyectar SanctionService (y en el servicio SanctionRepository, StudentRepository)
    - _Requirements: 1.2, 2.4, 3.3, 4.3_
  - [x] 11.2 Cobertura en paquete sanctions; corregir huecos
    - _Requirements: 4_

- [x] 12. Checkpoint final – Build verde, GET listado de sancionados, POST crear sanción, PATCH levantar sanción

- [x] 13. Integración con préstamos — `LoanService.Create` consulta sanciones `active`; **403** si aplica. Documentado en `sanctions/requirements.md` (Req. 5), `sanctions/design.md`, `loan/requirements.md`, `loan/design.md`, `functionalities/README.md`; Postman *LOANS MODULE* actualizado.
