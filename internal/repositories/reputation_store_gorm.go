package repositories

import (
	"context"
	"errors"

	"library_back/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type reputationStoreGorm struct {
	db *gorm.DB
}

// NewReputationStoreGorm builds a ReputationStore backed by GORM.
func NewReputationStoreGorm(db *gorm.DB) ReputationStore {
	return &reputationStoreGorm{db: db}
}

func (s *reputationStoreGorm) ExistsReturnEventForLoan(ctx context.Context, loanID string) (bool, error) {
	var n int64
	err := s.db.WithContext(ctx).Model(&models.ReputationEvent{}).
		Where("reference_id = ? AND event_type IN ?", loanID,
			[]models.ReputationEventType{
				models.ReputationEventBookReturnedOnTime,
				models.ReputationEventBookReturnedLate,
			}).
		Count(&n).Error
	return n > 0, err
}

func (s *reputationStoreGorm) CreateReputationEvent(ctx context.Context, ev *models.ReputationEvent) error {
	return s.db.WithContext(ctx).Create(ev).Error
}

func (s *reputationStoreGorm) IncrementStudentActiveLoans(ctx context.Context, studentID string, delta int) error {
	if delta == 0 {
		return nil
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var stats models.StudentStats
		err := tx.Where("student_id = ?", studentID).First(&stats).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			n := 0
			if delta > 0 {
				n = delta
			}
			return tx.Create(&models.StudentStats{
				StudentID:   studentID,
				ActiveLoans: n,
			}).Error
		}
		if err != nil {
			return err
		}
		newVal := stats.ActiveLoans + delta
		if newVal < 0 {
			newVal = 0
		}
		return tx.Model(&stats).Update("active_loans", newVal).Error
	})
}

func (s *reputationStoreGorm) ApplyStudentReturn(ctx context.Context, studentID string, points int, onTime bool) (newStreakDays int, err error) {
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var stats models.StudentStats
		e := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("student_id = ?", studentID).First(&stats).Error
		if errors.Is(e, gorm.ErrRecordNotFound) {
			streak := streakAfterReturn(0, onTime)
			newStreakDays = streak
			stats = models.StudentStats{
				StudentID:         studentID,
				TotalRead:         1,
				ActiveLoans:       0,
				TotalPoints:       points,
				CurrentStreakDays: streak,
				LongestStreakDays: streak,
			}
			if !onTime {
				stats.OverdueCount = 1
			}
			return tx.Create(&stats).Error
		}
		if e != nil {
			return e
		}

		newActive := stats.ActiveLoans - 1
		if newActive < 0 {
			newActive = 0
		}
		streak := streakAfterReturn(stats.CurrentStreakDays, onTime)
		newStreakDays = streak
		longest := stats.LongestStreakDays
		if streak > longest {
			longest = streak
		}
		overdueDelta := 0
		if !onTime {
			overdueDelta = 1
		}

		return tx.Model(&stats).Updates(map[string]interface{}{
			"total_read":          gorm.Expr("total_read + ?", 1),
			"active_loans":        newActive,
			"total_points":        gorm.Expr("total_points + ?", points),
			"current_streak_days": streak,
			"longest_streak_days": longest,
			"overdue_count":       gorm.Expr("overdue_count + ?", overdueDelta),
		}).Error
	})
	return newStreakDays, err
}

func (s *reputationStoreGorm) ExistsReputationEventByTypeAndReference(ctx context.Context, typ models.ReputationEventType, referenceID string) (bool, error) {
	var n int64
	err := s.db.WithContext(ctx).Model(&models.ReputationEvent{}).
		Where("event_type = ? AND reference_id = ?", typ, referenceID).
		Count(&n).Error
	return n > 0, err
}

func (s *reputationStoreGorm) CreateReputationEventWithPointsIfNew(ctx context.Context, studentID string, points int, typ models.ReputationEventType, referenceID string, description *string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var n int64
		if err := tx.Model(&models.ReputationEvent{}).
			Where("event_type = ? AND reference_id = ?", typ, referenceID).
			Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			return nil
		}
		ref := referenceID
		ev := &models.ReputationEvent{
			ID:          uuid.NewString(),
			StudentID:   studentID,
			EventType:   typ,
			Points:      points,
			Description: description,
			ReferenceID: &ref,
		}
		if err := tx.Create(ev).Error; err != nil {
			return err
		}

		var stats models.StudentStats
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("student_id = ?", studentID).First(&stats).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tp := points
			if tp < 0 {
				tp = 0
			}
			return tx.Create(&models.StudentStats{
				StudentID:   studentID,
				TotalPoints: tp,
			}).Error
		}
		if err != nil {
			return err
		}
		return tx.Model(&stats).Update("total_points", gorm.Expr("GREATEST(0, total_points + ?)", points)).Error
	})
}

func streakAfterReturn(prev int, onTime bool) int {
	if onTime {
		if prev < 0 {
			prev = 0
		}
		return prev + 1
	}
	return 0
}
