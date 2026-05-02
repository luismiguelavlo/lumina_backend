package book

import (
	"context"
	"testing"
	"time"

	"library_back/internal/models"
	"library_back/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakeBookRepo struct {
	existsISBN        bool
	existsCode        bool
	created           *models.Book
	createErr         error
	authorsByBook     map[string][]models.AuthorRef
	genresByBook      map[string][]models.GenreRef
	statusByBook      map[string]string
	updateErr         error
	deleteErr         error
	lastUpdateAuthors *[]string
	lastUpdateGenres  *[]string
}

func (f *fakeBookRepo) CreateWithRelations(ctx context.Context, b *models.Book, authorIDs, genreIDs []string) error {
	_ = ctx
	_ = authorIDs
	_ = genreIDs
	if f.createErr != nil {
		return f.createErr
	}
	f.created = b
	return nil
}

func (f *fakeBookRepo) GetByID(ctx context.Context, id string) (*models.Book, error) {
	_ = ctx
	if f.created != nil && f.created.ID == id {
		b := *f.created
		if b.CreatedAt.IsZero() {
			b.CreatedAt = time.Now().UTC()
			b.UpdatedAt = b.CreatedAt
		}
		return &b, nil
	}
	return nil, nil
}

func (f *fakeBookRepo) List(ctx context.Context, filter repositories.BookListFilter, limit, offset int) ([]models.BookListItem, int64, error) {
	return nil, 0, nil
}

func (f *fakeBookRepo) UpdateBook(ctx context.Context, b *models.Book, authorIDs, genreIDs *[]string) error {
	_ = ctx
	f.lastUpdateAuthors = authorIDs
	f.lastUpdateGenres = genreIDs
	return f.updateErr
}

func (f *fakeBookRepo) SoftDelete(ctx context.Context, id string) error {
	_ = ctx
	_ = id
	return f.deleteErr
}

func (f *fakeBookRepo) ExistsByISBN(ctx context.Context, isbn string, excludeID string) (bool, error) {
	_ = ctx
	_ = excludeID
	return f.existsISBN, nil
}

func (f *fakeBookRepo) ExistsByCatalogCode(ctx context.Context, code string, excludeID string) (bool, error) {
	_ = ctx
	_ = excludeID
	return f.existsCode, nil
}

func (f *fakeBookRepo) GetBookAuthors(ctx context.Context, bookID string) ([]models.AuthorRef, error) {
	_ = ctx
	if f.authorsByBook != nil {
		return f.authorsByBook[bookID], nil
	}
	return []models.AuthorRef{}, nil
}

func (f *fakeBookRepo) GetBookGenres(ctx context.Context, bookID string) ([]models.GenreRef, error) {
	_ = ctx
	if f.genresByBook != nil {
		return f.genresByBook[bookID], nil
	}
	return []models.GenreRef{}, nil
}

func (f *fakeBookRepo) GetBookStatus(ctx context.Context, bookID string) (string, error) {
	_ = ctx
	if f.statusByBook != nil {
		return f.statusByBook[bookID], nil
	}
	return "available", nil
}

type fakeAuthorRepo struct {
	exists map[string]bool
}

func (f *fakeAuthorRepo) List(ctx context.Context) ([]models.Author, error) { return nil, nil }

func (f *fakeAuthorRepo) Create(ctx context.Context, a *models.Author) error {
	_ = ctx
	_ = a
	return nil
}

func (f *fakeAuthorRepo) ExistsByID(ctx context.Context, id string) (bool, error) {
	return f.exists[id], nil
}

type fakeGenreRepo struct {
	exists map[string]bool
}

func (f *fakeGenreRepo) List(ctx context.Context) ([]models.Genre, error) { return nil, nil }

func (f *fakeGenreRepo) ExistsByID(ctx context.Context, id string) (bool, error) {
	return f.exists[id], nil
}

func TestBookService_CreateDuplicateISBN(t *testing.T) {
	s := NewBookService(
		&fakeBookRepo{existsISBN: true},
		&fakeAuthorRepo{exists: map[string]bool{}},
		&fakeGenreRepo{exists: map[string]bool{}},
	)
	_, err := s.Create(context.Background(), models.CreateBookRequest{
		Title: "T", ISBN: "1", CatalogCode: "C1",
	})
	if err != ErrDuplicateISBN {
		t.Fatalf("want ErrDuplicateISBN got %v", err)
	}
}

func TestBookService_CreateAuthorMissing(t *testing.T) {
	aid := uuid.NewString()
	s := NewBookService(
		&fakeBookRepo{},
		&fakeAuthorRepo{exists: map[string]bool{aid: false}},
		&fakeGenreRepo{exists: map[string]bool{}},
	)
	_, err := s.Create(context.Background(), models.CreateBookRequest{
		Title: "T", ISBN: "1", CatalogCode: "C1",
		AuthorIDs: []string{aid},
	})
	if err != ErrAuthorNotFound {
		t.Fatalf("want ErrAuthorNotFound got %v", err)
	}
}

func TestBookService_DeleteNotFound(t *testing.T) {
	s := NewBookService(
		&fakeBookRepo{deleteErr: gorm.ErrRecordNotFound},
		&fakeAuthorRepo{},
		&fakeGenreRepo{},
	)
	err := s.Delete(context.Background(), uuid.NewString())
	if err != ErrBookNotFound {
		t.Fatalf("want ErrBookNotFound got %v", err)
	}
}

func TestBookService_UpdatePatchesRelationsWhenSlicePresent(t *testing.T) {
	id := uuid.NewString()
	b := &models.Book{
		ID: id, Title: "Old", ISBN: "x", CatalogCode: "c",
		TotalCopies: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	fr := &fakeBookRepo{
		created: b,
		authorsByBook: map[string][]models.AuthorRef{
			id: {{ID: uuid.NewString(), Name: "A"}},
		},
		genresByBook: map[string][]models.GenreRef{},
		statusByBook: map[string]string{id: "available"},
	}
	aid := uuid.NewString()
	s := NewBookService(
		fr,
		&fakeAuthorRepo{exists: map[string]bool{aid: true}},
		&fakeGenreRepo{exists: map[string]bool{}},
	)
	empty := []string{}
	_, err := s.Update(context.Background(), id, models.UpdateBookRequest{AuthorIDs: empty})
	if err != nil {
		t.Fatal(err)
	}
	if fr.lastUpdateAuthors == nil || len(*fr.lastUpdateAuthors) != 0 {
		t.Fatalf("expected replace with empty authors")
	}
}
