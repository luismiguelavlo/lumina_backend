package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"library_back/internal/models"
	"library_back/internal/services/ranking"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockRankingSvc struct {
	top3            []models.Top3Item
	top3Err         error
	lbData          []models.LeaderboardItem
	lbCur           *models.CurrentUserRank
	lbErr           error
	studentRank     *models.CurrentUserRank
	studentRankErr  error
	lastStudentRank string
	lastOffset      int
	lastLimit       int
	lastStudentID   string
	lastEmail       string
	lastStudentCode string
}

func (m *mockRankingSvc) Top3(ctx context.Context) ([]models.Top3Item, error) {
	_ = ctx
	if m.top3Err != nil {
		return nil, m.top3Err
	}
	return m.top3, nil
}

func (m *mockRankingSvc) Leaderboard(ctx context.Context, offset, limit int, studentID, email, studentIDCode string) ([]models.LeaderboardItem, *models.CurrentUserRank, error) {
	_ = ctx
	m.lastOffset = offset
	m.lastLimit = limit
	m.lastStudentID = studentID
	m.lastEmail = email
	m.lastStudentCode = studentIDCode
	if m.lbErr != nil {
		return nil, nil, m.lbErr
	}
	return m.lbData, m.lbCur, nil
}

func (m *mockRankingSvc) StudentRank(ctx context.Context, studentID string) (*models.CurrentUserRank, error) {
	_ = ctx
	m.lastStudentRank = studentID
	if m.studentRankErr != nil {
		return nil, m.studentRankErr
	}
	return m.studentRank, nil
}

func TestRankingHandler_Top3_200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockRankingSvc{
		top3: []models.Top3Item{{Rank: 1, Name: "A B", Points: 10, StudentID: "s1"}},
	}
	h := &RankingHandler{svc: ms}
	r := gin.New()
	r.GET("/ranking/top3", h.Top3)
	req := httptest.NewRequest(http.MethodGet, "/ranking/top3", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d %s", w.Code, w.Body.String())
	}
	var body models.Top3Response
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) != 1 || body.Data[0].Rank != 1 {
		t.Fatalf("%+v", body)
	}
}

func TestRankingHandler_LeaderboardDefaults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockRankingSvc{}
	h := &RankingHandler{svc: ms}
	r := gin.New()
	r.GET("/ranking/leaderboard", h.Leaderboard)
	req := httptest.NewRequest(http.MethodGet, "/ranking/leaderboard", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || ms.lastOffset != 3 || ms.lastLimit != 3 {
		t.Fatalf("code=%d off=%d lim=%d", w.Code, ms.lastOffset, ms.lastLimit)
	}
	if ms.lastStudentID != "" || ms.lastEmail != "" || ms.lastStudentCode != "" {
		t.Fatalf("unexpected lookup %+v %+v %+v", ms.lastStudentID, ms.lastEmail, ms.lastStudentCode)
	}
}

func TestRankingHandler_Leaderboard400BadStudentID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockRankingSvc{}
	h := &RankingHandler{svc: ms}
	r := gin.New()
	r.GET("/ranking/leaderboard", h.Leaderboard)
	req := httptest.NewRequest(http.MethodGet, "/ranking/leaderboard?student_id=not-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestRankingHandler_Leaderboard400MultipleLookup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockRankingSvc{}
	h := &RankingHandler{svc: ms}
	r := gin.New()
	r.GET("/ranking/leaderboard", h.Leaderboard)
	req := httptest.NewRequest(http.MethodGet, "/ranking/leaderboard?email=a@b.com&student_id_code=X", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d %s", w.Code, w.Body.String())
	}
}

func TestRankingHandler_Leaderboard400BadEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockRankingSvc{}
	h := &RankingHandler{svc: ms}
	r := gin.New()
	r.GET("/ranking/leaderboard", h.Leaderboard)
	req := httptest.NewRequest(http.MethodGet, "/ranking/leaderboard?email=not-an-email", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestRankingHandler_LeaderboardWithEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockRankingSvc{
		lbData: []models.LeaderboardItem{{Rank: 4, StudentID: "x", Name: "P Q", Points: 5}},
		lbCur:  &models.CurrentUserRank{Rank: 12, StudentID: "y", Name: "Me", Points: 2},
	}
	h := &RankingHandler{svc: ms}
	r := gin.New()
	r.GET("/ranking/leaderboard", h.Leaderboard)
	req := httptest.NewRequest(http.MethodGet, "/ranking/leaderboard?offset=3&limit=5&email=user@example.com", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || ms.lastEmail != "user@example.com" || ms.lastStudentID != "" || ms.lastStudentCode != "" {
		t.Fatalf("code=%d email=%q id=%q code=%q", w.Code, ms.lastEmail, ms.lastStudentID, ms.lastStudentCode)
	}
}

func TestRankingHandler_LeaderboardWithStudentID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sid := uuid.NewString()
	ms := &mockRankingSvc{lbCur: &models.CurrentUserRank{Rank: 1, StudentID: sid}}
	h := &RankingHandler{svc: ms}
	r := gin.New()
	r.GET("/ranking/leaderboard", h.Leaderboard)
	req := httptest.NewRequest(http.MethodGet, "/ranking/leaderboard?student_id="+sid, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || ms.lastStudentID != sid {
		t.Fatalf("code=%d sid=%q", w.Code, ms.lastStudentID)
	}
}

func TestRankingHandler_StudentRank200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sid := uuid.NewString()
	ms := &mockRankingSvc{
		studentRank: &models.CurrentUserRank{Rank: 5, StudentID: sid, Name: "A B", Points: 99},
	}
	h := &RankingHandler{svc: ms}
	r := gin.New()
	r.GET("/ranking/students/:studentId", h.StudentRank)
	req := httptest.NewRequest(http.MethodGet, "/ranking/students/"+sid, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK || ms.lastStudentRank != sid {
		t.Fatalf("code=%d last=%q body=%s", w.Code, ms.lastStudentRank, w.Body.String())
	}
	var body models.StudentRankResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Rank != 5 || body.Points != 99 {
		t.Fatalf("%+v", body)
	}
}

func TestRankingHandler_StudentRank400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &RankingHandler{svc: &mockRankingSvc{}}
	r := gin.New()
	r.GET("/ranking/students/:studentId", h.StudentRank)
	req := httptest.NewRequest(http.MethodGet, "/ranking/students/not-uuid", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 got %d", w.Code)
	}
}

func TestRankingHandler_StudentRank404NotInBoard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := &mockRankingSvc{studentRankErr: ranking.ErrNotInLeaderboard}
	h := &RankingHandler{svc: ms}
	r := gin.New()
	r.GET("/ranking/students/:studentId", h.StudentRank)
	sid := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "/ranking/students/"+sid, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404 got %d", w.Code)
	}
}
