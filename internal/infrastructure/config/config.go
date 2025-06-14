package config

import (
	"os"

	"gorm.io/gorm"
)

type Config struct {
	Environment string
	Port        string
	Database    *DatabaseConfig
	DB          *gorm.DB
}

func NewConfig() (*Config, error) {
	config := &Config{
		Environment: os.Getenv("ENV"),
		Port:        os.Getenv("PORT"),
		Database:    NewDatabaseConfig(),
	}

	// Conectar a la base de datos
	db, err := config.Database.Connect()
	if err != nil {
		return nil, err
	}
	config.DB = db

	return config, nil
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}
