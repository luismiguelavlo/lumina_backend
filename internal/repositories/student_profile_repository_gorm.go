package repositories

import (
	"context"
	"errors"
	"strings"
	"time"

	"library_back/internal/models"

	"gorm.io/gorm"
)

// StudentProfileRepository loads aggregated profile data for a student.
type StudentProfileRepository interface {
	GetPersonalStats(ctx context.Context, studentID string) (models.PersonalStats, error)
	ListLoanHistory(ctx context.Context, studentID string, limit int) ([]models.LoanHistoryItem, error)
	ListBadgeGallery(ctx context.Context, studentID string) (models.BadgeGallery, error)
}

type studentProfileRepositoryGorm struct {
	db *gorm.DB
}

// NewStudentProfileRepositoryGorm returns a GORM-backed StudentProfileRepository.
func NewStudentProfileRepositoryGorm(db *gorm.DB) StudentProfileRepository {
	return &studentProfileRepositoryGorm{db: db}
}

func (r *studentProfileRepositoryGorm) GetPersonalStats(ctx context.Context, studentID string) (models.PersonalStats, error) {
	var st models.StudentStats
	err := r.db.WithContext(ctx).Where("student_id = ?", studentID).First(&st).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.PersonalStats{}, nil
		}
		return models.PersonalStats{}, err
	}
	return models.PersonalStats{
		TotalRead:         st.TotalRead,
		ActiveLoans:       st.ActiveLoans,
		OverdueCount:      st.OverdueCount,
		CurrentStreakDays: st.CurrentStreakDays,
		LongestStreakDays: st.LongestStreakDays,
		TotalPoints:       st.TotalPoints,
		GlobalRank:        st.GlobalRank,
	}, nil
}

type loanHistorySQLRow struct {
	LoanID     string     `gorm:"column:loan_id"`
	BookID     string     `gorm:"column:book_id"`
	Title      string     `gorm:"column:title"`
	CoverURL   *string    `gorm:"column:cover_url"`
	AuthorsAgg *string    `gorm:"column:authors_agg"`
	BorrowedAt time.Time  `gorm:"column:borrowed_at"`
	DueDate    time.Time  `gorm:"column:due_date"`
	ReturnedAt *time.Time `gorm:"column:returned_at"`
	Status     string     `gorm:"column:status"`
}

func (r *studentProfileRepositoryGorm) ListLoanHistory(ctx context.Context, studentID string, limit int) ([]models.LoanHistoryItem, error) {
	var rows []loanHistorySQLRow
	err := r.db.WithContext(ctx).Raw(`
SELECT
	l.id AS loan_id,
	l.book_id,
	b.title,
	b.cover_url,
	(SELECT string_agg(a.name, ', ' ORDER BY a.name)
	 FROM book_authors ba
	 JOIN authors a ON a.id = ba.author_id
	 WHERE ba.book_id = b.id) AS authors_agg,
	l.borrowed_at,
	l.due_date,
	l.returned_at,
	l.status::text AS status
FROM loans l
JOIN books b ON b.id = l.book_id AND b.deleted_at IS NULL
WHERE l.student_id = ?
ORDER BY l.borrowed_at DESC
LIMIT ?
`, studentID, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]models.LoanHistoryItem, 0, len(rows))
	for _, rw := range rows {
		var authors []string
		if rw.AuthorsAgg != nil && *rw.AuthorsAgg != "" {
			for _, p := range strings.Split(*rw.AuthorsAgg, ", ") {
				if t := strings.TrimSpace(p); t != "" {
					authors = append(authors, t)
				}
			}
		}
		if authors == nil {
			authors = []string{}
		}
		var ret *string
		if rw.ReturnedAt != nil {
			s := rw.ReturnedAt.UTC().Format(time.RFC3339)
			ret = &s
		}
		out = append(out, models.LoanHistoryItem{
			LoanID:     rw.LoanID,
			BookID:     rw.BookID,
			Title:      rw.Title,
			CoverURL:   rw.CoverURL,
			Authors:    authors,
			BorrowedAt: rw.BorrowedAt.UTC().Format(time.RFC3339),
			DueDate:    rw.DueDate.UTC().Format("2006-01-02"),
			ReturnedAt: ret,
			Status:     rw.Status,
		})
	}
	return out, nil
}

type badgeGallerySQLRow struct {
	ID          string     `gorm:"column:id"`
	Slug        string     `gorm:"column:slug"`
	Name        string     `gorm:"column:name"`
	Description *string    `gorm:"column:description"`
	IconURL     *string    `gorm:"column:icon_url"`
	Criteria    *string    `gorm:"column:criteria"`
	EarnedAt    *time.Time `gorm:"column:earned_at"`
}

func (r *studentProfileRepositoryGorm) ListBadgeGallery(ctx context.Context, studentID string) (models.BadgeGallery, error) {
	var rows []badgeGallerySQLRow
	err := r.db.WithContext(ctx).Raw(`
SELECT
	b.id,
	b.slug,
	b.name,
	b.description,
	b.icon_url,
	b.criteria,
	sb.earned_at
FROM badges b
LEFT JOIN student_badges sb ON sb.badge_id = b.id AND sb.student_id = ?
ORDER BY b.name ASC
`, studentID).Scan(&rows).Error
	if err != nil {
		return models.BadgeGallery{}, err
	}
	items := make([]models.BadgeGalleryItem, 0, len(rows))
	earned := 0
	for _, rw := range rows {
		en := rw.EarnedAt != nil
		if en {
			earned++
		}
		var eat *string
		if rw.EarnedAt != nil {
			s := rw.EarnedAt.UTC().Format(time.RFC3339)
			eat = &s
		}
		items = append(items, models.BadgeGalleryItem{
			ID:          rw.ID,
			Slug:        rw.Slug,
			Name:        rw.Name,
			Description: rw.Description,
			IconURL:     rw.IconURL,
			Criteria:    rw.Criteria,
			Earned:      en,
			EarnedAt:    eat,
		})
	}
	return models.BadgeGallery{
		TotalBadges: len(rows),
		EarnedCount: earned,
		Badges:      items,
	}, nil
}
