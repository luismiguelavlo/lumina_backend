package models

import "time"

// Loan is a book loan to a student.
// DueDate is stored as date in DB; use UTC midnight for the calendar day.
type Loan struct {
	ID         string     `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador del préstamo"`
	BookID     string     `json:"book_id" gorm:"type:uuid;not null;index;comment:Libro prestado"`
	StudentID  string     `json:"student_id" gorm:"type:uuid;not null;index;comment:Estudiante que recibe el libro"`
	IssuedBy   *string    `json:"issued_by,omitempty" gorm:"type:uuid;comment:Usuario staff que registró el préstamo"`
	BorrowedAt time.Time  `json:"borrowed_at" gorm:"autoCreateTime;comment:Momento en que salió el material"`
	DueDate    time.Time  `json:"due_date" gorm:"type:date;not null;comment:Fecha límite de devolución solo día"`
	ReturnedAt *time.Time `json:"returned_at,omitempty" gorm:"comment:Momento de devolución efectiva si ya se devolvió"`
	Status     LoanStatus `json:"status" gorm:"type:loan_status;not null;default:active;comment:Estado activo devuelto o moroso"`
	Notes      *string    `json:"notes,omitempty" gorm:"type:text;comment:Notas internas del bibliotecario"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime;comment:Alta del registro de préstamo"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime;comment:Última actualización del préstamo"`
}
