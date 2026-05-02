package repositories

import (
	"context"
	"errors"
	"strings"

	"library_back/internal/models"

	"gorm.io/gorm"
)

// StudentRepository persists students (active-only reads for staff API).
type StudentRepository interface {
	Create(ctx context.Context, s *models.Student) error
	GetActiveByID(ctx context.Context, id string) (*models.Student, error)
	List(ctx context.Context, filter models.StudentFilter, limit, offset int) ([]models.Student, int64, error)
	Update(ctx context.Context, s *models.Student) error
	Deactivate(ctx context.Context, id string) error
	ExistsByStudentIDCode(ctx context.Context, code string, excludeID string) (bool, error)
	ExistsByEmail(ctx context.Context, email string, excludeID string) (bool, error)
	EnsureStatsRow(ctx context.Context, studentID string) error
	// GetActiveStudentIDByEmail returns internal UUID for an active student with this email (case-insensitive). Empty if none.
	GetActiveStudentIDByEmail(ctx context.Context, email string) (string, error)
	// GetActiveStudentIDByStudentIDCode returns internal UUID for an active student with this matrícula (trimmed exact match). Empty if none.
	GetActiveStudentIDByStudentIDCode(ctx context.Context, studentIDCode string) (string, error)
}

type studentRepositoryGorm struct {
	db *gorm.DB
}

// NewStudentRepositoryGorm returns a GORM-backed StudentRepository.
func NewStudentRepositoryGorm(db *gorm.DB) StudentRepository {
	return &studentRepositoryGorm{db: db}
}

func (r *studentRepositoryGorm) Create(ctx context.Context, s *models.Student) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *studentRepositoryGorm) GetActiveByID(ctx context.Context, id string) (*models.Student, error) {
	var s models.Student
	err := r.db.WithContext(ctx).
		Preload("Dept").
		Where("id = ? AND is_active = ?", id, true).
		First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepositoryGorm) List(ctx context.Context, filter models.StudentFilter, limit, offset int) ([]models.Student, int64, error) {
	search := strings.TrimSpace(filter.Search)
	q := r.db.WithContext(ctx).Model(&models.Student{}).Where("is_active = ?", true)
	if search != "" {
		pattern := "%" + search + "%"
		q = q.Where(
			"(email IS NOT NULL AND email ILIKE ?) OR student_id_code ILIKE ? OR (first_name || ' ' || last_name) ILIKE ?",
			pattern, pattern, pattern,
		)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []models.Student
	err := q.Preload("Dept").Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error
	return rows, total, err
}

func (r *studentRepositoryGorm) Update(ctx context.Context, s *models.Student) error {
	res := r.db.WithContext(ctx).Model(&models.Student{}).
		Where("id = ? AND is_active = ?", s.ID, true).
		Updates(map[string]interface{}{
			"first_name":               s.FirstName,
			"last_name":                s.LastName,
			"student_id_code":          s.StudentIDCode,
			"email":                    s.Email,
			"avatar_url":               s.AvatarURL,
			"department_id":            s.DepartmentID,
			"degree_level":             s.DegreeLevel,
			"major":                    s.Major,
			"expected_graduation_year": s.ExpectedGraduationYear,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *studentRepositoryGorm) Deactivate(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Model(&models.Student{}).
		Where("id = ? AND is_active = ?", id, true).
		Update("is_active", false)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *studentRepositoryGorm) ExistsByStudentIDCode(ctx context.Context, code string, excludeID string) (bool, error) {
	q := r.db.WithContext(ctx).Model(&models.Student{}).Where("student_id_code = ?", code)
	if excludeID != "" {
		q = q.Where("id != ?", excludeID)
	}
	var n int64
	if err := q.Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *studentRepositoryGorm) ExistsByEmail(ctx context.Context, email string, excludeID string) (bool, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return false, nil
	}
	q := r.db.WithContext(ctx).Model(&models.Student{}).Where("email = ?", email)
	if excludeID != "" {
		q = q.Where("id != ?", excludeID)
	}
	var n int64
	if err := q.Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *studentRepositoryGorm) EnsureStatsRow(ctx context.Context, studentID string) error {
	stats := models.StudentStats{StudentID: studentID}
	return r.db.WithContext(ctx).Where("student_id = ?", studentID).FirstOrCreate(&stats).Error
}

func (r *studentRepositoryGorm) GetActiveStudentIDByEmail(ctx context.Context, email string) (string, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return "", nil
	}
	var id string
	err := r.db.WithContext(ctx).Model(&models.Student{}).
		Select("id").
		Where("is_active = ? AND email IS NOT NULL AND LOWER(TRIM(email)) = LOWER(?)", true, email).
		Limit(1).Scan(&id).Error
	return id, err
}

func (r *studentRepositoryGorm) GetActiveStudentIDByStudentIDCode(ctx context.Context, studentIDCode string) (string, error) {
	code := strings.TrimSpace(studentIDCode)
	if code == "" {
		return "", nil
	}
	var id string
	err := r.db.WithContext(ctx).Model(&models.Student{}).
		Select("id").
		Where("is_active = ? AND student_id_code = ?", true, code).
		Limit(1).Scan(&id).Error
	return id, err
}
