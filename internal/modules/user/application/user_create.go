package application

import (
	"time"

	"go-hexagonal-template/internal/modules/user/domain/model"
	"go-hexagonal-template/internal/modules/user/domain/port"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserUseCase struct {
	userRepository port.UserRepository
}

func NewCreateUserUseCase(userRepository port.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepository: userRepository,
	}
}

type CreateUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=6"`
}

func (uc *CreateUserUseCase) Execute(input CreateUserInput) (*model.User, error) {
	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:     input.Email,
		Name:      input.Name,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return uc.userRepository.Create(user)
}
