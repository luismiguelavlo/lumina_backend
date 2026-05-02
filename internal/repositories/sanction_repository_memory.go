package repositories

import (
	"context"
	"sort"
	"sync"
	"time"

	"library_back/internal/models"

	"github.com/google/uuid"
)

// SanctionRepositoryInMemory is an in-memory SanctionRepository for tests.
type SanctionRepositoryInMemory struct {
	mu        sync.Mutex
	sanctions []models.Sanction
}

// NewSanctionRepositoryInMemory returns an empty in-memory store.
func NewSanctionRepositoryInMemory() *SanctionRepositoryInMemory {
	return &SanctionRepositoryInMemory{sanctions: make([]models.Sanction, 0)}
}

func (r *SanctionRepositoryInMemory) Create(ctx context.Context, s *models.Sanction) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	if s.AppliedAt.IsZero() {
		s.AppliedAt = now
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
	}
	if s.Status == "" {
		s.Status = models.SanctionStatusActive
	}
	cp := *s
	r.sanctions = append(r.sanctions, cp)
	return nil
}

func (r *SanctionRepositoryInMemory) GetByID(ctx context.Context, id string) (*models.Sanction, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.sanctions {
		if r.sanctions[i].ID == id {
			cp := r.sanctions[i]
			return &cp, nil
		}
	}
	return nil, nil
}

func (r *SanctionRepositoryInMemory) ListActive(ctx context.Context, filter SanctionListFilter, limit, offset int) ([]SanctionListRow, int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	var match []models.Sanction
	for i := range r.sanctions {
		if r.sanctions[i].Status == models.SanctionStatusActive &&
			(filter.StudentID == "" || r.sanctions[i].StudentID == filter.StudentID) {
			match = append(match, r.sanctions[i])
		}
	}
	sort.Slice(match, func(i, j int) bool {
		return match[i].AppliedAt.After(match[j].AppliedAt)
	})
	total := int64(len(match))
	if offset > len(match) {
		offset = len(match)
	}
	end := offset + limit
	if end > len(match) {
		end = len(match)
	}
	page := match[offset:end]
	out := make([]SanctionListRow, 0, len(page))
	for i := range page {
		out = append(out, SanctionListRow{
			SanctionID:    page[i].ID,
			StudentID:     page[i].StudentID,
			StudentIDCode: "MEM-STUDENT",
			FirstName:     "Mem",
			LastName:      "Student",
			Email:         nil,
			Reason:        page[i].Reason,
			Status:        string(page[i].Status),
			AppliedAt:     page[i].AppliedAt,
			AppliedByID:   page[i].AppliedBy,
		})
	}
	return out, total, nil
}

func (r *SanctionRepositoryInMemory) CountActiveByStudentID(ctx context.Context, studentID string) (int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	var n int64
	for i := range r.sanctions {
		if r.sanctions[i].StudentID == studentID && r.sanctions[i].Status == models.SanctionStatusActive {
			n++
		}
	}
	return n, nil
}

func (r *SanctionRepositoryInMemory) Lift(ctx context.Context, id string, liftedAt time.Time, liftedBy string) (int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.sanctions {
		if r.sanctions[i].ID != id {
			continue
		}
		if r.sanctions[i].Status != models.SanctionStatusActive {
			return 0, nil
		}
		r.sanctions[i].Status = models.SanctionStatusLifted
		r.sanctions[i].LiftedAt = &liftedAt
		lb := liftedBy
		r.sanctions[i].LiftedBy = &lb
		return 1, nil
	}
	return 0, nil
}
