package ranking

import (
	"context"
	"strings"

	"library_back/internal/models"
	"library_back/internal/repositories"
)

// rankingStudentReader resolves students for ranking APIs (active reads + public identifiers).
type rankingStudentReader interface {
	GetActiveByID(ctx context.Context, id string) (*models.Student, error)
	GetActiveStudentIDByEmail(ctx context.Context, email string) (string, error)
	GetActiveStudentIDByStudentIDCode(ctx context.Context, studentIDCode string) (string, error)
}

// Service exposes reputation ranking read APIs.
type Service struct {
	repo     repositories.RankingRepository
	students rankingStudentReader
}

// NewService constructs RankingService. students may be nil only in tests that never need lookups.
func NewService(repo repositories.RankingRepository, students rankingStudentReader) *Service {
	return &Service{repo: repo, students: students}
}

// Top3 returns up to three leaderboard rows for the Top Readers section.
func (s *Service) Top3(ctx context.Context) ([]models.Top3Item, error) {
	rows, err := s.repo.GetTop3(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]models.Top3Item, 0, len(rows))
	for i := range rows {
		out = append(out, rowToTop3Item(&rows[i]))
	}
	return out, nil
}

// Leaderboard returns a slice of the global leaderboard and optionally the viewer's rank.
// studentID must be a valid UUID when non-empty. email and studentIDCode are alternative public lookups (only one should be used).
func (s *Service) Leaderboard(ctx context.Context, offset, limit int, studentID, email, studentIDCode string) ([]models.LeaderboardItem, *models.CurrentUserRank, error) {
	rows, err := s.repo.GetLeaderboardSlice(ctx, offset, limit)
	if err != nil {
		return nil, nil, err
	}
	data := make([]models.LeaderboardItem, 0, len(rows))
	for i := range rows {
		data = append(data, rowToLeaderboardItem(&rows[i]))
	}
	resolved, err := s.resolveStudentUUID(ctx, studentID, email, studentIDCode)
	if err != nil {
		return nil, nil, err
	}
	var cur *models.CurrentUserRank
	if resolved != "" {
		row, err := s.repo.GetRankByStudentID(ctx, resolved)
		if err != nil {
			return nil, nil, err
		}
		if row != nil {
			cur = rowToCurrentUser(row)
		}
	}
	return data, cur, nil
}

func (s *Service) resolveStudentUUID(ctx context.Context, studentID, email, studentIDCode string) (string, error) {
	if strings.TrimSpace(studentID) != "" {
		return strings.TrimSpace(studentID), nil
	}
	em := strings.TrimSpace(email)
	if em != "" {
		if s.students == nil {
			return "", nil
		}
		return s.students.GetActiveStudentIDByEmail(ctx, em)
	}
	code := strings.TrimSpace(studentIDCode)
	if code != "" {
		if s.students == nil {
			return "", nil
		}
		return s.students.GetActiveStudentIDByStudentIDCode(ctx, code)
	}
	return "", nil
}

// StudentRank returns the global rank row for an active student by internal UUID.
func (s *Service) StudentRank(ctx context.Context, studentID string) (*models.CurrentUserRank, error) {
	if s.students == nil {
		return nil, ErrRankingStudentNotFound
	}
	st, err := s.students.GetActiveByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, ErrRankingStudentNotFound
	}
	row, err := s.repo.GetRankByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNotInLeaderboard
	}
	return rowToCurrentUser(row), nil
}

func rowToTop3Item(r *repositories.LeaderboardRow) models.Top3Item {
	return models.Top3Item{
		Rank:      r.RankPosition,
		Name:      fullName(r.FirstName, r.LastName),
		Points:    r.TotalPoints,
		StudentID: r.StudentID,
		AvatarURL: r.AvatarURL,
	}
}

func rowToLeaderboardItem(r *repositories.LeaderboardRow) models.LeaderboardItem {
	return models.LeaderboardItem{
		Rank:              r.RankPosition,
		StudentID:         r.StudentID,
		Name:              fullName(r.FirstName, r.LastName),
		Points:            r.TotalPoints,
		AvatarURL:         r.AvatarURL,
		BooksRead:         r.BooksRead,
		CurrentStreakDays: r.CurrentStreakDays,
	}
}

func rowToCurrentUser(r *repositories.LeaderboardRow) *models.CurrentUserRank {
	return &models.CurrentUserRank{
		Rank:              r.RankPosition,
		StudentID:         r.StudentID,
		Name:              fullName(r.FirstName, r.LastName),
		Points:            r.TotalPoints,
		BooksRead:         r.BooksRead,
		CurrentStreakDays: r.CurrentStreakDays,
	}
}

func fullName(first, last string) string {
	return strings.TrimSpace(strings.TrimSpace(first) + " " + strings.TrimSpace(last))
}
