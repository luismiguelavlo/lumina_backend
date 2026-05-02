package models

import "time"

// StudentBadge links a student to an earned badge.
type StudentBadge struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador de la concesión"`
	StudentID string    `json:"student_id" gorm:"type:uuid;not null;index;uniqueIndex:uniq_student_badge;comment:Estudiante que obtuvo la insignia"`
	BadgeID   string    `json:"badge_id" gorm:"type:uuid;not null;index;uniqueIndex:uniq_student_badge;comment:Insignia concedida"`
	EarnedAt  time.Time `json:"earned_at" gorm:"autoCreateTime;comment:Momento en que se registró el logro"`
}
