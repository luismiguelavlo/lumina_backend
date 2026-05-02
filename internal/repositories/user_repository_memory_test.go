package repositories

import (
	"context"
	"testing"

	"library_back/internal/models"
)

func TestUserRepositoryInMemory_CreateFind(t *testing.T) {
	ctx := context.Background()
	r := NewUserRepositoryInMemory()
	u := &models.User{
		ID: "id-a", FirstName: "A", LastName: "B", Email: "a@b.com",
		PasswordHash: "h", Role: models.UserRoleAdmin, IsActive: false,
	}
	if err := r.Create(ctx, u); err != nil {
		t.Fatal(err)
	}
	got, err := r.FindByEmail(ctx, "a@b.com")
	if err != nil || got == nil || got.ID != "id-a" {
		t.Fatalf("FindByEmail: %+v err %v", got, err)
	}
	byID, err := r.FindByID(ctx, "id-a")
	if err != nil || byID == nil || byID.Email != "a@b.com" {
		t.Fatalf("FindByID: %+v err %v", byID, err)
	}
	miss, err := r.FindByEmail(ctx, "missing@test.com")
	if err != nil || miss != nil {
		t.Fatalf("FindByEmail miss: %+v err %v", miss, err)
	}
	missID, err := r.FindByID(ctx, "nope")
	if err != nil || missID != nil {
		t.Fatalf("FindByID miss: %+v err %v", missID, err)
	}
}

func TestUserRepositoryInMemory_DuplicateEmail(t *testing.T) {
	ctx := context.Background()
	r := NewUserRepositoryInMemory()
	u1 := &models.User{ID: "1", FirstName: "A", LastName: "B", Email: "dup@test.com", PasswordHash: "x", Role: models.UserRoleAdmin}
	u2 := &models.User{ID: "2", FirstName: "C", LastName: "D", Email: "dup@test.com", PasswordHash: "y", Role: models.UserRoleLibrarian}
	if err := r.Create(ctx, u1); err != nil {
		t.Fatal(err)
	}
	if err := r.Create(ctx, u2); err != ErrDuplicateEmail {
		t.Fatalf("want ErrDuplicateEmail got %v", err)
	}
}
