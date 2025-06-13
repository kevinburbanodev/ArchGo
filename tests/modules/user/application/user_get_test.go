package application_test

import (
	"testing"

	"go-hexagonal-template/internal/modules/user/application"
	"go-hexagonal-template/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestGetUserUseCase_Execute(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewUserRepositoryMock()
	useCase := application.NewGetUserUseCase(mockRepo)
	expectedID := uint(123)

	// Act
	user, err := useCase.Execute(expectedID)

	// Assert
	assert.NoError(t, err, "No debería haber error al obtener el usuario")
	assert.NotNil(t, user, "El usuario no debería ser nil")
	assert.Equal(t, expectedID, user.ID, "El ID no coincide")
	assert.Equal(t, "test@example.com", user.Email, "El email no coincide")
	assert.Equal(t, "Test", user.Name, "El nombre no coincide")
	assert.NotEmpty(t, user.CreatedAt, "La fecha de creación no debería estar vacía")
	assert.NotEmpty(t, user.UpdatedAt, "La fecha de actualización no debería estar vacía")
}

func TestGetUserUseCase_Execute_NotFound(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewUserRepositoryMock()
	useCase := application.NewGetUserUseCase(mockRepo)

	// Act
	user, err := useCase.Execute(9999)

	// Assert
	assert.Error(t, err, "Debería retornar un error cuando el usuario no existe")
	assert.Nil(t, user, "El usuario debería ser nil cuando no existe")
}
