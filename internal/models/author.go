package models

import "time"

// Author is a book author.
type Author struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador del autor"`
	Name      string    `json:"name" gorm:"size:200;not null;comment:Nombre completo o seudónimo del autor"`
	Bio       *string   `json:"bio,omitempty" gorm:"type:text;comment:Biografía o nota editorial opcional"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;comment:Fecha de alta del autor en catálogo"`
}
