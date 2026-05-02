# Prompt para Stitch (Google) — Diseño UI de Lumina Library

Copia el prompt de abajo en Stitch. Está pensado para generar un diseño completo y coherente que cubra todas las funcionalidades del backend.

---

## PROMPT

Diseña la interfaz completa de **Lumina Library**, una aplicación web de gestión para bibliotecas académicas. El usuario principal es un **administrador/bibliotecario** que gestiona libros, estudiantes, préstamos, multas, sanciones y consulta ranking y analíticas. Los **estudiantes no tienen login**; son registrados y consultados por el admin. El diseño debe ser moderno, limpio, profesional, con tema claro (opción dark mode), sidebar de navegación persistente, y estilo dashboard/panel de administración. Usa una paleta de colores inspirada en tonos azul oscuro, blanco y acentos dorados o ámbar para transmitir seriedad académica. Tipografía legible (sans-serif). Todos los listados deben tener paginación, barra de búsqueda cuando corresponda, y estados vacíos con ilustración. Los formularios deben mostrar validación inline con mensajes de error por campo. Diseña las siguientes pantallas:

---

### PANTALLA 1: Login
- Página centrada con logo "Lumina Library" arriba.
- Formulario: campo email, campo contraseña (con toggle mostrar/ocultar), botón "Iniciar sesión".
- Enlace inferior "¿No tienes cuenta? Regístrate".
- Validación inline: email inválido, contraseña vacía. Mensaje de error general si las credenciales son incorrectas ("Credenciales inválidas").

### PANTALLA 2: Registro
- Formulario: first_name, last_name, email, password (mínimo 8 caracteres), rol (selector: Admin / Librarian), avatar_url (opcional).
- Botón "Registrarse". Enlace "¿Ya tienes cuenta? Inicia sesión".
- Nota visible: "Tu cuenta será creada como inactiva. Un administrador debe activarla."
- Validación inline por cada campo.

### PANTALLA 3: Dashboard (página principal tras login)
- **Sidebar izquierdo** persistente en TODAS las pantallas con iconos y texto: Dashboard, Catálogo, Estudiantes, Préstamos, Multas, Sanciones, Ranking, y al fondo: nombre del usuario logueado con avatar y botón "Cerrar sesión".
- **Área principal:** 3 tarjetas métricas en fila superior:
  - "Total Books" (icono libro, número grande).
  - "Active Students" (icono personas, número grande).
  - "Overdue Fines" (icono moneda/alerta, monto en formato moneda ej. $1,250.00).
- **Sección "Most Borrowed Books":** lista o tarjetas horizontales con los 5 libros más prestados (portada mini, título, autor, número de préstamos).
- **Sección "Recent Activity":** lista vertical tipo timeline con los últimos 10 eventos (icono por tipo de evento, título, descripción breve, fecha relativa "hace 2 horas").

### PANTALLA 4: Catálogo — Listado de libros
- Barra superior: título "Catálogo", barra de búsqueda (placeholder "Buscar por título, autor o ISBN"), botón "+ Nuevo libro".
- Vista en tabla o grilla de tarjetas (toggle vista). Cada elemento muestra: portada (imagen), título, autor(es), ISBN, badge de status (verde "Available" / rojo "Borrowed").
- Paginación inferior (anterior/siguiente, número de página).
- Estado vacío: ilustración + "No hay libros en el catálogo. ¡Agrega el primero!".

### PANTALLA 5: Catálogo — Detalle del libro
- Header: portada grande a la izquierda, a la derecha: título (h1), autores (enlaces o texto), ISBN, catalog_code, status badge.
- Sección info: sinopsis, año de publicación, páginas, ubicación física ("Section A, Shelf 3, Row 2"), total de copias.
- Tags de géneros/categorías (badges/chips: "Fiction", "Science", etc.).
- Botones: "Editar" y "Eliminar" (con confirmación modal).

### PANTALLA 6: Catálogo — Formulario crear/editar libro
- Formulario con campos: título*, ISBN*, código de catálogo*, sinopsis (textarea), año de publicación, páginas, URL portada (con preview de imagen), ubicación, total de copias (número).
- Selector múltiple de autores (dropdown searchable con los autores existentes).
- Selector múltiple de géneros (dropdown searchable o chips seleccionables).
- Botón "Guardar" / "Crear libro". Validación inline.

### PANTALLA 7: Estudiantes — Listado
- Barra superior: título "Estudiantes", barra de búsqueda (placeholder "Buscar por nombre, email o matrícula"), botón "+ Nuevo estudiante".
- Tabla con columnas: Avatar (circle), Nombre completo, Matrícula (student_id_code), Email, Departamento, Estado (badge verde "Activo").
- Click en fila → ir al detalle/perfil. Paginación inferior.

### PANTALLA 8: Estudiantes — Formulario crear/editar
- Campos: first_name*, last_name*, student_id_code* (matrícula), email (opcional), avatar_url (con preview), departamento (dropdown de departamentos), degree_level (ej. "Undergraduate Student"), major, expected_graduation_year.
- Botón "Guardar". Validación inline (matrícula duplicada → error 409 mostrado).

### PANTALLA 9: Estudiantes — Perfil del estudiante (PANTALLA IMPORTANTE)
- **Header del perfil:** avatar grande, nombre completo, debajo: "Undergraduate Student - Computer Science" (degree_level + major), "Member since Sept 2021".
- **Sección "Personal Stats":** 4 tarjetas en fila:
  - "Books Read" (icono libro abierto, número).
  - "Active Loans" (icono reloj, número).
  - "Overdue" (icono alerta, número, rojo si > 0).
  - "Current Streak" (icono fuego/rayo, "12 days").
- **Sección "Loan History":** tabla o lista de tarjetas. Cada ítem: portada mini del libro, título, autor(es), fecha de préstamo, fecha límite, estado (badge: "Active" azul, "Returned" verde, "Overdue" rojo), fecha de devolución si aplica. Separar visualmente "Currently Reading" (activos) arriba y "Returned" abajo. Botón "View All" que lleva a la pantalla de préstamos filtrada por ese estudiante.
- **Sección "Badge Gallery":** encabezado "Insignias (8/12)" con barra de progreso. Grid de insignias: las ganadas se muestran con color e icono + nombre + fecha de obtención; las no ganadas aparecen en gris/bloqueadas con icono de candado. Botón "Otorgar insignia" (abre modal con dropdown de insignias disponibles).
- Botones en el header: "Editar" y "Desactivar" (con confirmación modal: "¿Seguro que quieres desactivar a este estudiante?").

### PANTALLA 10: Préstamos — Loan Control Panel
- Barra superior: título "Panel de Préstamos", filtros: selector de status (Todos / Active / Overdue / Returned), barra de búsqueda opcional, botón "+ Nuevo préstamo".
- Tabla con columnas: Libro (portada mini + título), ISBN, Prestatario (nombre del estudiante), Fecha límite (due_date), Tiempo restante (badge: verde "5 days", amarillo "1 day", rojo "Overdue -3 days"), Status (badge).
- Acción por fila: botón "Devolver" (solo si active/overdue) que abre modal de confirmación: "¿Registrar devolución de [título] por [estudiante]?" con botón "Confirmar devolución".
- Paginación inferior.

### PANTALLA 11: Préstamos — Formulario nuevo préstamo
- Formulario o modal: 
  - Buscador de estudiante (input que busca por nombre/matrícula y muestra sugerencias dropdown).
  - Buscador de libro (input que busca por título/ISBN, solo libros con copias disponibles, muestra sugerencias con portada mini).
  - Fecha de vencimiento (date picker, no permite fechas pasadas).
- Botón "Crear préstamo". Validación: estudiante no encontrado, libro no encontrado, libro sin copias disponibles (409), fecha pasada.

### PANTALLA 12: Multas — Listado
- Barra superior: título "Multas", filtros: selector de status (Todas / Pending / Paid / Waived), filtro por estudiante (buscador), botón "+ Nueva multa".
- Tabla con columnas: ID multa, Estudiante (nombre), Monto (formato moneda), Motivo, Status (badge: amarillo "Pending", verde "Paid", gris "Waived"), Fecha de creación, Fecha de pago (si aplica).
- Acciones por fila (solo si status = pending): botón "Marcar pagada" (icono check verde), botón "Condonar" (icono gift/corazón). Ambas con modal de confirmación.

### PANTALLA 13: Multas — Formulario nueva multa
- Formulario o modal:
  - Buscador de estudiante.
  - Selector de préstamo del estudiante (dropdown con los préstamos del estudiante seleccionado, mostrando título del libro y fecha).
  - Monto (input numérico, no negativo).
  - Motivo/razón (textarea opcional).
- Botón "Crear multa". Validación inline.

### PANTALLA 14: Sanciones — Listado de sancionados
- Barra superior: título "Sanciones", botón "+ Nueva sanción".
- Tabla con columnas: Estudiante (avatar + nombre + matrícula), Motivo, Fecha de aplicación, Aplicada por (nombre del admin).
- Acción por fila: botón "Levantar sanción" (icono unlock) con modal de confirmación: "¿Levantar sanción de [estudiante]? El estudiante podrá volver a solicitar préstamos."
- Solo se muestran sanciones activas. Paginación.

### PANTALLA 15: Sanciones — Formulario nueva sanción
- Formulario o modal:
  - Buscador de estudiante (input con sugerencias).
  - Motivo/razón* (textarea, obligatorio).
- Botón "Aplicar sanción". Validación: estudiante no encontrado, motivo vacío. Si ya tiene sanción activa, mostrar advertencia o error 409.

### PANTALLA 16: Ranking — Reputation Ranking
- **Sección Top 3 "Top Readers":** 3 tarjetas destacadas tipo podio (la del centro/1ro más grande). Cada tarjeta: posición (1, 2, 3 con medalla dorada/plateada/bronce), nombre del estudiante, avatar, total de puntos. Diseño visualmente atractivo.
- **Sección "Global Leaderboard":** tabla/lista desde la posición 4 en adelante. Columnas: Posición (#), Avatar, Nombre, Puntos. Resaltada la fila del usuario actual ("You") con fondo diferente si aparece. Paginación o scroll infinito.
- **Tarjetas laterales o superiores del usuario actual ("Your Stats"):** 4 mini-tarjetas:
  - "Your Rank" (#15 por ejemplo).
  - "Total Points" (1,250).
  - "Books Read" (23).
  - "Current Streak" (12 days).

### PANTALLA 17: Modal de confirmación (componente reutilizable)
- Overlay oscuro, tarjeta centrada con: icono (warning/check/trash según acción), título, mensaje descriptivo, dos botones: "Cancelar" (outline) y "Confirmar" (primario, color según acción: rojo para eliminar/desactivar, verde para confirmar devolución/pago).

### PANTALLA 18: Estado vacío (componente reutilizable)
- Ilustración sutil, texto principal "No hay [elementos] todavía", texto secundario con instrucción, botón de acción ("+ Crear primero").

### PANTALLA 19: Notificaciones/toasts (componente reutilizable)
- Toast en esquina superior derecha: éxito (verde), error (rojo), warning (amarillo). Icono + mensaje breve. Se cierra automáticamente en 4 segundos o con botón X.

---

**NOTAS DE COHERENCIA PARA EL DISEÑO:**
- El sidebar de navegación es IDÉNTICO en todas las pantallas (Dashboard, Catálogo, Estudiantes, Préstamos, Multas, Sanciones, Ranking). Incluye el logo "Lumina Library" arriba y el perfil del admin abajo.
- Todos los listados usan el MISMO componente de tabla con paginación, barra de búsqueda y filtros.
- Todos los formularios usan el MISMO estilo de inputs, labels, mensajes de error inline y botones.
- Los badges de status usan colores consistentes: verde (available/active/paid/returned a tiempo), rojo (borrowed/overdue), amarillo (pending), gris (waived/bloqueado), azul (active loan).
- Los modales de confirmación usan el MISMO componente con variaciones de color según la acción.
- La pantalla de perfil del estudiante es la más rica: combina stats, historial y badges. Es la pantalla showcase de la aplicación.
- El dashboard es la landing page tras login. El ranking puede ser accesible sin autenticación.

Genera todas las pantallas como un flujo de diseño conectado, donde la navegación del sidebar permite acceder a cada sección. Incluye estados hover en botones y filas de tabla. Diseña para escritorio (1440px) con layout responsive.
