package config

import (
	"fmt"
	"os"

	"go-hexagonal-template/internal/modules/user/domain/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	SSLMode     string
	LogLevel    string
	AutoMigrate bool
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		User:        os.Getenv("DB_USER"),
		Password:    os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		SSLMode:     os.Getenv("DB_SSL_MODE"),
		LogLevel:    os.Getenv("GORM_LOG_LEVEL"),
		AutoMigrate: os.Getenv("GORM_AUTO_MIGRATE") == "true",
	}
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

func (c *DatabaseConfig) getLogLevel() logger.LogLevel {
	switch c.LogLevel {
	case "debug":
		return logger.Info
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Silent
	}
}

func (c *DatabaseConfig) Connect() (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(c.getLogLevel()),
	}

	db, err := gorm.Open(postgres.Open(c.GetDSN()), gormConfig)
	if err != nil {
		return nil, err
	}

	// Auto-migrar las tablas si está configurado
	if c.AutoMigrate {
		if err := db.AutoMigrate(&model.User{}); err != nil {
			return nil, fmt.Errorf("error auto-migrando tablas: %v", err)
		}
	}

	return db, nil
}
