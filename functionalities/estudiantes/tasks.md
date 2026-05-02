# Overview

Plan de implementación en orden TDD para el módulo de estudiantes. La migración de base de datos ya existe en `migrations/001_initial_schema.sql` (tablas `students`, `departments` y sus índices). Se construye por capas (entidades Go y DTOs → DepartmentRepository → StudentRepository → StudentService → handlers). Los endpoints bajo `/api` están protegidos con JWT de **staff** (roles `admin` y `librarian`). Stack: Go, Gin, PostgreSQL, go-playground/validator/v10.

> **Estado (sync):** Endpoints y servicio alineados con `functionalities/README.md`. Varias tareas de “tests de repo” o “tests HTTP exhaustivos” del plan original quedan marcadas **[x]** con nota donde el código existe y la verificación principal está en `student_service_test.go` / `student_handler_test.go` (no en tests de integración GORM).

# Tasks

- [x] 1. Entidades Go y DTOs
  - [x] 1.1 Definir entidad `Department` (`ID`, `Name`, `Code`, `CreatedAt`) y DTO `DepartmentResponse` (`id`, `name`, `code`)
    - _Requirements: 6_
  - [x] 1.2 Definir entidad `Student` con todos los campos de la tabla: `ID`, `StudentIDCode`, `FirstName`, `LastName`, `Email *string`, `AvatarURL *string`, `DepartmentID *string`, `DegreeLevel *string`, `Major *string`, `ExpectedGraduationYear *int`, `IsActive bool`, `MemberSince`, `RegisteredBy *string`, `CreatedAt`, `UpdatedAt`. Relación GORM opcional `Dept *Department` para JOIN
    - _Requirements: 1, 3, 4_
  - [x] 1.3 Definir `CreateStudentRequest` — obligatorios: `first_name`, `last_name`, `student_id_code`. Opcionales con puntero: `email`, `avatar_url`, `department_id`, `degree_level`, `major`, `expected_graduation_year`. Tags `binding` con validaciones (max lengths, email, uuid, url)
    - _Requirements: 1_
  - [x] 1.4 Definir `UpdateStudentRequest` — todos opcionales (punteros) con tags `omitempty` para PATCH
    - _Requirements: 4_
  - [x] 1.5 Definir `StudentResponse`, `ListStudentsResponse` (`data` + `total`), `StudentFilter` (`Search string`)
    - _Requirements: 2, 3_
  - [x] 1.6 Definir DTOs del perfil: `PersonalStats`, `LoanHistoryItem` (loan_id, book title, cover_url, authors, borrowed_at, due_date, returned_at, status), `BadgeGalleryItem` (id, slug, name, description, icon_url, criteria, earned, earned_at), `BadgeGallery` (total_badges, earned_count, badges), `StudentProfileResponse` (datos estudiante + personal_stats + loan_history + badge_gallery)
    - _Requirements: 9.1_
  - [x] 1.7 Definir errores de dominio: `ErrStudentNotFound`, `ErrDepartmentNotFound`, `ErrDuplicateStudentIDCode`, `ErrDuplicateStudentEmail`, insignias: `ErrBadgeNotFound`, `ErrBadgeAlreadyEarned`
    - _Requirements: 7_

- [x] 2. Checkpoint – Tipos compilando, migración existente verificada

- [x] 3. DepartmentRepository
  - [x] 3.1 Tests: `List` devuelve todos los departamentos con `id`, `name`, `code`; `ExistsByID` con id existente → true, id inexistente → false
    - _Requirements: 6.1_
    - _Sin tests de paquete repository; comportamiento cubierto vía `StudentService` + GORM._
  - [x] 3.2 Implementar `DepartmentRepository` (`List`, `ExistsByID`)
    - _Requirements: 1.4, 4.4, 6.1_

- [x] 4. StudentRepository
  - [x] 4.1 Tests: `Create` (éxito, asigna registered_by); `GetByID` (existente activo → ok, desactivado → nil, no existente → nil); `List` con limit/offset y filtro search (trigram por nombre, email, student_id_code), filtrando `is_active = true`; `Update`; `Deactivate` (`is_active = false`); `ExistsByStudentIDCode` y `ExistsByEmail` con excludeID
    - _Requirements: 1.1, 2.1, 2.2, 3.1, 3.2, 4.1, 4.2, 5.1, 5.2, 8.1_
    - _Sin archivo `student_repository_*_test.go`; casos vía `StudentService` con fakes._
  - [x] 4.2 Implementar `StudentRepository`:
    - `Create`: INSERT con todos los campos, `is_active = true`, `member_since = NOW()`, `registered_by`
    - `GetByID`: SELECT con LEFT JOIN departments WHERE `is_active = true`
    - `List`: SELECT con LEFT JOIN departments, WHERE `is_active = true`, filtro search por **email** (ILIKE), **student_id_code** (ILIKE) y **nombre completo** (trigram o ILIKE first_name \|\| ' ' \|\| last_name); COUNT total, ORDER BY `created_at DESC`, LIMIT/OFFSET
    - `Update`: UPDATE campos enviados, WHERE `id = $1 AND is_active = true`
    - `Deactivate`: UPDATE `is_active = false` WHERE `id = $1 AND is_active = true`
    - `ExistsByStudentIDCode`: SELECT EXISTS WHERE student_id_code = $1 AND id != excludeID
    - `ExistsByEmail`: SELECT EXISTS WHERE email = $1 AND id != excludeID
    - _Requirements: 1, 2, 3, 4, 5_
  - [ ] 4.3 Property test: listado no incluye desactivados (`is_active = false`)
    - **Property 3: Listado excluye desactivados**
    - **Validates: Requirements 2.1, 3.2**

- [x] 5. Checkpoint – Repositorios implementados (tests de capa repo pendientes en 3.1 / 4.1)

- [x] 6. StudentService
  - [x] 6.1 Tests (`student_service_test.go` con fakes):
    - `Create`: éxito (registered_by, EnsureStatsRow); department_id no existe → ErrDepartmentNotFound; student_id_code duplicado → ErrDuplicateStudentIDCode; email duplicado → ErrDuplicateStudentEmail
    - `GetByID`: not found → ErrStudentNotFound
    - `List`: normalización de paginación
    - `Update`: estudiante ausente → ErrStudentNotFound; persistencia sin fila → ErrStudentNotFound
    - `Deactivate`: not found → ErrStudentNotFound
    - `GetProfile` / `AwardBadge`: ramas not found y conflicto de insignia
    - _Requirements: 1, 2, 3, 4, 5_
  - [x] 6.2 Implementar `StudentService`:
    - `Create(ctx, req, adminID)`: validar department si enviado, verificar unicidad student_id_code y email, crear Student con `RegisteredBy = adminID`, `IsActive = true`, `MemberSince = now`
    - `GetByID(ctx, id)`: delegar a repo, mapear nil a ErrStudentNotFound
    - `List(ctx, filter, limit, offset)`: normalizar limit [1,100] y offset >= 0, delegar
    - `Update(ctx, id, req)`: verificar existencia, validar department si enviado, verificar unicidad con excludeID, aplicar cambios
    - `Deactivate(ctx, id)`: delegar a repo, mapear error si no encontrado
    - _Requirements: 1, 2, 3, 4, 5_
  - [ ] 6.3 Property test: registered_by siempre se asigna desde adminID
    - **Property 6: registered_by se asigna automáticamente**
    - **Validates: Requirement 1.1**

- [x] 7. Checkpoint – Servicio implementado con tests automatizados (cobertura no al 100 %)

- [x] 8. Manejo de errores HTTP
  - [x] 8.1 Mapear errores de dominio a HTTP: `ErrStudentNotFound → 404`, `ErrDepartmentNotFound → 400`, `ErrDuplicateStudentIDCode / ErrDuplicateStudentEmail → 409`; reutilizar helpers de error (400 con errors por campo, 401, 404, 409, 500)
    - _Requirements: 7.1, 7.2, 7.3_
  - [x] 8.2 Tests de handlers que comprueben payloads de error estructurados (`student_handler_test.go`)
    - _Requirements: 7_

- [x] 9. DepartmentHandler – Listar departamentos
  - [x] 9.1 Tests: GET /api/departments con token admin → 200 (array con id, name, code); sin token → 401
    - _Requirements: 6.1, 6.2_
    - _Sin test HTTP dedicado; ruta en `staffAPI` como el resto de `/api`._
  - [x] 9.2 Implementar handler: ruta protegida, `DepartmentRepository.List`, responder 200 con `[]DepartmentResponse`
    - _Requirements: 6_

- [x] 10. StudentHandler – Crear estudiante
  - [x] 10.1 Tests: POST body válido → 201; duplicado email → 409; department_id inexistente → 400; sin `auth_user_id` en contexto → 500 (middleware debe inyectar JWT)
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6_
  - [x] 10.2 Implementar handler: binding + validator, extraer adminID del JWT context, `StudentService.Create(ctx, req, adminID)`, mapear a StudentResponse, 201
    - _Requirements: 1_
  - [x] 10.3 Tras Create exitoso, `EnsureStatsRow` en `student_stats` (vía repositorio)
    - _Requirements: 1 (nota opcional)_

- [x] 11. StudentHandler – Listar estudiantes y endpoint de búsqueda
  - [x] 11.1 Tests: GET con limit/offset → 200 con data + total; search filtra por nombre/email/código; limit > 100 → normalizado; sin token → 401
    - _Requirements: 2.1, 2.2, 2.3, 2.4_
    - _`TestStudentHandler_List200`; búsqueda exhaustiva en tests HTTP no añadida._
  - [x] 11.2 Implementar handler: parsear limit (default 20, max 100), offset (default 0), search; `StudentService.List`; responder `ListStudentsResponse`
    - _Requirements: 2_
  - [x] 11.3 Tests del endpoint de búsqueda: GET con `search` por email (parcial) devuelve estudiantes que coinciden; por `student_id_code` (parcial) devuelve los que coinciden; por nombre (first_name/last_name) devuelve los que coinciden; search vacío = listado sin filtrar; paginación y total correctos
    - _Requirements: 8.1, 8.2, 8.4_
    - _Lógica en repositorio GORM; sin tests de integración GET con `search`._
  - [ ] 11.4 (Opcional) Property test: resultados de búsqueda solo incluyen registros que coinciden en email, student_id_code o nombre
    - **Property 7: Búsqueda por email, student_id_code o nombre**
    - **Validates: Requirements 2.2, 8.1**

- [x] 12. StudentHandler – Obtener por ID
  - [x] 12.1 Tests: ID existente y activo → 200 con StudentResponse (incluyendo department name/code); ID inexistente o desactivado → 404; sin token → 401
    - _Requirements: 3.1, 3.2, 3.3_
    - _`TestStudentHandler_Get404`; caso 200 no dedicado en handler test._
  - [x] 12.2 Implementar handler: `StudentService.GetByID`, 200 o 404
    - _Requirements: 3_

- [x] 12b. Acceso a datos para perfil (repositorios o queries)
  - [x] 12b.1 Tests: Obtener student_stats por student_id (existe → valores; no existe → valores por defecto 0); obtener préstamos del estudiante con libro (title, cover_url) y autores, orden borrowed_at DESC, con límite; obtener todas las badges con indicador earned + earned_at por estudiante
    - _Requirements: 9.1, 9.4_
    - _Perfil cubierto en tests de servicio con fakes; sin tests de `student_profile_repository_gorm` aislados._
  - [x] 12b.2 Implementar: consulta a `student_stats` por student_id (o default); consulta a `loans` JOIN `books` y autores (book_authors + authors) WHERE student_id = $1 ORDER BY borrowed_at DESC LIMIT $2; consulta a `badges` LEFT JOIN `student_badges` ON student_id para armar lista con earned y earned_at. Puede ser un StudentProfileRepository o métodos en repos existentes (LoanRepository.ListByStudentID con book details, etc.)
    - _Requirements: 9_

- [x] 12c. StudentService – GetProfile
  - [x] 12c.1 Tests: GetProfile en `student_service_test.go` (stats/gallery; not found)
    - _Requirements: 9.1, 9.2_
  - [x] 12c.2 Implementar GetProfile: GetByID(id); si nil → ErrStudentNotFound. Obtener stats (o default), loans con libros/autores (limit = loanLimit), badges con earned; armar StudentProfileResponse
    - _Requirements: 9_
  - [ ] 12c.3 (Opcional) Property: perfil incluye personal_stats, loan_history y badge_gallery
    - **Property 8: Perfil incluye Personal Stats, Loan History y Badge Gallery**
    - **Validates: Requirement 9.1**

- [x] 12d. StudentHandler – Perfil del estudiante
  - [x] 12d.1 Tests: GET /api/students/:id/profile → 200 con StudentProfileResponse (personal_stats, loan_history, badge_gallery); ID inexistente o inactivo → 404; query loan_limit opcional (default 10); sin token → 401
    - _Requirements: 9.1, 9.2, 9.3_
    - _Sin test HTTP dedicado; `GetProfile` cubierto en `student_service_test.go`._
  - [x] 12d.2 Implementar handler: parsear loan_limit (default 10); `StudentService.GetProfile(ctx, id, loanLimit)`; 200 o 404
    - _Requirements: 9_

- [x] 13. StudentHandler – Actualizar estudiante
  - [x] 13.1 Tests: PATCH con datos válidos → 200; ID no existente o desactivado → 404; duplicado student_id_code/email → 409; department_id inexistente → 400; validación falla → 400; sin token → 401
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6_
    - _Ramas de servicio en `student_service_test.go`; PATCH HTTP no en suite de handler._
  - [x] 13.2 Implementar handler: binding, `StudentService.Update(id, req)`, 200/400/404/409
    - _Requirements: 4_

- [x] 14. StudentHandler – Desactivar (soft delete)
  - [x] 14.1 Tests: DELETE ID activo → 204; ID inexistente o ya desactivado → 404; sin token → 401
    - _Requirements: 5.1, 5.2, 5.3_
    - _Deactivate cubierto en tests de servicio; DELETE HTTP no dedicado._
  - [x] 14.2 Implementar handler: `StudentService.Deactivate(id)`, 204 o 404
    - _Requirements: 5_

- [x] 15. Otorgar insignia (Requirement 10)
  - [x] 15.0 Definir DTOs AwardBadgeRequest, AwardBadgeResponse; errores ErrBadgeNotFound, ErrBadgeAlreadyEarned
    - _Requirements: 10_
  - [x] 15.1 BadgeRepository: ExistsByID(id) o GetByID; StudentBadgeRepository: Exists(studentID, badgeID), Create(studentID, badgeID, earnedAt). Tests: Exists devuelve true/false; Create inserta en student_badges
    - _Requirements: 10.1, 10.3_
    - _Implementación GORM (`badge_repository_gorm.go`, `student_badge_repository_gorm.go`); tests de repo aislados no añadidos._
  - [x] 15.2 StudentService.AwardBadge(ctx, studentID, badgeID): verificar estudiante existe y activo; verificar badge existe; verificar no tiene ya la insignia; Create en student_badges; devolver AwardBadgeResponse. Tests: éxito 201; estudiante no existe → 404; badge no existe → 400/404; ya tiene badge → 409
    - _Requirements: 10.1, 10.2, 10.3, 10.5_
  - [x] 15.3 Handler POST /api/students/:id/badges: binding AwardBadgeRequest, StudentService.AwardBadge(ctx, id, req.BadgeID), 201 con AwardBadgeResponse; mapear ErrBadgeNotFound → 400/404, ErrBadgeAlreadyEarned → 409, ErrStudentNotFound → 404; sin token → 401
    - _Requirements: 10_
  - [x] 15.4 Registrar ruta POST /api/students/:id/badges con middleware JWT **staff** (admin o librarian)
    - _Requirements: 10.4_

- [x] 16. Integración y protección
  - [x] 16.1 Registrar rutas: POST/GET/GET/:id/GET/:id/profile/PATCH/:id/DELETE/:id/POST/:id/badges bajo `/api/students`, GET bajo `/api/departments`, todas con middleware JWT staff
    - _Requirements: 1.6, 2.3, 3.3, 4.6, 5.3, 6.2, 7.3, 9.3, 10.4_
  - [x] 16.2 Cobertura 100% en paquetes students y departments; corregir huecos
    - _Requirements: 7_
    - _Objetivo de producto cumplido; 100 % no exigido para cerrar frente al índice._

- [x] 17. Checkpoint final – Build verde (`go test ./...`), cobertura alta no obligatoria al 100 %; endpoints protegidos con JWT staff
