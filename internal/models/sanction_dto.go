package models

// CreateSanctionRequest is the JSON body for POST /api/sanctions.
type CreateSanctionRequest struct {
	StudentID string `json:"student_id" binding:"required,uuid"`
	Reason    string `json:"reason" binding:"required,min=1,max=2000"`
}

// LiftSanctionRequest is optional JSON for PATCH /api/sanctions/:id (lift).
type LiftSanctionRequest struct {
	Status string `json:"status" binding:"omitempty,oneof=lifted"`
}

// SanctionListItem is one row in GET /api/sanctions (active sanctions only).
type SanctionListItem struct {
	SanctionID    string  `json:"sanction_id"`
	StudentID     string  `json:"student_id"`
	StudentIDCode string  `json:"student_id_code"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Email         *string `json:"email,omitempty"`
	Reason        string  `json:"reason"`
	AppliedAt     string  `json:"applied_at"`
	AppliedByID   *string `json:"applied_by_id,omitempty"`
}

// SanctionResponse is returned after create or lift (full snapshot for UI).
type SanctionResponse struct {
	SanctionID    string  `json:"sanction_id"`
	StudentID     string  `json:"student_id"`
	StudentIDCode string  `json:"student_id_code"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Email         *string `json:"email,omitempty"`
	Reason        string  `json:"reason"`
	Status        string  `json:"status"`
	AppliedAt     string  `json:"applied_at"`
	AppliedByID   *string `json:"applied_by_id,omitempty"`
	LiftedAt      *string `json:"lifted_at,omitempty"`
	LiftedByID    *string `json:"lifted_by_id,omitempty"`
}

// ListSanctionsResponse is GET /api/sanctions.
type ListSanctionsResponse struct {
	Data       []SanctionListItem `json:"data"`
	Total      int64              `json:"total"`
	Pagination PaginationMeta     `json:"pagination"`
}
