package seeds

import (
	"context"
	"errors"
	"fmt"

	"library_back/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Summary describes what the public seed step did (best-effort counts).
type Summary struct {
	Departments int `json:"departments_upserted"`
	Genres      int `json:"genres_upserted"`
	Authors     int `json:"authors_inserted"`
	Books       int `json:"books_upserted"`
	BookLinks   int `json:"book_relations_inserted"`
	Badges      int `json:"badges_upserted"`
}

// Run inserts reference data that has no admin CRUD in the API (departments, authors,
// genres, demo books + links). Idempotent: safe to call multiple times.
func Run(ctx context.Context, db *gorm.DB) (*Summary, error) {
	var out Summary
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		out.Departments, err = seedDepartments(tx)
		if err != nil {
			return err
		}
		out.Genres, err = seedGenres(tx)
		if err != nil {
			return err
		}
		out.Authors, err = seedAuthors(tx)
		if err != nil {
			return err
		}
		out.Badges, err = seedBadges(tx)
		if err != nil {
			return err
		}
		out.Books, out.BookLinks, err = seedBooksAndLinks(tx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func seedDepartments(tx *gorm.DB) (int, error) {
	rows := []models.Department{
		{ID: "10000000-0000-4000-8000-000000000001", Name: "Computer Science", Code: "CS"},
		{ID: "10000000-0000-4000-8000-000000000002", Name: "English Literature", Code: "ENG"},
		{ID: "10000000-0000-4000-8000-000000000003", Name: "Mathematics", Code: "MATH"},
		{ID: "10000000-0000-4000-8000-000000000004", Name: "Physics", Code: "PHY"},
		{ID: "10000000-0000-4000-8000-000000000005", Name: "Fine Arts", Code: "ART"},
	}
	n := 0
	for _, d := range rows {
		r := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoNothing: true,
		}).Create(&d)
		if r.Error != nil {
			return n, r.Error
		}
		if r.RowsAffected > 0 {
			n++
		}
	}
	return n, nil
}

func seedGenres(tx *gorm.DB) (int, error) {
	rows := []models.Genre{
		{ID: "20000000-0000-4000-8000-000000000001", Name: "Fantasy", Code: "FIC"},
		{ID: "20000000-0000-4000-8000-000000000002", Name: "Science Fiction", Code: "SCI"},
		{ID: "20000000-0000-4000-8000-000000000003", Name: "Romance", Code: "ROM"},
		{ID: "20000000-0000-4000-8000-000000000004", Name: "Non-fiction", Code: "NF"},
	}
	n := 0
	for _, g := range rows {
		r := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoNothing: true,
		}).Create(&g)
		if r.Error != nil {
			return n, r.Error
		}
		if r.RowsAffected > 0 {
			n++
		}
	}
	return n, nil
}

func seedAuthors(tx *gorm.DB) (int, error) {
	rows := []models.Author{
		{ID: "30000000-0000-4000-8000-000000000001", Name: "Patrick Rothfuss"},
		{ID: "30000000-0000-4000-8000-000000000002", Name: "Ursula K. Le Guin"},
		{ID: "30000000-0000-4000-8000-000000000003", Name: "Isaac Asimov"},
	}
	n := 0
	for _, a := range rows {
		var existing models.Author
		err := tx.Where("id = ?", a.ID).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return n, err
		}
		if err := tx.Create(&a).Error; err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

func seedBadges(tx *gorm.DB) (int, error) {
	// Same slugs as migration 001; ON CONFLICT DO NOTHING keeps idempotent runs safe.
	rows := []models.Badge{
		{Slug: "punctual_reader", Name: "Punctual Reader", Description: ptr("Always returns books on time"), Criteria: ptr("Return 10+ books before due date with no overdue")},
		{Slug: "bookworm", Name: "Bookworm", Description: ptr("Avid reader with impressive stats"), Criteria: ptr("Read 50+ books")},
		{Slug: "explorer", Name: "Explorer", Description: ptr("Reads across diverse genres"), Criteria: ptr("Borrow books from 5+ different genres")},
		{Slug: "night_owl", Name: "Night Owl", Description: ptr("Late-night library patron"), Criteria: ptr("Borrow books during night hours 10+ times")},
		{Slug: "reviewer", Name: "Reviewer", Description: ptr("Active contributor to book reviews"), Criteria: ptr("Write 10+ book reviews")},
		{Slug: "collector", Name: "Collector", Description: ptr("Has read the most books in a category"), Criteria: ptr("Read 20+ books in a single genre")},
	}
	n := 0
	for _, b := range rows {
		r := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "slug"}},
			DoNothing: true,
		}).Create(&b)
		if r.Error != nil {
			return n, r.Error
		}
		if r.RowsAffected > 0 {
			n++
		}
	}
	return n, nil
}

func ptr(s string) *string { return &s }

func seedBooksAndLinks(tx *gorm.DB) (books int, links int, err error) {
	type bookSeed struct {
		Book         models.Book
		AuthorIDs    []string
		GenreCodes   []string
	}
	seeds := []bookSeed{
		{
			Book: models.Book{
				ID:              "40000000-0000-4000-8000-000000000001",
				Title:           "The Name of the Wind",
				ISBN:            "978-SEED-DEMO-0001",
				CatalogCode:     "FIC-SEED-001",
				Synopsis:        ptr("Demo title for catalog and loans."),
				PublicationYear: ptrInt(2007),
				Pages:           ptrInt(662),
				Location:        ptr("Section A, Shelf 1"),
				TotalCopies:     5,
			},
			AuthorIDs:  []string{"30000000-0000-4000-8000-000000000001"},
			GenreCodes: []string{"FIC"},
		},
		{
			Book: models.Book{
				ID:              "40000000-0000-4000-8000-000000000002",
				Title:           "The Left Hand of Darkness",
				ISBN:            "978-SEED-DEMO-0002",
				CatalogCode:     "SCI-SEED-001",
				Synopsis:        ptr("Demo science fiction title."),
				PublicationYear: ptrInt(1969),
				Pages:           ptrInt(320),
				Location:        ptr("Section B, Shelf 2"),
				TotalCopies:     3,
			},
			AuthorIDs:  []string{"30000000-0000-4000-8000-000000000002"},
			GenreCodes: []string{"SCI"},
		},
	}
	for _, s := range seeds {
		r := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "isbn"}},
			DoNothing: true,
		}).Create(&s.Book)
		if r.Error != nil {
			return books, links, r.Error
		}
		if r.RowsAffected > 0 {
			books++
		}

		for _, aid := range s.AuthorIDs {
			res := tx.Exec(`INSERT INTO book_authors (book_id, author_id) VALUES (?, ?) ON CONFLICT DO NOTHING`, s.Book.ID, aid)
			if res.Error != nil {
				return books, links, res.Error
			}
			links += int(res.RowsAffected)
		}
		for _, code := range s.GenreCodes {
			var gid string
			if err := tx.Model(&models.Genre{}).Select("id").Where("code = ?", code).Scan(&gid).Error; err != nil {
				return books, links, err
			}
			if gid == "" {
				return books, links, fmt.Errorf("genre code %q not found after seed", code)
			}
			res := tx.Exec(`INSERT INTO book_genres (book_id, genre_id) VALUES (?, ?) ON CONFLICT DO NOTHING`, s.Book.ID, gid)
			if res.Error != nil {
				return books, links, res.Error
			}
			links += int(res.RowsAffected)
		}
	}
	return books, links, nil
}

func ptrInt(i int) *int { return &i }
