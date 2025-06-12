package application_test

import (
	"testing"

	"go-hexagonal-template/internal/modules/user/application"
	"go-hexagonal-template/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestLoginUserUseCase_Execute_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock()
	useCase := application.NewLoginUserUseCase(mockRepo)
	input := application.LoginUserInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Act
	result, err := useCase.Execute(input)

	// Assert
	assert.NoError(t, err, "No debería haber error en el login exitoso")
	assert.NotNil(t, result, "El resultado no debería ser nil")
	assert.NotEmpty(t, result.Token, "El token no debería estar vacío")
	assert.NotNil(t, result.User, "El usuario no debería ser nil")
	assert.Equal(t, input.Email, result.User.Email, "El email no coincide")
}

func TestLoginUserUseCase_Execute_InvalidEmail(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock()
	useCase := application.NewLoginUserUseCase(mockRepo)
	input := application.LoginUserInput{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	// Act
	result, err := useCase.Execute(input)

	// Assert
	assert.Error(t, err, "Debería haber error con email inexistente")
	assert.Nil(t, result, "El resultado debería ser nil si el email no existe")
}

func TestLoginUserUseCase_Execute_InvalidPassword(t *testing.T) {
	mockRepo := mocks.NewUserRepositoryMock()
	useCase := application.NewLoginUserUseCase(mockRepo)
	input := application.LoginUserInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// Act
	result, err := useCase.Execute(input)

	// Assert
	assert.Error(t, err, "Debería haber error con contraseña incorrecta")
	assert.Nil(t, result, "El resultado debería ser nil si la contraseña es incorrecta")
}
