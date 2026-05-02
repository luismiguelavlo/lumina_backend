package ranking

import "errors"

var (
	// ErrRankingStudentNotFound is returned when the student id does not exist or the student is inactive.
	ErrRankingStudentNotFound = errors.New("ranking student not found")
	// ErrNotInLeaderboard is returned when the student is active but has no row in the leaderboard view (no stats).
	ErrNotInLeaderboard = errors.New("student not in leaderboard")
)
