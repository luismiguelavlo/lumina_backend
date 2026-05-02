package book

import "errors"

var (
	ErrBookNotFound         = errors.New("book not found")
	ErrAuthorNotFound       = errors.New("author not found")
	ErrGenreNotFound        = errors.New("genre not found")
	ErrDuplicateISBN        = errors.New("isbn already exists")
	ErrDuplicateCatalogCode = errors.New("catalog_code already exists")
)
