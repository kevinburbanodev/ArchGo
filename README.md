# Go Hexagonal Template

Este proyecto es una plantilla para crear aplicaciones en Go utilizando la arquitectura hexagonal (también conocida como puertos y adaptadores).

## Estructura del Proyecto

```
.
├── cmd/
│   └── server/         # Punto de entrada de la aplicación
├── internal/
│   ├── handlers/       # Manejadores HTTP
│   ├── infrastructure/ # Implementaciones concretas
│   ├── middleware/     # Middleware de la aplicación
│   └── modules/        # Módulos de la aplicación
│       └── user/       # Módulo de usuario
│           ├── application/    # Casos de uso
│           ├── domain/         # Modelos y puertos
│           └── infrastructure/ # Implementaciones
└── tests/              # Tests de la aplicación
```

## Requisitos

- Go 1.21 o superior
- PostgreSQL
- Make (opcional, para usar los comandos make)

## Configuración

1. Clona el repositorio
2. Copia el archivo `.env.example` a `.env` y configura las variables de entorno
3. Ejecuta `go mod download` para instalar las dependencias
4. Ejecuta `go run cmd/server/main.go` para iniciar el servidor

## Variables de Entorno

### Configuración de la Base de Datos
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_hexagonal
DB_SSL_MODE=disable
```

### Configuración de GORM
```
GORM_LOG_LEVEL=debug    # Niveles: debug, info, warn, error, silent
GORM_AUTO_MIGRATE=true  # true/false
```

### Configuración del Servidor
```
PORT=3000
JWT_SECRET=your-secret-key
ENV=development
```

### Explicación de Variables Específicas

#### DB_SSL_MODE
Configura el modo SSL para la conexión a PostgreSQL:
- `disable`: No usa SSL (recomendado para desarrollo local)
- `require`: Requiere conexión SSL
- `verify-ca`: Verifica que el certificado del servidor esté firmado por una CA confiable
- `verify-full`: Verifica el certificado y el nombre del host (más seguro)

#### GORM_LOG_LEVEL
Controla el nivel de logging de GORM:
- `debug`: Muestra todas las consultas SQL y detalles
- `info`: Muestra información general
- `warn`: Solo muestra advertencias
- `error`: Solo muestra errores
- `silent`: No muestra logs

#### GORM_AUTO_MIGRATE
Habilita/deshabilita la migración automática de la base de datos:
- `true`: GORM creará/actualizará las tablas automáticamente
- `false`: No se realizarán migraciones automáticas

## Documentación de la API (Swagger)

La documentación de la API está disponible a través de Swagger UI. Para acceder:

1. Inicia el servidor
2. Abre en tu navegador: `http://localhost:3000/swagger/index.html`

En Swagger UI podrás:
- Ver toda la documentación de la API
- Probar los endpoints directamente
- Ver los modelos de datos
- Ver los códigos de respuesta posibles
- Probar la autenticación con JWT

## Endpoints

### Usuarios

#### Crear Usuario
```bash
curl --location 'http://localhost:3000/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "password123"
}'
```

#### Login
```bash
curl --location 'http://localhost:3000/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "test@example.com",
    "password": "password123"
}'
```

#### Obtener Usuario (requiere autenticación)
```bash
curl --location 'http://localhost:3000/api/users/1' \
--header 'Authorization: Bearer <token>'
```

## Tests

Para ejecutar los tests:

```bash
go test ./...
```

## Comandos Make

- `make run`: Inicia el servidor
- `make test`: Ejecuta los tests
- `make build`: Compila la aplicación
- `make clean`: Limpia los archivos compilados

## Seguridad

### Rate Limiting

La API implementa rate limiting para prevenir ataques de fuerza bruta y DoS. Las características incluyen:

- **Límites por IP**: 100 peticiones por minuto
- **Headers de Respuesta**:
  - `X-RateLimit-Limit`: Límite total de peticiones (100)
  - `X-RateLimit-Remaining`: Peticiones restantes
  - `X-RateLimit-Reset`: Tiempo hasta que se reinicie el contador

#### Ejemplo de Respuesta con Rate Limit
```http
HTTP/1.1 200 OK
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1710288000
```

#### Ejemplo de Respuesta al Exceder el Límite
```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json

{
    "error": "Has excedido el límite de peticiones. Por favor, espera un momento."
}
```

### Otras Medidas de Seguridad

- **Autenticación JWT**: Todas las rutas protegidas requieren un token JWT válido
- **SSL/TLS**: Configurable a través de `DB_SSL_MODE` para conexiones seguras a la base de datos
- **Validación de Entrada**: Todos los datos de entrada son validados antes de ser procesados
- **Headers de Seguridad**: La API incluye headers de seguridad estándar



## Dockerización

### Requisitos
- Docker
- Docker Compose

### Estructura de Docker

#### Dockerfile
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 3000
CMD ["./main"]
```

#### docker-compose.yml
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=go_hexagonal
      - DB_SSL_MODE=disable
      - GORM_LOG_LEVEL=debug
      - GORM_AUTO_MIGRATE=true
      - JWT_SECRET=your-secret-key
    depends_on:
      - postgres
    networks:
      - app-network

  postgres:
    image: postgres:15-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=go_hexagonal
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - app-network

volumes:
  postgres_data:

networks:
  app-network:
    driver: bridge
```

### Características de la Dockerización

1. **Multi-stage Build**
   - Reduce el tamaño final de la imagen
   - Separa el proceso de compilación del entorno de ejecución
   - Incluye solo los archivos necesarios

2. **Servicios**
   - **app**: Servicio principal de la aplicación
     - Compilado desde el Dockerfile
     - Expone el puerto 3000
     - Configurado con variables de entorno
   - **postgres**: Base de datos PostgreSQL
     - Usa la imagen oficial de PostgreSQL
     - Persiste datos mediante volumen
     - Configurado con credenciales básicas

3. **Redes**
   - Red dedicada `app-network`
   - Aislamiento de servicios
   - Comunicación segura entre contenedores

4. **Volúmenes**
   - `postgres_data`: Persiste los datos de PostgreSQL
   - Evita pérdida de datos al reiniciar contenedores

### Comandos Docker

1. **Construir y ejecutar**
```bash
# Construir y levantar todos los servicios
docker-compose up --build

# Ejecutar en segundo plano
docker-compose up -d
```

2. **Gestión de contenedores**
```bash
# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down

# Reiniciar servicios
docker-compose restart
```

3. **Mantenimiento**
```bash
# Limpiar contenedores no utilizados
docker system prune

# Ver uso de recursos
docker stats
```

### Variables de Entorno

Las variables de entorno se pueden configurar de dos formas:

1. **En docker-compose.yml**
```yaml
environment:
  - DB_HOST=postgres
  - DB_PORT=5432
  # ... otras variables
```

2. **En archivo .env**
```env
DB_HOST=postgres
DB_PORT=5432
# ... otras variables
```

### Consideraciones de Seguridad

1. **Credenciales**
   - No incluir credenciales sensibles en el código
   - Usar variables de entorno o secrets de Docker
   - Cambiar las credenciales por defecto en producción

2. **Redes**
   - Usar redes Docker para aislar servicios
   - Exponer solo los puertos necesarios
   - Configurar políticas de red restrictivas

3. **Volúmenes**
   - Usar volúmenes nombrados para persistencia
   - Configurar permisos adecuados
   - Hacer backup regular de los datos 

   ## Licencia
    MIT 