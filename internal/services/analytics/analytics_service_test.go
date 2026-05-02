package analytics

import (
	"context"
	"errors"
	"testing"
	"time"

	"library_back/internal/models"
)

type fakeAnalyticsRepo struct {
	totalBooks     int64
	activeStudents int64
	fines          float64
	top            []models.MostBorrowedItem
	activity       []models.RecentActivityItem
	err            error
	lastMB         int
	lastRA         int
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

func TestService_GetDashboardDefaultsAndCaps(t *testing.T) {
	fr := &fakeAnalyticsRepo{
		totalBooks: 10, activeStudents: 3, fines: 12.345,
		top:      []models.MostBorrowedItem{{BookID: "b1", Title: "T", CatalogCode: "C", BorrowCount: 2}},
		activity: []models.RecentActivityItem{{ID: "a1", EventType: "x", Title: "y", CreatedAt: time.Now().UTC()}},
	}
	s := NewService(fr)
	out, err := s.GetDashboard(context.Background(), models.DashboardLimits{MostBorrowedLimit: 0, RecentActivityLimit: -1})
	if err != nil {
		t.Fatal(err)
	}
	if out.TotalBooks != 10 || out.ActiveStudents != 3 || out.OverdueFines != 12.35 {
		t.Fatalf("%+v", out)
	}
	if fr.lastMB != defaultMostBorrowed || fr.lastRA != defaultRecentActivity {
		t.Fatalf("limits mb=%d ra=%d", fr.lastMB, fr.lastRA)
	}

	fr2 := &fakeAnalyticsRepo{}
	s2 := NewService(fr2)
	_, _ = s2.GetDashboard(context.Background(), models.DashboardLimits{MostBorrowedLimit: 999, RecentActivityLimit: 999})
	if fr2.lastMB != maxMostBorrowed || fr2.lastRA != maxRecentActivity {
		t.Fatalf("cap mb=%d ra=%d", fr2.lastMB, fr2.lastRA)
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
