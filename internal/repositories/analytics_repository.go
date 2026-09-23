package repositories

import (
	"context"

	"library_back/internal/models"
)

// AnalyticsRepository runs read-only dashboard aggregations.
type AnalyticsRepository interface {
	TotalBooks(ctx context.Context) (int64, error)
	ActiveStudents(ctx context.Context) (int64, error)
	OverdueFinesTotal(ctx context.Context) (float64, error)
	PendingFinesCount(ctx context.Context) (int64, error)
	OverdueLoansCount(ctx context.Context) (int64, error)
	ActiveLoansCount(ctx context.Context) (int64, error)
	DueSoonCount(ctx context.Context) (int64, error)
	ActiveSanctionsCount(ctx context.Context) (int64, error)
	ReturnsThisWeekCount(ctx context.Context) (int64, error)
	NewStudentsThisMonthCount(ctx context.Context) (int64, error)
	CopyAvailabilityTotals(ctx context.Context) (available, checkedOut int64, err error)
	MostBorrowedBooks(ctx context.Context, limit int) ([]models.MostBorrowedItem, error)
	RecentActivity(ctx context.Context, limit int) ([]models.RecentActivityItem, error)
	TopOverdueLoans(ctx context.Context, limit int) ([]models.TopOverdueItem, error)
}
