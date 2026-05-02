package models

// CreateLoanRequest is the JSON body for POST /api/loans.
type CreateLoanRequest struct {
	StudentID string `json:"student_id" binding:"required,uuid"`
	BookID    string `json:"book_id" binding:"required,uuid"`
	DueDate   string `json:"due_date" binding:"required"` // YYYY-MM-DD
}

// LoanListItem is one row in GET /api/loans.
type LoanListItem struct {
	LoanID        string `json:"loan_id"`
	BookID        string `json:"book_id"`
	BookTitle     string `json:"book_title"`
	Borrower      string `json:"borrower"`
	DueDate       string `json:"due_date"`
	TimeRemaining int    `json:"time_remaining"`
	ISBN          string `json:"isbn"`
	Status        string `json:"status"`
}

// LoanDetailResponse is POST create / PATCH return payload.
type LoanDetailResponse struct {
	LoanID        string `json:"loan_id"`
	BookID        string `json:"book_id"`
	BookTitle     string `json:"book_title"`
	Borrower      string `json:"borrower"`
	BorrowedAt    string `json:"borrowed_at"`
	DueDate       string `json:"due_date"`
	TimeRemaining int    `json:"time_remaining"`
	ISBN          string `json:"isbn"`
	Status        string `json:"status"`
}

// ListLoansResponse is GET /api/loans.
type ListLoansResponse struct {
	Data       []LoanListItem `json:"data"`
	Total      int64          `json:"total"`
	Pagination PaginationMeta `json:"pagination"`
}
