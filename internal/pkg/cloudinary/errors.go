package cloudinary

import "errors"

var (
	ErrInvalidFolder = errors.New("invalid upload folder")
	ErrEmptyFile     = errors.New("empty file")
	ErrTooLarge      = errors.New("file exceeds maximum size")
	ErrProvider      = errors.New("cloudinary provider error")
	ErrNotConfigured = errors.New("cloudinary is not configured")
)
