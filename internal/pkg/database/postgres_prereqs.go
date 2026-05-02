package database

import (
	"fmt"

	"gorm.io/gorm"
)

// EnsurePostgreSQLPrereqs creates extensions and ENUM types needed by GORM models
// before AutoMigrate. Safe to run on every startup (idempotent).
func EnsurePostgreSQLPrereqs(db *gorm.DB) error {
	stmts := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,

		`DO $$ BEGIN
			CREATE TYPE user_role AS ENUM ('admin', 'librarian');
		EXCEPTION WHEN duplicate_object THEN NULL; END $$`,

		`DO $$ BEGIN
			CREATE TYPE loan_status AS ENUM ('active', 'returned', 'overdue');
		EXCEPTION WHEN duplicate_object THEN NULL; END $$`,

		`DO $$ BEGIN
			CREATE TYPE fine_status AS ENUM ('pending', 'paid', 'waived');
		EXCEPTION WHEN duplicate_object THEN NULL; END $$`,

		`DO $$ BEGIN
			CREATE TYPE sanction_status AS ENUM ('active', 'lifted');
		EXCEPTION WHEN duplicate_object THEN NULL; END $$`,

		`DO $$ BEGIN
			CREATE TYPE reputation_event_type AS ENUM (
				'book_returned_on_time',
				'book_returned_late',
				'badge_earned',
				'streak_milestone',
				'fine_paid',
				'sanction_received'
			);
		EXCEPTION WHEN duplicate_object THEN NULL; END $$`,
	}

	for _, sql := range stmts {
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("postgres prereq: %w", err)
		}
	}
	return nil
}

// EnsureBookAvailabilityView creates or replaces the catalog view used by
// CatalogBookRepository (list/detail status). Requires books and loans tables (run after AutoMigrate). Idempotent.
func EnsureBookAvailabilityView(db *gorm.DB) error {
	const sql = `CREATE OR REPLACE VIEW book_availability AS
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
WHERE b.deleted_at IS NULL`
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("book_availability view: %w", err)
	}
	return nil
}

// EnsureLeaderboardView creates or replaces the reputation ranking view (students + student_stats).
// Requires students and student_stats (run after AutoMigrate). Idempotent.
func EnsureLeaderboardView(db *gorm.DB) error {
	const sql = `CREATE OR REPLACE VIEW leaderboard AS
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
ORDER BY ss.total_points DESC`
	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("leaderboard view: %w", err)
	}
	return nil
}
