package analytics

import (
	"context"
	"math"

	"library_back/internal/models"
	"library_back/internal/repositories"
)

const (
	defaultMostBorrowed   = 5
	maxMostBorrowed       = 20
	defaultRecentActivity = 10
	maxRecentActivity     = 50
	defaultTopOverdue     = 5
	maxTopOverdue         = 20
)

// Service aggregates dashboard metrics (read-only).
type Service struct {
	repo repositories.AnalyticsRepository
}

// NewService constructs AnalyticsService.
func NewService(repo repositories.AnalyticsRepository) *Service {
	return &Service{repo: repo}
}

// GetDashboard loads all dashboard sections using normalized limits.
func (s *Service) GetDashboard(ctx context.Context, limits models.DashboardLimits) (*models.DashboardResponse, error) {
	mb := normalizeMostBorrowed(limits.MostBorrowedLimit)
	ra := normalizeRecentActivity(limits.RecentActivityLimit)
	to := normalizeTopOverdue(limits.TopOverdueLimit)

	totalBooks, err := s.repo.TotalBooks(ctx)
	if err != nil {
		return nil, err
	}
	active, err := s.repo.ActiveStudents(ctx)
	if err != nil {
		return nil, err
	}
	fines, err := s.repo.OverdueFinesTotal(ctx)
	if err != nil {
		return nil, err
	}
	pendingFinesCount, err := s.repo.PendingFinesCount(ctx)
	if err != nil {
		return nil, err
	}
	overdueLoans, err := s.repo.OverdueLoansCount(ctx)
	if err != nil {
		return nil, err
	}
	activeLoans, err := s.repo.ActiveLoansCount(ctx)
	if err != nil {
		return nil, err
	}
	dueSoon, err := s.repo.DueSoonCount(ctx)
	if err != nil {
		return nil, err
	}
	activeSanctions, err := s.repo.ActiveSanctionsCount(ctx)
	if err != nil {
		return nil, err
	}
	returnsWeek, err := s.repo.ReturnsThisWeekCount(ctx)
	if err != nil {
		return nil, err
	}
	newStudents, err := s.repo.NewStudentsThisMonthCount(ctx)
	if err != nil {
		return nil, err
	}
	available, checkedOut, err := s.repo.CopyAvailabilityTotals(ctx)
	if err != nil {
		return nil, err
	}
	top, err := s.repo.MostBorrowedBooks(ctx, mb)
	if err != nil {
		return nil, err
	}
	act, err := s.repo.RecentActivity(ctx, ra)
	if err != nil {
		return nil, err
	}
	overdueList, err := s.repo.TopOverdueLoans(ctx, to)
	if err != nil {
		return nil, err
	}

	return &models.DashboardResponse{
		TotalBooks:        totalBooks,
		ActiveStudents:    active,
		OverdueFines:      roundMoney(fines),
		PendingFinesCount: pendingFinesCount,
		OverdueLoans:      overdueLoans,
		ActiveLoans:       activeLoans,
		DueSoon:           dueSoon,
		ActiveSanctions:   activeSanctions,
		ReturnsThisWeek:   returnsWeek,
		NewStudentsMonth:  newStudents,
		AvailableCopies:   available,
		CheckedOutCopies:  checkedOut,
		MostBorrowedBooks: top,
		RecentActivity:    act,
		TopOverdue:        overdueList,
	}, nil
}

func normalizeMostBorrowed(n int) int {
	if n <= 0 {
		return defaultMostBorrowed
	}
	if n > maxMostBorrowed {
		return maxMostBorrowed
	}
	return n
}

func normalizeRecentActivity(n int) int {
	if n <= 0 {
		return defaultRecentActivity
	}
	if n > maxRecentActivity {
		return maxRecentActivity
	}
	return n
}

func normalizeTopOverdue(n int) int {
	if n <= 0 {
		return defaultTopOverdue
	}
	if n > maxTopOverdue {
		return maxTopOverdue
	}
	return n
}

func roundMoney(v float64) float64 {
	return math.Round(v*100) / 100
}
