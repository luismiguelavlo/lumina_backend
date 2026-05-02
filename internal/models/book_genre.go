package models

// BookGenre is the books ↔ genres junction (composite PK).
type BookGenre struct {
	BookID  string `json:"book_id" gorm:"type:uuid;primaryKey;comment:Libro clasificado"`
	GenreID string `json:"genre_id" gorm:"type:uuid;primaryKey;comment:Género asignado al libro"`
}
