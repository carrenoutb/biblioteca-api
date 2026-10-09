package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"biblioteca-api/internal/middleware"
)

// SetupRoutes registra todas las rutas de la API sobre el router de Gin.
func SetupRoutes(
	router *gin.Engine,
	authHandler *AuthHandler,
	bookHandler *BookHandler,
	loanHandler *LoanHandler,
	jwtSecret string,
) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api")

	// Rutas públicas de autenticación
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// A partir de aquí, todas las rutas requieren un token JWT válido
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(jwtSecret))
	{
		protected.GET("/auth/me", authHandler.Me)

		protected.GET("/books", bookHandler.List)
		protected.GET("/books/:id", bookHandler.Get)

		protected.POST("/loans", loanHandler.Borrow)
		protected.GET("/loans/me", loanHandler.MyLoans)
		protected.PUT("/loans/:id/return", loanHandler.Return)
	}

	// Rutas exclusivas para administradores
	admin := api.Group("")
	admin.Use(middleware.AuthRequired(jwtSecret), middleware.AdminRequired())
	{
		admin.POST("/books", bookHandler.Create)
		admin.PUT("/books/:id", bookHandler.Update)
		admin.DELETE("/books/:id", bookHandler.Delete)
		admin.GET("/loans", loanHandler.AllLoans)
	}
}
