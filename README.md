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

## Licencia

MIT 