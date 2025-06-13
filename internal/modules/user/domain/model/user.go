package model

import (
	"time"

	"gorm.io/gorm"
)

// User representa un usuario en el sistema
// @Description Modelo de usuario del sistema
type User struct {
	// @Description ID único del usuario
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	// @Description Nombre del usuario
	Name string `json:"name" binding:"required"`

	// @Description Correo electrónico del usuario (debe ser único)
	Email string `json:"email" binding:"required,email" gorm:"unique"`

	// @Description Contraseña del usuario (hasheada)
	Password string `json:"-" binding:"required"`

	// @Description Fecha de creación del usuario
	CreatedAt time.Time `json:"created_at"`

	// @Description Fecha de última actualización del usuario
	UpdatedAt time.Time `json:"updated_at"`

	// @Description Fecha de eliminación del usuario (soft delete)
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
