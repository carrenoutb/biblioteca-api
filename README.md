# Biblioteca API

API REST en Go para un sistema de biblioteca, **proyecto del curso de Arquitectura de Software**. El repositorio implementa una **arquitectura por capas** (`handler` → `service` → `repository`) sobre MySQL, con autenticación JWT, documentación Swagger y está **dockerizado** (MySQL + API) para ejecutarse en cualquier máquina con un solo comando.

Permite: registro/login con JWT, catálogo de libros (CRUD para administradores), y préstamos de libros para estudiantes.

---

## Stack

- **Go 1.22** + **Gin** (router)
- **sqlx** + **MySQL 8** (persistencia)
- **JWT** (autenticación) + **bcrypt** (hash de contraseñas)
- **Swagger** (documentación interactiva de la API)
- **Docker Compose** (MySQL + API)

---

# 🎓 Guía de inicio rápido (para estudiantes)

Esta guía te lleva de **cero a la API corriendo en tu máquina**, paso a paso y explicado. Está pensada para que la sigas aunque **nunca hayas usado Docker** y aunque **no tengas Go ni MySQL instalados** (todo vive dentro de los contenedores).

## 0. ¿Qué necesitas?

- Una computadora con **Windows, macOS o Linux**.
- **Docker** instalado (la sección 1 te dice cómo y cómo comprobar que funcionó).
- Conexión a **internet** (solo la primera vez: Docker descarga las imágenes).

> ✅ **Lo que NO necesitas:** ni Go, ni MySQL, ni ninguna otra herramienta. La API y la base de datos corren dentro de contenedores.

## 1. Instalar Docker (según tu sistema operativo)

### Windows / macOS — Docker Desktop

Instala **[Docker Desktop](https://www.docker.com/products/docker-desktop/)** (el instalador "todo-en-uno": trae el motor de Docker, Compose y una interfaz gráfica). Ábrelo y espera a que el icono indique que está corriendo.

### Linux — Docker Engine + plugin Compose

Con el gestor de paquetes de tu distribución, por ejemplo en Debian/Ubuntu:

```bash
sudo apt update
sudo apt install -y docker.io docker-compose-plugin
sudo systemctl enable --now docker
```

### Verifica la instalación (obligatorio para continuar)

```bash
docker --version
docker compose version
```

Debes ver algo como (las versiones exactas pueden variar):

```
Docker version 29.5.2, build ...
Docker Compose version 5.5.1
```

---

## 2. Las piezas del proyecto (para entender antes de ejecutar)

| Archivo | ¿Qué es? |
|---|---|
| `Dockerfile` | La **receta** para construir la imagen de la API: toma una imagen base de Go, descarga dependencias, compila el binario y deja lista la aplicación (imagen *multi-stage*: un paso para compilar, otro ligero para ejecutar). |
| `docker-compose.yml` | El **director de orquesta**: define que la app son 2 contenedores — `mysql` (base de datos MySQL 8) y `api` (nuestra API) —, los conecta en una red propia, expone puertos y lanza las migraciones automáticamente. |
| `.env.example` → `.env` | Las **variables de entorno**: configuración (usuario/clave de la BD, secreto JWT…) que no debe ir "quemada" en el código. Se copia a `.env` (ver paso 3.2). |
| `migrations/` | Los **scripts SQL** que crean el esquema y cargan datos de ejemplo; MySQL los ejecuta solo, la primera vez que arranca con la base vacía. |
| `docs/` | La documentación Swagger ya generada (no se edita a mano). |

**Resumen:** `Dockerfile` crea la **imagen** de la API → `docker compose` ejecuta esa imagen junto a MySQL como **contenedores** → las **imágenes base** (`golang`, `mysql`) se descargan del **registro** (Docker Hub).

---

## 3. Paso a paso

### 3.1 — Obtener el código

Clona el repositorio (o copia la carpeta del proyecto) y entra en ella:

```bash
git clone <url-del-repositorio> biblioteca-api
cd biblioteca-api
```

### 3.2 — Crear tu archivo `.env`

> ⚠️ **Importante:** el repositorio **no contiene un archivo `.env`**. En su lugar trae una plantilla llamada **`.env.example`** con las mismas variables y valores de ejemplo. El `.env` se crea en **tu máquina** copiando esa plantilla:

```bash
cp .env.example .env
```

**¿Por qué el `.env` no está en el repositorio?** Porque contiene datos sensibles (usuario/contraseña de la base de datos y la clave de firma de los tokens `JWT_SECRET`). Si se subiera al repo (público en GitHub), cualquiera podría ver esas credenciales. Por eso el `.gitignore` lo excluye y en el repo solo viaja `.env.example`.

**¿Qué contiene?** Cada variable la lee la API al arrancar:

| Variable | Qué es | Valor de ejemplo |
|---|---|---|
| `DB_HOST` | Dónde vive la base de datos (ver nota 👇) | `mysql` |
| `DB_PORT` | Puerto de MySQL | `3306` |
| `DB_USER` | Usuario de la base de datos | `biblioteca_user` |
| `DB_PASSWORD` | Contraseña de la base de datos | `biblioteca_pass` |
| `DB_NAME` | Nombre de la base de datos | `biblioteca` |
| `JWT_SECRET` | Clave con la que se firman los tokens (¡cámbiala si vas más allá de la clase!) | `TOKEN_BIBLIOTECA` |
| `JWT_EXPIRATION_HOURS` | Horas de validez de un token | `24` |
| `APP_PORT` | Puerto donde escucha la API | `8080` |

> **Nota sobre `DB_HOST`:** `mysql` **no es una dirección de internet** — es el *nombre del servicio* de la base de datos dentro de la red de Docker (lo inyecta `docker-compose.yml`). No lo cambies a menos que sepas lo que haces.

> ⚠️ Si **no** creas el `.env`, `docker compose up --build` fallará (Compose está configurado con `env_file: .env` y necesita el archivo). Este paso es obligatorio.

### 3.3 — Levantar todo con Docker

```bash
docker compose up --build
```

**Qué está pasando internamente, en orden:**

1. Docker lee el `Dockerfile` y **construye la imagen de la API** (descarga `golang:...-alpine` desde Docker Hub, ejecuta `go mod download` y `go build`).
2. Descarga la imagen de **MySQL 8**.
3. Arranca **MySQL primero** y espera a que esté sano (el *healthcheck* le "hace ping" cada 5 segundos).
4. Con MySQL sano, ejecuta **automáticamente** las migraciones (`migrations/001_init.sql`) y el seed en la base vacía.
5. Arranca la **API** en el puerto `8080`.

> 💡 La primera ejecución tarda unos minutos (descarga imágenes). Para correrlo en segundo plano y seguir usando la terminal: `docker compose up -d`. Para detener sin borrar: `Ctrl + C`.

### 3.4 — Verificar que funciona

**a) Estado de los servicios** (en otra terminal):

```bash
docker compose ps
```

Ambos contenedores deben estar `Up` y MySQL `(healthy)`.

**b) La documentación interactiva (Swagger):** abre en el navegador

```
http://localhost:8080/swagger/index.html
```

**c) Prueba real del flujo completo con REST Client** (`register` → `login` → `books`):

En lugar de usar Postman o `curl`, esta guía usa la extensión **REST Client** de VS Code (gratuita y sin salir del editor). Sigue estos pasos:

**Paso 1 — Instalar la extensión**
1. Abre VS Code y ve a la vista de extensiones: `Cmd/Ctrl + Shift + X`.
2. Busca **"REST Client"** (autor: *Huachao Mao*) e instálala.
3. Verás el panel de la extensión: cualquier archivo `.http` (o `.rest`) se detecta automáticamente y muestra un botón **`Send Request`** sobre cada petición.

**Paso 2 — Crear el archivo `.http`**

Crea un archivo nuevo en la raíz del proyecto llamado `pruebas-api.http`. En el **paso 3** vas a escribir cada petición **una por una** (copia cada *request* en tu archivo y pruébala); al final del paso 3 tienes el **archivo completo** por si prefieres copiarlo todo de una vez.

**Paso 3 — Escribir y probar cada petición (request → respuesta)**

### 3.1 — Estado de la API (Health)

Escribe esto en `pruebas-api.http`:

```http
### 1. Estado de la API
GET http://localhost:8080/health
```

Haz clic en **`Send Request`** → la **respuesta esperada** es:

```json
{"status":"ok"}
```

**¿Qué comprueba?** Que la API esté viva y respondiendo. Si ves esto, el contenedor `api` está corriendo bien.

### 3.2 — Registrar un usuario nuevo

Agrega al final del archivo:

```http
### 2. Registrar un usuario nuevo
# @name register
POST http://localhost:8080/api/auth/register
Content-Type: application/json

{
  "name": "estudiante-taller",
  "email": "taller-practica@example.com",
  "password": "pass@123"
}
```

**Send Request** → la **respuesta esperada** es (el `id` puede variar):

```json
{"token":"eyJhbGciOi...","user":{"id":1,"name":"estudiante-taller","email":"taller-practica@example.com","role":"student","created_at":"2026-10-09T..."}}
```

**¿Qué pasó?** Se creó el usuario con rol `student` y la API devolvió su **token** automáticamente. Fíjate en `# @name register`: le pone un nombre a esta petición (lo usaremos como referencia).

### 3.3 — Iniciar sesión (obtener el token)

Agrega:

```http
### 3. Iniciar sesión (obtener el token)
# @name login
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "email": "taller-practica@example.com",
  "password": "pass@123"
}
```

**Send Request** → **respuesta esperada**:

```json
{"token":"eyJhbGciOi...","user":{"id":1,"name":"estudiante-taller","email":"taller-practica@example.com","role":"student",...}}
```

**¿Qué pasó?** Con `# @name login` guardamos la respuesta de esta petición con nombre `login`. Así, más adelante podremos decirle a otra petición: *"usa el `token` de la respuesta de `login`"* — sin copiar y pegar el token a mano.

### 3.4 — Listar libros (con el token del login)

Agrega la última petición:

```http
### 4. Listar libros (usa el token del login automáticamente)
GET http://localhost:8080/api/books
Authorization: Bearer {{login.response.body.token}}
```

> ⚠️ Antes de ejecutar esta petición, ejecuta la **3.3** (login) para que la variable `{{login.response.body.token}}` exista.

**Send Request** → **respuesta esperada** (los libros del seed):

```json
[
  {"id":1,"title":"Cien años de soledad","author":"Gabriel García Márquez","isbn":"9780307474728","category":"Novela","total_copies":3,"available_copies":3},
  {"id":2,"title":"Clean Code","author":"Robert C. Martin","isbn":"9780132350884","category":"Tecnología","total_copies":2,"available_copies":2}
]
```

**¿Qué pasó?** La expresión `{{login.response.body.token}}` tomó automáticamente el campo `token` de la **respuesta** de la petición `login` (3.3) y lo puso en la cabecera `Authorization`. REST Client hace el *encadenado* por ti.

---

### 🧩 El archivo completo (`pruebas-api.http`)

Si prefieres copiar todo de una vez (o verificar que tu archivo está bien), esta es la versión completa de lo que acabamos de escribir:

```http
# ============================================================
# Pruebas de la API Biblioteca — REST Client (VS Code)
# Requisito: la API corriendo (docker compose up --build)
# ============================================================

# Dirección base de la API
@baseUrl = http://localhost:8080

### 1. Estado de la API
GET {{baseUrl}}/health

### 2. Registrar un usuario nuevo
# @name register
POST {{baseUrl}}/api/auth/register
Content-Type: application/json

{
  "name": "estudiante-taller",
  "email": "taller-practica@example.com",
  "password": "pass@123"
}

### 3. Iniciar sesión (obtener el token)
# @name login
POST {{baseUrl}}/api/auth/login
Content-Type: application/json

{
  "email": "taller-practica@example.com",
  "password": "pass@123"
}

### 4. Listar libros (usa el token del login automáticamente)
GET {{baseUrl}}/api/books
Authorization: Bearer {{login.response.body.token}}
```

**Resumen de la sintaxis** que acabas de usar:

| Elemento | Qué es |
|---|---|
| `###` | Separa una petición de la siguiente (obligatorio entre peticiones) |
| `#` | Comentario (no se envía) |
| `@baseUrl = ...` | Define una **variable** reutilizable con `{{baseUrl}}` |
| `# @name login` | **Nombra** la petición para usar su respuesta más adelante |
| (línea en blanco) | Separa las cabeceras del **cuerpo** (obligatorio en JSON) |
| `{{login.response.body.token}}` | Toma el campo `token` de la **respuesta** de la petición llamada `login` |

> 💡 El repositorio ya incluye `req-example.http` como referencia — ahí verás la variante *manual* (copiando el token literal en `@token`).

> Si la lista de libros aparece → **la API y la base de datos están funcionando** 🎉

---

# 📎 Anexos — Material de consulta

> ℹ️ **Estos anexos NO forman parte de la guía de inicio** (secciones 0 a 3). Son material de consulta: comandos útiles, problemas comunes, glosario y la referencia de la API. Para poner en marcha el proyecto solo necesitas las **secciones 0 a 3**.

---

## 4. Comandos útiles

| Comando | Para qué sirve |
|---|---|
| `docker compose ps` | Ver el estado de los contenedores del proyecto |
| `docker compose logs -f api` | Ver los logs de la API en vivo |
| `docker compose logs mysql` | Ver los logs de MySQL (útil cuando "no conecta") |
| `docker compose exec api sh` | Abrir una terminal *dentro* del contenedor de la API |
| `docker compose stop` / `start` | Detener / volver a arrancar los contenedores (sin borrarlos) |
| `docker compose down` | Detener y **eliminar** los contenedores (conserva los datos) |
| `docker compose down -v` | Igual que `down` y además **borra el volumen de datos** (⚠️ el siguiente `up` recrea la BD con las migraciones) |

---

## 5. Problemas comunes y soluciones

| Problema | Solución |
|---|---|
| `docker: command not found` | Docker no está instalado o no corre: revisa la sección 1 |
| `... port is already allocated` (3306 o 8080) | Otra app usa ese puerto: ciérrala o cambia el mapeo en `docker-compose.yml` (ej. `"3307:3306"` / `"8081:8080"`) |
| El build falla descargando imágenes | Sucede sin internet o con red corporativa/proxy |
| La API responde error de conexión a MySQL | `docker compose logs api` y `docker compose logs mysql`; confirma que MySQL esté `(healthy)` antes que la API y que en tu `.env` el `DB_HOST=mysql` |
| Quieres los datos de ejemplo originales | ⚠️ `docker compose down -v` y luego `docker compose up --build` |
| Aviso `the attribute 'version' is obsolete` | Solo es un aviso de Compose; puedes borrar la línea `version: "3.9"` del compose |
| `404` al abrir `http://localhost:8080/` | Normal: la API no tiene ruta raíz. Usa `/health` o `/swagger/index.html` |

---

## 6. Glosario rápido

| Término | En una línea |
|---|---|
| **Imagen** | Plantilla de solo lectura con todo lo necesario; se descarga o se construye |
| **Contenedor** | Una imagen *corriendo*: un proceso aislado |
| **Registro / Docker Hub** | El "almacén" público de imágenes |
| **Docker Compose** | Define y lanza varios contenedores juntos desde un solo archivo |
| **Mapeo de puertos** | `"8080:80"` = puerto de tu máquina `:` puerto del contenedor |
| **Volumen** | Disco persistente: los datos sobreviven aunque el contenedor se borre |
| **Healthcheck** | Prueba automática de que un contenedor está "sano" |

---

# 📚 Referencia de la API

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
docker compose down -v
docker compose up --build
```

El flag `-v` borra el volumen de MySQL, por lo que las migraciones y el seed se vuelven a ejecutar.