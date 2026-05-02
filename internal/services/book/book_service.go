package book

import (
	"context"
	"errors"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BookService implements catalog book use cases.
type BookService struct {
	books   repositories.CatalogBookRepository
	authors repositories.AuthorCatalogRepository
	genres  repositories.GenreCatalogRepository
}

// NewBookService constructs BookService.
func NewBookService(
	books repositories.CatalogBookRepository,
	authors repositories.AuthorCatalogRepository,
	genres repositories.GenreCatalogRepository,
) *BookService {
	return &BookService{books: books, authors: authors, genres: genres}
}

// Create validates relations and uniqueness, then persists the book.
func (s *BookService) Create(ctx context.Context, req models.CreateBookRequest) (*models.BookDetailResponse, error) {
	if err := s.ensureAuthorIDs(ctx, req.AuthorIDs); err != nil {
		return nil, err
	}
	if err := s.ensureGenreIDs(ctx, req.GenreIDs); err != nil {
		return nil, err
	}
	dup, err := s.books.ExistsByISBN(ctx, req.ISBN, "")
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, ErrDuplicateISBN
	}
	dup, err = s.books.ExistsByCatalogCode(ctx, req.CatalogCode, "")
	if err != nil {
		return nil, err
	}
	if dup {
		return nil, ErrDuplicateCatalogCode
	}

	copies := 1
	if req.TotalCopies != nil {
		copies = *req.TotalCopies
	}
	b := &models.Book{
		ID:              uuid.NewString(),
		Title:           req.Title,
		ISBN:            req.ISBN,
		CatalogCode:     req.CatalogCode,
		Synopsis:        req.Synopsis,
		PublicationYear: req.PublicationYear,
		Pages:           req.Pages,
		CoverURL:        req.CoverURL,
		Location:        req.Location,
		TotalCopies:     copies,
	}
	if err := s.books.CreateWithRelations(ctx, b, req.AuthorIDs, req.GenreIDs); err != nil {
		return nil, err
	}
	return s.buildDetail(ctx, b.ID)
}

// GetByID returns book detail or ErrBookNotFound.
func (s *BookService) GetByID(ctx context.Context, id string) (*models.BookDetailResponse, error) {
	return s.buildDetail(ctx, id)
}

// List returns paginated summarized books.
func (s *BookService) List(ctx context.Context, filter repositories.BookListFilter, limit, offset int) ([]models.BookListItem, int64, error) {
	return s.books.List(ctx, filter, limit, offset)
}

// Update merges fields and optionally replaces author/genre links.
func (s *BookService) Update(ctx context.Context, id string, req models.UpdateBookRequest) (*models.BookDetailResponse, error) {
	b, err := s.books.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrBookNotFound
	}

	if req.ISBN != nil {
		dup, err := s.books.ExistsByISBN(ctx, *req.ISBN, id)
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, ErrDuplicateISBN
		}
		b.ISBN = *req.ISBN
	}
	if req.CatalogCode != nil {
		dup, err := s.books.ExistsByCatalogCode(ctx, *req.CatalogCode, id)
		if err != nil {
			return nil, err
		}
		if dup {
			return nil, ErrDuplicateCatalogCode
		}
		b.CatalogCode = *req.CatalogCode
	}
	if req.Title != nil {
		b.Title = *req.Title
	}
	if req.Synopsis != nil {
		b.Synopsis = req.Synopsis
	}
	if req.PublicationYear != nil {
		b.PublicationYear = req.PublicationYear
	}
	if req.Pages != nil {
		b.Pages = req.Pages
	}
	if req.CoverURL != nil {
		b.CoverURL = req.CoverURL
	}
	if req.Location != nil {
		b.Location = req.Location
	}
	if req.TotalCopies != nil {
		b.TotalCopies = *req.TotalCopies
	}

	var authorPatch *[]string
	if req.AuthorIDs != nil {
		if err := s.ensureAuthorIDs(ctx, req.AuthorIDs); err != nil {
			return nil, err
		}
		authorPatch = &req.AuthorIDs
	}
	var genrePatch *[]string
	if req.GenreIDs != nil {
		if err := s.ensureGenreIDs(ctx, req.GenreIDs); err != nil {
			return nil, err
		}
		genrePatch = &req.GenreIDs
	}

	if err := s.books.UpdateBook(ctx, b, authorPatch, genrePatch); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}
	return s.buildDetail(ctx, id)
}

// Delete soft-deletes a book.
func (s *BookService) Delete(ctx context.Context, id string) error {
	err := s.books.SoftDelete(ctx, id)
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrBookNotFound
	}
	return err
}

func (s *BookService) ensureAuthorIDs(ctx context.Context, ids []string) error {
	for _, id := range ids {
		ok, err := s.authors.ExistsByID(ctx, id)
		if err != nil {
			return err
		}
		if !ok {
			return ErrAuthorNotFound
		}
	}
	return nil
}

func (s *BookService) ensureGenreIDs(ctx context.Context, ids []string) error {
	for _, id := range ids {
		ok, err := s.genres.ExistsByID(ctx, id)
		if err != nil {
			return err
		}
		if !ok {
			return ErrGenreNotFound
		}
	}
	return nil
}

func (s *BookService) buildDetail(ctx context.Context, bookID string) (*models.BookDetailResponse, error) {
	b, err := s.books.GetByID(ctx, bookID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, ErrBookNotFound
	}
	authors, err := s.books.GetBookAuthors(ctx, bookID)
	if err != nil {
		return nil, err
	}
	genres, err := s.books.GetBookGenres(ctx, bookID)
	if err != nil {
		return nil, err
	}
	status, err := s.books.GetBookStatus(ctx, bookID)
	if err != nil {
		return nil, err
	}
	if authors == nil {
		authors = []models.AuthorRef{}
	}
	if genres == nil {
		genres = []models.GenreRef{}
	}
	return &models.BookDetailResponse{
		ID:              b.ID,
		Title:           b.Title,
		ISBN:            b.ISBN,
		CatalogCode:     b.CatalogCode,
		Synopsis:        b.Synopsis,
		PublicationYear: b.PublicationYear,
		Pages:           b.Pages,
		CoverURL:        b.CoverURL,
		Location:        b.Location,
		TotalCopies:     b.TotalCopies,
		Status:          status,
		Authors:         authors,
		Genres:          genres,
		CreatedAt:       b.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:       b.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}
