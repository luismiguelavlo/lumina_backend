-- ============================================================
-- Lumina Library — Rollback Initial Schema
-- ============================================================

DROP TRIGGER IF EXISTS trg_student_stats_updated_at ON student_stats;
DROP TRIGGER IF EXISTS trg_loans_updated_at ON loans;
DROP TRIGGER IF EXISTS trg_books_updated_at ON books;
DROP TRIGGER IF EXISTS trg_students_updated_at ON students;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP VIEW IF EXISTS leaderboard;
DROP VIEW IF EXISTS book_availability;

DROP TABLE IF EXISTS revoked_tokens;
DROP TABLE IF EXISTS activity_log;
DROP TABLE IF EXISTS reputation_events;
DROP TABLE IF EXISTS student_stats;
DROP TABLE IF EXISTS student_badges;
DROP TABLE IF EXISTS badges;
DROP TABLE IF EXISTS sanctions;
DROP TABLE IF EXISTS fines;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS book_genres;
DROP TABLE IF EXISTS book_authors;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS genres;
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS departments;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS reputation_event_type;
DROP TYPE IF EXISTS sanction_status;
DROP TYPE IF EXISTS fine_status;
DROP TYPE IF EXISTS loan_status;
DROP TYPE IF EXISTS user_role;
