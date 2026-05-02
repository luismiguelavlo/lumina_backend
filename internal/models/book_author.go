package models

// BookAuthor is the books ↔ authors junction (composite PK).
type BookAuthor struct {
	BookID   string `json:"book_id" gorm:"type:uuid;primaryKey;comment:Libro en la relación autor"`
	AuthorID string `json:"author_id" gorm:"type:uuid;primaryKey;comment:Autor asociado al libro"`
}
