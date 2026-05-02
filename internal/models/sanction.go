package models

import "time"

// Sanction restricts a student (e.g. after repeated violations).
type Sanction struct {
	ID        string         `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador de la sanción"`
	StudentID string         `json:"student_id" gorm:"type:uuid;not null;index;comment:Estudiante sancionado"`
	Reason    string         `json:"reason" gorm:"type:text;not null;comment:Motivo textual de la restricción"`
	Status    SanctionStatus `json:"status" gorm:"type:sanction_status;not null;default:active;comment:Activa o levantada"`
	AppliedAt time.Time      `json:"applied_at" gorm:"autoCreateTime;comment:Momento en que entró en vigor"`
	AppliedBy *string        `json:"applied_by,omitempty" gorm:"type:uuid;comment:Usuario staff que aplicó la sanción"`
	LiftedAt  *time.Time     `json:"lifted_at,omitempty" gorm:"comment:Momento en que se levantó si aplica"`
	LiftedBy  *string        `json:"lifted_by,omitempty" gorm:"type:uuid;comment:Usuario staff que levantó la sanción"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;comment:Alta del registro"`
}
