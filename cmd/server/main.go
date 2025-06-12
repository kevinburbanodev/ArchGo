package main

import (
	"fmt"
	"go-hexagonal-template/internal/handlers"
	"go-hexagonal-template/internal/infrastructure/config"
	"go-hexagonal-template/internal/middleware"
	"go-hexagonal-template/internal/modules/user/infrastructure/persistence"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Cargar configuración
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Configurar el modo de Gin según el entorno
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Crear una instancia de Gin
	r := gin.Default()

	// Inicializar handlers
	healthHandler := handlers.NewHealthHandler()
	userRepo := persistence.NewUserRepositoryImpl(cfg.DB)
	userHandler := handlers.NewUserHandler(userRepo)

	// Definir rutas públicas
	r.GET("/healthy", healthHandler.HealthCheck)
	r.POST("/users", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)

	// Definir rutas protegidas
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users/:id", userHandler.GetUser)
	}

	// Obtener el puerto de la variable de entorno o usar 3000 por defecto
	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	// Iniciar el servidor en el puerto configurado
	r.Run(fmt.Sprintf(":%s", port))
}
