# Funcionalidades Lumina Library — Índice y cobertura

Documento de referencia: qué está definido en cada módulo y qué faltaba (y dónde se cerró).

---

## 1. Lo que ya está definido (por módulo)

| Módulo | requirements | design | tasks | Cubre (resumen) |
|--------|--------------|--------|-------|------------------|
| **autenticación** | ✓ | ✓ | ✓ | Registro, login, refresh, logout, blacklist, cron limpieza. Tablas: users, revoked_tokens. |
| **catalog** | ✓ | ✓ | ✓ | CRUD libros, soft delete (deleted_at), listado resumido, detalle, autores/géneros, book_availability. Tablas: books, authors, genres, book_authors, book_genres. |
| **estudiantes** | ✓ | ✓ | ✓ | CRUD estudiantes, listado + búsqueda (email, student_id_code, nombre), GET por ID, PATCH, DELETE (soft), departamentos, **perfil** (personal stats, loan history, badge gallery), **otorgar insignia** (POST :id/badges). Opcional: crear student_stats al crear estudiante. Tablas: students, departments, student_badges, badges. |
| **loan** | ✓ | ✓ | ✓ | Listar préstamos (filtro status, opcional student_id para “View All”), crear préstamo, **devolver libro** (PATCH return; escribe `activity_log` `book_returned` con actor JWT). **Bloqueo:** no crear préstamo si el estudiante tiene sanción `active` (**403**; ver `sanctions` + `loan/design.md`). **Job en aplicación:** `active` → `overdue` cuando `due_date < CURRENT_DATE` (ticker + ejecución al arranque; ver `loan/design.md`). Vista `book_availability` cuenta préstamos `active` y `overdue` como prestadas. Integración: RecordLoanCreated y RecordReturn (reputation); opcional: multa automática si devolución tardía. Tabla: loans. |
| **fines** | ✓ | ✓ | ✓ | Listar multas (filtros student_id, status), crear multa, GET por ID, marcar pagada, marcar condonada. Tabla: fines. |
| **sanctions** | ✓ | ✓ | ✓ | **API:** `GET/POST /api/sanctions`, `PATCH /api/sanctions/:id` (levantar). Listado solo **active** con datos del estudiante; `applied_by` / `lifted_by` desde JWT. Tabla: sanctions. |
| **ranking** | ✓ | ✓ | ✓ | **Público (sin JWT):** `GET /api/ranking/top3`, `GET /api/ranking/leaderboard`. `current_user` con **uno** de: `student_id`, `email`, `student_id_code` (resolución en `StudentRepository`). Vista `leaderboard` + `EnsureLeaderboardView`. |
| **reputation** | ✓ | ✓ | ✓ | **Servicio interno** (`internal/services/reputation`): **RecordReturn** / **RecordLoanCreated** vía **`ReputationStore`**. **Activity log:** `book_returned` escrito desde **`LoanService`** (`ActivityLogRepository` GORM) tras devolución, con `actor_id` JWT. Idempotencia reputación por `reference_id` = `loan_id`; puntos **+10** / **-5**; racha y `overdue_count` en tardía. Errores de reputación / activity log no revierten préstamo. Detalle: `reputation/design.md`. Backlog: eventos `fine_paid` / `sanction_received` / badges en `reputation_events`. Tablas: `reputation_events`, `student_stats`, `activity_log` (escritura en return). |
| **analytics** | ✓ | ✓ | ✓ | Dashboard: total_books, active_students, overdue_fines, most_borrowed_books, recent_activity. Solo lectura. Tablas: books, students, fines, loans, activity_log. |

Todas las tablas del schema (`001_initial_schema.sql`) quedan cubiertas por al menos un módulo, salvo las extensiones opcionales indicadas abajo.

---

## 2. Huecos que se cerraron en este documento

- **Préstamo bloqueado por sanción activa:** `LoanService` consulta sanciones `active` antes de crear un préstamo; **403** si aplica. Detalle en `loan/requirements.md` (Req. 2.7), `loan/design.md`, `sanctions/requirements.md` (Req. 5) y `sanctions/design.md`.

Se añadieron requisitos/diseño en los módulos existentes para dejar el flujo completo definido antes de pasar a código:

- **active_loans al crear préstamo** → Definido en **reputation** (RecordLoanCreated) y **loan** (invocar tras Create).
- **Multa automática al devolver tarde** → Definido en **loan** (requirement opcional: crear multa al devolver préstamo en mora); implementación opcional vía FineService desde LoanService.
- **Listar préstamos por estudiante (View All)** → Definido en **loan** (filtro opcional `student_id` en GET /api/loans).
- **Crear student_stats al crear estudiante** → Definido en **estudiantes** (opcional: crear fila con 0s).
- **Activity log al devolver** → Ya estaba como opcional en **reputation** (Requirement 3).
- **Otorgar insignia (badge) a estudiante** → Definido en **estudiantes** (Requirement 10: endpoint admin para otorgar badge).
- **Préstamos en mora (`overdue`)** → Job periódico en el servidor (actualización en BD) documentado en **loan** (Requirement 5); no requiere pg_cron.

Detalle de cada uno en los requirements/design de los módulos indicados.

---

## 3. Resumen de “qué falta por definir”

Tras las actualizaciones anteriores, **no queda ningún hueco obligatorio** para poder pasar a código con flujo coherente:

- **Opcionales** (ya documentados como “opcional” o “fase posterior” en requirements):
  - ~~Activity log al devolver~~ — implementado desde `LoanService` + `ActivityLogRepository` (ver reputation Req. 3).
  - Actualizar `global_rank` en student_stats (ranking se calcula con RANK() en la vista; no es obligatorio persistir rank).
  - CRUD completo de badges (crear/editar definiciones de insignias); hoy solo se usan las del seed y “otorgar” al estudiante (estudiantes Req 10).

---

## 4. Orden sugerido de implementación

1. **autenticación** (base de JWT y usuarios).
2. **catalog** (libros y disponibilidad).
3. **estudiantes** (sin perfil completo al inicio; perfil cuando existan loan + reputation).
4. **loan** (listar, crear, devolver + llamada a ReputationService; job `active`→`overdue` según `loan/design.md`).
5. **reputation** (RecordReturn + RecordLoanCreated; opcional activity_log).
6. **fines** (CRUD; opcional: creación automática desde reputation/loan en devolución tardía).
7. **sanctions**, **ranking**, **analytics** (en el orden que prefieras).

Cada módulo tiene en su carpeta: `requirements.md`, `design.md`, `tasks.md`. La plantilla de artefactos está en `plaintext.md`.

---

## 5. Sincronización `tasks.md` ↔ código (2026-04)

La tabla de la sección 1 marca **requirements / design / tasks** como cubiertos (✓) a nivel de **módulo entregado**. Algunos `functionalities/*/tasks.md` conservaban casillas `[ ]` en subtareas pese a que el comportamiento ya existía en el repositorio (p. ej. catálogo en `internal/services/book` y repos `catalog_*`).

Se actualizaron los checklists para que coincidan con ese criterio:

| Archivo | Qué se hizo |
|---------|-------------|
| `catalog/tasks.md` | Marcado **[x]** en el plan de entrega (handlers, servicio, repos, wiring). Quedan **[ ]** solo ítems **opcionales** de property-based tests (4.3, 6.3, 11.3). Overview: `deleted_at` en `001_initial_schema.sql`, rutas y paquetes reales. |
| `autenticación/tasks.md` | Checkpoints **8**, **15.2** y **16** pasan a **[x]** con nota: implementación cerrada; medición/CI al 100 % de cobertura sigue siendo mejora opcional. |
| `estudiantes/tasks.md` | Bloque de nota **Estado (sync)**; subtareas de tests/repo marcadas **[x]** donde hay implementación + verificación vía servicio/handler (con nota si falta test HTTP o GORM dedicado). Opcionales (p. ej. 4.3, 6.3, 11.4, 12c.3) siguen **[ ]**. |
| `analytics/tasks.md` | **3.1** y **9.2** en **[x]** con la misma filosofía (SQL en GORM; tests con fake/mocks). Opcionales 3.3 y 7.3 sin cambio. |

**Regla práctica:** `[x]` = comportamiento acordado implementado (y, cuando aplica, tests suficientes para el repo actual). `[ ]` = backlog u opcional explícito en el propio `tasks.md`.
