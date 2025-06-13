package main

import (
	"fmt"
	"go-hexagonal-template/internal/handlers"
	"go-hexagonal-template/internal/infrastructure/config"
	"go-hexagonal-template/internal/middleware"
	"go-hexagonal-template/internal/modules/user/infrastructure/persistence"
	"log"

	_ "go-hexagonal-template/docs" // Esto es importante para la documentación Swagger

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Go Hexagonal Template API
// @version         1.0
// @description     API Template con arquitectura hexagonal en Go
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load("../../.env"); err != nil {
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

	// Configurar Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Obtener el puerto de la variable de entorno o usar 3000 por defecto
	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	// Iniciar el servidor en el puerto configurado
	r.Run(fmt.Sprintf(":%s", port))
}
