package author

import (
	"context"
	"testing"

	"library_back/internal/models"
)

type fakeAuthorRepo struct {
	last *models.Author
}

func (f *fakeAuthorRepo) Create(ctx context.Context, a *models.Author) error {
	_ = ctx
	f.last = a
	return nil
}

func TestService_Create(t *testing.T) {
	fr := &fakeAuthorRepo{}
	s := NewService(fr)
	out, err := s.Create(context.Background(), models.CreateAuthorRequest{Name: "  Ada Lovelace  "})
	if err != nil {
		t.Fatal(err)
	}
	if out.Name != "Ada Lovelace" || out.ID == "" || fr.last == nil || fr.last.Name != "Ada Lovelace" {
		t.Fatalf("out=%+v last=%+v", out, fr.last)
	}
}
