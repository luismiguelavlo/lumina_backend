package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Fine is a monetary penalty linked to a loan.
type Fine struct {
	ID        string          `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4();comment:Identificador de la multa"`
	LoanID    string          `json:"loan_id" gorm:"type:uuid;not null;index;comment:Préstamo que originó o vincula la multa"`
	StudentID string          `json:"student_id" gorm:"type:uuid;not null;index;comment:Estudiante responsable del pago"`
	Amount    decimal.Decimal `json:"amount" gorm:"type:numeric(10,2);not null;comment:Importe en moneda local mayor o igual que cero"`
	Status    FineStatus      `json:"status" gorm:"type:fine_status;not null;default:pending;comment:Estado pendiente pagada o condonada"`
	Reason    *string         `json:"reason,omitempty" gorm:"type:text;comment:Motivo o detalle mostrado al usuario"`
	CreatedAt time.Time       `json:"created_at" gorm:"autoCreateTime;comment:Cuándo se generó la multa"`
	PaidAt    *time.Time      `json:"paid_at,omitempty" gorm:"comment:Marca temporal cuando el estado pasó a pagada"`
}
