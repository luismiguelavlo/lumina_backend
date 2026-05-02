package models

import "time"

// ReputationEvent records a points change for gamification / reputation.
type ReputationEvent struct {
	ID          string              `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador del evento de puntos"`
	StudentID   string              `json:"student_id" gorm:"type:uuid;not null;index;comment:Estudiante afectado"`
	EventType   ReputationEventType `json:"event_type" gorm:"type:reputation_event_type;not null;comment:Tipo de hecho que originó el movimiento"`
	Points      int                 `json:"points" gorm:"not null;comment:Delta de puntos positivo o negativo"`
	Description *string             `json:"description,omitempty" gorm:"type:text;comment:Texto libre opcional para auditoría"`
	ReferenceID *string             `json:"reference_id,omitempty" gorm:"type:uuid;comment:ID externo para idempotencia ej préstamo origen"`
	CreatedAt   time.Time           `json:"created_at" gorm:"autoCreateTime;comment:Momento del evento"`
}
