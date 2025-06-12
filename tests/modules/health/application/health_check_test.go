package application_test

import (
	"testing"

	"go-hexagonal-template/internal/modules/health/application"
	"go-hexagonal-template/internal/modules/health/domain/model"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckUseCase_Execute(t *testing.T) {
	// Arrange
	useCase := application.NewHealthCheckUseCase()
	expectedHealth := &model.Health{
		Status:  "UP",
		Message: "¡Healthy!",
	}

	// Act
	health := useCase.Execute()

	// Assert
	assert.Equal(t, expectedHealth.Status, health.Status)
	assert.Equal(t, expectedHealth.Message, health.Message)
}
