package models

// Top3Item is one row in the Top Readers reputation block.
type Top3Item struct {
	Rank      int     `json:"rank"`
	Name      string  `json:"name"`
	Points    int     `json:"points"`
	StudentID string  `json:"student_id"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// Top3Response is the body for GET /api/ranking/top3.
type Top3Response struct {
	Data []Top3Item `json:"data"`
}

// LeaderboardItem is one slice of the global leaderboard.
type LeaderboardItem struct {
	Rank              int     `json:"rank"`
	StudentID         string  `json:"student_id"`
	Name              string  `json:"name"`
	Points            int     `json:"points"`
	AvatarURL         *string `json:"avatar_url,omitempty"`
	BooksRead         int     `json:"books_read,omitempty"`
	CurrentStreakDays int     `json:"current_streak_days,omitempty"`
}

// CurrentUserRank is the authenticated / selected student's row ("Your rank").
type CurrentUserRank struct {
	Rank              int    `json:"rank"`
	StudentID         string `json:"student_id"`
	Name              string `json:"name"`
	Points            int    `json:"points"`
	BooksRead         int    `json:"books_read,omitempty"`
	CurrentStreakDays int    `json:"current_streak_days,omitempty"`
}

// LeaderboardResponse is the body for GET /api/ranking/leaderboard.
type LeaderboardResponse struct {
	Data        []LeaderboardItem `json:"data"`
	CurrentUser *CurrentUserRank  `json:"current_user"`
}

// StudentRankResponse is the body for GET /api/ranking/students/:studentId (same shape as CurrentUserRank).
type StudentRankResponse = CurrentUserRank
