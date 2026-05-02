package models

import "time"

// Department is an academic or organizational unit.
type Department struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador del departamento o facultad"`
	Name      string    `json:"name" gorm:"size:150;not null;uniqueIndex;comment:Nombre legible del departamento único"`
	Code      string    `json:"code" gorm:"size:20;not null;uniqueIndex;comment:Código corto único para referencias internas"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment:Fecha de creación del registro"`
}
