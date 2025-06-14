package config

import (
	"fmt"
	"os"

	"go-hexagonal-template/internal/modules/user/domain/model"

	"gorm.io/gorm"

	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Environment string
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	SSLMode     string
	LogLevel    string
	AutoMigrate bool
	RefreshDB   bool
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Environment: os.Getenv("ENV"),
		Host:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		User:        os.Getenv("DB_USER"),
		Password:    os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		SSLMode:     os.Getenv("DB_SSL_MODE"),
		LogLevel:    os.Getenv("GORM_LOG_LEVEL"),
		AutoMigrate: os.Getenv("GORM_AUTO_MIGRATE") == "true",
		RefreshDB:   os.Getenv("GORM_REFRESH_DB") == "true",
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

	// Solo en desarrollo y si está configurado, hacer refresh de las tablas
	if c.Environment == "dev" && c.RefreshDB {
		if err := db.Migrator().DropTable(&model.User{}); err != nil {
			return nil, fmt.Errorf("error eliminando tablas: %v", err)
		}
	}

	// Auto-migrar las tablas si está configurado y si no estamos en producción
	if c.AutoMigrate && c.Environment != "prod" {
		if err := db.AutoMigrate(&model.User{}); err != nil {
			return nil, fmt.Errorf("error auto-migrando tablas: %v", err)
		}
	}

	return db, nil
}
