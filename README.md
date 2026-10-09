# Biblioteca API

API REST en Go para un sistema de biblioteca, pensada como backend del proyecto de curso de **Front-End Development con React**.

Permite: registro/login con JWT, catálogo de libros (CRUD para administradores), y préstamos de libros para estudiantes.

---

## Stack

- **Go 1.22** + **Gin** (router)
- **sqlx** + **MySQL 8** (persistencia)
- **JWT** (autenticación) + **bcrypt** (hash de contraseñas)
- **Swagger** (documentación interactiva de la API)
- **Docker Compose** (MySQL + API)

---

## Requisitos previos

- [Docker](https://www.docker.com/) y Docker Compose instalados
- Go 1.22+ instalado localmente **solo si vas a generar la documentación Swagger o correr el proyecto fuera de Docker**

---

## Puesta en marcha (primera vez)

### 1. Configurar variables de entorno

```bash
cp .env.example .env
```

Puedes dejar los valores por defecto para desarrollo local; solo cambia `JWT_SECRET` si vas a exponer esto más allá de la clase.

### 2. Generar go.sum y la documentación Swagger

Como el proyecto se generó sin acceso a internet, falta un paso único antes de levantar todo:

```bash
# Descarga las dependencias y genera go.sum
go mod tidy

# Instala la herramienta swag (una sola vez)
go install github.com/swaggo/swag/cmd/swag@latest

# Genera la documentación a partir de los comentarios @Summary/@Router del código
swag init -g cmd/api/main.go -o docs
```

Esto crea `docs/docs.go`, `docs/swagger.json` y `docs/swagger.yaml`.

### 3. Levantar todo con Docker

```bash
docker-compose up --build
```

Esto levanta:
- **MySQL** en el puerto `3306`, con el esquema y datos de ejemplo cargados automáticamente (ver `migrations/001_init.sql`)
- **La API** en el puerto `8080`

### 4. Explorar la API

Abre en tu navegador:

```
http://localhost:8080/swagger/index.html
```

Ahí puedes ver todos los endpoints, probar peticiones directamente, y ver los modelos de datos — sin necesitar Postman.

---

## Endpoints principales

| Método | Ruta                  | Descripción                          | Auth requerido |
|--------|-----------------------|---------------------------------------|:--------------:|
| POST   | `/api/auth/register`  | Crear cuenta                          | No             |
| POST   | `/api/auth/login`     | Iniciar sesión (devuelve JWT)         | No             |
| GET    | `/api/auth/me`        | Perfil del usuario autenticado        | Sí             |
| GET    | `/api/books`          | Listar libros (filtros: category, search) | Sí         |
| GET    | `/api/books/:id`      | Detalle de un libro                   | Sí             |
| POST   | `/api/books`          | Crear libro                           | Sí (admin)     |
| PUT    | `/api/books/:id`      | Actualizar libro                      | Sí (admin)     |
| DELETE | `/api/books/:id`      | Eliminar libro                        | Sí (admin)     |
| POST   | `/api/loans`          | Pedir prestado un libro               | Sí             |
| GET    | `/api/loans/me`       | Mis préstamos                         | Sí             |
| PUT    | `/api/loans/:id/return` | Devolver un libro                   | Sí             |
| GET    | `/api/loans`          | Ver todos los préstamos               | Sí (admin)     |

### Cómo autenticarse desde React

1. `POST /api/auth/register` o `/api/auth/login` → obtienes un `token`
2. Guarda el token (ej. en memoria o `localStorage`, según lo que enseñes)
3. Envíalo en cada petición protegida:

```js
fetch("http://localhost:8080/api/books", {
  headers: {
    Authorization: `Bearer ${token}`
  }
})
```

---

## Crear un usuario administrador

Por defecto, todo usuario que se registra vía `/api/auth/register` queda con rol `student`. Para crear un admin (y poder gestionar el catálogo de libros), conéctate directamente a MySQL y actualiza el rol:

```sql
UPDATE users SET role = 'admin' WHERE email = 'tu_correo@ejemplo.com';
```

---

## Estructura del proyecto

```
cmd/api/main.go         → punto de entrada
internal/
  config/                → carga de variables de entorno
  db/                     → conexión a MySQL
  models/                 → structs de dominio (User, Book, Loan)
  repository/             → acceso a datos (SQL)
  service/                → lógica de negocio (auth, préstamos, etc.)
  handler/                → controladores HTTP + rutas
  middleware/             → JWT auth, permisos de admin
migrations/               → esquema SQL + datos de ejemplo
docs/                     → documentación Swagger (generada, no editar a mano)
```

---

## Reiniciar la base de datos desde cero

Si necesitas volver a los datos de ejemplo originales:

```bash
docker-compose down -v
docker-compose up --build
```

El flag `-v` borra el volumen de MySQL, por lo que las migraciones y el seed se vuelven a ejecutar.
