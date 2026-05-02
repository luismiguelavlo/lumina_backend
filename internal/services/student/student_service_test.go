package student

import (
	"context"
	"strings"
	"testing"
	"time"

	"library_back/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakeStudentRepo struct {
	createErr        error
	lastCreated      *models.Student
	getActive        *models.Student
	getErr           error
	list             []models.Student
	listTotal        int64
	listErr          error
	updateErr        error
	deactivateErr    error
	existsCode       bool
	existsEmail      bool
	ensureStatsCalls int
}

func (f *fakeStudentRepo) Create(ctx context.Context, s *models.Student) error {
	_ = ctx
	if f.createErr != nil {
		return f.createErr
	}
	f.lastCreated = s
	f.getActive = s
	return nil
}

func (f *fakeStudentRepo) GetActiveByID(ctx context.Context, id string) (*models.Student, error) {
	_ = ctx
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.getActive != nil && f.getActive.ID == id {
		return f.getActive, nil
	}
	return nil, nil
}

func (f *fakeStudentRepo) List(ctx context.Context, filter models.StudentFilter, limit, offset int) ([]models.Student, int64, error) {
	_ = ctx
	_ = filter
	_ = limit
	_ = offset
	if f.listErr != nil {
		return nil, 0, f.listErr
	}
	return f.list, f.listTotal, nil
}

func (f *fakeStudentRepo) Update(ctx context.Context, s *models.Student) error {
	_ = ctx
	_ = s
	return f.updateErr
}

func (f *fakeStudentRepo) Deactivate(ctx context.Context, id string) error {
	_ = ctx
	_ = id
	return f.deactivateErr
}

func (f *fakeStudentRepo) ExistsByStudentIDCode(ctx context.Context, code string, excludeID string) (bool, error) {
	_ = ctx
	_ = code
	_ = excludeID
	return f.existsCode, nil
}

func (f *fakeStudentRepo) ExistsByEmail(ctx context.Context, email string, excludeID string) (bool, error) {
	_ = ctx
	_ = email
	_ = excludeID
	return f.existsEmail, nil
}

func (f *fakeStudentRepo) EnsureStatsRow(ctx context.Context, studentID string) error {
	_ = ctx
	_ = studentID
	f.ensureStatsCalls++
	return nil
}

func (f *fakeStudentRepo) GetActiveStudentIDByEmail(ctx context.Context, email string) (string, error) {
	_ = ctx
	_ = email
	return "", nil
}

func (f *fakeStudentRepo) GetActiveStudentIDByStudentIDCode(ctx context.Context, studentIDCode string) (string, error) {
	_ = ctx
	_ = studentIDCode
	return "", nil
}

type fakeDeptRepo struct {
	exists bool
	err    error
}

func (f *fakeDeptRepo) List(ctx context.Context) ([]models.Department, error) { return nil, nil }

func (f *fakeDeptRepo) ExistsByID(ctx context.Context, id string) (bool, error) {
	_ = ctx
	_ = id
	if f.err != nil {
		return false, f.err
	}
	return f.exists, nil
}

type fakeProfileRepo struct{}

func (f *fakeProfileRepo) GetPersonalStats(ctx context.Context, studentID string) (models.PersonalStats, error) {
	_ = ctx
	_ = studentID
	return models.PersonalStats{TotalRead: 1}, nil
}

func (f *fakeProfileRepo) ListLoanHistory(ctx context.Context, studentID string, limit int) ([]models.LoanHistoryItem, error) {
	_ = ctx
	_ = studentID
	_ = limit
	return []models.LoanHistoryItem{}, nil
}

func (f *fakeProfileRepo) ListBadgeGallery(ctx context.Context, studentID string) (models.BadgeGallery, error) {
	_ = ctx
	_ = studentID
	return models.BadgeGallery{TotalBadges: 2, EarnedCount: 0, Badges: []models.BadgeGalleryItem{}}, nil
}

type fakeBadgeRepo struct {
	badge *models.Badge
}

func (f *fakeBadgeRepo) GetByID(ctx context.Context, id string) (*models.Badge, error) {
	_ = ctx
	_ = id
	return f.badge, nil
}

func (f *fakeBadgeRepo) GetBySlug(ctx context.Context, slug string) (*models.Badge, error) {
	_ = ctx
	_ = slug
	return nil, nil
}

func (f *fakeBadgeRepo) ListAll(ctx context.Context) ([]models.Badge, error) {
	_ = ctx
	return nil, nil
}

func (f *fakeBadgeRepo) Create(ctx context.Context, b *models.Badge) error {
	_ = ctx
	_ = b
	return nil
}

func (f *fakeBadgeRepo) Save(ctx context.Context, b *models.Badge) error {
	_ = ctx
	_ = b
	return nil
}

func (f *fakeBadgeRepo) Delete(ctx context.Context, id string) error {
	_ = ctx
	_ = id
	return nil
}

func (f *fakeBadgeRepo) CountStudentBadgesByBadgeID(ctx context.Context, badgeID string) (int64, error) {
	_ = ctx
	_ = badgeID
	return 0, nil
}

type fakeStudBadgeRepo struct {
	exists       bool
	err          error
	deleteRows   int64
	deleteErr    error
}

func (f *fakeStudBadgeRepo) Exists(ctx context.Context, studentID, badgeID string) (bool, error) {
	_ = ctx
	_ = studentID
	_ = badgeID
	return f.exists, f.err
}

func (f *fakeStudBadgeRepo) Create(ctx context.Context, sb *models.StudentBadge) error {
	_ = ctx
	sb.EarnedAt = time.Now().UTC()
	return nil
}

func (f *fakeStudBadgeRepo) DeleteByStudentAndBadge(ctx context.Context, studentID, badgeID string) (int64, error) {
	_ = ctx
	_ = studentID
	_ = badgeID
	if f.deleteErr != nil {
		return 0, f.deleteErr
	}
	if f.deleteRows > 0 {
		return f.deleteRows, nil
	}
	if f.exists {
		return 1, nil
	}
	return 0, nil
}

func TestStudentService_CreateDepartmentMissing(t *testing.T) {
	did := uuid.NewString()
	s := NewStudentService(
		&fakeStudentRepo{},
		&fakeDeptRepo{exists: false},
		&fakeProfileRepo{},
		&fakeBadgeRepo{},
		&fakeStudBadgeRepo{},
		nil,
	)
	_, err := s.Create(context.Background(), models.CreateStudentRequest{
		FirstName: "A", LastName: "B", StudentIDCode: "X-1",
		DepartmentID: &did,
	}, "admin")
	if err != ErrDepartmentNotFound {
		t.Fatalf("want ErrDepartmentNotFound got %v", err)
	}
}

func TestStudentService_CreateDuplicateCode(t *testing.T) {
	s := NewStudentService(
		&fakeStudentRepo{existsCode: true},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{},
		&fakeStudBadgeRepo{},
		nil,
	)
	_, err := s.Create(context.Background(), models.CreateStudentRequest{
		FirstName: "A", LastName: "B", StudentIDCode: "DUP",
	}, "admin")
	if err != ErrDuplicateStudentIDCode {
		t.Fatalf("want ErrDuplicateStudentIDCode got %v", err)
	}
}

func TestStudentService_CreateDuplicateEmail(t *testing.T) {
	em := "x@test.com"
	s := NewStudentService(
		&fakeStudentRepo{existsEmail: true},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{},
		&fakeStudBadgeRepo{},
		nil,
	)
	_, err := s.Create(context.Background(), models.CreateStudentRequest{
		FirstName: "A", LastName: "B", StudentIDCode: "U-1",
		Email: &em,
	}, "admin")
	if err != ErrDuplicateStudentEmail {
		t.Fatalf("want ErrDuplicateStudentEmail got %v", err)
	}
}

func TestStudentService_CreateSetsRegisteredByAndStats(t *testing.T) {
	fr := &fakeStudentRepo{}
	st := NewStudentService(
		fr,
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{},
		&fakeStudBadgeRepo{},
		nil,
	)
	out, err := st.Create(context.Background(), models.CreateStudentRequest{
		FirstName: "A", LastName: "B", StudentIDCode: "L-1",
	}, "staff-user-1")
	if err != nil {
		t.Fatal(err)
	}
	if fr.lastCreated == nil || out.ID != fr.lastCreated.ID {
		t.Fatalf("response id")
	}
	if fr.lastCreated.RegisteredBy == nil || *fr.lastCreated.RegisteredBy != "staff-user-1" {
		t.Fatalf("registered_by not set from adminID")
	}
	if !fr.lastCreated.IsActive {
		t.Fatalf("is_active")
	}
	if fr.ensureStatsCalls != 1 {
		t.Fatalf("EnsureStatsRow calls = %d", fr.ensureStatsCalls)
	}
}

func TestStudentService_CreateAutoGeneratesStudentIDCode(t *testing.T) {
	fr := &fakeStudentRepo{}
	st := NewStudentService(
		fr,
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{},
		&fakeStudBadgeRepo{},
		nil,
	)
	out, err := st.Create(context.Background(), models.CreateStudentRequest{
		FirstName: "A", LastName: "B",
	}, "admin")
	if err != nil {
		t.Fatal(err)
	}
	if fr.lastCreated == nil || fr.lastCreated.StudentIDCode == "" {
		t.Fatal("expected generated student_id_code")
	}
	if !strings.HasPrefix(fr.lastCreated.StudentIDCode, "AUTO-") {
		t.Fatalf("want AUTO- prefix, got %q", fr.lastCreated.StudentIDCode)
	}
	if out.StudentIDCode != fr.lastCreated.StudentIDCode {
		t.Fatalf("response code mismatch")
	}
}

func TestStudentService_GetByIDNotFound(t *testing.T) {
	s := NewStudentService(
		&fakeStudentRepo{},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{},
		&fakeStudBadgeRepo{},
		nil,
	)
	_, err := s.GetByID(context.Background(), uuid.NewString())
	if err != ErrStudentNotFound {
		t.Fatalf("want ErrStudentNotFound got %v", err)
	}
}

func TestStudentService_ListNormalizesPagination(t *testing.T) {
	s := NewStudentService(
		&fakeStudentRepo{list: []models.Student{}, listTotal: 0},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{},
		&fakeStudBadgeRepo{},
		nil,
	)
	_, _, err := s.List(context.Background(), models.StudentFilter{}, 0, -1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStudentService_UpdateStudentMissing(t *testing.T) {
	s := NewStudentService(
		&fakeStudentRepo{},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{},
		&fakeStudBadgeRepo{},
		nil,
	)
	nm := "N"
	_, err := s.Update(context.Background(), uuid.NewString(), models.UpdateStudentRequest{FirstName: &nm})
	if err != ErrStudentNotFound {
		t.Fatalf("want ErrStudentNotFound got %v", err)
	}
}

func TestStudentService_UpdatePersistNotFound(t *testing.T) {
	s := NewStudentService(
		&fakeStudentRepo{
			getActive: &models.Student{
				ID: "x", StudentIDCode: "c", FirstName: "a", LastName: "b",
				IsActive: true, MemberSince: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
			},
			updateErr: gorm.ErrRecordNotFound,
		},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{},
		&fakeStudBadgeRepo{},
		nil,
	)
	nm := "N"
	_, err := s.Update(context.Background(), "x", models.UpdateStudentRequest{FirstName: &nm})
	if err != ErrStudentNotFound {
		t.Fatalf("want ErrStudentNotFound got %v", err)
	}
}

func TestStudentService_DeactivateNotFound(t *testing.T) {
	s := NewStudentService(
		&fakeStudentRepo{deactivateErr: gorm.ErrRecordNotFound},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{},
		&fakeStudBadgeRepo{},
		nil,
	)
	err := s.Deactivate(context.Background(), uuid.NewString())
	if err != ErrStudentNotFound {
		t.Fatalf("want ErrStudentNotFound got %v", err)
	}
}

func TestStudentService_AwardBadgeNotFound(t *testing.T) {
	s := NewStudentService(
		&fakeStudentRepo{getActive: nil},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{badge: nil},
		&fakeStudBadgeRepo{},
		nil,
	)
	_, err := s.AwardBadge(context.Background(), uuid.NewString(), uuid.NewString())
	if err != ErrStudentNotFound {
		t.Fatalf("want ErrStudentNotFound got %v", err)
	}
}

func TestStudentService_AwardBadgeBadgeMissing(t *testing.T) {
	sid := uuid.NewString()
	s := NewStudentService(
		&fakeStudentRepo{getActive: &models.Student{ID: sid, StudentIDCode: "s", FirstName: "a", LastName: "b", IsActive: true, MemberSince: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{badge: nil},
		&fakeStudBadgeRepo{},
		nil,
	)
	_, err := s.AwardBadge(context.Background(), sid, uuid.NewString())
	if err != ErrBadgeNotFound {
		t.Fatalf("want ErrBadgeNotFound got %v", err)
	}
}

func TestStudentService_AwardBadgeAlreadyEarned(t *testing.T) {
	sid := uuid.NewString()
	bid := uuid.NewString()
	s := NewStudentService(
		&fakeStudentRepo{getActive: &models.Student{ID: sid, StudentIDCode: "s", FirstName: "a", LastName: "b", IsActive: true, MemberSince: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{badge: &models.Badge{ID: bid, Slug: "x", Name: "X"}},
		&fakeStudBadgeRepo{exists: true},
		nil,
	)
	_, err := s.AwardBadge(context.Background(), sid, bid)
	if err != ErrBadgeAlreadyEarned {
		t.Fatalf("want ErrBadgeAlreadyEarned got %v", err)
	}
}

func TestStudentService_RevokeBadgeSuccess(t *testing.T) {
	sid := uuid.NewString()
	bid := uuid.NewString()
	s := NewStudentService(
		&fakeStudentRepo{getActive: &models.Student{ID: sid, StudentIDCode: "s", FirstName: "a", LastName: "b", IsActive: true, MemberSince: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{badge: &models.Badge{ID: bid, Slug: "x", Name: "X"}},
		&fakeStudBadgeRepo{deleteRows: 1},
		nil,
	)
	if err := s.RevokeBadge(context.Background(), sid, bid); err != nil {
		t.Fatal(err)
	}
}

func TestStudentService_RevokeBadgeNotEarned(t *testing.T) {
	sid := uuid.NewString()
	bid := uuid.NewString()
	s := NewStudentService(
		&fakeStudentRepo{getActive: &models.Student{ID: sid, StudentIDCode: "s", FirstName: "a", LastName: "b", IsActive: true, MemberSince: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{badge: &models.Badge{ID: bid, Slug: "x", Name: "X"}},
		&fakeStudBadgeRepo{},
		nil,
	)
	err := s.RevokeBadge(context.Background(), sid, bid)
	if err != ErrStudentDoesNotHaveBadge {
		t.Fatalf("want ErrStudentDoesNotHaveBadge got %v", err)
	}
}

func TestStudentService_GetProfile(t *testing.T) {
	sid := uuid.NewString()
	s := NewStudentService(
		&fakeStudentRepo{getActive: &models.Student{
			ID: sid, StudentIDCode: "s", FirstName: "a", LastName: "b",
			IsActive: true, MemberSince: time.Now(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{},
		&fakeStudBadgeRepo{},
		nil,
	)
	p, err := s.GetProfile(context.Background(), sid, 5)
	if err != nil {
		t.Fatal(err)
	}
	if p.PersonalStats.TotalRead != 1 {
		t.Fatalf("stats")
	}
	if p.BadgeGallery.TotalBadges != 2 {
		t.Fatalf("gallery")
	}
}

func TestStudentService_GetProfileNotFound(t *testing.T) {
	s := NewStudentService(
		&fakeStudentRepo{},
		&fakeDeptRepo{},
		&fakeProfileRepo{},
		&fakeBadgeRepo{},
		&fakeStudBadgeRepo{},
		nil,
	)
	_, err := s.GetProfile(context.Background(), uuid.NewString(), 10)
	if err != ErrStudentNotFound {
		t.Fatalf("want ErrStudentNotFound got %v", err)
	}
}

func TestNormalizeStudentListParams(t *testing.T) {
	l, o := normalizeStudentListParams(0, -5)
	if l != 20 || o != 0 {
		t.Fatalf("defaults")
	}
	l, o = normalizeStudentListParams(200, 0)
	if l != 100 {
		t.Fatalf("max cap")
	}
}
