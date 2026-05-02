package repositories

import (
	"context"
	"errors"
	"strings"
	"time"

	"library_back/internal/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type fineRepositoryGorm struct {
	db *gorm.DB
}

// NewFineRepositoryGorm returns a GORM-backed FineRepository.
func NewFineRepositoryGorm(db *gorm.DB) FineRepository {
	return &fineRepositoryGorm{db: db}
}

func (r *fineRepositoryGorm) Create(ctx context.Context, f *models.Fine) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *fineRepositoryGorm) GetByID(ctx context.Context, id string) (*models.Fine, error) {
	var f models.Fine
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&f).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *fineRepositoryGorm) List(ctx context.Context, filter FineListFilter, limit, offset int) ([]FineWithStudent, int64, error) {
	where := "1=1"
	args := []interface{}{}
	if sid := strings.TrimSpace(filter.StudentID); sid != "" {
		where += " AND f.student_id = ?"
		args = append(args, sid)
	}
	if st := strings.TrimSpace(strings.ToLower(filter.Status)); st != "" {
		where += " AND f.status::text = ?"
		args = append(args, st)
	}

	var total int64
	countSQL := `SELECT COUNT(*) FROM fines f WHERE ` + where
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	type row struct {
		ID        string
		LoanID    string
		StudentID string
		Amount    decimal.Decimal `gorm:"column:amount"`
		Status    string
		Reason    *string
		CreatedAt time.Time
		PaidAt    *time.Time
		Name      string `gorm:"column:student_full_name"`
	}
	listSQL := `
SELECT
	f.id,
	f.loan_id,
	f.student_id,
	f.amount,
	f.status::text AS status,
	f.reason,
	f.created_at,
	f.paid_at,
	BTRIM(COALESCE(s.first_name, '') || ' ' || COALESCE(s.last_name, '')) AS student_full_name
FROM fines f
LEFT JOIN students s ON s.id = f.student_id
WHERE ` + where + `
ORDER BY f.created_at DESC
LIMIT ? OFFSET ?`
	argsWithPage := append(append([]interface{}{}, args...), limit, offset)
	var rows []row
	if err := r.db.WithContext(ctx).Raw(listSQL, argsWithPage...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	out := make([]FineWithStudent, 0, len(rows))
	for i := range rows {
		f := models.Fine{
			ID:        rows[i].ID,
			LoanID:    rows[i].LoanID,
			StudentID: rows[i].StudentID,
			Amount:    rows[i].Amount,
			Status:    models.FineStatus(rows[i].Status),
			Reason:    rows[i].Reason,
			CreatedAt: rows[i].CreatedAt,
			PaidAt:    rows[i].PaidAt,
		}
		out = append(out, FineWithStudent{Fine: f, StudentName: rows[i].Name})
	}
	return out, total, nil
}

func (r *fineRepositoryGorm) UpdateStatus(ctx context.Context, id string, status models.FineStatus, paidAt *time.Time) (int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Fine{}).
		Where("id = ? AND status = ?", id, models.FineStatusPending)
	var res *gorm.DB
	if status == models.FineStatusPaid {
		res = q.Updates(map[string]interface{}{
			"status":  status,
			"paid_at": paidAt,
		})
	} else {
		res = q.Update("status", status)
	}
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

func (r *fineRepositoryGorm) CountPendingByLoanID(ctx context.Context, loanID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Fine{}).
		Where("loan_id = ? AND status = ?", loanID, models.FineStatusPending).
		Count(&n).Error
	return n, err
}
