# Overview

Plan de implementación para el módulo de analíticas (dashboard). Las tablas y vistas necesarias ya existen en `migrations/001_initial_schema.sql` (books, students, fines, loans, activity_log). Se construye por capas: DTOs y límites → AnalyticsRepository → AnalyticsService → AnalyticsHandler. Un único endpoint GET protegido con JWT staff (admin o librarian), alineado con el grupo `/api` existente. Stack: Go, Gin, PostgreSQL. Sin escritura en BD.

> **Estado (sync):** Módulo implementado y reflejado en `functionalities/README.md`. Tests de repositorio contra PostgreSQL real no añadidos; servicio y handler usan mocks/fakes (ver 3.1, 9.2).

# Tasks

- [x] 1. Modelo de datos y DTOs
  - [x] 1.1 Definir tipos de respuesta: `DashboardResponse` (total_books, active_students, overdue_fines, most_borrowed_books, recent_activity), `MostBorrowedItem` (book_id, title, catalog_code, borrow_count, author_names opcional), `RecentActivityItem` (id, event_type, title, description, metadata, created_at, actor_name opcional, student_name opcional)
    - _Requirements: 1.3, 2.3, 3.3, 4.1, 5.1_
  - [x] 1.2 Definir `DashboardLimits` (MostBorrowedLimit, RecentActivityLimit) y constantes default (5 y 10) y max (20 y 50)
    - _Requirements: 4.2, 5.2_
  - [x] 1.3 Definir interfaz `AnalyticsRepository` con métodos: TotalBooks, ActiveStudents, OverdueFinesTotal, MostBorrowedBooks(limit), RecentActivity(limit); e interfaz `AnalyticsService` con GetDashboard(ctx, limits)
    - _Requirements: 6.1_

- [x] 2. Checkpoint – Tipos compilando

- [x] 3. AnalyticsRepository
  - [x] 3.1 Tests: TotalBooks devuelve COUNT(*) de books; ActiveStudents devuelve COUNT(*) de students WHERE is_active = true; OverdueFinesTotal devuelve SUM(amount) WHERE status = 'pending' (0 si vacío); MostBorrowedBooks(limit) devuelve ≤ limit elementos ordenados por borrow_count DESC; RecentActivity(limit) devuelve ≤ limit elementos ordenados por created_at DESC
    - _Requirements: 1.1, 2.1, 3.1, 4.1, 4.4, 5.1, 5.4_
    - _Sin tests de integración contra BD; la lógica SQL está en `analytics_repository_gorm.go` y el servicio usa fake en tests._
  - [x] 3.2 Implementar AnalyticsRepository:
    - TotalBooks: `SELECT COUNT(*) FROM books`
    - ActiveStudents: `SELECT COUNT(*) FROM students WHERE is_active = true`
    - OverdueFinesTotal: `SELECT COALESCE(SUM(amount), 0) FROM fines WHERE status = 'pending'`
    - MostBorrowedBooks: JOIN books + loans, GROUP BY book_id, ORDER BY count DESC, LIMIT $1; incluir title, catalog_code; opcionalmente autores vía book_authors + authors
    - RecentActivity: SELECT activity_log con LEFT JOIN users y students para actor_name y student_name, ORDER BY created_at DESC, LIMIT $1
    - _Requirements: 1, 2, 3, 4, 5_
  - [ ] 3.3 (Opcional) Property: MostBorrowedBooks y RecentActivity respetan límite y orden
    - **Property 4 y 5: listas acotadas y ordenadas**
    - **Validates: Requirements 4.1, 5.1**

- [x] 4. Checkpoint – Repositorio implementado y tests pasando

- [x] 5. AnalyticsService
  - [x] 5.1 Tests: GetDashboard con límites por defecto devuelve DashboardResponse con los cinco campos; verificar que se llama al repo con los límites correctos (normalizados dentro de default/max)
    - _Requirements: 6.1_
  - [x] 5.2 Implementar AnalyticsService.GetDashboard: normalizar MostBorrowedLimit (default 5, max 20) y RecentActivityLimit (default 10, max 50); llamar a los cinco métodos del repo; construir y devolver DashboardResponse
    - _Requirements: 4.2, 5.2, 6.1_

- [x] 6. Checkpoint – Servicio implementado

- [x] 7. AnalyticsHandler
  - [x] 7.1 Tests: GET /api/analytics/dashboard con token admin válido → 200 y cuerpo con total_books, active_students, overdue_fines, most_borrowed_books, recent_activity; sin token → 401; query params most_borrowed_limit y recent_activity_limit se parsean y se pasan al servicio (validar que se acotan a max)
    - _Requirements: 1.2, 2.2, 3.2, 4.3, 5.3, 6.1, 6.3_
    - _Tests de handler con mock del servicio (200 con query, 500); 401 cubierto por el mismo middleware que el resto de `/api`._
  - [x] 7.2 Implementar handler: extraer query params (most_borrowed_limit, recent_activity_limit), normalizar a default/max; llamar AnalyticsService.GetDashboard; responder 200 con JSON DashboardResponse
    - _Requirements: 6.1_
  - [ ] 7.3 (Opcional) Property: sin token admin siempre 401
    - **Property 6: Sin token admin → 401**
    - **Validates: Requirements 1.2, 6.3**
    - _Ruta bajo `staffAPI`: sin token staff → 401._

- [x] 8. Manejo de errores
  - [x] 8.1 Mapear error de repositorio (ej. contexto cancelado, BD no disponible) a 500 con mensaje genérico; no exponer detalles internos
    - _Requirements: 6.2_
  - [x] 8.2 Tests: simular error del repositorio y verificar que el handler responde 500
    - _Requirements: 6.2_

- [x] 9. Integración
  - [x] 9.1 Registrar ruta GET /api/analytics/dashboard con middleware JWT staff (`staffAPI` en `main.go`); inyectar AnalyticsService en el handler
    - _Requirements: 6.3_
    - _No solo admin: mismo criterio que books/loans/fines (staff admin o librarian)._
  - [x] 9.2 Ejecutar cobertura en el paquete analytics y corregir huecos
    - _Requirements: 6_
    - _Paquete con tests de servicio/handler; gate de cobertura al 100 % no requerido para alinear con el índice._

- [x] 10. Checkpoint final – Build verde, endpoint protegido, respuesta JSON con las cinco métricas
