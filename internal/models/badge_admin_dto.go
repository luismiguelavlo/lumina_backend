package models

// CreateBadgeDefinitionRequest is POST /api/badges.
type CreateBadgeDefinitionRequest struct {
	Slug        string  `json:"slug" binding:"required,min=1,max=50"`
	Name        string  `json:"name" binding:"required,min=1,max=100"`
	IconURL     *string `json:"icon_url,omitempty" binding:"omitempty,max=2048"`
	Description *string `json:"description,omitempty"`
	Criteria    *string `json:"criteria,omitempty"`
}

// PatchBadgeDefinitionRequest is PATCH /api/badges/:id (all fields optional).
type PatchBadgeDefinitionRequest struct {
	Slug        *string `json:"slug,omitempty" binding:"omitempty,min=1,max=50"`
	Name        *string `json:"name,omitempty" binding:"omitempty,min=1,max=100"`
	IconURL     *string `json:"icon_url,omitempty" binding:"omitempty,max=2048"`
	Description *string `json:"description,omitempty"`
	Criteria    *string `json:"criteria,omitempty"`
}

// ListBadgeDefinitionsResponse is GET /api/badges.
type ListBadgeDefinitionsResponse struct {
	Data []BadgeDefinitionResponse `json:"data"`
}

// BadgeDefinitionResponse is one badge row for admin APIs.
type BadgeDefinitionResponse struct {
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	IconURL     *string `json:"icon_url,omitempty"`
	Description *string `json:"description,omitempty"`
	Criteria    *string `json:"criteria,omitempty"`
	CreatedAt   string  `json:"created_at"`
}
