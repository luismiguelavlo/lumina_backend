package repositories

import (
	"context"
	"fmt"
	"strings"

	"library_back/internal/models"

	"gorm.io/gorm"
)

// BookListFilter drives GET /api/books search.
type BookListFilter struct {
	Search  string
	GenreID string
}

// CatalogBookRepository persists books and M:N links for the catalog API.
type CatalogBookRepository interface {
	CreateWithRelations(ctx context.Context, b *models.Book, authorIDs, genreIDs []string) error
	GetByID(ctx context.Context, id string) (*models.Book, error)
	List(ctx context.Context, filter BookListFilter, limit, offset int) ([]models.BookListItem, int64, error)
	UpdateBook(ctx context.Context, b *models.Book, authorIDs, genreIDs *[]string) error
	SoftDelete(ctx context.Context, id string) error
	ExistsByISBN(ctx context.Context, isbn string, excludeID string) (bool, error)
	ExistsByCatalogCode(ctx context.Context, code string, excludeID string) (bool, error)
	GetBookAuthors(ctx context.Context, bookID string) ([]models.AuthorRef, error)
	GetBookGenres(ctx context.Context, bookID string) ([]models.GenreRef, error)
	GetBookStatus(ctx context.Context, bookID string) (string, error)
}

type catalogBookRepositoryGorm struct {
	db *gorm.DB
}

// NewCatalogBookRepositoryGorm returns a GORM-backed CatalogBookRepository.
func NewCatalogBookRepositoryGorm(db *gorm.DB) CatalogBookRepository {
	return &catalogBookRepositoryGorm{db: db}
}

func (r *catalogBookRepositoryGorm) CreateWithRelations(ctx context.Context, b *models.Book, authorIDs, genreIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(b).Error; err != nil {
			return err
		}
		if err := replaceBookAuthors(tx, b.ID, authorIDs); err != nil {
			return err
		}
		return replaceBookGenres(tx, b.ID, genreIDs)
	})
}

func (r *catalogBookRepositoryGorm) GetByID(ctx context.Context, id string) (*models.Book, error) {
	var b models.Book
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&b).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

func (r *catalogBookRepositoryGorm) List(ctx context.Context, filter BookListFilter, limit, offset int) ([]models.BookListItem, int64, error) {
	search := strings.TrimSpace(filter.Search)
	genreID := strings.TrimSpace(filter.GenreID)
	pattern := "%" + search + "%"

	baseWhere := "b.deleted_at IS NULL"
	args := []interface{}{}
	if search != "" {
		baseWhere += ` AND (
			b.title ILIKE ? OR
			b.isbn ILIKE ? OR
			EXISTS (
				SELECT 1 FROM book_authors ba
				JOIN authors a ON a.id = ba.author_id
				WHERE ba.book_id = b.id AND a.name ILIKE ?
			)
		)`
		args = append(args, pattern, pattern, pattern)
	}
	if genreID != "" {
		baseWhere += ` AND EXISTS (
			SELECT 1
			FROM book_genres bg
			WHERE bg.book_id = b.id AND bg.genre_id = ?
		)`
		args = append(args, genreID)
	}

	var total int64
	countSQL := fmt.Sprintf(`SELECT COUNT(*) FROM books b WHERE %s`, baseWhere)
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	type row struct {
		ID       string
		CoverURL *string
		Title    string
		ISBN     string
		Status   string
		Author   string
	}
	var rows []row
	listSQL := fmt.Sprintf(`
SELECT
	b.id,
	b.cover_url,
	b.title,
	b.isbn,
	COALESCE(ba.status, 'available') AS status,
	COALESCE((
		SELECT string_agg(a.name, ', ' ORDER BY a.name)
		FROM book_authors x
		JOIN authors a ON a.id = x.author_id
		WHERE x.book_id = b.id
	), '') AS author
FROM books b
LEFT JOIN book_availability ba ON ba.book_id = b.id
WHERE %s
ORDER BY b.created_at DESC
LIMIT ? OFFSET ?
`, baseWhere)
	listArgs := append(append([]interface{}{}, args...), limit, offset)
	if err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	out := make([]models.BookListItem, 0, len(rows))
	for _, rw := range rows {
		out = append(out, models.BookListItem{
			ID:       rw.ID,
			CoverURL: rw.CoverURL,
			Title:    rw.Title,
			ISBN:     rw.ISBN,
			Status:   rw.Status,
			Author:   rw.Author,
		})
	}
	return out, total, nil
}

func (r *catalogBookRepositoryGorm) UpdateBook(ctx context.Context, b *models.Book, authorIDs, genreIDs *[]string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"title":             b.Title,
			"isbn":              b.ISBN,
			"catalog_code":      b.CatalogCode,
			"synopsis":          b.Synopsis,
			"publication_year":  b.PublicationYear,
			"pages":             b.Pages,
			"cover_url":         b.CoverURL,
			"location":          b.Location,
			"total_copies":      b.TotalCopies,
		}
		res := tx.Model(&models.Book{}).Where("id = ? AND deleted_at IS NULL", b.ID).Updates(updates)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if authorIDs != nil {
			if err := replaceBookAuthors(tx, b.ID, *authorIDs); err != nil {
				return err
			}
		}
		if genreIDs != nil {
			if err := replaceBookGenres(tx, b.ID, *genreIDs); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *catalogBookRepositoryGorm) SoftDelete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Model(&models.Book{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()"))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *catalogBookRepositoryGorm) ExistsByISBN(ctx context.Context, isbn string, excludeID string) (bool, error) {
	q := r.db.WithContext(ctx).Model(&models.Book{}).Where("isbn = ? AND deleted_at IS NULL", isbn)
	if excludeID != "" {
		q = q.Where("id != ?", excludeID)
	}
	var n int64
	if err := q.Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *catalogBookRepositoryGorm) ExistsByCatalogCode(ctx context.Context, code string, excludeID string) (bool, error) {
	q := r.db.WithContext(ctx).Model(&models.Book{}).Where("catalog_code = ? AND deleted_at IS NULL", code)
	if excludeID != "" {
		q = q.Where("id != ?", excludeID)
	}
	var n int64
	if err := q.Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *catalogBookRepositoryGorm) GetBookAuthors(ctx context.Context, bookID string) ([]models.AuthorRef, error) {
	var refs []models.AuthorRef
	err := r.db.WithContext(ctx).
		Table("authors").
		Select("authors.id AS id, authors.name AS name").
		Joins("JOIN book_authors ON book_authors.author_id = authors.id").
		Where("book_authors.book_id = ?", bookID).
		Order("authors.name ASC").
		Scan(&refs).Error
	return refs, err
}

func (r *catalogBookRepositoryGorm) GetBookGenres(ctx context.Context, bookID string) ([]models.GenreRef, error) {
	var refs []models.GenreRef
	err := r.db.WithContext(ctx).
		Table("genres").
		Select("genres.id AS id, genres.name AS name, genres.code AS code").
		Joins("JOIN book_genres ON book_genres.genre_id = genres.id").
		Where("book_genres.book_id = ?", bookID).
		Order("genres.name ASC").
		Scan(&refs).Error
	return refs, err
}

func (r *catalogBookRepositoryGorm) GetBookStatus(ctx context.Context, bookID string) (string, error) {
	var st string
	err := r.db.WithContext(ctx).
		Raw(`SELECT status FROM book_availability WHERE book_id = ?`, bookID).
		Scan(&st).Error
	if err != nil {
		return "", err
	}
	if st == "" {
		return "available", nil
	}
	return st, nil
}

func replaceBookAuthors(tx *gorm.DB, bookID string, authorIDs []string) error {
	if err := tx.Where("book_id = ?", bookID).Delete(&models.BookAuthor{}).Error; err != nil {
		return err
	}
	for _, aid := range authorIDs {
		if err := tx.Create(&models.BookAuthor{BookID: bookID, AuthorID: aid}).Error; err != nil {
			return err
		}
	}
	return nil
}

func replaceBookGenres(tx *gorm.DB, bookID string, genreIDs []string) error {
	if err := tx.Where("book_id = ?", bookID).Delete(&models.BookGenre{}).Error; err != nil {
		return err
	}
	for _, gid := range genreIDs {
		if err := tx.Create(&models.BookGenre{BookID: bookID, GenreID: gid}).Error; err != nil {
			return err
		}
	}
	return nil
}
