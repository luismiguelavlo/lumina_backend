package models

import (
	"encoding/json"
	"time"
)

// ActivityLog is an audit / activity feed entry.
type ActivityLog struct {
	ID          string          `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador de la entrada de actividad"`
	EventType   string          `json:"event_type" gorm:"size:50;not null;index;comment:Código de tipo ej book_returned o student_registered"`
	Title       string          `json:"title" gorm:"size:200;not null;comment:Título breve para listados del dashboard"`
	Description *string         `json:"description,omitempty" gorm:"type:text;comment:Detalle extendido opcional"`
	StudentID   *string         `json:"student_id,omitempty" gorm:"type:uuid;index;comment:Estudiante relacionado si aplica"`
	ActorID     *string         `json:"actor_id,omitempty" gorm:"type:uuid;index;comment:Usuario staff que originó el evento"`
	Metadata    json.RawMessage `json:"metadata,omitempty" gorm:"type:jsonb;serializer:json;comment:Payload JSON arbitrario según event_type"`
	CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime;index;comment:Orden típico descendente para feed"`
}

// TableName keeps GORM using the singular table name from the SQL schema.
func (ActivityLog) TableName() string {
	return "activity_log"
}
