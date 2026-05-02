package models

import "time"

// Badge is an earnable achievement definition.
type Badge struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador de la insignia"`
	Slug        string    `json:"slug" gorm:"size:50;not null;uniqueIndex;comment:Clave estable para lógica ej first-loan"`
	Name        string    `json:"name" gorm:"size:100;not null;comment:Nombre visible de la insignia"`
	IconURL     *string   `json:"icon_url,omitempty" gorm:"type:text;comment:URL del icono en interfaz"`
	Description *string   `json:"description,omitempty" gorm:"type:text;comment:Texto corto para el usuario"`
	Criteria    *string   `json:"criteria,omitempty" gorm:"type:text;comment:Descripción de cómo se obtiene"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;comment:Alta de la definición de insignia"`
}
