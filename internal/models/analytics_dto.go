package models

import "time"

// DashboardLimits controls list sizes for GET /api/analytics/dashboard.
type DashboardLimits struct {
	MostBorrowedLimit   int
	RecentActivityLimit int
	TopOverdueLimit     int
}

// DashboardResponse is the JSON body for the analytics dashboard.
type DashboardResponse struct {
	TotalBooks         int64                `json:"total_books"`
	ActiveStudents     int64                `json:"active_students"`
	OverdueFines       float64              `json:"overdue_fines"` // pending fines amount (legacy field name)
	PendingFinesCount  int64                `json:"pending_fines_count"`
	OverdueLoans       int64                `json:"overdue_loans"`
	ActiveLoans        int64                `json:"active_loans"`
	DueSoon            int64                `json:"due_soon"`
	ActiveSanctions    int64                `json:"active_sanctions"`
	ReturnsThisWeek    int64                `json:"returns_this_week"`
	NewStudentsMonth   int64                `json:"new_students_month"`
	AvailableCopies    int64                `json:"available_copies"`
	CheckedOutCopies   int64                `json:"checked_out_copies"`
	MostBorrowedBooks  []MostBorrowedItem   `json:"most_borrowed_books"`
	RecentActivity     []RecentActivityItem `json:"recent_activity"`
	TopOverdue         []TopOverdueItem     `json:"top_overdue"`
}

// MostBorrowedItem is one row in the "most borrowed books" list.
type MostBorrowedItem struct {
	BookID      string   `json:"book_id"`
	Title       string   `json:"title"`
	CatalogCode string   `json:"catalog_code"`
	BorrowCount int64    `json:"borrow_count"`
	AuthorNames []string `json:"author_names,omitempty"`
}

// RecentActivityItem is one row from activity_log for the dashboard feed.
type RecentActivityItem struct {
	ID          string         `json:"id"`
	EventType   string         `json:"event_type"`
	Title       string         `json:"title"`
	Description *string        `json:"description,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	ActorName   *string        `json:"actor_name,omitempty"`
	StudentName *string        `json:"student_name,omitempty"`
}

// TopOverdueItem is one overdue loan for the librarian action queue.
type TopOverdueItem struct {
	LoanID      string    `json:"loan_id"`
	StudentID   string    `json:"student_id"`
	StudentName string    `json:"student_name"`
	StudentCode string    `json:"student_code,omitempty"`
	BookID      string    `json:"book_id"`
	BookTitle   string    `json:"book_title"`
	DueDate     time.Time `json:"due_date"`
	DaysOverdue int       `json:"days_overdue"`
}
