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

func (r *analyticsRepositoryGorm) PendingFinesCount(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Fine{}).
		Where("status = ?", models.FineStatusPending).
		Count(&n).Error
	return n, err
}

func (r *analyticsRepositoryGorm) OverdueLoansCount(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Loan{}).
		Where("status = ?", models.LoanStatusOverdue).
		Count(&n).Error
	return n, err
}

func (r *analyticsRepositoryGorm) ActiveLoansCount(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Loan{}).
		Where("status IN ?", []models.LoanStatus{models.LoanStatusActive, models.LoanStatusOverdue}).
		Count(&n).Error
	return n, err
}

func (r *analyticsRepositoryGorm) DueSoonCount(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Raw(`
SELECT COUNT(*) FROM loans
WHERE status = ?
  AND due_date::date >= CURRENT_DATE
  AND due_date::date <= CURRENT_DATE + INTERVAL '7 days'
`, models.LoanStatusActive).Scan(&n).Error
	return n, err
}

func (r *analyticsRepositoryGorm) ActiveSanctionsCount(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Sanction{}).
		Where("status = ?", models.SanctionStatusActive).
		Count(&n).Error
	return n, err
}

func (r *analyticsRepositoryGorm) ReturnsThisWeekCount(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Raw(`
SELECT COUNT(*) FROM loans
WHERE status = ?
  AND returned_at IS NOT NULL
  AND returned_at >= NOW() - INTERVAL '7 days'
`, models.LoanStatusReturned).Scan(&n).Error
	return n, err
}

func (r *analyticsRepositoryGorm) NewStudentsThisMonthCount(ctx context.Context) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Raw(`
SELECT COUNT(*) FROM students
WHERE created_at >= date_trunc('month', NOW())
`).Scan(&n).Error
	return n, err
}

func (r *analyticsRepositoryGorm) CopyAvailabilityTotals(ctx context.Context) (available, checkedOut int64, err error) {
	type row struct {
		Available  int64 `gorm:"column:available"`
		CheckedOut int64 `gorm:"column:checked_out"`
	}
	var out row
	err = r.db.WithContext(ctx).Raw(`
SELECT
  COALESCE(SUM(ba.available_copies), 0) AS available,
  COALESCE(SUM(ba.checked_out), 0) AS checked_out
FROM book_availability ba
JOIN books b ON b.id = ba.book_id
WHERE b.deleted_at IS NULL
`).Scan(&out).Error
	return out.Available, out.CheckedOut, err
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

type topOverdueScan struct {
	LoanID      string    `gorm:"column:loan_id"`
	StudentID   string    `gorm:"column:student_id"`
	StudentName string    `gorm:"column:student_name"`
	StudentCode string    `gorm:"column:student_code"`
	BookID      string    `gorm:"column:book_id"`
	BookTitle   string    `gorm:"column:book_title"`
	DueDate     time.Time `gorm:"column:due_date"`
	DaysOverdue int       `gorm:"column:days_overdue"`
}

func (r *analyticsRepositoryGorm) TopOverdueLoans(ctx context.Context, limit int) ([]models.TopOverdueItem, error) {
	const q = `
SELECT
  l.id::text AS loan_id,
  s.id::text AS student_id,
  TRIM(COALESCE(s.first_name, '') || ' ' || COALESCE(s.last_name, '')) AS student_name,
  COALESCE(s.student_id_code, '') AS student_code,
  b.id::text AS book_id,
  b.title AS book_title,
  l.due_date,
  GREATEST(0, (CURRENT_DATE - l.due_date::date))::int AS days_overdue
FROM loans l
JOIN students s ON s.id = l.student_id
JOIN books b ON b.id = l.book_id
WHERE l.status = ?
ORDER BY l.due_date ASC, l.id ASC
LIMIT ?`

	var rows []topOverdueScan
	if err := r.db.WithContext(ctx).Raw(q, models.LoanStatusOverdue, limit).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]models.TopOverdueItem, 0, len(rows))
	for i := range rows {
		out = append(out, models.TopOverdueItem{
			LoanID:      rows[i].LoanID,
			StudentID:   rows[i].StudentID,
			StudentName: strings.TrimSpace(rows[i].StudentName),
			StudentCode: strings.TrimSpace(rows[i].StudentCode),
			BookID:      rows[i].BookID,
			BookTitle:   rows[i].BookTitle,
			DueDate:     rows[i].DueDate.UTC(),
			DaysOverdue: rows[i].DaysOverdue,
		})
	}
	return out, nil
}
