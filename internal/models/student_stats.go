package models

import "time"

// StudentStats holds aggregated counters for a student (1:1 with students).
type StudentStats struct {
	ID                string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador de la fila de estadísticas"`
	StudentID         string    `json:"student_id" gorm:"type:uuid;not null;uniqueIndex;comment:Estudiante dueño métricas 1 a 1"`
	TotalRead         int       `json:"total_read" gorm:"not null;default:0;comment:Libros leídos o devueltos contabilizados"`
	ActiveLoans       int       `json:"active_loans" gorm:"not null;default:0;comment:Préstamos vigentes en este momento"`
	OverdueCount      int       `json:"overdue_count" gorm:"not null;default:0;comment:Veces o ítems en mora acumulados"`
	CurrentStreakDays int       `json:"current_streak_days" gorm:"not null;default:0;comment:Racha actual de devoluciones puntuales en días"`
	LongestStreakDays int       `json:"longest_streak_days" gorm:"not null;default:0;comment:Mejor racha histórica en días"`
	TotalPoints       int       `json:"total_points" gorm:"not null;default:0;comment:Puntos de reputación acumulados"`
	GlobalRank        *int      `json:"global_rank,omitempty" gorm:"comment:Posición en ranking global calculada externamente"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:Última actualización de contadores"`
}
