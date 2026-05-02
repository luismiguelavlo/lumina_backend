package repositories

import (
	"context"
	"errors"
	"time"

	"library_back/internal/models"

	"gorm.io/gorm"
)

type sanctionRepositoryGorm struct {
	db *gorm.DB
}

// NewSanctionRepositoryGorm returns a GORM-backed SanctionRepository.
func NewSanctionRepositoryGorm(db *gorm.DB) SanctionRepository {
	return &sanctionRepositoryGorm{db: db}
}

func (r *sanctionRepositoryGorm) Create(ctx context.Context, s *models.Sanction) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *sanctionRepositoryGorm) GetByID(ctx context.Context, id string) (*models.Sanction, error) {
	var s models.Sanction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *sanctionRepositoryGorm) ListActive(ctx context.Context, filter SanctionListFilter, limit, offset int) ([]SanctionListRow, int64, error) {
	where := `sc.status::text = 'active' AND st.is_active = true`
	args := make([]interface{}, 0, 3)
	if filter.StudentID != "" {
		where += " AND sc.student_id = ?"
		args = append(args, filter.StudentID)
	}
	var total int64
	countSQL := `SELECT COUNT(*) FROM sanctions sc JOIN students st ON st.id = sc.student_id WHERE ` + where
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	type row struct {
		SanctionID    string
		StudentID     string
		StudentIDCode string
		FirstName     string
		LastName      string
		Email         *string
		Reason        string
		Status        string
		AppliedAt     time.Time
		AppliedByID   *string `gorm:"column:applied_by_id"`
	}
	listSQL := `
SELECT
	sc.id AS sanction_id,
	st.id AS student_id,
	st.student_id_code,
	st.first_name,
	st.last_name,
	st.email,
	sc.reason,
	sc.status::text AS status,
	sc.applied_at,
	sc.applied_by AS applied_by_id
FROM sanctions sc
JOIN students st ON st.id = sc.student_id
WHERE ` + where + `
ORDER BY sc.applied_at DESC
LIMIT ? OFFSET ?`
	var rows []row
	listArgs := append(args, limit, offset)
	if err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]SanctionListRow, 0, len(rows))
	for i := range rows {
		out = append(out, SanctionListRow{
			SanctionID:    rows[i].SanctionID,
			StudentID:     rows[i].StudentID,
			StudentIDCode: rows[i].StudentIDCode,
			FirstName:     rows[i].FirstName,
			LastName:      rows[i].LastName,
			Email:         rows[i].Email,
			Reason:        rows[i].Reason,
			Status:        rows[i].Status,
			AppliedAt:     rows[i].AppliedAt,
			AppliedByID:   rows[i].AppliedByID,
		})
	}
	return out, total, nil
}

func (r *sanctionRepositoryGorm) CountActiveByStudentID(ctx context.Context, studentID string) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Sanction{}).
		Where("student_id = ? AND status = ?", studentID, models.SanctionStatusActive).
		Count(&n).Error
	return n, err
}

func (r *sanctionRepositoryGorm) Lift(ctx context.Context, id string, liftedAt time.Time, liftedBy string) (int64, error) {
	res := r.db.WithContext(ctx).Model(&models.Sanction{}).
		Where("id = ? AND status = ?", id, models.SanctionStatusActive).
		Updates(map[string]interface{}{
			"status":    models.SanctionStatusLifted,
			"lifted_at": liftedAt,
			"lifted_by": liftedBy,
		})
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}
