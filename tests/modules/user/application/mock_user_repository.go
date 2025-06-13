package application_test

import (
	"time"

	"go-hexagonal-template/internal/modules/user/domain/model"

	"github.com/stretchr/testify/assert"
)

// MockUserRepository es un mock del repositorio para testing
type MockUserRepository struct{}

func (m *MockUserRepository) Create(user *model.User) (*model.User, error) {
	if user.ID == 0 {
		user.ID = 1 // ID fijo para testing
	}
	return user, nil
}

func (m *MockUserRepository) GetByID(id uint) (*model.User, error) {
	if id == 9999 {
		return nil, assert.AnError
	}
	now := time.Now()
	return &model.User{
		ID:        id,
		Email:     "test@example.com",
		Name:      "Test",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
