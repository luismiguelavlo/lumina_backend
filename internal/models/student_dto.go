package models

// CreateStudentRequest is the JSON body for POST /api/students.
type CreateStudentRequest struct {
	FirstName              string  `json:"first_name" binding:"required,max=100"`
	LastName               string  `json:"last_name" binding:"required,max=100"`
	StudentIDCode          string  `json:"student_id_code" binding:"omitempty,max=50"` // omit → server assigns AUTO-YYYYMMDD-<hex>
	Email                  *string `json:"email" binding:"omitempty,email,max=255"`
	AvatarURL              *string `json:"avatar_url" binding:"omitempty,url"`
	DepartmentID           *string `json:"department_id" binding:"omitempty,uuid"`
	DegreeLevel            *string `json:"degree_level" binding:"omitempty,max=100"`
	Major                  *string `json:"major" binding:"omitempty,max=150"`
	ExpectedGraduationYear *int    `json:"expected_graduation_year"`
}

// UpdateStudentRequest is the JSON body for PATCH /api/students/:id.
type UpdateStudentRequest struct {
	FirstName              *string `json:"first_name" binding:"omitempty,max=100"`
	LastName               *string `json:"last_name" binding:"omitempty,max=100"`
	StudentIDCode          *string `json:"student_id_code" binding:"omitempty,max=50"`
	Email                  *string `json:"email" binding:"omitempty,email,max=255"`
	AvatarURL              *string `json:"avatar_url" binding:"omitempty,url"`
	DepartmentID           *string `json:"department_id" binding:"omitempty,uuid"`
	DegreeLevel            *string `json:"degree_level" binding:"omitempty,max=100"`
	Major                  *string `json:"major" binding:"omitempty,max=150"`
	ExpectedGraduationYear *int    `json:"expected_graduation_year"`
}

// StudentResponse is the public student payload for list/detail/create/update.
type StudentResponse struct {
	ID                     string  `json:"id"`
	StudentIDCode          string  `json:"student_id_code"`
	FirstName              string  `json:"first_name"`
	LastName               string  `json:"last_name"`
	Email                  *string `json:"email,omitempty"`
	AvatarURL              *string `json:"avatar_url,omitempty"`
	DepartmentID           *string `json:"department_id,omitempty"`
	DepartmentName         *string `json:"department_name,omitempty"`
	DepartmentCode         *string `json:"department_code,omitempty"`
	DegreeLevel            *string `json:"degree_level,omitempty"`
	Major                  *string `json:"major,omitempty"`
	ExpectedGraduationYear *int    `json:"expected_graduation_year,omitempty"`
	IsActive               bool    `json:"is_active"`
	MemberSince            string  `json:"member_since"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              string  `json:"updated_at"`
}

// ListStudentsResponse is GET /api/students.
type ListStudentsResponse struct {
	Data       []StudentResponse `json:"data"`
	Total      int64             `json:"total"`
	Pagination PaginationMeta    `json:"pagination"`
}

// DepartmentResponse is an item in GET /api/departments.
type DepartmentResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// StudentFilter drives list/search for students.
type StudentFilter struct {
	Search string
}

// PersonalStats is the student_stats slice of the profile.
type PersonalStats struct {
	TotalRead         int  `json:"total_read"`
	ActiveLoans       int  `json:"active_loans"`
	OverdueCount      int  `json:"overdue_count"`
	CurrentStreakDays int  `json:"current_streak_days"`
	LongestStreakDays int  `json:"longest_streak_days,omitempty"`
	TotalPoints       int  `json:"total_points,omitempty"`
	GlobalRank        *int `json:"global_rank,omitempty"`
}

// LoanHistoryItem is one loan row in the student profile.
type LoanHistoryItem struct {
	LoanID     string   `json:"loan_id"`
	BookID     string   `json:"book_id"`
	Title      string   `json:"title"`
	CoverURL   *string  `json:"cover_url,omitempty"`
	Authors    []string `json:"authors"`
	BorrowedAt string   `json:"borrowed_at"`
	DueDate    string   `json:"due_date"`
	ReturnedAt *string  `json:"returned_at,omitempty"`
	Status     string   `json:"status"`
}

// BadgeGalleryItem is one badge in the gallery (earned or locked).
type BadgeGalleryItem struct {
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IconURL     *string `json:"icon_url,omitempty"`
	Criteria    *string `json:"criteria,omitempty"`
	Earned      bool    `json:"earned"`
	EarnedAt    *string `json:"earned_at,omitempty"`
}

// BadgeGallery groups badge progress for the profile.
type BadgeGallery struct {
	TotalBadges int                `json:"total_badges"`
	EarnedCount int                `json:"earned_count"`
	Badges      []BadgeGalleryItem `json:"badges"`
}

// StudentProfileResponse is GET /api/students/:id/profile.
type StudentProfileResponse struct {
	ID                     string            `json:"id"`
	StudentIDCode          string            `json:"student_id_code"`
	FirstName              string            `json:"first_name"`
	LastName               string            `json:"last_name"`
	AvatarURL              *string           `json:"avatar_url,omitempty"`
	DegreeLevel            *string           `json:"degree_level,omitempty"`
	Major                  *string           `json:"major,omitempty"`
	MemberSince            string            `json:"member_since"`
	DepartmentName         *string           `json:"department_name,omitempty"`
	DepartmentCode         *string           `json:"department_code,omitempty"`
	Email                  *string           `json:"email,omitempty"`
	ExpectedGraduationYear *int              `json:"expected_graduation_year,omitempty"`
	PersonalStats          PersonalStats     `json:"personal_stats"`
	LoanHistory            []LoanHistoryItem `json:"loan_history"`
	BadgeGallery           BadgeGallery      `json:"badge_gallery"`
}

// AwardBadgeRequest is POST /api/students/:id/badges.
type AwardBadgeRequest struct {
	BadgeID string `json:"badge_id" binding:"required,uuid"`
}

// AwardBadgeResponse is returned after awarding a badge.
type AwardBadgeResponse struct {
	ID       string `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	EarnedAt string `json:"earned_at"`
}
