package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"library_back/internal/models"
	"library_back/internal/services/ranking"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	msgInvalidStudentIDRanking = "Parámetro student_id no válido"
	msgRankingMultipleLookup   = "Indique solo uno de: student_id, email o student_id_code"
	msgRankingInvalidEmail     = "Email no válido"
	msgRankingStudentCodeLong  = "student_id_code demasiado largo"
	msgRankingStudentNotFound  = "Estudiante no encontrado"
	msgRankingNotInLeaderboard = "El estudiante no figura en el ranking"
)

type rankingServiceAPI interface {
	Top3(ctx context.Context) ([]models.Top3Item, error)
	Leaderboard(ctx context.Context, offset, limit int, studentID, email, studentIDCode string) ([]models.LeaderboardItem, *models.CurrentUserRank, error)
	StudentRank(ctx context.Context, studentID string) (*models.CurrentUserRank, error)
}

// RankingHandler exposes GET /api/ranking/* (público, sin JWT).
type RankingHandler struct {
	svc rankingServiceAPI
}

// NewRankingHandler constructs RankingHandler.
func NewRankingHandler(svc *ranking.Service) *RankingHandler {
	return &RankingHandler{svc: svc}
}

// Top3 GET /api/ranking/top3
func (h *RankingHandler) Top3(c *gin.Context) {
	items, err := h.svc.Top3(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	c.JSON(http.StatusOK, models.Top3Response{Data: items})
}

// Leaderboard GET /api/ranking/leaderboard
// Query: offset, limit (paginación del listado); opcionalmente uno solo de student_id (UUID), email, student_id_code para current_user.
func (h *RankingHandler) Leaderboard(c *gin.Context) {
	offset, limit := normalizeLeaderboardOffsetLimit(c.Query("offset"), c.Query("limit"))

	rawID := strings.TrimSpace(c.Query("student_id"))
	rawEmail := strings.TrimSpace(c.Query("email"))
	rawCode := strings.TrimSpace(c.Query("student_id_code"))

	nSet := 0
	if rawID != "" {
		nSet++
	}
	if rawEmail != "" {
		nSet++
	}
	if rawCode != "" {
		nSet++
	}
	if nSet > 1 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: msgRankingMultipleLookup,
			Errors:  map[string]string{"query": msgRankingMultipleLookup},
		})
		return
	}

	if rawID != "" {
		if _, err := uuid.Parse(rawID); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Message: msgInvalidStudentIDRanking,
				Errors:  map[string]string{"student_id": msgInvalidStudentIDRanking},
			})
			return
		}
	}

	if rawEmail != "" {
		if len(rawEmail) > 255 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Message: msgRankingInvalidEmail,
				Errors:  map[string]string{"email": msgRankingInvalidEmail},
			})
			return
		}
		if _, err := mail.ParseAddress(rawEmail); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Message: msgRankingInvalidEmail,
				Errors:  map[string]string{"email": msgRankingInvalidEmail},
			})
			return
		}
	}

	if rawCode != "" && len(rawCode) > 50 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: msgRankingStudentCodeLong,
			Errors:  map[string]string{"student_id_code": msgRankingStudentCodeLong},
		})
		return
	}

	data, cur, err := h.svc.Leaderboard(c.Request.Context(), offset, limit, rawID, rawEmail, rawCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		return
	}
	c.JSON(http.StatusOK, models.LeaderboardResponse{Data: data, CurrentUser: cur})
}

// StudentRank GET /api/ranking/students/:studentId — global rank for one active student.
func (h *RankingHandler) StudentRank(c *gin.Context) {
	rawID := strings.TrimSpace(c.Param("studentId"))
	if rawID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: msgInvalidStudentIDRanking,
			Errors:  map[string]string{"student_id": msgInvalidStudentIDRanking},
		})
		return
	}
	if _, err := uuid.Parse(rawID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: msgInvalidStudentIDRanking,
			Errors:  map[string]string{"student_id": msgInvalidStudentIDRanking},
		})
		return
	}
	out, err := h.svc.StudentRank(c.Request.Context(), rawID)
	if err != nil {
		switch {
		case errors.Is(err, ranking.ErrRankingStudentNotFound):
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgRankingStudentNotFound})
		case errors.Is(err, ranking.ErrNotInLeaderboard):
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: msgRankingNotInLeaderboard})
		default:
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: msgInternal})
		}
		return
	}
	c.JSON(http.StatusOK, out)
}

func normalizeLeaderboardOffsetLimit(offsetStr, limitStr string) (offset, limit int) {
	offset = 3
	limit = 3
	if offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil {
			if v < 0 {
				offset = 0
			} else {
				offset = v
			}
		}
	}
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			switch {
			case v < 1:
				limit = 3
			case v > 50:
				limit = 50
			default:
				limit = v
			}
		}
	}
	return offset, limit
}

var _ rankingServiceAPI = (*ranking.Service)(nil)
