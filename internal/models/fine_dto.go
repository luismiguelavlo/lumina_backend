package models

// CreateFineRequest is the JSON body for POST /api/fines.
type CreateFineRequest struct {
	LoanID    string   `json:"loan_id" binding:"required,uuid"`
	StudentID string   `json:"student_id" binding:"required,uuid"`
	Amount    float64  `json:"amount" binding:"required,gte=0"`
	Reason    *string  `json:"reason,omitempty" binding:"omitempty,max=500"`
}

// FineListItem is one row in GET /api/fines.
type FineListItem struct {
	ID          string  `json:"id"`
	LoanID      string  `json:"loan_id"`
	StudentID   string  `json:"student_id"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	Reason      *string `json:"reason,omitempty"`
	CreatedAt   string  `json:"created_at"`
	PaidAt      *string `json:"paid_at,omitempty"`
	StudentName *string `json:"student_name,omitempty"`
}

// FineResponse is detail, create, or status-update payload.
type FineResponse struct {
	ID        string   `json:"id"`
	LoanID    string   `json:"loan_id"`
	StudentID string   `json:"student_id"`
	Amount    float64  `json:"amount"`
	Status    string   `json:"status"`
	Reason    *string  `json:"reason,omitempty"`
	CreatedAt string   `json:"created_at"`
	PaidAt    *string  `json:"paid_at,omitempty"`
}

// ListFinesResponse is GET /api/fines.
type ListFinesResponse struct {
	Data       []FineListItem `json:"data"`
	Total      int64          `json:"total"`
	Pagination PaginationMeta `json:"pagination"`
}
