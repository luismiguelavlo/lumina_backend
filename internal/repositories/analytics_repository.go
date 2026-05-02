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
	MostBorrowedBooks(ctx context.Context, limit int) ([]models.MostBorrowedItem, error)
	RecentActivity(ctx context.Context, limit int) ([]models.RecentActivityItem, error)
}
