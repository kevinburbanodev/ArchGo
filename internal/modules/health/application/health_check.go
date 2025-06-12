package application

import (
	"go-hexagonal-template/internal/modules/health/domain/model"
)

type HealthCheckUseCase struct{}

func NewHealthCheckUseCase() *HealthCheckUseCase {
	return &HealthCheckUseCase{}
}

func (uc *HealthCheckUseCase) Execute() *model.Health {
	return &model.Health{
		Status:  "UP",
		Message: "¡Healthy!",
	}
}
