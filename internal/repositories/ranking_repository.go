package repositories

import "context"

// LeaderboardRow is one row from the leaderboard view (read model).
type LeaderboardRow struct {
	StudentID         string  `gorm:"column:student_id"`
	FirstName         string  `gorm:"column:first_name"`
	LastName          string  `gorm:"column:last_name"`
	AvatarURL         *string `gorm:"column:avatar_url"`
	TotalPoints       int     `gorm:"column:total_points"`
	BooksRead         int     `gorm:"column:books_read"`
	CurrentStreakDays int     `gorm:"column:current_streak_days"`
	RankPosition      int     `gorm:"column:rank_position"`
}

// RankingRepository reads the leaderboard view for reputation ranking APIs.
type RankingRepository interface {
	GetTop3(ctx context.Context) ([]LeaderboardRow, error)
	GetLeaderboardSlice(ctx context.Context, offset, limit int) ([]LeaderboardRow, error)
	// GetRankByStudentID returns nil, nil when the student is not in the leaderboard (no stats or inactive).
	GetRankByStudentID(ctx context.Context, studentID string) (*LeaderboardRow, error)
}
