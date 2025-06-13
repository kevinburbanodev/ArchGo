package model_test

import (
	"testing"
	"time"

	"go-hexagonal-template/internal/modules/user/domain/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserModel(t *testing.T) {
	t.Run("debería crear un usuario válido", func(t *testing.T) {
		// Arrange
		now := time.Now()
		user := &model.User{
			Name:      "Test User",
			Email:     "test@example.com",
			Password:  "hashedPassword123",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Assert
		assert.NotNil(t, user)
		assert.Equal(t, "Test User", user.Name)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "hashedPassword123", user.Password)
		assert.Equal(t, now, user.CreatedAt)
		assert.Equal(t, now, user.UpdatedAt)
		assert.Empty(t, user.DeletedAt)
	})

	t.Run("debería manejar un usuario con ID", func(t *testing.T) {
		// Arrange
		user := &model.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		}

		// Assert
		assert.Equal(t, uint(1), user.ID)
	})

	t.Run("debería manejar un usuario con soft delete", func(t *testing.T) {
		// Arrange
		now := time.Now()
		user := &model.User{
			Name:      "Test User",
			Email:     "test@example.com",
			DeletedAt: gorm.DeletedAt{Time: now, Valid: true},
		}

		// Assert
		assert.True(t, user.DeletedAt.Valid)
		assert.Equal(t, now, user.DeletedAt.Time)
	})
}
