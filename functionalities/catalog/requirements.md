# Introduction

Gestión del catálogo de libros (Catalog Management) de Lumina Library. Permite a los administradores crear libros, listarlos con una vista resumida (imagen, título, autor, ISBN, status), consultar el detalle completo de un libro con sus relaciones (autor(es), categoría/género(s), ubicación), actualizar un libro y eliminarlo en forma lógica (soft delete). Los datos se almacenan en las tablas `books`, `authors`, `genres`, `book_authors` y `book_genres` del esquema PostgreSQL. El status de disponibilidad (available/borrowed) se deriva de la vista `book_availability` o del conteo de préstamos aún no devueltos (`active` y `overdue`; ver módulo loan). Stack: Go, Gin, PostgreSQL, arquitectura en capas (handlers → services → repositories). Todos los endpoints están protegidos con JWT (rol admin).

# Glossary

- **Book (Libro):** Registro en la tabla `books`. Campos: title, isbn, catalog_code, synopsis, publication_year, pages, cover_url, location, total_copies. Soft delete mediante columna `deleted_at`.
- **Author (Autor):** Tabla `authors` (id, name, bio). Relación M:N con libros vía `book_authors`. Un libro puede tener varios autores.
- **Genre / Category (Género o categoría):** Tabla `genres` (id, name, code). Relación M:N con libros vía `book_genres`. Un libro puede tener uno o más géneros.
- **Location:** Campo `books.location` (VARCHAR). Ubicación física en la biblioteca (ej. "Section A, Shelf 3, Row 2").
- **Status:** Estado de disponibilidad del libro. Valores: `available` (hay copias disponibles) o `borrowed` (todas las copias prestadas). Se calcula a partir de `total_copies` y el número de préstamos no devueltos —`active` y `overdue`— (vista `book_availability` o equivalente).
- **Catalog code:** Código único del libro en el catálogo (ej. FIC-001, SCI-001). Columna `books.catalog_code` UNIQUE.
- **Soft delete:** Marcar un libro como eliminado con `deleted_at = NOW()` sin borrar el registro. Los listados y el detalle excluyen libros con `deleted_at` no nulo.

# Requirements

### Requirement 1: Crear libro

**User Story:** Como administrador, quiero crear un libro en el catálogo con sus datos básicos, autores, géneros y ubicación para que esté disponible en el sistema.

#### Acceptance Criteria

1. WHEN el administrador envía POST con los campos obligatorios (title, isbn, catalog_code) y opcionales válidos (synopsis, publication_year, pages, cover_url, location, total_copies, author_ids, genre_ids), THE sistema SHALL crear el libro con `deleted_at = null`, insertar en `book_authors` y `book_genres` cuando se envíen author_ids y genre_ids, y devolver 201 con los datos del libro creado (incluyendo relaciones si se desea).
2. WHEN `isbn` o `catalog_code` ya existen en la base de datos, THE sistema SHALL responder 409 con mensaje de conflicto.
3. WHEN se envía `author_ids` o `genre_ids` y algún ID no existe en su tabla, THE sistema SHALL responder 400 con mensaje de validación.
4. WHEN falta algún campo obligatorio o la validación falla (ej. isbn vacío, total_copies negativo), THE sistema SHALL responder 400 con objeto de errores de validación estructurado.
5. WHEN no se envía token válido de administrador, THE sistema SHALL responder 401.

### Requirement 2: Listar libros (vista resumida)

**User Story:** Como administrador, quiero listar los libros del catálogo mostrando solo imagen, título, autor, ISBN y status para una vista de gestión rápida.

#### Acceptance Criteria

1. WHEN el administrador solicita GET al endpoint de listado con paginación (limit, offset), THE sistema SHALL devolver solo libros no eliminados (`deleted_at` nulo), ordenados de forma consistente (ej. created_at DESC), con total de registros. Cada elemento SHALL incluir únicamente: cover_url (imagen), title, author (o authors: nombres concatenados o array), isbn, status (available | borrowed).
2. WHEN se envía un filtro de búsqueda (ej. search por título, autor o ISBN), THE sistema SHALL filtrar según diseño (trigram o LIKE) y devolver solo los que coincidan.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.
4. WHEN limit u offset son inválidos, THE sistema SHALL normalizar a valores por defecto (limit default 20, max 100; offset >= 0).

### Requirement 3: Detalle completo del libro

**User Story:** Como administrador, quiero ver el detalle completo de un libro incluyendo sus relaciones con autor(es), categoría(s)/género(s) y ubicación.

#### Acceptance Criteria

1. WHEN el administrador solicita GET por ID de un libro existente y no eliminado, THE sistema SHALL devolver 200 con todos los campos del libro (id, title, isbn, catalog_code, synopsis, publication_year, pages, cover_url, location, total_copies, created_at, updated_at), el status (available/borrowed), y las relaciones: autor(es) (id y name), género(s)/categoría(s) (id, name, code), y location (campo del libro).
2. WHEN el ID no existe o el libro tiene `deleted_at` no nulo, THE sistema SHALL responder 404.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.

### Requirement 4: Actualizar libro

**User Story:** Como administrador, quiero actualizar la información de un libro existente (campos básicos y relaciones).

#### Acceptance Criteria

1. WHEN el administrador envía PUT o PATCH con datos válidos para un ID existente y no eliminado, THE sistema SHALL actualizar el libro y, según diseño, las tablas `book_authors` y `book_genres` (reemplazar relaciones si se envían author_ids/genre_ids), y devolver 200 con los datos actualizados.
2. WHEN el ID no existe o el libro está eliminado, THE sistema SHALL responder 404.
3. WHEN tras la actualización se violaría unicidad de `isbn` o `catalog_code` (otro libro ya tiene ese valor), THE sistema SHALL responder 409.
4. WHEN author_ids o genre_ids contienen IDs inexistentes, THE sistema SHALL responder 400.
5. WHEN la validación falla, THE sistema SHALL responder 400 con errores estructurados.
6. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.

### Requirement 5: Eliminar libro (soft delete)

**User Story:** Como administrador, quiero dar de baja un libro del catálogo para que deje de aparecer en listados y búsquedas sin borrar su historial de préstamos.

#### Acceptance Criteria

1. WHEN el administrador envía DELETE para un ID existente y no eliminado, THE sistema SHALL marcar `deleted_at = NOW()` y responder 204.
2. WHEN el ID no existe o el libro ya está eliminado, THE sistema SHALL responder 404.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.

### Requirement 6: Catálogos de autores y géneros

**User Story:** Como administrador, quiero obtener las listas de autores y de géneros para usarlas al crear o editar libros.

#### Acceptance Criteria

1. WHEN el administrador solicita GET al endpoint de autores, THE sistema SHALL devolver la lista de autores (id, name; opcionalmente bio) con 200.
2. WHEN el administrador solicita GET al endpoint de géneros, THE sistema SHALL devolver la lista de géneros (id, name, code) con 200.
3. WHEN no hay token válido de administrador, THE sistema SHALL responder 401.

### Requirement 7: Validación y errores

**User Story:** Como desarrollador, quiero respuestas HTTP consistentes (400, 401, 404, 409, 500) y validación con mensajes por campo.

#### Acceptance Criteria

1. WHEN la entrada no cumple reglas de validación (campos requeridos, formatos, rangos), THE sistema SHALL responder 400 con payload estructurado (message + errors por campo).
2. WHEN ocurre un error interno no esperado, THE sistema SHALL responder 500 con mensaje genérico.
3. WHEN se accede a cualquier endpoint del catálogo sin token válido de admin, THE sistema SHALL responder 401.
