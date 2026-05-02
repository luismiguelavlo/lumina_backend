package repositories

import (
	"context"
	"errors"
	"strings"
	"time"

	"library_back/internal/models"

	"gorm.io/gorm"
)

type loanRepositoryGorm struct {
	db *gorm.DB
}

// NewLoanRepositoryGorm returns a GORM-backed LoanRepository.
func NewLoanRepositoryGorm(db *gorm.DB) LoanRepository {
	return &loanRepositoryGorm{db: db}
}

func (r *loanRepositoryGorm) Create(ctx context.Context, loan *models.Loan) error {
	return r.db.WithContext(ctx).Create(loan).Error
}

func (r *loanRepositoryGorm) GetByID(ctx context.Context, id string) (*models.Loan, error) {
	var l models.Loan
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&l).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *loanRepositoryGorm) GetJoinByID(ctx context.Context, id string) (*LoanJoinRow, error) {
	var row LoanJoinRow
	err := r.db.WithContext(ctx).Raw(`
SELECT
	l.id AS loan_id,
	l.book_id,
	b.title AS book_title,
	b.isbn AS book_isbn,
	l.student_id,
	s.first_name AS borrower_first,
	s.last_name AS borrower_last,
	l.due_date,
	l.status::text AS status,
	l.borrowed_at,
	l.returned_at
FROM loans l
JOIN books b ON b.id = l.book_id AND b.deleted_at IS NULL
JOIN students s ON s.id = l.student_id
WHERE l.id = ?
`, id).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.LoanID == "" {
		return nil, nil
	}
	return &row, nil
}

func (r *loanRepositoryGorm) List(ctx context.Context, filter LoanListFilter, limit, offset int) ([]LoanJoinRow, int64, error) {
	status := strings.TrimSpace(strings.ToLower(filter.Status))
	if status == "" {
		status = string(models.LoanStatusActive)
	}
	studentID := strings.TrimSpace(filter.StudentID)

	where := "l.status::text = ?"
	args := []interface{}{status}
	if studentID != "" {
		where += " AND l.student_id = ?"
		args = append(args, studentID)
	}

	var total int64
	countSQL := `
SELECT COUNT(*)
FROM loans l
JOIN books b ON b.id = l.book_id AND b.deleted_at IS NULL
JOIN students s ON s.id = l.student_id
WHERE ` + where
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	listSQL := `
SELECT
	l.id AS loan_id,
	l.book_id,
	b.title AS book_title,
	b.isbn AS book_isbn,
	l.student_id,
	s.first_name AS borrower_first,
	s.last_name AS borrower_last,
	l.due_date,
	l.status::text AS status,
	l.borrowed_at,
	l.returned_at
FROM loans l
JOIN books b ON b.id = l.book_id AND b.deleted_at IS NULL
JOIN students s ON s.id = l.student_id
WHERE ` + where + `
ORDER BY l.due_date ASC, l.id ASC
LIMIT ? OFFSET ?
`
	listArgs := append(append([]interface{}{}, args...), limit, offset)
	var rows []LoanJoinRow
	if err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *loanRepositoryGorm) Return(ctx context.Context, id string, returnedAt time.Time) error {
	res := r.db.WithContext(ctx).Model(&models.Loan{}).
		Where("id = ? AND status IN ?", id, []models.LoanStatus{models.LoanStatusActive, models.LoanStatusOverdue}).
		Updates(map[string]interface{}{
			"status":      models.LoanStatusReturned,
			"returned_at": returnedAt,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *loanRepositoryGorm) CountCheckedOutByBookID(ctx context.Context, bookID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Loan{}).
		Where("book_id = ? AND status IN ?", bookID, []models.LoanStatus{models.LoanStatusActive, models.LoanStatusOverdue}).
		Count(&n).Error
	return n, err
}

func (r *loanRepositoryGorm) MarkActiveLoansOverdue(ctx context.Context) (int64, error) {
	res := r.db.WithContext(ctx).Model(&models.Loan{}).
		Where("status = ? AND due_date < CURRENT_DATE", models.LoanStatusActive).
		Update("status", models.LoanStatusOverdue)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}
