package analytics

import (
	"context"
	"errors"
	"testing"
	"time"

	"library_back/internal/models"
)

type fakeAnalyticsRepo struct {
	totalBooks       int64
	activeStudents   int64
	fines            float64
	pendingFines     int64
	overdueLoans     int64
	activeLoans      int64
	dueSoon          int64
	activeSanctions  int64
	returnsWeek      int64
	newStudentsMonth int64
	availableCopies  int64
	checkedOut       int64
	top              []models.MostBorrowedItem
	activity         []models.RecentActivityItem
	topOverdue       []models.TopOverdueItem
	err              error
	lastMB           int
	lastRA           int
	lastTO           int
}

func (f *fakeAnalyticsRepo) TotalBooks(ctx context.Context) (int64, error) {
	_ = ctx
	if f.err != nil {
		return 0, f.err
	}
	return f.totalBooks, nil
}

func (f *fakeAnalyticsRepo) ActiveStudents(ctx context.Context) (int64, error) {
	_ = ctx
	if f.err != nil {
		return 0, f.err
	}
	return f.activeStudents, nil
}

func (f *fakeAnalyticsRepo) OverdueFinesTotal(ctx context.Context) (float64, error) {
	_ = ctx
	if f.err != nil {
		return 0, f.err
	}
	return f.fines, nil
}

func (f *fakeAnalyticsRepo) PendingFinesCount(ctx context.Context) (int64, error) {
	_ = ctx
	if f.err != nil {
		return 0, f.err
	}
	return f.pendingFines, nil
}

func (f *fakeAnalyticsRepo) OverdueLoansCount(ctx context.Context) (int64, error) {
	_ = ctx
	if f.err != nil {
		return 0, f.err
	}
	return f.overdueLoans, nil
}

func (f *fakeAnalyticsRepo) ActiveLoansCount(ctx context.Context) (int64, error) {
	_ = ctx
	if f.err != nil {
		return 0, f.err
	}
	return f.activeLoans, nil
}

func (f *fakeAnalyticsRepo) DueSoonCount(ctx context.Context) (int64, error) {
	_ = ctx
	if f.err != nil {
		return 0, f.err
	}
	return f.dueSoon, nil
}

func (f *fakeAnalyticsRepo) ActiveSanctionsCount(ctx context.Context) (int64, error) {
	_ = ctx
	if f.err != nil {
		return 0, f.err
	}
	return f.activeSanctions, nil
}

func (f *fakeAnalyticsRepo) ReturnsThisWeekCount(ctx context.Context) (int64, error) {
	_ = ctx
	if f.err != nil {
		return 0, f.err
	}
	return f.returnsWeek, nil
}

func (f *fakeAnalyticsRepo) NewStudentsThisMonthCount(ctx context.Context) (int64, error) {
	_ = ctx
	if f.err != nil {
		return 0, f.err
	}
	return f.newStudentsMonth, nil
}

func (f *fakeAnalyticsRepo) CopyAvailabilityTotals(ctx context.Context) (available, checkedOut int64, err error) {
	_ = ctx
	if f.err != nil {
		return 0, 0, f.err
	}
	return f.availableCopies, f.checkedOut, nil
}

func (f *fakeAnalyticsRepo) MostBorrowedBooks(ctx context.Context, limit int) ([]models.MostBorrowedItem, error) {
	_ = ctx
	f.lastMB = limit
	if f.err != nil {
		return nil, f.err
	}
	return f.top, nil
}

func (f *fakeAnalyticsRepo) RecentActivity(ctx context.Context, limit int) ([]models.RecentActivityItem, error) {
	_ = ctx
	f.lastRA = limit
	if f.err != nil {
		return nil, f.err
	}
	return f.activity, nil
}

func (f *fakeAnalyticsRepo) TopOverdueLoans(ctx context.Context, limit int) ([]models.TopOverdueItem, error) {
	_ = ctx
	f.lastTO = limit
	if f.err != nil {
		return nil, f.err
	}
	return f.topOverdue, nil
}

func TestService_GetDashboardDefaultsAndCaps(t *testing.T) {
	fr := &fakeAnalyticsRepo{
		totalBooks: 10, activeStudents: 3, fines: 12.345,
		pendingFines: 2, overdueLoans: 4, activeLoans: 9, dueSoon: 3,
		activeSanctions: 1, returnsWeek: 7, newStudentsMonth: 5,
		availableCopies: 40, checkedOut: 12,
		top:        []models.MostBorrowedItem{{BookID: "b1", Title: "T", CatalogCode: "C", BorrowCount: 2}},
		activity:   []models.RecentActivityItem{{ID: "a1", EventType: "x", Title: "y", CreatedAt: time.Now().UTC()}},
		topOverdue: []models.TopOverdueItem{{LoanID: "l1", StudentName: "Ana", BookTitle: "Book", DaysOverdue: 3}},
	}
	s := NewService(fr)
	out, err := s.GetDashboard(context.Background(), models.DashboardLimits{MostBorrowedLimit: 0, RecentActivityLimit: -1})
	if err != nil {
		t.Fatal(err)
	}
	if out.TotalBooks != 10 || out.ActiveStudents != 3 || out.OverdueFines != 12.35 {
		t.Fatalf("%+v", out)
	}
	if out.PendingFinesCount != 2 || out.OverdueLoans != 4 || out.ActiveLoans != 9 || out.DueSoon != 3 {
		t.Fatalf("ops metrics: %+v", out)
	}
	if out.ActiveSanctions != 1 || out.ReturnsThisWeek != 7 || out.NewStudentsMonth != 5 {
		t.Fatalf("engagement metrics: %+v", out)
	}
	if out.AvailableCopies != 40 || out.CheckedOutCopies != 12 || len(out.TopOverdue) != 1 {
		t.Fatalf("inventory/overdue: %+v", out)
	}
	if fr.lastMB != defaultMostBorrowed || fr.lastRA != defaultRecentActivity || fr.lastTO != defaultTopOverdue {
		t.Fatalf("limits mb=%d ra=%d to=%d", fr.lastMB, fr.lastRA, fr.lastTO)
	}

	fr2 := &fakeAnalyticsRepo{}
	s2 := NewService(fr2)
	_, _ = s2.GetDashboard(context.Background(), models.DashboardLimits{
		MostBorrowedLimit: 999, RecentActivityLimit: 999, TopOverdueLimit: 999,
	})
	if fr2.lastMB != maxMostBorrowed || fr2.lastRA != maxRecentActivity || fr2.lastTO != maxTopOverdue {
		t.Fatalf("cap mb=%d ra=%d to=%d", fr2.lastMB, fr2.lastRA, fr2.lastTO)
	}
}

func TestService_GetDashboardPropagatesError(t *testing.T) {
	fr := &fakeAnalyticsRepo{err: errors.New("db")}
	s := NewService(fr)
	_, err := s.GetDashboard(context.Background(), models.DashboardLimits{})
	if err == nil {
		t.Fatal("want error")
	}
}
