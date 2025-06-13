package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimiterMiddleware crea un middleware para limitar las peticiones
func RateLimiterMiddleware() gin.HandlerFunc {
	// Crear un rate limiter que permita 100 peticiones por minuto
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}

	// Usar memoria como almacenamiento para el rate limiter
	store := memory.NewStore()

	// Crear el limiter
	limiterInstance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Obtener el IP del cliente
		ip := c.ClientIP()

		// Obtener el contexto del limiter
		context, err := limiterInstance.Get(c.Request.Context(), ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al verificar el límite de peticiones",
			})
			c.Abort()
			return
		}

		// Agregar headers con información del rate limit
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", string(context.Remaining))
		c.Header("X-RateLimit-Reset", string(context.Reset))

		// Si se excedió el límite, retornar error
		if context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Has excedido el límite de peticiones. Por favor, espera un momento.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
