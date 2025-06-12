# Go Hexagonal Template

Este es un template para proyectos Go que implementa una arquitectura hexagonal (también conocida como puertos y adaptadores) con vertical slicing, utilizando Gin como framework web, GORM como ORM y PostgreSQL como base de datos por defecto.

## Características Principales

- **Arquitectura Hexagonal**: Implementación limpia de la arquitectura hexagonal para una clara separación de responsabilidades.
- **Vertical Slicing**: Organización del código por funcionalidades (módulos) en lugar de capas técnicas.
- **Autenticación JWT**: Sistema de autenticación completo con tokens JWT.
- **Middleware de Autenticación**: Middleware para proteger rutas que requieren autenticación.
- **GORM + PostgreSQL**: ORM con soporte para PostgreSQL por defecto.
- **Testing**: Tests unitarios y de integración incluidos.
- **Configuración por Variables de Entorno**: Uso de `.env` para configuración.

## Estructura del Proyecto

```
.
├── cmd/
│   └── server/          # Punto de entrada de la aplicación
├── internal/
│   ├── infrastructure/  # Implementaciones concretas (adaptadores)
│   │   ├── auth/       # Implementación de autenticación JWT
│   │   ├── config/     # Configuración de la aplicación
│   │   └── persistence/# Implementaciones de repositorios
│   ├── middleware/     # Middlewares de la aplicación
│   └── modules/        # Módulos de la aplicación (vertical slicing)
│       └── user/       # Módulo de usuario
│           ├── application/  # Casos de uso
│           ├── domain/      # Entidades y puertos
│           └── infrastructure/# Implementaciones específicas
└── tests/              # Tests de integración
```

## Módulo de Usuario

El template incluye un módulo de usuario completo con las siguientes funcionalidades:

- Creación de usuarios
- Autenticación (login)
- Obtención de usuario por ID
- Validación de datos
- Hashing de contraseñas
- Manejo de errores

## API Endpoints

### Rutas Públicas
- `POST /users` - Crear usuario
- `POST /login` - Autenticación de usuario

### Rutas Protegidas (requieren JWT)
- `GET /api/users/:id` - Obtener usuario por ID

## Configuración

El proyecto utiliza variables de entorno para su configuración. Crea un archivo `.env` en la raíz del proyecto con las siguientes variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSL_MODE=disable
JWT_SECRET_KEY=your_secret_key
PORT=3000
ENV=development
```

## Cambiar la Base de Datos

El template está configurado para usar PostgreSQL por defecto, pero puede ser fácilmente adaptado para usar otras bases de datos soportadas por GORM.

### Usando SQLite

1. Modifica el archivo `go.mod` para incluir el driver de SQLite:
```go
require (
    gorm.io/driver/sqlite v1.5.5
)
```

2. Modifica el archivo `internal/infrastructure/config/database.go`:
```go
import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return db, nil
}
```

### Usando MySQL

1. Modifica el archivo `go.mod` para incluir el driver de MySQL:
```go
require (
    gorm.io/driver/mysql v1.5.2
)
```

2. Modifica el archivo `internal/infrastructure/config/database.go`:
```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return db, nil
}
```

## Ejecutar el Proyecto

1. Clona el repositorio
2. Instala las dependencias:
```bash
go mod download
```

3. Configura las variables de entorno en el archivo `.env`

4. Ejecuta la aplicación:
```bash
go run cmd/server/main.go
```

## Ejecutar Tests

```bash
go test ./...
```

## Estructura de la Arquitectura Hexagonal

### Dominio (Domain)
- Contiene las entidades y reglas de negocio
- Define los puertos (interfaces) para la comunicación con el exterior
- No tiene dependencias externas

### Aplicación (Application)
- Implementa los casos de uso
- Orquesta las operaciones entre el dominio y los adaptadores
- Depende solo del dominio

### Infraestructura (Infrastructure)
- Implementa los adaptadores para bases de datos, servicios externos, etc.
- Implementa las interfaces definidas en el dominio
- Contiene la configuración de la aplicación

## Middleware de Autenticación

El middleware de autenticación verifica:
1. Presencia del token en el header Authorization
2. Formato correcto del token (Bearer)
3. Validez del token JWT
4. Almacena la información del usuario en el contexto

## Contribuir

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles. 