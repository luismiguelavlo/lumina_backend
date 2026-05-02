# Overview

Plan de implementación para Catalog Management. La tabla `books` y las tablas relacionadas (`authors`, `genres`, `book_authors`, `book_genres`) y la vista `book_availability` ya existen en `migrations/001_initial_schema.sql`, **incluyendo** `books.deleted_at` para soft delete (no hay migración `002` separada en este repo). Implementación en código: `internal/models` (`book`, `catalog_dto`), `internal/repositories` (`catalog_book_repository`, `catalog_author_repository`, `catalog_genre_repository`), `internal/services/book`, `internal/handlers` (`book_handler`, `author_handler`, `genre_handler`). Endpoints bajo `/api` con JWT staff. Stack: Go, Gin, PostgreSQL, go-playground/validator/v10.

> **Estado (sync):** Funcionalidad de catálogo alineada con este índice (`functionalities/README.md`). Tests automatizados principales en `book_service_test.go`; no hay suite dedicada de tests HTTP de libros ni tests de capa repositorio con memoria (ver tareas 4.1, 8.2 y bloque 9–14).

# Tasks

- [x] 1. Migración y modelo de datos
  - [x] 1.1 Columna `deleted_at` e índices en `migrations/001_initial_schema.sql` (no archivo `002` aparte).
    - _Requirements: 5.1_
  - [x] 1.2 Definir entidad `Book` con todos los campos de la tabla (incluyendo `DeletedAt *time.Time`); entidades `Author` y `Genre`; DTOs `CreateBookRequest`, `UpdateBookRequest`, `BookListItem`, `BookDetailResponse`, `AuthorRef`, `GenreRef`, `AuthorResponse`, `GenreResponse`, `ListBooksResponse`
    - _Requirements: 1, 2, 3, 4_
  - [x] 1.3 Definir errores de dominio: `ErrBookNotFound`, `ErrAuthorNotFound`, `ErrGenreNotFound`, `ErrDuplicateISBN`, `ErrDuplicateCatalogCode`
    - _Requirements: 7_
  - [x] 1.4 Definir interfaces `BookRepository`, `AuthorRepository`, `GenreRepository` y `BookService` según design.md
    - _Requirements: 1, 2, 3, 4, 5, 6_

- [x] 2. Checkpoint – Migración aplicable, tipos compilando

- [x] 3. AuthorRepository y GenreRepository
  - [x] 3.1 Tests: AuthorRepository List devuelve todos; ExistsByID(id existente) true, (inexistente) false. GenreRepository igual.
    - _Requirements: 6.1, 6.2_
    - _Cobertura parcial vía `book_service_test.go` y fakes; sin archivo `*_repository_*_test.go` dedicado._
  - [x] 3.2 Implementar AuthorRepository (List, ExistsByID) y GenreRepository (List, ExistsByID)
    - _Requirements: 1.3, 4.4, 6_

- [x] 4. BookRepository
  - [x] 4.1 Tests: Create; GetByID (existente no eliminado → ok, eliminado o inexistente → nil); List con limit/offset y filtro search (WHERE deleted_at IS NULL), verificando que cada ítem tiene solo campos resumidos (cover_url, title, author, isbn, status); Update; SoftDelete (deleted_at = NOW()); ExistsByISBN y ExistsByCatalogCode con excludeID; SetBookAuthors/SetBookGenres (reemplazo); GetBookAuthors, GetBookGenres; GetBookStatus (available/borrowed desde vista o subconsulta)
    - _Requirements: 1.1, 2.1, 2.2, 3.1, 4.1, 5.1_
    - _Lógica en `catalog_book_repository.go`; pruebas directas del repo no añadidas; comportamiento transitivamente cubierto por servicio._
  - [x] 4.2 Implementar BookRepository: todas las operaciones filtrando deleted_at IS NULL donde corresponda; listado con JOIN para autores (agregar nombres) y status (vista book_availability o COUNT de loans activos)
    - _Requirements: 2.1, 3.1, 4.1, 5.1_
  - [ ] 4.3 (Opcional) Property: listado no incluye libros eliminados
    - **Property 3: Listado excluye eliminados**
    - **Validates: Requirements 2.1, 3.2**

- [x] 5. Checkpoint – Repositorios implementados y tests pasando

- [x] 6. BookService
  - [x] 6.1 Tests: Create (éxito con author_ids y genre_ids; isbn duplicado → 409; catalog_code duplicado → 409; author_id inexistente → 400; genre_id inexistente → 400); GetByID (ok con relaciones; not found → ErrBookNotFound); List (paginación, filtro search); Update (ok; not found; duplicado isbn/catalog_code; author/genre inexistente); Delete (ok; not found)
    - _Requirements: 1, 2, 3, 4, 5_
    - _Ver `internal/services/book/book_service_test.go` (subconjunto de casos)._
  - [x] 6.2 Implementar BookService: Create validando autores y géneros y unicidad isbn/catalog_code; GetByID construyendo BookDetailResponse con autores, géneros, location, status; List delegando al repo; Update validando existencia, unicidad y author_ids/genre_ids; Delete con soft delete
    - _Requirements: 1, 2, 3, 4, 5_
  - [ ] 6.3 (Opcional) Property: detalle incluye autores, géneros y location
    - **Property 5: Detalle incluye relaciones**
    - **Validates: Requirement 3.1**

- [x] 7. Checkpoint – Servicio implementado, cobertura en servicio

- [x] 8. Manejo de errores HTTP
  - [x] 8.1 Mapear ErrBookNotFound → 404, ErrAuthorNotFound/ErrGenreNotFound → 400, ErrDuplicateISBN/ErrDuplicateCatalogCode → 409; helpers 400 con errors por campo, 401, 404, 409, 500
    - _Requirements: 7.1, 7.2, 7.3_
  - [x] 8.2 Tests de handlers que comprueben 400 estructurado, 401 sin token, 404 y 409 según caso
    - _Requirements: 7_
    - _Sin `book_handler_test.go`; errores HTTP cubiertos en flujo real + servicio._

- [x] 9. AuthorHandler y GenreHandler (catálogos)
  - [x] 9.1 Tests: GET /api/authors y GET /api/genres con token admin → 200 y lista; sin token → 401
    - _Requirements: 6.1, 6.2, 6.3_
    - _401 igual patrón `staffAPI` que otros `/api`; tests HTTP dedicados no añadidos._
  - [x] 9.2 Implementar handlers: rutas protegidas, llamar AuthorRepository.List y GenreRepository.List, responder con []AuthorResponse y []GenreResponse
    - _Requirements: 6_

- [x] 10. BookHandler – Crear libro
  - [x] 10.1 Tests: POST con body válido y author_ids/genre_ids existentes → 201; isbn o catalog_code duplicado → 409; author_id o genre_id inexistente → 400; validación falla → 400; sin token → 401
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_
    - _Comportamiento vía servicio; tests de handler POST no dedicados._
  - [x] 10.2 Implementar handler: binding + validator, BookService.Create, mapeo a BookDetailResponse, 201
    - _Requirements: 1_

- [x] 11. BookHandler – Listar libros
  - [x] 11.1 Tests: GET con limit/offset → 200 con data y total; cada ítem tiene solo id, cover_url, title, author, isbn, status; search filtra; limit/offset inválidos se normalizan; sin token → 401
    - _Requirements: 2.1, 2.2, 2.3, 2.4_
    - _Tests de handler GET no dedicados._
  - [x] 11.2 Implementar handler: parsear limit (default 20, max 100), offset, search; BookService.List; responder ListBooksResponse
    - _Requirements: 2_
  - [ ] 11.3 (Opcional) Property: listado solo campos resumidos
    - **Property 4: Listado solo campos resumidos**
    - **Validates: Requirement 2.1**

- [x] 12. BookHandler – Detalle del libro
  - [x] 12.1 Tests: ID existente y no eliminado → 200 con BookDetailResponse (autores, géneros, location, status); ID inexistente o eliminado → 404; sin token → 401
    - _Requirements: 3.1, 3.2, 3.3_
    - _Tests de handler no dedicados._
  - [x] 12.2 Implementar handler: BookService.GetByID, 200 o 404
    - _Requirements: 3_

- [x] 13. BookHandler – Actualizar libro
  - [x] 13.1 Tests: PUT/PATCH con datos válidos → 200; ID no existente o eliminado → 404; duplicado isbn/catalog_code → 409; author_id o genre_id inexistente → 400; validación → 400; sin token → 401
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6_
    - _Parcialmente cubierto en `book_service_test.go` (UpdatePatchesRelationsWhenSlicePresent)._
  - [x] 13.2 Implementar handler: binding, BookService.Update(id, req), 200/400/404/409
    - _Requirements: 4_

- [x] 14. BookHandler – Eliminar (soft delete)
  - [x] 14.1 Tests: DELETE ID existente y no eliminado → 204; ID inexistente o ya eliminado → 404; sin token → 401
    - _Requirements: 5.1, 5.2, 5.3_
    - _Servicio: `TestBookService_DeleteNotFound`._
  - [x] 14.2 Implementar handler: BookService.Delete(id), 204 o 404
    - _Requirements: 5_

- [x] 15. Integración y protección
  - [x] 15.1 Registrar rutas: POST/GET/GET/:id/PUT/:id/PATCH/:id/DELETE/:id bajo /api/books, GET /api/authors, GET /api/genres, todas con middleware JWT admin
    - _Requirements: 1.5, 2.3, 3.3, 4.6, 5.3, 6.3, 7.3_
    - _Staff (admin/librarian) en `main.go`, coherente con el resto de `/api`._
  - [x] 15.2 Cobertura en paquetes catalog (books, authors, genres); corregir huecos
    - _Requirements: 7_
    - _Cobertura no al 100 %; build y tests principales verdes._

- [x] 16. Checkpoint final – Build verde, todos los endpoints protegidos, listado solo con imagen/título/autor/isbn/status, detalle con autor/categoría/location
