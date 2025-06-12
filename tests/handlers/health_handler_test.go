package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-hexagonal-template/internal/handlers"
	"go-hexagonal-template/internal/modules/health/domain/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupHealthTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	healthHandler := handlers.NewHealthHandler()
	router.GET("/health", healthHandler.HealthCheck)
	return router
}

func TestHealthHandler_HealthCheck(t *testing.T) {
	// Arrange
	router := setupHealthTestRouter()
	expectedResponse := &model.Health{
		Status:  "UP",
		Message: "¡Healthy!",
	}

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", bytes.NewBuffer(nil))
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response model.Health
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Status, response.Status)
	assert.Equal(t, expectedResponse.Message, response.Message)
}

func TestHealthHandler_HealthCheck_WrongMethod(t *testing.T) {
	// Arrange
	router := setupHealthTestRouter()

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/health", bytes.NewBuffer(nil))
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}
