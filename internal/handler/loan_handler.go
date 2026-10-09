package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"biblioteca-api/internal/models"
	"biblioteca-api/internal/service"
)

type LoanHandler struct {
	loanService *service.LoanService
}

func NewLoanHandler(loanService *service.LoanService) *LoanHandler {
	return &LoanHandler{loanService: loanService}
}

// Borrow godoc
// @Summary      Pedir prestado un libro
// @Tags         loans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.LoanRequest true "ID del libro a pedir prestado"
// @Success      201 {object} models.Loan
// @Failure      400 {object} map[string]string
// @Router       /loans [post]
func (h *LoanHandler) Borrow(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req models.LoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan, err := h.loanService.Borrow(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, loan)
}

// MyLoans godoc
// @Summary      Mis préstamos
// @Description  Lista los préstamos del usuario autenticado, incluyendo activos e histórico
// @Tags         loans
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.LoanDetail
// @Router       /loans/me [get]
func (h *LoanHandler) MyLoans(c *gin.Context) {
	userID := c.GetInt("user_id")

	loans, err := h.loanService.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener los préstamos"})
		return
	}

	c.JSON(http.StatusOK, loans)
}

// Return godoc
// @Summary      Devolver un libro
// @Tags         loans
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID del préstamo"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /loans/{id}/return [put]
func (h *LoanHandler) Return(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	if err := h.loanService.Return(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "libro devuelto correctamente"})
}

// AllLoans godoc
// @Summary      Ver todos los préstamos (solo admin)
// @Tags         loans
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.LoanDetail
// @Router       /loans [get]
func (h *LoanHandler) AllLoans(c *gin.Context) {
	loans, err := h.loanService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron obtener los préstamos"})
		return
	}

	c.JSON(http.StatusOK, loans)
}
