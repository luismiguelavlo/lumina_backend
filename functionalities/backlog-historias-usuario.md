# Backlog de historias de usuario — Lumina Library

Fuente: `library_back/functionalities/*/requirements.md`

Cada historia está en **su propia tabla** con el formato solicitado: cabecera con título, **Descripción**, **Validación** (criterios de aceptación) y pie con **Id**, **Prioridad**, **Estimación** y **Dependencia**.

- **Estimación:** story points (SP), salvo que indiques otro criterio (por ejemplo horas).
- **Dependencia:** IDs de historias previas o "—" si no aplica.

---

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">AUTH-01 — Registro de usuario (admin / bibliotecario)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como sistema, quiero permitir el registro de un nuevo usuario con nombre, apellido, email, contraseña y rol, para que pueda ser activado manualmente y luego iniciar sesión.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando se envía POST con first_name, last_name, email, password y opcionalmente role (por defecto admin) y avatar_url, el sistema debe crear el usuario con is_active = false, almacenar la contraseña hasheada (bcrypt) y devolver 201 con los datos del usuario (sin password_hash).</li><li>Cuando el email ya existe en la base de datos, el sistema debe responder 409 con mensaje genérico de conflicto.</li><li>Cuando falta algún campo obligatorio (first_name, last_name, email, password) o la validación falla, el sistema debe responder 400 con un objeto de errores de validación estructurado.</li><li>Cuando la petición no es JSON válido o el cuerpo está vacío, el sistema debe responder 400 con mensaje genérico.</li><li>Cuando se envía un role que no es admin ni librarian, el sistema debe responder 400 con error de validación.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">AUTH-01</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">—</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">AUTH-02 — Inicio de sesión</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador o bibliotecario, quiero iniciar sesión con email y contraseña para obtener tokens y acceder al sistema.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando se envían email y contraseña correctos y el usuario existe con is_active = true, el sistema debe responder 200 con access_token y refresh_token (JWT) con tiempo de vida de 10 horas.</li><li>Cuando las credenciales son incorrectas o el usuario no existe, el sistema debe responder 401 con mensaje genérico (&quot;credenciales inválidas&quot;).</li><li>Cuando el usuario existe pero is_active = false, el sistema debe responder 401 con el mismo mensaje genérico.</li><li>Cuando falta email o password o la validación falla, el sistema debe responder 400 con objeto de errores de validación.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">AUTH-02</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">AUTH-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">AUTH-03 — Renovación de tokens (refresh)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como usuario autenticado, quiero renovar mis tokens usando el refresh token para mantener la sesión.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando se envía un refresh token válido (no expirado, no revocado, usuario activo), el sistema debe responder 200 con un nuevo access_token y refresh_token.</li><li>Cuando el refresh token está expirado, revocado o es inválido, el sistema debe responder 401.</li><li>Cuando el usuario asociado al token tiene is_active = false, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">AUTH-03</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">AUTH-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">AUTH-04 — Cierre de sesión (revocación de tokens)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como usuario autenticado, quiero cerrar sesión para que mis tokens actuales dejen de ser válidos.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando se envía POST con access token válido en Authorization y opcionalmente refresh_token en el cuerpo, el sistema debe insertar los JTI en revoked_tokens y responder 204.</li><li>Cuando no se envía token o el token es inválido o expirado, el sistema debe responder 401.</li><li>Cuando un token está en revoked_tokens, el sistema debe considerarlo inválido en todas las operaciones.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">AUTH-04</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">AUTH-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">AUTH-05 — Validación y manejo de errores (autenticación)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como desarrollador, quiero respuestas HTTP consistentes y validación estructurada.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando cualquier entrada no cumple las reglas de validación (formato email, longitud de password &gt;= 8, campos requeridos, longitudes máximas), el sistema debe responder 400 con payload estructurado (message + errors por campo).</li><li>Cuando se produce un error interno no esperado, el sistema debe responder 500 con mensaje genérico sin exponer detalles.</li><li>Cuando se intenta acceder a un recurso protegido sin token o con token inválido, el sistema debe responder 401 con mensaje genérico.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">AUTH-05</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">AUTH-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">AUTH-06 — Limpieza de la lista negra (cron)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como sistema, quiero limpiar periódicamente la tabla revoked_tokens para no acumular registros indefinidamente.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando han pasado 2 días desde la última ejecución, el sistema debe ejecutar un job que elimine de revoked_tokens todos los registros con revoked_at &lt; NOW() - INTERVAL &#x27;24 hours&#x27;.</li><li>Cuando el job se ejecuta, debe hacerlo sin bloquear las peticiones HTTP.</li><li>Cuando el job falla, el sistema debe registrar el error y continuar; la siguiente ejecución intentará de nuevo.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">AUTH-06</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">2</td>
<td style="border:1px solid #ccc;padding:10px;">AUTH-04</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">CAT-01 — Crear libro en el catálogo</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero crear un libro en el catálogo con sus datos básicos, autores, géneros y ubicación para que esté disponible en el sistema.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía POST con los campos obligatorios (title, isbn, catalog_code) y opcionales válidos (synopsis, publication_year, pages, cover_url, location, total_copies, author_ids, genre_ids), el sistema debe crear el libro con deleted_at = null, insertar en book_authors y book_genres cuando se envíen author_ids y genre_ids, y devolver 201 con los datos del libro creado (incluyendo relaciones si se desea).</li><li>Cuando isbn o catalog_code ya existen en la base de datos, el sistema debe responder 409 con mensaje de conflicto.</li><li>Cuando se envía author_ids o genre_ids y algún ID no existe en su tabla, el sistema debe responder 400 con mensaje de validación.</li><li>Cuando falta algún campo obligatorio o la validación falla (por ejemplo isbn vacío, total_copies negativo), el sistema debe responder 400 con objeto de errores de validación estructurado.</li><li>Cuando no se envía token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">CAT-01</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">8</td>
<td style="border:1px solid #ccc;padding:10px;">AUTH-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">CAT-02 — Listar libros (vista resumida)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero listar los libros del catálogo mostrando solo imagen, título, autor, ISBN y status para una vista de gestión rápida.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita GET al endpoint de listado con paginación (limit, offset), el sistema debe devolver solo libros no eliminados (deleted_at nulo), ordenados de forma consistente (por ejemplo created_at DESC), con el total de registros. Cada elemento debe incluir únicamente: cover_url (imagen), title, author (o authors: nombres concatenados o array), isbn, status (available | borrowed).</li><li>Cuando se envía un filtro de búsqueda (por ejemplo search por título, autor o ISBN), el sistema debe filtrar según diseño (trigram o LIKE) y devolver solo los que coincidan.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li><li>Cuando limit u offset son inválidos, el sistema debe normalizar a valores por defecto (limit por defecto 20, máximo 100; offset &gt;= 0).</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">CAT-02</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">CAT-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">CAT-03 — Detalle completo del libro</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero ver el detalle completo de un libro incluyendo sus relaciones con autor(es), categoría(s) / género(s) y ubicación.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita GET por ID de un libro existente y no eliminado, el sistema debe devolver 200 con todos los campos del libro (id, title, isbn, catalog_code, synopsis, publication_year, pages, cover_url, location, total_copies, created_at, updated_at), el status (available/borrowed), y las relaciones: autor(es) (id y name), género(s) / categoría(s) (id, name, code), y location (campo del libro).</li><li>Cuando el ID no existe o el libro tiene deleted_at no nulo, el sistema debe responder 404.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">CAT-03</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">CAT-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">CAT-04 — Actualizar libro</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero actualizar la información de un libro existente (campos básicos y relaciones).</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía PUT o PATCH con datos válidos para un ID existente y no eliminado, el sistema debe actualizar el libro y, según diseño, las tablas book_authors y book_genres (reemplazar relaciones si se envían author_ids/genre_ids), y devolver 200 con los datos actualizados.</li><li>Cuando el ID no existe o el libro está eliminado, el sistema debe responder 404.</li><li>Cuando tras la actualización se violaría unicidad de isbn o catalog_code (otro libro ya tiene ese valor), el sistema debe responder 409.</li><li>Cuando author_ids o genre_ids contienen IDs inexistentes, el sistema debe responder 400.</li><li>Cuando la validación falla, el sistema debe responder 400 con errores estructurados.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">CAT-04</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">CAT-01, CAT-03</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">CAT-05 — Eliminar libro (soft delete)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero dar de baja un libro del catálogo para que deje de aparecer en listados y búsquedas sin borrar su historial de préstamos.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía DELETE para un ID existente y no eliminado, el sistema debe marcar deleted_at = NOW() y responder 204.</li><li>Cuando el ID no existe o el libro ya está eliminado, el sistema debe responder 404.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">CAT-05</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">CAT-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">CAT-06 — Catálogos de autores y géneros</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero obtener las listas de autores y de géneros para usarlas al crear o editar libros.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita GET al endpoint de autores, el sistema debe devolver la lista de autores (id, name; opcionalmente bio) con 200.</li><li>Cuando el administrador solicita GET al endpoint de géneros, el sistema debe devolver la lista de géneros (id, name, code) con 200.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">CAT-06</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">CAT-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">CAT-07 — Validación y errores (catálogo)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como desarrollador, quiero respuestas HTTP consistentes (400, 401, 404, 409, 500) y validación con mensajes por campo.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando la entrada no cumple reglas de validación (campos requeridos, formatos, rangos), el sistema debe responder 400 con payload estructurado (message + errors por campo).</li><li>Cuando ocurre un error interno no esperado, el sistema debe responder 500 con mensaje genérico.</li><li>Cuando se accede a cualquier endpoint del catálogo sin token válido de admin, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">CAT-07</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">CAT-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">STU-01 — Crear estudiante</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero registrar un estudiante con sus datos personales y académicos para gestionar quién solicita libros.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía POST con los campos obligatorios (first_name, last_name, student_id_code) y campos opcionales válidos (email, department_id, degree_level, major, expected_graduation_year, avatar_url), el sistema debe crear el estudiante con is_active = true, member_since = NOW(), registered_by = admin_id (extraído del JWT) y devolver 201 con los datos del estudiante. (Opcional) el sistema puede crear una fila en student_stats con valores en 0; si no se implementa, el módulo reputation hace upsert al registrar préstamos/devoluciones.</li><li>Cuando student_id_code ya existe en la base de datos, el sistema debe responder 409 con mensaje de conflicto.</li><li>Cuando email ya existe en la base de datos (y no es nulo), el sistema debe responder 409 con mensaje de conflicto.</li><li>Cuando se envía department_id y no existe en la tabla departments, el sistema debe responder 400 con mensaje de validación.</li><li>Cuando falta algún campo obligatorio o la validación falla (por ejemplo email con formato inválido, student_id_code vacío), el sistema debe responder 400 con objeto de errores de validación estructurado.</li><li>Cuando no se envía token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">STU-01</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">8</td>
<td style="border:1px solid #ccc;padding:10px;">AUTH-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">STU-02 — Listar estudiantes activos</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero listar estudiantes activos con paginación y búsqueda por nombre o email.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita GET con paginación (limit, offset), el sistema debe devolver solo estudiantes con is_active = true, ordenados por created_at DESC, con el total de registros que cumplen el filtro.</li><li>Cuando se envía el query param search, el sistema debe filtrar por coincidencia parcial en nombre completo (first_name || &#x27; &#x27; || last_name) usando búsqueda trigram, o por email o student_id_code.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li><li>Cuando limit u offset son inválidos (por ejemplo limit &gt; 100, negativo), el sistema debe normalizar a valores por defecto (limit por defecto 20, máximo 100; offset por defecto 0).</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">STU-02</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">STU-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">STU-03 — Obtener estudiante por ID</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero ver el detalle completo de un estudiante.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el ID existe y el estudiante tiene is_active = true, el sistema debe devolver 200 con todos los datos del estudiante (incluyendo department name/code si tiene department_id).</li><li>Cuando el ID no existe o el estudiante tiene is_active = false, el sistema debe responder 404.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">STU-03</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">STU-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">STU-04 — Actualizar estudiante</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero actualizar la información de un estudiante existente.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía PATCH con datos válidos para un ID existente y activo (is_active = true), el sistema debe actualizar solo los campos enviados y devolver 200 con los datos actualizados.</li><li>Cuando el ID no existe o el estudiante tiene is_active = false, el sistema debe responder 404.</li><li>Cuando la actualización violaría unicidad de student_id_code o email (otro estudiante ya tiene ese valor), el sistema debe responder 409.</li><li>Cuando se envía department_id y no existe en departments, el sistema debe responder 400.</li><li>Cuando la validación falla (por ejemplo email mal formado), el sistema debe responder 400 con errores estructurados.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">STU-04</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">STU-01, STU-03</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">STU-05 — Desactivar estudiante (soft delete)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero desactivar un estudiante para que deje de aparecer en listados sin borrar su historial de préstamos.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía DELETE para un ID existente y activo, el sistema debe marcar is_active = false y responder 204.</li><li>Cuando el ID no existe o el estudiante ya está desactivado (is_active = false), el sistema debe responder 404.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">STU-05</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">STU-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">STU-06 — Listar departamentos (catálogo)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero obtener la lista de departamentos para usarla al crear o editar estudiantes.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita GET al endpoint de departamentos, el sistema debe devolver la lista de todos los departamentos (id, name, code) con 200.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">STU-06</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">2</td>
<td style="border:1px solid #ccc;padding:10px;">STU-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">STU-07 — Validación y errores (estudiantes)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como desarrollador, quiero respuestas HTTP consistentes y validación con mensajes por campo.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando la entrada no cumple reglas de validación (formato email, campos requeridos, max length), el sistema debe responder 400 con payload estructurado (message + errors por campo).</li><li>Cuando ocurre un error interno no esperado, el sistema debe responder 500 con mensaje genérico.</li><li>Cuando se accede a cualquier endpoint de estudiantes o departamentos sin token válido de admin, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">STU-07</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">STU-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">STU-08 — Búsqueda de estudiantes</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero buscar estudiantes por email, código de matrícula (student_id_code) o por nombre para localizar rápidamente a un estudiante.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía GET al endpoint de estudiantes con el query param search (término de búsqueda), el sistema debe devolver solo estudiantes activos que coincidan por al menos uno de: (a) email (coincidencia parcial, case-insensitive), (b) student_id_code (coincidencia parcial o exacta), (c) nombre completo — first_name y/o last_name (coincidencia parcial con búsqueda trigram o ILIKE).</li><li>Cuando se envía search vacío o se omite, el sistema debe comportarse como listado normal (lista paginada sin filtrar por búsqueda).</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li><li>La paginación (limit, offset) debe aplicarse igual que en el listado (limit por defecto 20, máximo 100; offset &gt;= 0). La respuesta debe incluir total con el número de registros que cumplen el criterio de búsqueda.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">STU-08</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">STU-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">STU-09 — Perfil del estudiante (stats, préstamos, insignias)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador o usuario de la aplicación, quiero ver el perfil completo de un estudiante con sus estadísticas personales, historial de préstamos e insignias para tener una vista unificada del estudiante (dashboard de perfil).</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando se solicita GET al endpoint de perfil del estudiante (por ejemplo GET /api/students/:id/profile) para un ID existente y activo (is_active = true), el sistema debe devolver 200 con un objeto que incluya: (a) datos básicos del estudiante para el encabezado del perfil (id, first_name, last_name, avatar_url, degree_level, major, member_since, department name/code si aplica); (b) Personal Stats: total_read, active_loans, overdue_count, current_streak_days y opcionalmente longest_streak_days, total_points, global_rank (desde student_stats; si no hay fila, valores por defecto 0 o null según diseño); (c) Loan History: lista de préstamos del estudiante (orden reciente primero), cada uno con loan id, book (title, cover_url, authors como texto o array), borrowed_at, due_date, returned_at, status (active | returned | overdue). Se puede limitar la cantidad en el perfil (por ejemplo últimos 10 o 20) con opción de &quot;View All&quot; vía endpoint de loans por estudiante; (d) Badge Gallery: total_badges, earned_count y lista de todas las insignias con id, slug, name, icon_url, description/criteria, earned (boolean), earned_at (si earned es true).</li><li>Cuando el ID no existe o el estudiante tiene is_active = false, el sistema debe responder 404.</li><li>Cuando no hay token válido de administrador (o el perfil está protegido), el sistema debe responder 401.</li><li>El historial de préstamos debe incluir datos del libro (título, portada, autores) mediante JOIN o consultas a books y tablas de autores; los préstamos devueltos (status = returned) y activos (status = active) deben poder distinguirse para mostrar &quot;Currently Reading&quot; y &quot;Returned&quot; en la interfaz.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">STU-09</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">8</td>
<td style="border:1px solid #ccc;padding:10px;">STU-03, LOAN-01, REP-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">STU-10 — Otorgar insignia a un estudiante</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero otorgar (asignar) una insignia a un estudiante para que aparezca en su perfil y en la galería de insignias.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía POST al endpoint de otorgar insignia (por ejemplo POST /api/students/:id/badges) con body { &quot;badge_id&quot;: &quot;uuid&quot; } para un estudiante existente y activo y una insignia existente (badges.id), el sistema debe insertar un registro en student_badges (student_id, badge_id, earned_at = NOW()) y devolver 201 (o 200 con el badge otorgado). Si el estudiante ya tiene esa insignia (combinación student_id + badge_id ya existe), el sistema debe responder 409 con mensaje de conflicto o 200 idempotente según criterio de negocio.</li><li>Cuando el student_id no existe o el estudiante tiene is_active = false, el sistema debe responder 404.</li><li>Cuando el badge_id no existe en la tabla badges, el sistema debe responder 400 o 404 con mensaje de validación.</li><li>Cuando no se envía token válido de administrador, el sistema debe responder 401.</li><li>La respuesta debe incluir al menos el badge otorgado (id, slug, name, earned_at) para mostrar en la UI.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">STU-10</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">STU-03</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">LOAN-01 — Listar préstamos (panel de control)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero ver la lista de préstamos con nombre del libro, ID del libro, prestatario, fecha de vencimiento, tiempo restante e ISBN para gestionar los préstamos activos.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita GET al endpoint de préstamos con paginación (limit, offset), el sistema debe devolver préstamos ordenados de forma consistente (por ejemplo due_date ASC para ver primero los más urgentes). Cada elemento debe incluir: loan id, book id, book title, borrower (nombre del estudiante: first_name + last_name o student_id_code según diseño), due_date, time_remaining (días restantes hasta due_date; 0 o negativo si ya venció), isbn del libro.</li><li>Cuando se envía el query param status (por ejemplo active, overdue, returned), el sistema debe filtrar por ese estado. Si no se envía, por defecto debe devolver solo préstamos con status active (o todos según criterio de diseño; se recomienda default active para el panel de control).</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li><li>Cuando limit u offset son inválidos, el sistema debe normalizar (limit por defecto 20, máximo 100; offset &gt;= 0). La respuesta debe incluir el total de registros que cumplen el filtro.</li><li>Cuando se envía el query param opcional student_id (UUID), el sistema debe filtrar los préstamos por ese estudiante, permitiendo listar todos los préstamos de un estudiante (por ejemplo para &quot;View All&quot; en el perfil del estudiante).</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">LOAN-01</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">AUTH-02, STU-01, CAT-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">LOAN-02 — Crear nuevo préstamo</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero registrar un nuevo préstamo indicando el estudiante, el libro y la fecha de vencimiento para que quede registrado en el sistema.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía POST con student_id, book_id y due_date válidos, el sistema debe crear el préstamo con status active, borrowed_at = NOW(), issued_by = admin_id (extraído del JWT), y devolver 201 con los datos del préstamo creado (incluyendo al menos loan id, book, borrower, due_date, time_remaining, isbn si se desea consistencia con el listado).</li><li>Cuando student_id o book_id no existen en la base de datos, el sistema debe responder 400 o 404 con mensaje claro (por ejemplo &quot;student not found&quot;, &quot;book not found&quot;).</li><li>Cuando due_date es anterior a la fecha actual, el sistema debe responder 400 con mensaje de validación (la fecha de vencimiento no puede ser en el pasado).</li><li>Cuando falta algún campo obligatorio o la validación falla, el sistema debe responder 400 con objeto de errores de validación estructurado.</li><li>Cuando no se envía token válido de administrador, el sistema debe responder 401.</li><li>Cuando el libro no tiene copias disponibles (todas prestadas), el sistema debe responder 409 o 400 con mensaje de conflicto (opcional según regla de negocio; si se permite múltiples préstamos del mismo libro a distintos estudiantes mientras no se excedan copias, la validación se hace en servicio).</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">LOAN-02</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">8</td>
<td style="border:1px solid #ccc;padding:10px;">LOAN-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">LOAN-03 — Devolver libro (registrar devolución)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero registrar la devolución de un libro para que el préstamo se marque como devuelto y la copia quede disponible nuevamente.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía PATCH (o POST) al endpoint de devolución para un loan_id existente con status active o overdue, el sistema debe actualizar status = returned, returned_at = NOW() y devolver 200 con los datos del préstamo actualizado.</li><li>Cuando el loan_id no existe, el sistema debe responder 404.</li><li>Cuando el préstamo ya tiene status returned, el sistema debe responder 400 o 409 con mensaje &quot;loan already returned&quot;.</li><li>Cuando no se envía token válido de administrador, el sistema debe responder 401.</li><li>(Opcional) Cuando la devolución es de un préstamo que estaba en mora (returned_at &gt; due_date), el sistema puede crear automáticamente una multa (ver módulo fines) con loan_id, student_id, amount y reason según regla de negocio; si no se implementa, la multa se crea manualmente vía POST /api/fines.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">LOAN-03</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">8</td>
<td style="border:1px solid #ccc;padding:10px;">LOAN-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">LOAN-04 — Validación y errores (préstamos)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como desarrollador, quiero respuestas HTTP consistentes (400, 401, 404, 409, 500) y validación con mensajes por campo.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando la entrada no cumple reglas de validación (due_date en el pasado, UUIDs inválidos, campos requeridos), el sistema debe responder 400 con payload estructurado (message + errors por campo).</li><li>Cuando ocurre un error interno no esperado, el sistema debe responder 500 con mensaje genérico.</li><li>Cuando se accede a cualquier endpoint de préstamos sin token válido de admin, el sistema debe responder 401.</li><li>Cuando se intenta devolver un préstamo ya devuelto, el sistema debe responder 400 o 409.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">LOAN-04</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">LOAN-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">REP-01 — Registrar devolución (RecordReturn)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como sistema, cuando un préstamo se marca como devuelto, quiero registrar el evento de reputación y actualizar las estadísticas del estudiante para que el ranking y el perfil reflejen la actividad.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el LoanService (u otro llamador) invoca ReputationService.RecordReturn(ctx, loan) con un préstamo ya en status returned (returned_at set), el servicio debe crear un registro en reputation_events con: student_id del préstamo, event_type = book_returned_on_time si returned_at &lt;= due_date (o en la fecha de vencimiento), o event_type = book_returned_late si returned_at &gt; due_date; points según regla de negocio (por ejemplo +10 on_time, -5 late; valores configurables); description opcional; reference_id = loan.id.</li><li>Cuando se invoca RecordReturn, el servicio debe actualizar la fila del estudiante en student_stats (o crearla si no existe, upsert): incrementar total_read en 1; decrementar active_loans en 1 (si el préstamo estaba activo); actualizar total_points sumando (o restando) los points del evento recién creado; opcionalmente actualizar current_streak_days y longest_streak_days (por ejemplo devolución a tiempo +1 día de racha, devolución tardía racha a 0).</li><li>Cuando el estudiante no tiene fila en student_stats, el servicio debe crear una con total_read=1, active_loans decrementado según lógica (si se conoce el conteo actual), total_points = points del evento, y el resto en 0 o por defecto.</li><li>La operación debe ser idempotente respecto a múltiples llamadas para el mismo loan_id: no crear eventos duplicados para el mismo préstamo (por ejemplo verificar si ya existe un reputation_event con reference_id = loan.id y tipo returned_on_time/late; si existe, no insertar de nuevo y opcionalmente no re-restar/sumar stats).</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">REP-01</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">8</td>
<td style="border:1px solid #ccc;padding:10px;">LOAN-03</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">REP-02 — Integración con LoanService (reputación)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como desarrollador, quiero que el flujo de devolución de un libro dispare automáticamente el registro de reputación.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el LoanService procesa una devolución (Return) con éxito, el LoanService debe invocar ReputationService.RecordReturn(ctx, loan) después de actualizar el préstamo a status returned (y setear returned_at), pasando el préstamo ya actualizado (o sus datos: student_id, due_date, returned_at, id).</li><li>Cuando RecordReturn falla (por ejemplo error de BD), el comportamiento debe ser según diseño: fallar la devolución (transacción) o registrar el error y continuar (devolución exitosa pero sin actualizar reputación); se recomienda no bloquear la devolución y loguear el fallo para reintento o corrección manual.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">REP-02</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">REP-01, LOAN-03</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">REP-03 — Activity log en devoluciones (opcional)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como sistema, quiero que las devoluciones queden registradas en el activity_log para el feed &quot;Recent activity&quot; del dashboard.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el diseño lo contempla, RecordReturn o el LoanService debe insertar una fila en activity_log con event_type (por ejemplo &quot;book_returned&quot;), title, description, student_id, actor_id (admin que registró la devolución) y metadata opcional, para que GET /api/analytics/dashboard incluya la actividad en recent_activity.</li><li>Si no se implementa en esta fase, el diseño debe dejar documentada la extensión para una fase posterior.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">REP-03</td>
<td style="border:1px solid #ccc;padding:10px;">Baja</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">REP-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">REP-04 — Reglas de puntos y racha</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como producto, quiero puntos y rachas coherentes para gamificación.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Los puntos por evento (book_returned_on_time, book_returned_late) deben ser configurables (constantes o config) y documentados en el design. Ejemplo: on_time +10, late -5.</li><li>La racha (current_streak_days) debe actualizarse según regla documentada: por ejemplo devolución a tiempo en el mismo día o día siguiente incrementa o mantiene racha; devolución tardía resetea a 0. longest_streak_days debe ser el máximo histórico de current_streak_days.</li><li>global_rank en student_stats puede actualizarse con un job periódico o trigger; no es obligatorio que RecordReturn calcule el rank (la vista leaderboard usa RANK() en la consulta).</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">REP-04</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">REP-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">REP-05 — Registrar creación de préstamo (RecordLoanCreated)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como sistema, cuando se crea un préstamo, quiero actualizar active_loans del estudiante para que las estadísticas y el perfil sean correctos.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el LoanService (u otro llamador) invoca ReputationService.RecordLoanCreated(ctx, studentID) tras crear un préstamo con éxito, el servicio debe actualizar la fila del estudiante en student_stats incrementando active_loans en 1. Si el estudiante no tiene fila, el servicio debe crear una con active_loans=1 y el resto en 0 (o por defecto).</li><li>La operación debe ser idempotente respecto a múltiples llamadas para el mismo evento si se pasa un identificador de préstamo (opcional); en su forma mínima, solo se incrementa active_loans sin comprobar duplicados.</li><li>Cuando RecordLoanCreated falla, el crear préstamo debe poder continuar (no bloquear); se recomienda loguear el error.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">REP-05</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">LOAN-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">FIN-01 — Listar multas</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero listar las multas con filtros por estudiante y por estado para gestionar cobros y condonaciones.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita GET al endpoint de multas con paginación (limit, offset), el sistema debe devolver multas ordenadas de forma consistente (por ejemplo created_at DESC), con el total de registros. Cada elemento debe incluir: id, loan_id, student_id, amount, status, reason, created_at, paid_at (si aplica), y opcionalmente datos del estudiante (nombre) o del préstamo según diseño.</li><li>Cuando se envía el query param student_id (UUID), el sistema debe filtrar por ese estudiante. Cuando se envía status (pending, paid, waived), el sistema debe filtrar por ese estado.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li><li>Cuando limit u offset son inválidos, el sistema debe normalizar (limit por defecto 20, máximo 100; offset &gt;= 0). La respuesta debe incluir el total de registros que cumplen el filtro.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">FIN-01</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">AUTH-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">FIN-02 — Crear multa</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero registrar una multa asociada a un préstamo y un estudiante (por ejemplo por devolución tardía) para que quede pendiente de cobro.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía POST con loan_id, student_id, amount (&gt;= 0) y opcionalmente reason, el sistema debe crear la multa con status pending y devolver 201 con los datos de la multa creada.</li><li>Cuando loan_id o student_id no existen, o el loan no pertenece al student_id, el sistema debe responder 400 o 404 con mensaje claro.</li><li>Cuando amount es negativo, el sistema debe responder 400 con objeto de errores de validación.</li><li>Cuando falta algún campo obligatorio (loan_id, student_id, amount), el sistema debe responder 400 con payload estructurado.</li><li>Cuando no se envía token válido de administrador, el sistema debe responder 401.</li><li>(Opcional) Cuando ya existe una multa pendiente para el mismo loan_id, el sistema debe responder 409 o permitir según regla de negocio; el diseño debe documentar el criterio.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">FIN-02</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">LOAN-02, STU-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">FIN-03 — Obtener multa por ID</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero ver el detalle de una multa.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el ID existe, el sistema debe devolver 200 con los datos de la multa (id, loan_id, student_id, amount, status, reason, created_at, paid_at) y opcionalmente datos del estudiante y del préstamo.</li><li>Cuando el ID no existe, el sistema debe responder 404.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">FIN-03</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">FIN-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">FIN-04 — Marcar multa como pagada</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero marcar una multa como pagada cuando el estudiante realiza el pago.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía PATCH al endpoint de marcar como pagada (por ejemplo PATCH /api/fines/:id/paid) para un id de multa existente con status pending, el sistema debe actualizar status = paid, paid_at = NOW() y devolver 200 con los datos actualizados.</li><li>Cuando el id no existe o la multa ya está paid o waived, el sistema debe responder 404 o 400 según diseño.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">FIN-04</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">FIN-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">FIN-05 — Marcar multa como condonada (waived)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero condonar una multa para que el estudiante no deba pagarla.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía PATCH al endpoint de condonar (por ejemplo PATCH /api/fines/:id/waived) para un id de multa existente con status pending, el sistema debe actualizar status = waived y devolver 200 (paid_at permanece null).</li><li>Cuando el id no existe o la multa ya está paid o waived, el sistema debe responder 404 o 400.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">FIN-05</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">FIN-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">FIN-06 — Validación y errores (multas)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como desarrollador, quiero respuestas HTTP consistentes (400, 401, 404, 409, 500).</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando la entrada no cumple reglas de validación (amount negativo, UUIDs inválidos), el sistema debe responder 400 con payload estructurado.</li><li>Cuando ocurre un error interno no esperado, el sistema debe responder 500 con mensaje genérico.</li><li>Cuando se accede a cualquier endpoint de multas sin token válido de admin, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">FIN-06</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">2</td>
<td style="border:1px solid #ccc;padding:10px;">FIN-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">SAN-01 — Listar estudiantes sancionados</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero ver la lista de estudiantes que están actualmente sancionados para gestionar las sanciones activas.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita GET al endpoint de sanciones con paginación (limit, offset), el sistema debe devolver solo sanciones con status active, ordenadas de forma consistente (por ejemplo applied_at DESC), con el total de registros. Cada elemento debe incluir datos de la sanción (id, reason, applied_at, applied_by si está disponible) y datos del estudiante (id, student_id_code, first_name, last_name, email según diseño) para identificar al estudiante sancionado.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li><li>Cuando limit u offset son inválidos, el sistema debe normalizar (limit por defecto 20, máximo 100; offset &gt;= 0). La respuesta debe incluir el total de registros.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">SAN-01</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">AUTH-02, STU-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">SAN-02 — Agregar estudiante a la lista de sanciones</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero agregar un estudiante a la lista de sanciones indicando el motivo para registrar la sanción.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador envía POST con student_id y reason válidos, el sistema debe crear la sanción con status active, applied_at = NOW(), applied_by = admin_id (extraído del JWT) y devolver 201 con los datos de la sanción creada (incluyendo datos del estudiante si se desea consistencia con el listado).</li><li>Cuando student_id no existe en la tabla students, el sistema debe responder 404 con mensaje claro (por ejemplo &quot;student not found&quot;).</li><li>Cuando falta reason o está vacío, el sistema debe responder 400 con objeto de errores de validación estructurado.</li><li>Cuando no se envía token válido de administrador, el sistema debe responder 401.</li><li>(Opcional) Cuando el estudiante ya tiene una sanción activa, el sistema debe responder 409 o permitir múltiples sanciones activas según regla de negocio; el diseño debe documentar el criterio elegido.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">SAN-02</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">SAN-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">SAN-03 — Levantar sanción</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero levantar una sanción para que el estudiante deje de figurar como sancionado.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita levantar la sanción (por ejemplo PATCH o POST al recurso de la sanción) para un id de sanción existente con status active, el sistema debe actualizar status = lifted, lifted_at = NOW(), lifted_by = admin_id (extraído del JWT) y devolver 200 (o 204 según diseño).</li><li>Cuando el id de sanción no existe o la sanción ya está lifted, el sistema debe responder 404.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">SAN-03</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">SAN-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">SAN-04 — Validación y errores (sanciones)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como desarrollador, quiero respuestas HTTP consistentes (400, 401, 404, 409, 500) y validación con mensajes por campo.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando la entrada no cumple reglas de validación (reason vacío, student_id inválido), el sistema debe responder 400 con payload estructurado (message + errors por campo).</li><li>Cuando ocurre un error interno no esperado, el sistema debe responder 500 con mensaje genérico.</li><li>Cuando se accede a cualquier endpoint de sanciones sin token válido de admin, el sistema debe responder 401.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">SAN-04</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">2</td>
<td style="border:1px solid #ccc;padding:10px;">SAN-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">RANK-01 — Top 3 estudiantes (nombre y puntos)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como usuario de la aplicación, quiero ver el Top 3 del Reputation Ranking con nombre y puntos para la sección &quot;Top Readers&quot;.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando se solicita GET al endpoint de top 3, el sistema debe devolver exactamente los tres primeros puestos del ranking (rank 1, 2, 3), cada uno con: posición (rank), nombre del estudiante (first_name, last_name o nombre completo según diseño), y total_points. El orden debe ser por posición ascendente (1, 2, 3).</li><li>Cuando hay menos de 3 estudiantes en el leaderboard (por ejemplo 0, 1 o 2), el sistema debe devolver solo los que existan (array de 0 a 3 elementos).</li><li>La fuente de datos debe ser la vista leaderboard (o equivalente: students activos con student_stats), ordenada por total_points DESC. Solo estudiantes con is_active = true y con fila en student_stats.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">RANK-01</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">REP-01, REP-05</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">RANK-02 — Leaderboard global y puesto del estudiante actual</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como usuario, quiero ver un tramo del leaderboard global desde el puesto 4 (por ejemplo 3 estudiantes: puestos 4, 5, 6) y además ver mi propio puesto y puntos aunque no esté en ese tramo.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando se solicita GET al endpoint de leaderboard con offset y limit (por ejemplo offset=3, limit=3), el sistema debe devolver una lista de estudiantes desde la posición (offset+1) con hasta limit elementos (por ejemplo posiciones 4, 5, 6), cada uno con: rank (posición), student_id, nombre (first_name, last_name o nombre completo), total_points, y opcionalmente avatar_url, books_read, current_streak_days según diseño.</li><li>Cuando se envía el parámetro opcional student_id (UUID del estudiante actual), el sistema debe incluir en la respuesta un objeto current_user (o equivalente) con: rank (posición de ese estudiante en el ranking global), student_id, nombre, total_points; y opcionalmente books_read, current_streak_days para las tarjetas &quot;Your Rank&quot;, &quot;Total Points&quot;, &quot;Books Read&quot;, &quot;Current Streak&quot;. Si el estudiante no tiene entrada en el leaderboard (no tiene student_stats o no está activo), current_user debe indicar que no tiene puesto (por ejemplo null o rank null).</li><li>Cuando no se envía student_id, el sistema debe devolver solo la lista del tramo (data) sin current_user, o current_user null.</li><li>Cuando offset o limit son inválidos (negativos, limit excesivo), el sistema debe normalizar (offset &gt;= 0, limit por defecto 3, máximo según diseño por ejemplo 50).</li><li>La fuente de datos debe ser la vista leaderboard (o students + student_stats) con ranking consistente con el Top 3 (mismo orden por total_points DESC).</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">RANK-02</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">RANK-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">RANK-03 — Validación y errores (ranking)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como desarrollador, quiero respuestas HTTP consistentes y manejo de parámetros inválidos.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando los query params (offset, limit, student_id) son inválidos (por ejemplo student_id no UUID), el sistema debe responder 400 o normalizar cuando sea razonable (offset/limit a valores por defecto).</li><li>Cuando ocurre un error interno no esperado, el sistema debe responder 500 con mensaje genérico.</li><li>Si los endpoints están protegidos por JWT, cuando no hay token válido, el sistema debe responder 401; si son públicos, no se exige autenticación.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">RANK-03</td>
<td style="border:1px solid #ccc;padding:10px;">Baja</td>
<td style="border:1px solid #ccc;padding:10px;">2</td>
<td style="border:1px solid #ccc;padding:10px;">RANK-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">ANL-01 — Métrica total de libros (total books)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero ver el número total de libros en el catálogo para tener una visión del tamaño de la colección.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita la métrica de total books (como parte del dashboard o en un endpoint dedicado), el sistema debe devolver el conteo de filas en la tabla books (o la suma de total_copies si se define como criterio de negocio; por defecto conteo de filas).</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li><li>Cuando la consulta se ejecuta correctamente, la respuesta debe incluir un valor numérico entero (total_books).</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">ANL-01</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">2</td>
<td style="border:1px solid #ccc;padding:10px;">CAT-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">ANL-02 — Métrica de estudiantes activos</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero ver cuántos estudiantes activos hay registrados en el sistema.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita la métrica de active students, el sistema debe devolver el conteo de filas en students donde is_active = true.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li><li>Cuando la consulta se ejecuta correctamente, la respuesta debe incluir un valor numérico entero (active_students).</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">ANL-02</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">2</td>
<td style="border:1px solid #ccc;padding:10px;">STU-01</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">ANL-03 — Métrica overdue fines (multas pendientes)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero ver el monto total de multas pendientes de pago (overdue fines).</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita la métrica de overdue fines, el sistema debe devolver la suma de fines.amount para todas las filas en fines donde status = &#x27;pending&#x27;. Si no hay multas pendientes, la suma debe ser 0.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li><li>Cuando la consulta se ejecuta correctamente, la respuesta debe incluir un valor numérico (decimal o entero según tipo en BD) con dos decimales para representar moneda (overdue_fines).</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">ANL-03</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">2</td>
<td style="border:1px solid #ccc;padding:10px;">FIN-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">ANL-04 — Libros más prestados (top N)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero ver los libros más prestados (top N) para conocer la demanda del catálogo.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita most borrowed books con un límite N (por ejemplo 5 o 10), el sistema debe devolver una lista ordenada por número de préstamos (count de loans por book_id) descendente, con hasta N elementos. Cada elemento debe incluir al menos: identificador del libro, título, cantidad de préstamos; y opcionalmente autor(es) o catalog_code.</li><li>Cuando el límite N no se envía o es inválido, el sistema debe usar un valor por defecto (por ejemplo 5) o un máximo (por ejemplo 20) según diseño.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li><li>Cuando no hay préstamos, el sistema debe devolver una lista vacía o los libros con 0 préstamos según criterio de diseño (recomendado: solo libros con al menos un préstamo, ordenados por count descendente).</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">ANL-04</td>
<td style="border:1px solid #ccc;padding:10px;">Media</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">LOAN-02</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">ANL-05 — Actividad reciente (recent activity)</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como administrador, quiero ver las actividades recientes del sistema (devoluciones, altas de estudiantes, avisos de mora, etc.) para tener un resumen de lo que ha ocurrido.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el administrador solicita recent activity con un límite N (por ejemplo 10 o 20), el sistema debe devolver las últimas N filas de activity_log ordenadas por created_at DESC, incluyendo al menos: id, event_type, title, description, created_at; y opcionalmente metadata, actor_id/actor_name, student_id/student_name.</li><li>Cuando el límite N no se envía o es inválido, el sistema debe usar un valor por defecto (por ejemplo 10) o un máximo (por ejemplo 50) según diseño.</li><li>Cuando no hay token válido de administrador, el sistema debe responder 401.</li><li>Cuando no hay registros en activity_log, el sistema debe devolver una lista vacía.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">ANL-05</td>
<td style="border:1px solid #ccc;padding:10px;">Baja</td>
<td style="border:1px solid #ccc;padding:10px;">3</td>
<td style="border:1px solid #ccc;padding:10px;">REP-03</td>
</tr>
</tbody>
</table>

<table style="width:100%;border-collapse:collapse;margin-bottom:28px;border:1px solid #ccc;">
<thead>
<tr><th colspan="4" style="background-color:#E65100;color:#fff;padding:12px 14px;text-align:left;font-size:1.05em;">ANL-06 — Dashboard de analíticas y errores</th></tr>
</thead>
<tbody>
<tr><td colspan="4" style="border:1px solid #ccc;padding:14px;vertical-align:top;background:#fff;">
<p style="margin:0 0 8px 0;"><strong>Descripción:</strong></p>
<p style="margin:0;">Como desarrollador, quiero un contrato API claro y respuestas HTTP consistentes para las analíticas.</p>
<p style="margin:14px 0 6px 0;"><strong>Validación:</strong></p>
<ul style="margin:0;padding-left:22px;"><li>Cuando el diseño define un único endpoint de dashboard (por ejemplo GET /api/analytics/dashboard), la respuesta debe incluir al menos total_books, active_students, overdue_fines, most_borrowed_books (array) y recent_activity (array), con códigos 200 o 401/500 según corresponda.</li><li>Cuando ocurre un error interno (por ejemplo fallo de BD), el sistema debe responder 500 con mensaje genérico y no exponer detalles internos.</li><li>Cuando se accede sin token válido de administrador, el sistema debe responder 401 en cualquier endpoint de analíticas.</li></ul>
</td></tr>
<tr style="font-weight:bold;background:#f0f0f0;">
<td style="border:1px solid #ccc;padding:10px;width:22%;">Id</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Prioridad</td>
<td style="border:1px solid #ccc;padding:10px;width:22%;">Estimación</td>
<td style="border:1px solid #ccc;padding:10px;width:34%;">Dependencia</td>
</tr>
<tr>
<td style="border:1px solid #ccc;padding:10px;">ANL-06</td>
<td style="border:1px solid #ccc;padding:10px;">Alta</td>
<td style="border:1px solid #ccc;padding:10px;">5</td>
<td style="border:1px solid #ccc;padding:10px;">ANL-01, ANL-02, ANL-03, ANL-04, ANL-05</td>
</tr>
</tbody>
</table>

## Resumen

- **Total historias:** 51
