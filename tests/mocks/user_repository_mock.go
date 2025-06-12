package mocks

import (
	"time"

	"go-hexagonal-template/internal/modules/user/domain/model"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// UserRepositoryMock es un mock del repositorio para testing
type UserRepositoryMock struct{}

func NewUserRepositoryMock() *UserRepositoryMock {
	return &UserRepositoryMock{}
}

func (m *UserRepositoryMock) Create(user *model.User) (*model.User, error) {
	if user.ID == "" {
		user.ID = "test-id-123" // ID fijo para testing
	}
	return user, nil
}

func (m *UserRepositoryMock) GetByID(id string) (*model.User, error) {
	if id == "nonexistent" {
		return nil, assert.AnError
	}
	now := time.Now().Format(time.RFC3339)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	return &model.User{
		ID:        id,
		Email:     "test@example.com",
		Name:      "Test",
		LastName:  "User",
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (m *UserRepositoryMock) GetByEmail(email string) (*model.User, error) {
	if email == "nonexistent@example.com" {
		return nil, assert.AnError
	}
	now := time.Now().Format(time.RFC3339)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	return &model.User{
		ID:        "test-id-123",
		Email:     email,
		Name:      "Test",
		LastName:  "User",
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
