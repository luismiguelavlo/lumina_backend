# Introduction

Reputation Ranking de Lumina Library. Permite mostrar el ranking de reputación de estudiantes por puntos: un endpoint devuelve el **Top 3** (nombre y puntos) para la sección "Reputation Ranking - Top Readers"; otro endpoint devuelve un tramo del **leaderboard global** desde la posición 4 (por defecto 3 estudiantes: puestos 4, 5, 6) y además la **posición y datos del estudiante actual** ("Your Rank" / "You" en el leaderboard). Los datos provienen de la vista `leaderboard` (students + student_stats: total_points, rank_position). Stack: Go, Gin, PostgreSQL, arquitectura en capas (handlers → services → repositories). Los endpoints **`GET /api/ranking/top3`** y **`GET /api/ranking/leaderboard`** son **públicos** (sin JWT): cualquier cliente puede consultar el ranking.

# Glossary

- **Reputation Ranking:** Clasificación de estudiantes por puntos de reputación (total_points en student_stats). Se muestra en la sección "Reputation Ranking - Top Readers" y en el "Global Leaderboard".
- **Top 3:** Los tres primeros puestos del ranking (rank_position 1, 2, 3). Cada uno con nombre del estudiante y total de puntos.
- **Leaderboard global:** Lista ordenada por total_points descendente; cada fila tiene rank (posición), estudiante (nombre, avatar opcional) y puntos. Solo incluye estudiantes activos con registro en student_stats (vista `leaderboard`).
- **Current user / Tu puesto:** Datos del estudiante en el ranking global (rank, puntos, etc.). Como los estudiantes no tienen login de app, se identifica con **exactamente uno** de los query params opcionales: **`student_id`** (UUID interno), **`email`** (correo del perfil, comparación sin distinguir mayúsculas) o **`student_id_code`** (matrícula / código institucional). Si no tiene fila en `student_stats` o no está activo, `current_user` es null.
- **total_points:** Puntos de reputación del estudiante; columna en `student_stats`. El ranking se ordena por total_points DESC.
- **rank_position:** Posición en el ranking (1 = primero). Se calcula con RANK() OVER (ORDER BY total_points DESC) en la vista `leaderboard`.

# Requirements

### Requirement 1: Top 3 estudiantes (nombre y puntos)

**User Story:** Como usuario de la aplicación, quiero ver el Top 3 del Reputation Ranking con nombre y puntos para la sección "Top Readers".

#### Acceptance Criteria

1. WHEN se solicita GET al endpoint de top 3, THE sistema SHALL devolver exactamente los tres primeros puestos del ranking (rank 1, 2, 3), cada uno con: posición (rank), nombre del estudiante (first_name, last_name o nombre completo según diseño), y total_points. Orden SHALL ser por posición ascendente (1, 2, 3).
2. WHEN hay menos de 3 estudiantes en el leaderboard (p. ej. 0, 1 o 2), THE sistema SHALL devolver solo los que existan (array de 0 a 3 elementos).
3. La fuente de datos SHALL ser la vista `leaderboard` (o equivalente: students activos con student_stats), ordenada por total_points DESC. Solo estudiantes con is_active = true y con fila en student_stats.

### Requirement 2: Leaderboard desde posición 4 + puesto del estudiante actual

**User Story:** Como usuario, quiero ver un tramo del leaderboard global desde el puesto 4 (p. ej. 3 estudiantes: puestos 4, 5, 6) y además ver mi propio puesto y puntos aunque no esté en ese tramo.

#### Acceptance Criteria

1. WHEN se solicita GET al endpoint de leaderboard con offset y limit (ej. offset=3, limit=3), THE sistema SHALL devolver una lista de estudiantes desde la posición (offset+1) con hasta `limit` elementos (ej. posiciones 4, 5, 6), cada uno con: rank (posición), student_id, nombre (first_name, last_name o nombre completo), total_points, y opcionalmente avatar_url, books_read, current_streak_days según diseño.
2. WHEN se envía **exactamente uno** de los parámetros opcionales `student_id` (UUID), `email` o `student_id_code`, THE sistema SHALL resolver al estudiante activo correspondiente y, si tiene fila en el leaderboard, incluir `current_user` con rank, student_id, nombre, total_points y opcionalmente books_read, current_streak_days. Si se envían dos o más a la vez, THE sistema SHALL responder **400**. Si el identificador no corresponde a ningún activo o no está en el leaderboard, `current_user` SHALL ser **null**.
3. WHEN no se envía ninguno de esos parámetros, THE lista del tramo (data) se devuelve y `current_user` SHALL ser **null**.
4. WHEN offset o limit son inválidos (negativos, limit excesivo), THE sistema SHALL normalizar (offset >= 0, limit por defecto 3, máximo según diseño p. ej. 50).
5. La fuente de datos SHALL ser la vista `leaderboard` (o students + student_stats) con ranking consistente con el Top 3 (mismo orden por total_points DESC).

### Requirement 3: Validación y errores

**User Story:** Como desarrollador, quiero respuestas HTTP consistentes y manejo de parámetros inválidos.

#### Acceptance Criteria

1. WHEN `student_id` no es UUID válido, o `email` no es un correo aceptable por RFC (parsing estricto), o `student_id_code` supera la longitud permitida, THE sistema SHALL responder **400** donde aplique; offset/limit se normalizan cuando es razonable.
2. WHEN ocurre un error interno no esperado, THE sistema SHALL responder 500 con mensaje genérico.
3. Los endpoints de ranking **no** exigen autenticación.
