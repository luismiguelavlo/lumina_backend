package models

import "time"

// DashboardLimits controls list sizes for GET /api/analytics/dashboard.
type DashboardLimits struct {
	MostBorrowedLimit   int
	RecentActivityLimit int
}

// DashboardResponse is the JSON body for the analytics dashboard.
type DashboardResponse struct {
	TotalBooks        int64                `json:"total_books"`
	ActiveStudents    int64                `json:"active_students"`
	OverdueFines      float64              `json:"overdue_fines"`
	MostBorrowedBooks []MostBorrowedItem   `json:"most_borrowed_books"`
	RecentActivity    []RecentActivityItem `json:"recent_activity"`
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
