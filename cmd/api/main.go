// Package main es el punto de entrada de la API de la Biblioteca.
//
// @title           Biblioteca API
// @version         1.0
// @description     API REST para el sistema de biblioteca (proyecto de curso Front-End con React).
// @BasePath        /api
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Escribe: Bearer <tu_token>
package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "biblioteca-api/docs" // Documentos generados por swag init

	"biblioteca-api/internal/config"
	"biblioteca-api/internal/db"
	"biblioteca-api/internal/handler"
	"biblioteca-api/internal/repository"
	"biblioteca-api/internal/service"
)

func main() {
	cfg := config.Load()

	database := db.Connect(cfg.DSN())
	defer database.Close()

	// Repositorios
	userRepo := repository.NewUserRepository(database)
	bookRepo := repository.NewBookRepository(database)
	loanRepo := repository.NewLoanRepository(database)

	// Servicios
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpHours)
	bookService := service.NewBookService(bookRepo)
	loanService := service.NewLoanService(loanRepo, bookRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	bookHandler := handler.NewBookHandler(bookService)
	loanHandler := handler.NewLoanHandler(loanService)

	router := gin.Default()

	// CORS abierto para que el frontend en React (otro puerto/origen) pueda consumir la API.
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	handler.SetupRoutes(router, authHandler, bookHandler, loanHandler, cfg.JWTSecret)

	log.Printf("Servidor escuchando en el puerto %s", cfg.AppPort)
	log.Printf("Documentación Swagger disponible en http://localhost:%s/swagger/index.html", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
