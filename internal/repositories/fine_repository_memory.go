package repositories

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"library_back/internal/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// FineRepositoryInMemory is an in-memory FineRepository for tests.
type FineRepositoryInMemory struct {
	mu    sync.Mutex
	fines []models.Fine
}

// NewFineRepositoryInMemory returns an empty in-memory store.
func NewFineRepositoryInMemory() *FineRepositoryInMemory {
	return &FineRepositoryInMemory{fines: make([]models.Fine, 0)}
}

func (r *FineRepositoryInMemory) Create(ctx context.Context, f *models.Fine) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	if f.CreatedAt.IsZero() {
		f.CreatedAt = now
	}
	if f.Status == "" {
		f.Status = models.FineStatusPending
	}
	cp := *f
	r.fines = append(r.fines, cp)
	return nil
}

func (r *FineRepositoryInMemory) GetByID(ctx context.Context, id string) (*models.Fine, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.fines {
		if r.fines[i].ID == id {
			cp := r.fines[i]
			return &cp, nil
		}
	}
	return nil, nil
}

func (r *FineRepositoryInMemory) List(ctx context.Context, filter FineListFilter, limit, offset int) ([]FineWithStudent, int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	sid := strings.TrimSpace(filter.StudentID)
	st := strings.TrimSpace(strings.ToLower(filter.Status))
	var match []models.Fine
	for i := range r.fines {
		f := r.fines[i]
		if sid != "" && f.StudentID != sid {
			continue
		}
		if st != "" && strings.ToLower(string(f.Status)) != st {
			continue
		}
		match = append(match, f)
	}
	sort.Slice(match, func(i, j int) bool {
		return match[i].CreatedAt.After(match[j].CreatedAt)
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
	out := make([]FineWithStudent, 0, len(page))
	for i := range page {
		out = append(out, FineWithStudent{Fine: page[i], StudentName: ""})
	}
	return out, total, nil
}

func (r *FineRepositoryInMemory) UpdateStatus(ctx context.Context, id string, status models.FineStatus, paidAt *time.Time) (int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.fines {
		if r.fines[i].ID != id {
			continue
		}
		if r.fines[i].Status != models.FineStatusPending {
			return 0, nil
		}
		r.fines[i].Status = status
		if status == models.FineStatusPaid && paidAt != nil {
			t := *paidAt
			r.fines[i].PaidAt = &t
		}
		if status == models.FineStatusWaived {
			r.fines[i].PaidAt = nil
		}
		return 1, nil
	}
	return 0, nil
}

func (r *FineRepositoryInMemory) CountPendingByLoanID(ctx context.Context, loanID string) (int64, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	var n int64
	for i := range r.fines {
		if r.fines[i].LoanID == loanID && r.fines[i].Status == models.FineStatusPending {
			n++
		}
	}
	return n, nil
}

// SeedFine inserts a fine for tests (optional helper).
func (r *FineRepositoryInMemory) SeedFine(f models.Fine) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if f.ID == "" {
		f.ID = uuid.NewString()
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now().UTC()
	}
	if f.Status == "" {
		f.Status = models.FineStatusPending
	}
	if f.Amount.IsZero() {
		f.Amount = decimal.Zero
	}
	r.fines = append(r.fines, f)
}
