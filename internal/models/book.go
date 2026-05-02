package models

import "time"

// Book is a catalog title (not a physical copy row).
type Book struct {
	ID              string     `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador del título en catálogo"`
	Title           string     `json:"title" gorm:"size:300;not null;comment:Título del libro"`
	ISBN            string     `json:"isbn" gorm:"size:20;not null;uniqueIndex;comment:ISBN único del ejemplar conceptual"`
	CatalogCode     string     `json:"catalog_code" gorm:"size:20;not null;uniqueIndex;comment:Código interno de estantería ej FIC-001"`
	Synopsis        *string    `json:"synopsis,omitempty" gorm:"type:text;comment:Resumen o sinopsis opcional"`
	PublicationYear *int       `json:"publication_year,omitempty" gorm:"comment:Año de publicación"`
	Pages           *int       `json:"pages,omitempty" gorm:"comment:Número de páginas"`
	CoverURL        *string    `json:"cover_url,omitempty" gorm:"type:text;comment:URL de la portada"`
	Location        *string    `json:"location,omitempty" gorm:"size:200;comment:Ubicación física en biblioteca ej estantería A"`
	TotalCopies     int        `json:"total_copies" gorm:"not null;default:1;comment:Ejemplares totales disponibles para préstamo"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" gorm:"comment:Marca de borrado lógico no nulo si el título fue retirado"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime;comment:Fecha de alta del título"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime;comment:Última modificación de metadatos"`
}
