package models

// CreateBookRequest is the JSON body for POST /api/books.
type CreateBookRequest struct {
	Title            string   `json:"title" binding:"required,max=300"`
	ISBN             string   `json:"isbn" binding:"required,max=20"`
	CatalogCode      string   `json:"catalog_code" binding:"required,max=20"`
	Synopsis         *string  `json:"synopsis"`
	PublicationYear  *int     `json:"publication_year"`
	Pages            *int     `json:"pages" binding:"omitempty,min=0"`
	CoverURL         *string  `json:"cover_url" binding:"omitempty,url"`
	Location         *string  `json:"location" binding:"omitempty,max=200"`
	TotalCopies      *int     `json:"total_copies" binding:"omitempty,min=0"`
	AuthorIDs        []string `json:"author_ids" binding:"omitempty,dive,uuid"`
	GenreIDs         []string `json:"genre_ids" binding:"omitempty,dive,uuid"`
}

// UpdateBookRequest is the JSON body for PUT/PATCH /api/books/:id (partial via pointers / omitempty).
type UpdateBookRequest struct {
	Title            *string  `json:"title" binding:"omitempty,max=300"`
	ISBN             *string  `json:"isbn" binding:"omitempty,max=20"`
	CatalogCode      *string  `json:"catalog_code" binding:"omitempty,max=20"`
	Synopsis         *string  `json:"synopsis"`
	PublicationYear  *int     `json:"publication_year"`
	Pages            *int     `json:"pages" binding:"omitempty,min=0"`
	CoverURL         *string  `json:"cover_url" binding:"omitempty,url"`
	Location         *string  `json:"location" binding:"omitempty,max=200"`
	TotalCopies      *int     `json:"total_copies" binding:"omitempty,min=0"`
	AuthorIDs        []string `json:"author_ids" binding:"omitempty,dive,uuid"`
	GenreIDs         []string `json:"genre_ids" binding:"omitempty,dive,uuid"`
}

// BookListItem is a summarized row for GET /api/books.
type BookListItem struct {
	ID       string  `json:"id"`
	CoverURL *string `json:"cover_url,omitempty"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	ISBN     string  `json:"isbn"`
	Status   string  `json:"status"`
}

// AuthorRef is an author on a book detail response.
type AuthorRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GenreRef is a genre on a book detail response.
type GenreRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// BookDetailResponse is GET /api/books/:id and create/update payloads.
type BookDetailResponse struct {
	ID               string      `json:"id"`
	Title            string      `json:"title"`
	ISBN             string      `json:"isbn"`
	CatalogCode      string      `json:"catalog_code"`
	Synopsis         *string     `json:"synopsis,omitempty"`
	PublicationYear  *int        `json:"publication_year,omitempty"`
	Pages            *int        `json:"pages,omitempty"`
	CoverURL         *string     `json:"cover_url,omitempty"`
	Location         *string     `json:"location,omitempty"`
	TotalCopies      int         `json:"total_copies"`
	Status           string      `json:"status"`
	Authors          []AuthorRef `json:"authors"`
	Genres           []GenreRef  `json:"genres"`
	CreatedAt        string      `json:"created_at"`
	UpdatedAt        string      `json:"updated_at"`
}

// AuthorResponse is an item in GET /api/authors.
type AuthorResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CreateAuthorRequest is the JSON body for POST /api/authors (staff).
type CreateAuthorRequest struct {
	Name string  `json:"name" binding:"required,min=1,max=200"`
	Bio  *string `json:"bio,omitempty" binding:"omitempty,max=5000"`
}

// AuthorCreatedResponse is returned after POST /api/authors.
type AuthorCreatedResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Bio       *string `json:"bio,omitempty"`
	CreatedAt string  `json:"created_at"`
}

// GenreResponse is an item in GET /api/genres.
type GenreResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// ListBooksResponse is GET /api/books.
type ListBooksResponse struct {
	Data       []BookListItem   `json:"data"`
	Total      int64            `json:"total"`
	Pagination PaginationMeta   `json:"pagination"`
}

// PaginationMeta is paging context for list endpoints.
type PaginationMeta struct {
	Limit       int   `json:"limit"`
	Offset      int   `json:"offset"`
	Count       int   `json:"count"`
	Total       int64 `json:"total"`
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}
