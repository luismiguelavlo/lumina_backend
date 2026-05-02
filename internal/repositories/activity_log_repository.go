package repositories

import (
	"context"

	"library_back/internal/models"
)

// ActivityLogRepository persists dashboard / audit feed rows.
type ActivityLogRepository interface {
	Create(ctx context.Context, entry *models.ActivityLog) error
}
