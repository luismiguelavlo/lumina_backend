package ranking

import (
	"context"
	"strings"
	"testing"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
)

type fakeRankingRepo struct {
	top3  []repositories.LeaderboardRow
	slice []repositories.LeaderboardRow
	byID  *repositories.LeaderboardRow
	err   error
}

func (f *fakeRankingRepo) GetTop3(ctx context.Context) ([]repositories.LeaderboardRow, error) {
	_ = ctx
	if f.err != nil {
		return nil, f.err
	}
	return f.top3, nil
}

func (f *fakeRankingRepo) GetLeaderboardSlice(ctx context.Context, offset, limit int) ([]repositories.LeaderboardRow, error) {
	_ = ctx
	_ = offset
	_ = limit
	if f.err != nil {
		return nil, f.err
	}
	return f.slice, nil
}

func (f *fakeRankingRepo) GetRankByStudentID(ctx context.Context, studentID string) (*repositories.LeaderboardRow, error) {
	_ = ctx
	_ = studentID
	if f.err != nil {
		return nil, f.err
	}
	return f.byID, nil
}

type noopStudentLookup struct{}

func (noopStudentLookup) GetActiveByID(ctx context.Context, id string) (*models.Student, error) {
	_ = ctx
	_ = id
	return nil, nil
}

func (noopStudentLookup) GetActiveStudentIDByEmail(ctx context.Context, email string) (string, error) {
	_ = ctx
	_ = email
	return "", nil
}

func (noopStudentLookup) GetActiveStudentIDByStudentIDCode(ctx context.Context, studentIDCode string) (string, error) {
	_ = ctx
	_ = studentIDCode
	return "", nil
}

type fakeStudentLookup struct {
	byEmail map[string]string
	byCode  map[string]string
	// activeIDs lists student UUIDs that count as active for GetActiveByID (minimal stub rows).
	activeIDs map[string]bool
}

func (f *fakeStudentLookup) GetActiveByID(ctx context.Context, id string) (*models.Student, error) {
	_ = ctx
	if f.activeIDs == nil || !f.activeIDs[id] {
		return nil, nil
	}
	now := time.Now().UTC()
	return &models.Student{
		ID: id, StudentIDCode: "X", FirstName: "A", LastName: "B",
		IsActive: true, MemberSince: now, CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (f *fakeStudentLookup) GetActiveStudentIDByEmail(ctx context.Context, email string) (string, error) {
	_ = ctx
	if f.byEmail == nil {
		return "", nil
	}
	return f.byEmail[strings.ToLower(strings.TrimSpace(email))], nil
}

func (f *fakeStudentLookup) GetActiveStudentIDByStudentIDCode(ctx context.Context, studentIDCode string) (string, error) {
	_ = ctx
	if f.byCode == nil {
		return "", nil
	}
	return f.byCode[strings.TrimSpace(studentIDCode)], nil
}

func TestService_Top3MapsNames(t *testing.T) {
	f := &fakeRankingRepo{
		top3: []repositories.LeaderboardRow{
			{StudentID: "a", FirstName: "Ana", LastName: "López", TotalPoints: 100, RankPosition: 1},
			{StudentID: "b", FirstName: "Bob", LastName: "Smith", TotalPoints: 90, RankPosition: 2},
		},
	}
	s := NewService(f, noopStudentLookup{})
	out, err := s.Top3(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 || out[0].Name != "Ana López" || out[0].Points != 100 || out[0].Rank != 1 {
		t.Fatalf("%+v", out)
	}
}

func TestService_LeaderboardWithoutStudent(t *testing.T) {
	f := &fakeRankingRepo{
		slice: []repositories.LeaderboardRow{
			{StudentID: "x", FirstName: "X", LastName: "Y", TotalPoints: 10, BooksRead: 2, CurrentStreakDays: 1, RankPosition: 4},
		},
	}
	s := NewService(f, noopStudentLookup{})
	data, cur, err := s.Leaderboard(context.Background(), 3, 3, "", "", "")
	if err != nil || cur != nil || len(data) != 1 || data[0].Rank != 4 {
		t.Fatalf("data=%+v cur=%v err=%v", data, cur, err)
	}
}

func TestService_LeaderboardWithStudentInBoard(t *testing.T) {
	sid := uuid.NewString()
	row := &repositories.LeaderboardRow{
		StudentID: sid, FirstName: "Z", LastName: "W", TotalPoints: 200, RankPosition: 7, BooksRead: 5, CurrentStreakDays: 3,
	}
	f := &fakeRankingRepo{
		slice: []repositories.LeaderboardRow{},
		byID:  row,
	}
	s := NewService(f, noopStudentLookup{})
	data, cur, err := s.Leaderboard(context.Background(), 3, 3, sid, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 0 || cur == nil || cur.Rank != 7 || cur.Points != 200 || cur.Name != "Z W" {
		t.Fatalf("data=%+v cur=%+v", data, cur)
	}
}

func TestService_LeaderboardWithStudentNotInBoard(t *testing.T) {
	sid := uuid.NewString()
	f := &fakeRankingRepo{
		slice: []repositories.LeaderboardRow{},
		byID:  nil,
	}
	s := NewService(f, noopStudentLookup{})
	_, cur, err := s.Leaderboard(context.Background(), 3, 3, sid, "", "")
	if err != nil || cur != nil {
		t.Fatalf("cur=%v err=%v", cur, err)
	}
}

func TestService_LeaderboardByEmailResolves(t *testing.T) {
	sid := uuid.NewString()
	row := &repositories.LeaderboardRow{
		StudentID: sid, FirstName: "E", LastName: "M", TotalPoints: 42, RankPosition: 10,
	}
	f := &fakeRankingRepo{slice: []repositories.LeaderboardRow{}, byID: row}
	st := &fakeStudentLookup{byEmail: map[string]string{"e@test.com": sid}}
	s := NewService(f, st)
	_, cur, err := s.Leaderboard(context.Background(), 3, 3, "", "e@test.com", "")
	if err != nil || cur == nil || cur.Rank != 10 {
		t.Fatalf("cur=%+v err=%v", cur, err)
	}
}

func TestService_LeaderboardByStudentIDCode(t *testing.T) {
	sid := uuid.NewString()
	row := &repositories.LeaderboardRow{StudentID: sid, FirstName: "A", LastName: "B", TotalPoints: 1, RankPosition: 99}
	f := &fakeRankingRepo{slice: nil, byID: row}
	st := &fakeStudentLookup{byCode: map[string]string{"LUM-001": sid}}
	s := NewService(f, st)
	_, cur, err := s.Leaderboard(context.Background(), 0, 3, "", "", "LUM-001")
	if err != nil || cur == nil || cur.Rank != 99 {
		t.Fatalf("cur=%+v err=%v", cur, err)
	}
}

func TestService_StudentRankOK(t *testing.T) {
	sid := uuid.NewString()
	row := &repositories.LeaderboardRow{
		StudentID: sid, FirstName: "Ana", LastName: "Ruiz", TotalPoints: 150,
		RankPosition: 3, BooksRead: 10, CurrentStreakDays: 2,
	}
	f := &fakeRankingRepo{byID: row}
	st := &fakeStudentLookup{activeIDs: map[string]bool{sid: true}}
	s := NewService(f, st)
	out, err := s.StudentRank(context.Background(), sid)
	if err != nil {
		t.Fatal(err)
	}
	if out.Rank != 3 || out.Points != 150 || out.Name != "Ana Ruiz" || out.StudentID != sid {
		t.Fatalf("%+v", out)
	}
}

func TestService_StudentRankNotFound(t *testing.T) {
	s := NewService(&fakeRankingRepo{}, &fakeStudentLookup{})
	_, err := s.StudentRank(context.Background(), uuid.NewString())
	if err != ErrRankingStudentNotFound {
		t.Fatalf("got %v", err)
	}
}

func TestService_StudentRankNotInLeaderboard(t *testing.T) {
	sid := uuid.NewString()
	f := &fakeRankingRepo{byID: nil}
	st := &fakeStudentLookup{activeIDs: map[string]bool{sid: true}}
	s := NewService(f, st)
	_, err := s.StudentRank(context.Background(), sid)
	if err != ErrNotInLeaderboard {
		t.Fatalf("got %v", err)
	}
}
