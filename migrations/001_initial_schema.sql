-- ============================================================
-- Lumina Library — Initial Database Schema (PostgreSQL)
-- ============================================================
-- users  = administradores / bibliotecarios (CON credenciales)
-- students = estudiantes / patrons (SIN credenciales, registrados por admins)
-- ============================================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ============================================================
-- 1. ENUMS
-- ============================================================

CREATE TYPE user_role AS ENUM ('admin', 'librarian');
CREATE TYPE loan_status AS ENUM ('active', 'returned', 'overdue');
CREATE TYPE fine_status AS ENUM ('pending', 'paid', 'waived');
CREATE TYPE sanction_status AS ENUM ('active', 'lifted');
CREATE TYPE reputation_event_type AS ENUM (
    'book_returned_on_time',
    'book_returned_late',
    'badge_earned',
    'streak_milestone',
    'fine_paid',
    'sanction_received'
);

-- ============================================================
-- 2. USERS (admin / librarian — with credentials)
-- ============================================================

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role          user_role NOT NULL DEFAULT 'admin',
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    avatar_url    TEXT,
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_role ON users (role);

-- ============================================================
-- 3. DEPARTMENTS
-- ============================================================

CREATE TABLE departments (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(150) NOT NULL UNIQUE,
    code       VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- 4. STUDENTS (patrons — NO credentials, registered by admins)
-- ============================================================

CREATE TABLE students (
    id                       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id_code          VARCHAR(50) NOT NULL UNIQUE, -- e.g. LUM-2024-001
    first_name               VARCHAR(100) NOT NULL,
    last_name                VARCHAR(100) NOT NULL,
    email                    VARCHAR(255) UNIQUE, -- contact email, not for login
    avatar_url               TEXT,
    department_id            UUID REFERENCES departments(id) ON DELETE SET NULL,
    degree_level             VARCHAR(100), -- e.g. "Undergraduate Student"
    major                    VARCHAR(150),
    expected_graduation_year INT,
    is_active                BOOLEAN NOT NULL DEFAULT TRUE,
    member_since             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    registered_by            UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_students_code ON students (student_id_code);
CREATE INDEX idx_students_email ON students (email);
CREATE INDEX idx_students_name ON students USING gin (
    (first_name || ' ' || last_name) gin_trgm_ops
);

-- ============================================================
-- 5. CATALOG: AUTHORS, GENRES, BOOKS
-- ============================================================

CREATE TABLE authors (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(200) NOT NULL,
    bio        TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_authors_name ON authors USING gin (name gin_trgm_ops);

CREATE TABLE genres (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(100) NOT NULL UNIQUE,
    code       VARCHAR(10) NOT NULL UNIQUE, -- FIC, SCI, ROM
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE books (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title            VARCHAR(300) NOT NULL,
    isbn             VARCHAR(20) NOT NULL UNIQUE,
    catalog_code     VARCHAR(20) NOT NULL UNIQUE, -- FIC-001, SCI-001
    synopsis         TEXT,
    publication_year INT,
    pages            INT,
    cover_url        TEXT,
    location         VARCHAR(200), -- "Section A, Shelf 3, Row 2"
    total_copies     INT NOT NULL DEFAULT 1 CHECK (total_copies >= 0),
    deleted_at       TIMESTAMPTZ, -- soft delete: NULL = active, NOT NULL = deleted
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_books_title ON books USING gin (title gin_trgm_ops);
CREATE INDEX idx_books_isbn ON books (isbn);
CREATE INDEX idx_books_catalog_code ON books (catalog_code);

CREATE TABLE book_authors (
    book_id   UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

CREATE TABLE book_genres (
    book_id  UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    genre_id UUID NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, genre_id)
);

-- ============================================================
-- 6. LOANS (student borrows a book)
-- ============================================================

CREATE TABLE loans (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    book_id     UUID NOT NULL REFERENCES books(id) ON DELETE RESTRICT,
    student_id  UUID NOT NULL REFERENCES students(id) ON DELETE RESTRICT,
    issued_by   UUID REFERENCES users(id) ON DELETE SET NULL, -- admin who created the loan
    borrowed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    due_date    DATE NOT NULL,
    returned_at TIMESTAMPTZ,
    status      loan_status NOT NULL DEFAULT 'active',
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_loans_student ON loans (student_id);
CREATE INDEX idx_loans_book ON loans (book_id);
CREATE INDEX idx_loans_status ON loans (status);
CREATE INDEX idx_loans_due_date ON loans (due_date) WHERE status = 'active';

-- ============================================================
-- 7. FINES (linked to student via loan)
-- ============================================================

CREATE TABLE fines (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    loan_id    UUID NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    amount     DECIMAL(10,2) NOT NULL CHECK (amount >= 0),
    status     fine_status NOT NULL DEFAULT 'pending',
    reason     TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    paid_at    TIMESTAMPTZ
);

CREATE INDEX idx_fines_student ON fines (student_id);
CREATE INDEX idx_fines_status ON fines (status) WHERE status = 'pending';

-- ============================================================
-- 8. SANCTIONS (applied to students by admins)
-- ============================================================

CREATE TABLE sanctions (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    reason     TEXT NOT NULL,
    status     sanction_status NOT NULL DEFAULT 'active',
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    applied_by UUID REFERENCES users(id) ON DELETE SET NULL, -- admin who sanctioned
    lifted_at  TIMESTAMPTZ,
    lifted_by  UUID REFERENCES users(id) ON DELETE SET NULL, -- admin who lifted
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sanctions_student ON sanctions (student_id);
CREATE INDEX idx_sanctions_active ON sanctions (student_id) WHERE status = 'active';

-- ============================================================
-- 9. BADGES & STUDENT BADGES
-- ============================================================

CREATE TABLE badges (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug        VARCHAR(50) NOT NULL UNIQUE,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    icon_url    TEXT,
    criteria    TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE student_badges (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    badge_id   UUID NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    earned_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (student_id, badge_id)
);

CREATE INDEX idx_student_badges_student ON student_badges (student_id);

-- ============================================================
-- 10. STUDENT STATS (cached / materialized per student)
-- ============================================================

CREATE TABLE student_stats (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id          UUID NOT NULL UNIQUE REFERENCES students(id) ON DELETE CASCADE,
    total_read          INT NOT NULL DEFAULT 0,
    active_loans        INT NOT NULL DEFAULT 0,
    overdue_count       INT NOT NULL DEFAULT 0,
    current_streak_days INT NOT NULL DEFAULT 0,
    longest_streak_days INT NOT NULL DEFAULT 0,
    total_points        INT NOT NULL DEFAULT 0,
    global_rank         INT,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_student_stats_rank ON student_stats (total_points DESC);

-- ============================================================
-- 11. REPUTATION EVENTS (point ledger for students)
-- ============================================================

CREATE TABLE reputation_events (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id   UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    event_type   reputation_event_type NOT NULL,
    points       INT NOT NULL,
    description  TEXT,
    reference_id UUID,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reputation_student ON reputation_events (student_id);
CREATE INDEX idx_reputation_created ON reputation_events (created_at DESC);

-- ============================================================
-- 12. ACTIVITY LOG (system-wide feed)
-- ============================================================

CREATE TABLE activity_log (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_id   UUID REFERENCES users(id) ON DELETE SET NULL, -- admin who performed
    student_id UUID REFERENCES students(id) ON DELETE SET NULL, -- related student
    event_type VARCHAR(50) NOT NULL,
    title      VARCHAR(200) NOT NULL,
    description TEXT,
    metadata   JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_activity_created ON activity_log (created_at DESC);
CREATE INDEX idx_activity_type ON activity_log (event_type);

-- ============================================================
-- 13. SEED: DEFAULT BADGES
-- ============================================================

INSERT INTO badges (slug, name, description, criteria) VALUES
    ('punctual_reader', 'Punctual Reader', 'Always returns books on time', 'Return 10+ books before due date with no overdue'),
    ('bookworm', 'Bookworm', 'Avid reader with impressive stats', 'Read 50+ books'),
    ('explorer', 'Explorer', 'Reads across diverse genres', 'Borrow books from 5+ different genres'),
    ('night_owl', 'Night Owl', 'Late-night library patron', 'Borrow books during night hours 10+ times'),
    ('reviewer', 'Reviewer', 'Active contributor to book reviews', 'Write 10+ book reviews'),
    ('collector', 'Collector', 'Has read the most books in a category', 'Read 20+ books in a single genre');

-- ============================================================
-- 14. VIEW: available copies per book
-- ============================================================

CREATE VIEW book_availability AS
SELECT
    b.id AS book_id,
    b.title,
    b.catalog_code,
    b.total_copies,
    COALESCE(active.checked_out, 0) AS checked_out,
    b.total_copies - COALESCE(active.checked_out, 0) AS available_copies,
    CASE
        WHEN b.total_copies - COALESCE(active.checked_out, 0) > 0 THEN 'available'
        ELSE 'borrowed'
    END AS status
FROM books b
LEFT JOIN (
    SELECT book_id, COUNT(*) AS checked_out
    FROM loans
    WHERE status IN ('active', 'overdue')
    GROUP BY book_id
) active ON active.book_id = b.id
WHERE b.deleted_at IS NULL;

-- ============================================================
-- 15. VIEW: student leaderboard
-- ============================================================

CREATE VIEW leaderboard AS
SELECT
    s.id AS student_id,
    s.first_name,
    s.last_name,
    s.avatar_url,
    ss.total_points,
    ss.total_read AS books_read,
    ss.current_streak_days,
    RANK() OVER (ORDER BY ss.total_points DESC) AS rank_position
FROM students s
JOIN student_stats ss ON ss.student_id = s.id
WHERE s.is_active = TRUE
ORDER BY ss.total_points DESC;

-- ============================================================
-- 16. REVOKED TOKENS (JWT blacklist for logout)
-- ============================================================

CREATE TABLE revoked_tokens (
    jti        VARCHAR(36) PRIMARY KEY,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_revoked_tokens_expires ON revoked_tokens (expires_at);
CREATE INDEX idx_revoked_tokens_revoked ON revoked_tokens (revoked_at);

-- ============================================================
-- 17. updated_at TRIGGER
-- ============================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_students_updated_at BEFORE UPDATE ON students FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_books_updated_at BEFORE UPDATE ON books FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_loans_updated_at BEFORE UPDATE ON loans FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER trg_student_stats_updated_at BEFORE UPDATE ON student_stats FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
