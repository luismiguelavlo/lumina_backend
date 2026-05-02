package models

import "time"

// Genre is a book category.
type Genre struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador del género literario"`
	Name      string    `json:"name" gorm:"size:100;not null;uniqueIndex;comment:Nombre legible del género"`
	Code      string    `json:"code" gorm:"size:10;not null;uniqueIndex;comment:Código corto único ej FIC o SCI"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment:Fecha de creación del género"`
}
