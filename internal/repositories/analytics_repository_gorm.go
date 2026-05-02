package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"library_back/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type analyticsRepositoryGorm struct {
	db *gorm.DB
}

// NewAnalyticsRepositoryGorm builds an AnalyticsRepository backed by GORM/PostgreSQL.
func NewAnalyticsRepositoryGorm(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepositoryGorm{db: db}
}

func (r *analyticsRepositoryGorm) TotalBooks(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Book{}).
		Where("deleted_at IS NULL").
		Count(&n).Error
	return n, err
}

func (r *analyticsRepositoryGorm) ActiveStudents(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Student{}).
		Where("is_active = ?", true).
		Count(&n).Error
	return n, err
}

func (r *analyticsRepositoryGorm) OverdueFinesTotal(ctx context.Context) (float64, error) {
	var sum decimal.Decimal
	err := r.db.WithContext(ctx).Raw(
		`SELECT COALESCE(SUM(amount), 0) FROM fines WHERE status = ?`,
		models.FineStatusPending,
	).Scan(&sum).Error
	if err != nil {
		return 0, err
	}
	f, _ := sum.Float64()
	return f, nil
}

type mostBorrowedScan struct {
	BookID      string `gorm:"column:book_id"`
	Title       string `gorm:"column:title"`
	CatalogCode string `gorm:"column:catalog_code"`
	BorrowCount int64  `gorm:"column:borrow_count"`
	AuthorNames sql.NullString
}

func (r *analyticsRepositoryGorm) MostBorrowedBooks(ctx context.Context, limit int) ([]models.MostBorrowedItem, error) {
	const q = `
SELECT b.id AS book_id, b.title, b.catalog_code, COUNT(l.id) AS borrow_count,
       (SELECT STRING_AGG(a.name, ', ' ORDER BY a.name)
        FROM book_authors ba
        JOIN authors a ON a.id = ba.author_id
        WHERE ba.book_id = b.id) AS author_names
FROM books b
LEFT JOIN loans l ON l.book_id = b.id
WHERE b.deleted_at IS NULL
GROUP BY b.id, b.title, b.catalog_code
HAVING COUNT(l.id) > 0
ORDER BY borrow_count DESC
LIMIT ?`

	var rows []mostBorrowedScan
	if err := r.db.WithContext(ctx).Raw(q, limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]models.MostBorrowedItem, 0, len(rows))
	for i := range rows {
		item := models.MostBorrowedItem{
			BookID:      rows[i].BookID,
			Title:       rows[i].Title,
			CatalogCode: rows[i].CatalogCode,
			BorrowCount: rows[i].BorrowCount,
		}
		if rows[i].AuthorNames.Valid && strings.TrimSpace(rows[i].AuthorNames.String) != "" {
			parts := strings.Split(rows[i].AuthorNames.String, ", ")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					item.AuthorNames = append(item.AuthorNames, p)
				}
			}
		}
		out = append(out, item)
	}
	return out, nil
}

type recentActivityScan struct {
	ID          string `gorm:"column:id"`
	EventType   string `gorm:"column:event_type"`
	Title       string `gorm:"column:title"`
	Description *string
	Metadata    []byte `gorm:"column:metadata"`
	CreatedAt   time.Time
	ActorName   sql.NullString `gorm:"column:actor_name"`
	StudentName sql.NullString `gorm:"column:student_name"`
}

func (r *analyticsRepositoryGorm) RecentActivity(ctx context.Context, limit int) ([]models.RecentActivityItem, error) {
	const q = `
SELECT a.id::text AS id, a.event_type, a.title, a.description, a.metadata, a.created_at,
       NULLIF(TRIM(COALESCE(u.first_name, '') || ' ' || COALESCE(u.last_name, '')), '') AS actor_name,
       NULLIF(TRIM(COALESCE(s.first_name, '') || ' ' || COALESCE(s.last_name, '')), '') AS student_name
FROM activity_log a
LEFT JOIN users u ON u.id = a.actor_id
LEFT JOIN students s ON s.id = a.student_id
ORDER BY a.created_at DESC
LIMIT ?`

	var rows []recentActivityScan
	if err := r.db.WithContext(ctx).Raw(q, limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]models.RecentActivityItem, 0, len(rows))
	for i := range rows {
		item := models.RecentActivityItem{
			ID:        rows[i].ID,
			EventType: rows[i].EventType,
			Title:     rows[i].Title,
			CreatedAt: rows[i].CreatedAt.UTC(),
		}
		if rows[i].Description != nil && strings.TrimSpace(*rows[i].Description) != "" {
			d := strings.TrimSpace(*rows[i].Description)
			item.Description = &d
		}
		if len(rows[i].Metadata) > 0 && string(rows[i].Metadata) != "{}" && string(rows[i].Metadata) != "null" {
			var meta map[string]any
			if err := json.Unmarshal(rows[i].Metadata, &meta); err == nil && len(meta) > 0 {
				item.Metadata = meta
			}
		}
		if rows[i].ActorName.Valid && strings.TrimSpace(rows[i].ActorName.String) != "" {
			a := strings.TrimSpace(rows[i].ActorName.String)
			item.ActorName = &a
		}
		if rows[i].StudentName.Valid && strings.TrimSpace(rows[i].StudentName.String) != "" {
			sn := strings.TrimSpace(rows[i].StudentName.String)
			item.StudentName = &sn
		}
		out = append(out, item)
	}
	return out, nil
}
