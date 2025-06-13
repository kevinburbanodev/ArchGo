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
	now := time.Now()
	if user.ID == 0 {
		user.ID = 1 // ID fijo para testing
	}
	user.CreatedAt = now
	user.UpdatedAt = now
	return user, nil
}

func (m *UserRepositoryMock) GetByID(id uint) (*model.User, error) {
	if id == 9999 {
		return nil, assert.AnError
	}
	now := time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	return &model.User{
		ID:        id,
		Email:     "test@example.com",
		Name:      "Test",
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (m *UserRepositoryMock) GetByEmail(email string) (*model.User, error) {
	if email == "nonexistent@example.com" {
		return nil, assert.AnError
	}
	now := time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	return &model.User{
		ID:        1,
		Email:     email,
		Name:      "Test",
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
