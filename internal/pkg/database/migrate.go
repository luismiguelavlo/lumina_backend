package database

import (
	"fmt"

	"library_back/internal/models"

	"gorm.io/gorm"
)

// Migrate runs prerequisite SQL (extensions, ENUM types) then GORM AutoMigrate.
// If you already applied migrations/001_initial_schema.sql, prereqs are no-ops.
func Migrate(db *gorm.DB) error {
	if err := EnsurePostgreSQLPrereqs(db); err != nil {
		return err
	}
	err := db.AutoMigrate(
		&models.User{},
		&models.RevokedToken{},
		&models.Department{},
		&models.Student{},
		&models.Author{},
		&models.Genre{},
		&models.Book{},
		&models.BookAuthor{},
		&models.BookGenre{},
		&models.Loan{},
		&models.Fine{},
		&models.Sanction{},
		&models.Badge{},
		&models.StudentBadge{},
		&models.StudentStats{},
		&models.ReputationEvent{},
		&models.ActivityLog{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	if err := EnsureBookAvailabilityView(db); err != nil {
		return err
	}
	if err := EnsureLeaderboardView(db); err != nil {
		return err
	}
	return nil
}
