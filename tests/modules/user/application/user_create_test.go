package application_test

import (
	"testing"

	"go-hexagonal-template/internal/modules/user/application"
	"go-hexagonal-template/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserUseCase_Execute(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewUserRepositoryMock()
	useCase := application.NewCreateUserUseCase(mockRepo)
	input := application.CreateUserInput{
		Email:    "test@example.com",
		Name:     "Test",
		LastName: "User",
	}

	// Act
	user, err := useCase.Execute(input)

	// Assert
	assert.NoError(t, err, "No debería haber error al crear el usuario")
	assert.NotNil(t, user, "El usuario no debería ser nil")
	assert.Equal(t, input.Email, user.Email, "El email no coincide")
	assert.Equal(t, input.Name, user.Name, "El nombre no coincide")
	assert.Equal(t, input.LastName, user.LastName, "El apellido no coincide")
	assert.NotEmpty(t, user.CreatedAt, "La fecha de creación no debería estar vacía")
	assert.NotEmpty(t, user.UpdatedAt, "La fecha de actualización no debería estar vacía")
}
