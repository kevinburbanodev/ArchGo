# Go Hexagonal Template

[English](README.md) | Español

Una plantilla robusta para construir aplicaciones en Go utilizando arquitectura hexagonal (también conocida como puertos y adaptadores).

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
└── tests/              # Pruebas de la aplicación
```

## Requisitos

- Go 1.21 o superior
- PostgreSQL
- Make (opcional, para comandos make)

## Configuración

1. Clonar el repositorio
2. Copiar `.env.example` a `.env` y configurar las variables de entorno
3. Ejecutar `go mod download` para instalar las dependencias
4. Ejecutar `go run cmd/server/main.go` para iniciar el servidor

## Variables de Entorno

### Configuración de Base de Datos
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
Configura el modo SSL para la conexión PostgreSQL:
- `disable`: Sin SSL (recomendado para desarrollo local)
- `require`: Requiere conexión SSL
- `verify-ca`: Verifica que el certificado del servidor esté firmado por una CA confiable
- `verify-full`: Verifica certificado y nombre de host (más seguro)

#### GORM_LOG_LEVEL
Controla el nivel de registro de GORM:
- `debug`: Muestra todas las consultas SQL y detalles
- `info`: Muestra información general
- `warn`: Muestra solo advertencias
- `error`: Muestra solo errores
- `silent`: Sin registros

#### GORM_AUTO_MIGRATE
Habilita/deshabilita la migración automática de la base de datos:
- `true`: GORM creará/actualizará tablas automáticamente
- `false`: No se realizarán migraciones automáticas

## Documentación de la API (Swagger)

La documentación de la API está disponible a través de Swagger UI. Para acceder:

1. Iniciar el servidor
2. Abrir en el navegador: `http://localhost:3000/swagger/index.html`

En Swagger UI puedes:
- Ver la documentación completa de la API
- Probar endpoints directamente
- Ver modelos de datos
- Ver códigos de respuesta
- Probar autenticación JWT

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

#### Iniciar Sesión
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

## Seguridad

### Límite de Tasa (Rate Limiting)

La API implementa límites de tasa para prevenir ataques de fuerza bruta y DoS. Características incluyen:

- **Límites basados en IP**: 100 solicitudes por minuto
- **Encabezados de Respuesta**:
  - `X-RateLimit-Limit`: Límite total de solicitudes (100)
  - `X-RateLimit-Remaining`: Solicitudes restantes
  - `X-RateLimit-Reset`: Tiempo hasta el reinicio del contador

#### Ejemplo de Respuesta con Límite de Tasa
```http
HTTP/1.1 200 OK
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1710288000
```

#### Ejemplo de Respuesta Cuando se Excede el Límite
```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json

{
    "error": "Has excedido el límite de solicitudes. Por favor, espera un momento."
}
```

### Otras Medidas de Seguridad

- **Autenticación JWT**: Todas las rutas protegidas requieren un token JWT válido
- **SSL/TLS**: Configurable a través de `DB_SSL_MODE` para conexiones seguras a la base de datos
- **Validación de Entrada**: Todos los datos de entrada son validados antes del procesamiento
- **Encabezados de Seguridad**: La API incluye encabezados de seguridad estándar

## Docker

### Requisitos
- Docker
- Docker Compose

### Estructura Docker

#### Dockerfile
```dockerfile
# Etapa de construcción
FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Etapa final
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

### Características de Docker

1. **Construcción Multi-etapa**
   - Reduce el tamaño de la imagen final
   - Separa el proceso de construcción del entorno de ejecución
   - Incluye solo los archivos necesarios

2. **Servicios**
   - **app**: Servicio principal de la aplicación
     - Construido desde Dockerfile
     - Expone el puerto 3000
     - Configurado con variables de entorno
   - **postgres**: Base de datos PostgreSQL
     - Usa la imagen oficial de PostgreSQL
     - Persiste datos a través de volumen
     - Configurado con credenciales básicas

3. **Redes**
   - Red dedicada `app-network`
   - Aislamiento de servicios
   - Comunicación segura entre contenedores

4. **Volúmenes**
   - `postgres_data`: Persiste datos de PostgreSQL
   - Previene pérdida de datos al reiniciar contenedores

### Comandos Docker

1. **Construir y Ejecutar**
```bash
# Construir e iniciar todos los servicios
docker-compose up --build

# Ejecutar en segundo plano
docker-compose up -d
```

2. **Gestión de Contenedores**
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

Las variables de entorno se pueden configurar de dos maneras:

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
   - Usar variables de entorno o secretos de Docker
   - Cambiar credenciales por defecto en producción

2. **Redes**
   - Usar redes Docker para aislar servicios
   - Exponer solo puertos necesarios
   - Configurar políticas de red restrictivas

3. **Volúmenes**
   - Usar volúmenes nombrados para persistencia
   - Configurar permisos apropiados
   - Respaldo regular de datos

## Contacto

### Autor
**Kevin Fernando Burbano Aragón**  
Ingeniero en Sistemas y Desarrollador de Software Senior con amplia experiencia en desarrollo de software.

### Información de Contacto
- **Email**: [burbanokevin1997@gmail.com](mailto:burbanokevin1997@gmail.com)
- **GitHub**: [@kevinburbanodev](https://github.com/kevinburbanodev)
- **LinkedIn**: [Kevin Fernando Burbano Aragón](https://www.linkedin.com/in/kevin-fernando-burbano-arag%C3%B3n-78b3871a0/)

### Experiencia
- Ingeniero en Sistemas con experiencia senior en desarrollo de software
- Especializado en arquitecturas limpias y patrones de diseño
- Experto en desarrollo de APIs y microservicios
- Amplio conocimiento en múltiples tecnologías y frameworks

## Licencia
MIT 