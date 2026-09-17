package service

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Mosteben/hotel-booking-system/internal/user/model"
	"github.com/google/uuid"
)

type mockUserLister struct {
	users []model.User
	err   error
}

func (m *mockUserLister) GetAll() ([]model.User, error) {
	return m.users, m.err
}

func TestGetAllUsers_MapsSafeFieldsOnly(t *testing.T) {
	id := uuid.New()
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	repo := &mockUserLister{
		users: []model.User{
			{
				ID:         id,
				FirstName:  "Ada",
				LastName:   "Lovelace",
				Email:      "ada@example.com",
				Password:   "super-secret-hash",
				Phone:      "+201000000000",
				Role:       "customer",
				IsActive:   true,
				IsVerified: false,
				CreatedAt:  createdAt,
			},
		},
	}

	svc := NewUserService(repo)

	summaries, err := svc.GetAllUsers()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}

	got := summaries[0]

	if got.ID != id.String() {
		t.Errorf("expected id %q, got %q", id.String(), got.ID)
	}

	if got.Email != "ada@example.com" {
		t.Errorf("expected email to be mapped, got %q", got.Email)
	}

	if got.CreatedAt != createdAt.Format(time.RFC3339) {
		t.Errorf("expected created_at %q, got %q", createdAt.Format(time.RFC3339), got.CreatedAt)
	}

	// AdminUserSummary has no password field at all, but this guards
	// against that ever silently changing in the future.
	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("failed to marshal summary: %v", err)
	}
	if strings.Contains(string(encoded), "super-secret-hash") {
		t.Errorf("admin user summary must never serialize the password hash")
	}
}

func TestGetAllUsers_PropagatesRepositoryError(t *testing.T) {
	repo := &mockUserLister{err: errors.New("db unavailable")}

	svc := NewUserService(repo)

	_, err := svc.GetAllUsers()

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestGetAllUsers_EmptyList(t *testing.T) {
	svc := NewUserService(&mockUserLister{users: []model.User{}})

	summaries, err := svc.GetAllUsers()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(summaries) != 0 {
		t.Errorf("expected an empty slice, got %d items", len(summaries))
	}
}
