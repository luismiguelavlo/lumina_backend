package badgedef

import "errors"

var (
	ErrNotFound  = errors.New("badge not found")
	ErrInUse     = errors.New("badge has been awarded to students")
	ErrSlugTaken = errors.New("badge slug already exists")
)
