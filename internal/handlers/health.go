package handlers

import (
	"net/http"

	"go-hexagonal-template/internal/modules/health/application"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	healthCheckUseCase *application.HealthCheckUseCase
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		healthCheckUseCase: application.NewHealthCheckUseCase(),
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	health := h.healthCheckUseCase.Execute()
	c.JSON(http.StatusOK, health)
}
