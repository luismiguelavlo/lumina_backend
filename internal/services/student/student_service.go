package student

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// badgeEarnedReputationHook is optional; when set, AwardBadge records a badge_earned reputation event once per student_badge row.
type badgeEarnedReputationHook interface {
	RecordBadgeEarned(ctx context.Context, studentBadgeID, studentID string)
}

type Service struct {
	students   repositories.StudentRepository
	depts      repositories.StudentDepartmentRepository
	profile    repositories.StudentProfileRepository
	badges     repositories.BadgeRepository
	studBadges repositories.StudentBadgeRepository
	rep        badgeEarnedReputationHook
	activityLog repositories.ActivityLogRepository
}

// NewStudentService constructs StudentService. rep may be nil (tests).
func NewStudentService(
	students repositories.StudentRepository,
	depts repositories.StudentDepartmentRepository,
	profile repositories.StudentProfileRepository,
	badges repositories.BadgeRepository,
	studBadges repositories.StudentBadgeRepository,
	rep badgeEarnedReputationHook,
) *Service {
	return &Service{
		students:   students,
		depts:      depts,
		profile:    profile,
		badges:     badges,
		studBadges: studBadges,
		rep:        rep,
	}
}

// SetActivityLog configures optional activity log persistence for dashboard feed events.
func (s *Service) SetActivityLog(activityLog repositories.ActivityLogRepository) {
	s.activityLog = activityLog
}

// Create registers a new active student.
func (s *Service) Create(ctx context.Context, req models.CreateStudentRequest, adminID string) (*models.StudentResponse, error) {
	if req.DepartmentID != nil {
		ok, err := s.depts.ExistsByID(ctx, *req.DepartmentID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrDepartmentNotFound
		}
	}
	code := strings.TrimSpace(req.StudentIDCode)
	if code == "" {
		var err error
		code, err = s.allocateAutoStudentIDCode(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		dup, err := s.students.ExistsByStudentIDCode(ctx, code, "")
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, ErrDuplicateStudentIDCode
		}
	}
	if req.Email != nil {
		em := strings.TrimSpace(*req.Email)
		if em != "" {
			dupEmail, err := s.students.ExistsByEmail(ctx, em, "")
			if err != nil {
				return nil, err
			}
			if dupEmail {
				return nil, ErrDuplicateStudentEmail
			}
		}
	}

	adminPtr := adminID
	st := &models.Student{
		ID:                     uuid.NewString(),
		StudentIDCode:          code,
		FirstName:              strings.TrimSpace(req.FirstName),
		LastName:               strings.TrimSpace(req.LastName),
		Email:                  normalizeOptionalEmail(req.Email),
		AvatarURL:              req.AvatarURL,
		DepartmentID:           req.DepartmentID,
		DegreeLevel:            req.DegreeLevel,
		Major:                  req.Major,
		ExpectedGraduationYear: req.ExpectedGraduationYear,
		IsActive:               true,
		MemberSince:            time.Now().UTC(),
		RegisteredBy:           &adminPtr,
	}
	if err := s.students.Create(ctx, st); err != nil {
		return nil, err
	}
	_ = s.students.EnsureStatsRow(ctx, st.ID)
	s.logStudentCreated(ctx, st, adminID)
	return s.GetByID(ctx, st.ID)
}

func (s *Service) logStudentCreated(ctx context.Context, st *models.Student, actorUserID string) {
	if s.activityLog == nil || st == nil {
		return
	}
	meta, err := json.Marshal(map[string]any{
		"student_id":      st.ID,
		"student_id_code": st.StudentIDCode,
		"created_at":      st.CreatedAt.UTC().Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("student activity_log metadata: %v", err)
		return
	}
	entry := &models.ActivityLog{
		ID:        uuid.NewString(),
		EventType: "student_created",
		Title:     "Estudiante creado",
		Metadata:  meta,
	}
	desc := "Nuevo estudiante registrado"
	entry.Description = &desc
	sid := st.ID
	entry.StudentID = &sid
	aid := strings.TrimSpace(actorUserID)
	if aid != "" {
		entry.ActorID = &aid
	}
	if err := s.activityLog.Create(ctx, entry); err != nil {
		log.Printf("student activity_log create: %v", err)
	}
}

// GetByID returns an active student or ErrStudentNotFound.
func (s *Service) GetByID(ctx context.Context, id string) (*models.StudentResponse, error) {
	st, err := s.students.GetActiveByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, ErrStudentNotFound
	}
	return studentToResponse(st), nil
}

// List returns paginated active students.
func (s *Service) List(ctx context.Context, filter models.StudentFilter, limit, offset int) ([]models.StudentResponse, int64, error) {
	limit, offset = normalizeStudentListParams(limit, offset)
	rows, total, err := s.students.List(ctx, filter, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out := make([]models.StudentResponse, 0, len(rows))
	for i := range rows {
		out = append(out, *studentToResponse(&rows[i]))
	}
	return out, total, nil
}

// Update patches an active student.
func (s *Service) Update(ctx context.Context, id string, req models.UpdateStudentRequest) (*models.StudentResponse, error) {
	st, err := s.students.GetActiveByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, ErrStudentNotFound
	}

	if req.DepartmentID != nil {
		ok, err := s.depts.ExistsByID(ctx, *req.DepartmentID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrDepartmentNotFound
		}
		st.DepartmentID = req.DepartmentID
	}
	if req.FirstName != nil {
		st.FirstName = strings.TrimSpace(*req.FirstName)
	}
	if req.LastName != nil {
		st.LastName = strings.TrimSpace(*req.LastName)
	}
	if req.StudentIDCode != nil {
		code := strings.TrimSpace(*req.StudentIDCode)
		dup, err := s.students.ExistsByStudentIDCode(ctx, code, id)
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, ErrDuplicateStudentIDCode
		}
		st.StudentIDCode = code
	}
	if req.Email != nil {
		em := strings.TrimSpace(*req.Email)
		if em == "" {
			st.Email = nil
		} else {
			dup, err := s.students.ExistsByEmail(ctx, em, id)
			if err != nil {
				return nil, err
			}
			if dup {
				return nil, ErrDuplicateStudentEmail
			}
			st.Email = &em
		}
	}
	if req.AvatarURL != nil {
		st.AvatarURL = req.AvatarURL
	}
	if req.DegreeLevel != nil {
		st.DegreeLevel = req.DegreeLevel
	}
	if req.Major != nil {
		st.Major = req.Major
	}
	if req.ExpectedGraduationYear != nil {
		st.ExpectedGraduationYear = req.ExpectedGraduationYear
	}

	if err := s.students.Update(ctx, st); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudentNotFound
		}
		return nil, err
	}
	return s.GetByID(ctx, id)
}

// Deactivate sets is_active = false.
func (s *Service) Deactivate(ctx context.Context, id string) error {
	err := s.students.Deactivate(ctx, id)
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrStudentNotFound
	}
	return err
}

// GetProfile returns dashboard data for an active student.
func (s *Service) GetProfile(ctx context.Context, id string, loanLimit int) (*models.StudentProfileResponse, error) {
	st, err := s.students.GetActiveByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, ErrStudentNotFound
	}
	if loanLimit <= 0 {
		loanLimit = 10
	}
	if loanLimit > 50 {
		loanLimit = 50
	}

	stats, err := s.profile.GetPersonalStats(ctx, id)
	if err != nil {
		return nil, err
	}
	loans, err := s.profile.ListLoanHistory(ctx, id, loanLimit)
	if err != nil {
		return nil, err
	}
	if loans == nil {
		loans = []models.LoanHistoryItem{}
	}
	gallery, err := s.profile.ListBadgeGallery(ctx, id)
	if err != nil {
		return nil, err
	}

	var deptName, deptCode *string
	if st.Dept != nil {
		deptName = &st.Dept.Name
		deptCode = &st.Dept.Code
	}

	return &models.StudentProfileResponse{
		ID:                     st.ID,
		StudentIDCode:          st.StudentIDCode,
		FirstName:              st.FirstName,
		LastName:               st.LastName,
		AvatarURL:              st.AvatarURL,
		DegreeLevel:            st.DegreeLevel,
		Major:                  st.Major,
		MemberSince:            st.MemberSince.UTC().Format(time.RFC3339),
		DepartmentName:         deptName,
		DepartmentCode:         deptCode,
		Email:                  st.Email,
		ExpectedGraduationYear: st.ExpectedGraduationYear,
		PersonalStats:          stats,
		LoanHistory:            loans,
		BadgeGallery:           gallery,
	}, nil
}

// AwardBadge assigns a badge to an active student.
func (s *Service) AwardBadge(ctx context.Context, studentID, badgeID string) (*models.AwardBadgeResponse, error) {
	st, err := s.students.GetActiveByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, ErrStudentNotFound
	}
	b, err := s.badges.GetByID(ctx, badgeID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrBadgeNotFound
	}
	has, err := s.studBadges.Exists(ctx, studentID, badgeID)
	if err != nil {
		return nil, err
	}
	if has {
		return nil, ErrBadgeAlreadyEarned
	}
	sb := &models.StudentBadge{
		ID:        uuid.NewString(),
		StudentID: studentID,
		BadgeID:   badgeID,
	}
	if err := s.studBadges.Create(ctx, sb); err != nil {
		return nil, err
	}
	if s.rep != nil {
		s.rep.RecordBadgeEarned(ctx, sb.ID, studentID)
	}
	earnedAt := sb.EarnedAt
	if earnedAt.IsZero() {
		earnedAt = time.Now().UTC()
	}
	return &models.AwardBadgeResponse{
		ID:       b.ID,
		Slug:     b.Slug,
		Name:     b.Name,
		EarnedAt: earnedAt.UTC().Format(time.RFC3339),
	}, nil
}

// RevokeBadge removes an earned badge from an active student.
func (s *Service) RevokeBadge(ctx context.Context, studentID, badgeID string) error {
	st, err := s.students.GetActiveByID(ctx, studentID)
	if err != nil {
		return err
	}
	if st == nil {
		return ErrStudentNotFound
	}
	b, err := s.badges.GetByID(ctx, badgeID)
	if err != nil {
		return err
	}
	if b == nil {
		return ErrBadgeNotFound
	}
	n, err := s.studBadges.DeleteByStudentAndBadge(ctx, studentID, badgeID)
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrStudentDoesNotHaveBadge
	}
	return nil
}

const autoStudentIDMaxAttempts = 10

func nextAutoStudentIDCode() string {
	d := time.Now().UTC().Format("20060102")
	raw := strings.ReplaceAll(uuid.NewString(), "-", "")
	return "AUTO-" + d + "-" + raw[:8]
}

func (s *Service) allocateAutoStudentIDCode(ctx context.Context) (string, error) {
	for i := 0; i < autoStudentIDMaxAttempts; i++ {
		c := nextAutoStudentIDCode()
		dup, err := s.students.ExistsByStudentIDCode(ctx, c, "")
		if err != nil {
			return "", err
		}
		if !dup {
			return c, nil
		}
	}
	return "", ErrStudentIDGenerationFailed
}

func normalizeOptionalEmail(e *string) *string {
	if e == nil {
		return nil
	}
	t := strings.TrimSpace(*e)
	if t == "" {
		return nil
	}
	return &t
}

func studentToResponse(st *models.Student) *models.StudentResponse {
	var deptName, deptCode *string
	if st.Dept != nil {
		deptName = &st.Dept.Name
		deptCode = &st.Dept.Code
	}
	return &models.StudentResponse{
		ID:                     st.ID,
		StudentIDCode:          st.StudentIDCode,
		FirstName:              st.FirstName,
		LastName:               st.LastName,
		Email:                  st.Email,
		AvatarURL:              st.AvatarURL,
		DepartmentID:           st.DepartmentID,
		DepartmentName:         deptName,
		DepartmentCode:         deptCode,
		DegreeLevel:            st.DegreeLevel,
		Major:                  st.Major,
		ExpectedGraduationYear: st.ExpectedGraduationYear,
		IsActive:               st.IsActive,
		MemberSince:            st.MemberSince.UTC().Format(time.RFC3339),
		CreatedAt:              st.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:              st.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func normalizeStudentListParams(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
