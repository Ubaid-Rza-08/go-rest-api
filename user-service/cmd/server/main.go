// cmd/server/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Ubaid-Rza-08/go-rest-api/internal/config"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/database"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/handlers"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/middleware"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/repository"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/service"
)

func main() {
	// --------------------------------------------------
	// Load configuration
	// --------------------------------------------------
	cfg := config.Load()

	// --------------------------------------------------
	// Logger
	// --------------------------------------------------
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// --------------------------------------------------
	// Gin mode
	// --------------------------------------------------
	gin.SetMode(cfg.Server.GinMode)

	// --------------------------------------------------
	// Database connection
	// --------------------------------------------------
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}
	defer db.Close()

	logger.Info("database connected")

	// --------------------------------------------------
	// Dependency Injection
	// --------------------------------------------------
	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(
		userRepo,
		cfg,
	)

	authHandler := handlers.NewAuthHandler(userService)
	userHandler := handlers.NewUserHandler(userService)

	// --------------------------------------------------
	// Router
	// --------------------------------------------------
	router := gin.New()

	// --------------------------------------------------
	// Global middleware
	// --------------------------------------------------
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))

	// --------------------------------------------------
	// Health check
	// --------------------------------------------------
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "go-rest-api",
		})
	})

	// --------------------------------------------------
	// API v1
	// --------------------------------------------------
	api := router.Group("/api/v1")

	// ---------------- AUTH ROUTES ----------------
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// ---------------- PROTECTED ROUTES ----------------
	protected := api.Group("/")
	protected.Use(middleware.Authenticate(cfg.JWT.Secret))
	{
		// Profile
		protected.GET("/profile", userHandler.GetProfile)

		// User CRUD
		protected.GET(
			"/users",
			middleware.RequireRole("admin"),
			userHandler.GetAll,
		)

		protected.GET("/users/:id", userHandler.GetByID)

		protected.PUT("/users/:id", userHandler.Update)

		protected.DELETE("/users/:id", userHandler.Delete)
	}

	// --------------------------------------------------
	// HTTP server configuration
	// --------------------------------------------------
	server := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// --------------------------------------------------
	// Start server
	// --------------------------------------------------
	go func() {
		logger.Info(
			"server started",
			zap.String("port", cfg.Server.Port),
			zap.String("env", cfg.App.Env),
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	// --------------------------------------------------
	// Graceful shutdown
	// --------------------------------------------------
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	logger.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("graceful shutdown failed", zap.Error(err))
	}

	logger.Info("server exited properly")
}