package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

const leaderboardOrder = "rank_position ASC, total_points DESC, student_id ASC"

type rankingRepositoryGorm struct {
	db *gorm.DB
}

// NewRankingRepositoryGorm builds a RankingRepository backed by PostgreSQL.
func NewRankingRepositoryGorm(db *gorm.DB) RankingRepository {
	return &rankingRepositoryGorm{db: db}
}

func (r *rankingRepositoryGorm) GetTop3(ctx context.Context) ([]LeaderboardRow, error) {
	var rows []LeaderboardRow
	err := r.db.WithContext(ctx).Table("leaderboard").
		Order(leaderboardOrder).
		Limit(3).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *rankingRepositoryGorm) GetLeaderboardSlice(ctx context.Context, offset, limit int) ([]LeaderboardRow, error) {
	var rows []LeaderboardRow
	err := r.db.WithContext(ctx).Table("leaderboard").
		Order(leaderboardOrder).
		Offset(offset).
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *rankingRepositoryGorm) GetRankByStudentID(ctx context.Context, studentID string) (*LeaderboardRow, error) {
	var row LeaderboardRow
	err := r.db.WithContext(ctx).Table("leaderboard").
		Where("student_id = ?", studentID).
		Take(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}
