package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"biblioteca-api/internal/models"
	"biblioteca-api/internal/service"
)

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

// List godoc
// @Summary      Listar libros
// @Description  Devuelve el catálogo de libros, con filtros opcionales por categoría y búsqueda de título
// @Tags         books
// @Produce      json
// @Security     BearerAuth
// @Param        category query string false "Filtrar por categoría"
// @Param        search   query string false "Buscar por título"
// @Success      200 {array} models.Book
// @Router       /books [get]
func (h *BookHandler) List(c *gin.Context) {
	category := c.Query("category")
	search := c.Query("search")

	books, err := h.bookService.List(category, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo obtener el catálogo"})
		return
	}
	c.JSON(http.StatusOK, books)
}

// Get godoc
// @Summary      Obtener un libro por ID
// @Tags         books
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del libro"
// @Success      200 {object} models.Book
// @Failure      404 {object} map[string]string
// @Router       /books/{id} [get]
func (h *BookHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	book, err := h.bookService.Get(id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "libro no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al buscar el libro"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// Create godoc
// @Summary      Crear un libro (solo admin)
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.BookRequest true "Datos del libro"
// @Success      201 {object} models.Book
// @Failure      400 {object} map[string]string
// @Router       /books [post]
func (h *BookHandler) Create(c *gin.Context) {
	var req models.BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := h.bookService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear el libro"})
		return
	}
	c.JSON(http.StatusCreated, book)
}

// Update godoc
// @Summary      Actualizar un libro (solo admin)
// @Tags         books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del libro"
// @Param        request body models.BookRequest true "Datos del libro"
// @Success      200 {object} models.Book
// @Failure      404 {object} map[string]string
// @Router       /books/{id} [put]
func (h *BookHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	var req models.BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book, err := h.bookService.Update(id, req)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "libro no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo actualizar el libro"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// Delete godoc
// @Summary      Eliminar un libro (solo admin)
// @Tags         books
// @Security     BearerAuth
// @Param        id path int true "ID del libro"
// @Success      204
// @Router       /books/{id} [delete]
func (h *BookHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	if err := h.bookService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo eliminar el libro"})
		return
	}

	c.Status(http.StatusNoContent)
}
