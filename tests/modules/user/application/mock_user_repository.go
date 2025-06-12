package application_test

import (
	"time"

	"go-hexagonal-template/internal/modules/user/domain/model"

	"github.com/stretchr/testify/assert"
)

// MockUserRepository es un mock del repositorio para testing
type MockUserRepository struct{}

func (m *MockUserRepository) Create(user *model.User) (*model.User, error) {
	if user.ID == "" {
		user.ID = "test-id-123" // ID fijo para testing
	}
	return user, nil
}

func (m *MockUserRepository) GetByID(id string) (*model.User, error) {
	if id == "nonexistent" {
		return nil, assert.AnError
	}
	now := time.Now().Format(time.RFC3339)
	return &model.User{
		ID:        id,
		Email:     "test@example.com",
		Name:      "Test",
		LastName:  "User",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
