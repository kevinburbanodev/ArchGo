package application

import (
	"errors"
	"fmt"

	"go-hexagonal-template/internal/infrastructure/auth"
	"go-hexagonal-template/internal/modules/user/domain/model"
	"go-hexagonal-template/internal/modules/user/domain/port"

	"golang.org/x/crypto/bcrypt"
)

type LoginUserUseCase struct {
	userRepository port.UserRepository
}

func NewLoginUserUseCase(userRepository port.UserRepository) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepository: userRepository,
	}
}

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func (uc *LoginUserUseCase) Execute(input LoginUserInput) (*LoginUserOutput, error) {
	// Buscar el usuario por email
	user, err := uc.userRepository.GetByEmail(input.Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Verificar la contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar el token JWT
	token, err := auth.GenerateToken(fmt.Sprintf("%d", user.ID), user.Email)
	if err != nil {
		return nil, err
	}

	return &LoginUserOutput{
		Token: token,
		User:  user,
	}, nil
}
