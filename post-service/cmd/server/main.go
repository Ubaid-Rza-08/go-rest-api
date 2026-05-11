package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Ubaid-Rza-08/post-service/internal/config"
	"github.com/Ubaid-Rza-08/post-service/internal/database"
	"github.com/Ubaid-Rza-08/post-service/internal/handlers"
	"github.com/Ubaid-Rza-08/post-service/internal/middleware"
	"github.com/Ubaid-Rza-08/post-service/internal/repository"
	"github.com/Ubaid-Rza-08/post-service/internal/service"
)

func main() {

	// --------------------------------------------------
	// Load configuration
	// --------------------------------------------------
	cfg := config.Load()

	// --------------------------------------------------
	// Database connection
	// --------------------------------------------------
	db, err := database.NewPostgres(cfg.DBURL)
	if err != nil {
		log.Fatal("DATABASE CONNECTION ERROR:", err)
	}

	defer db.Close()

	log.Println("DATABASE CONNECTED")

	// --------------------------------------------------
	// Dependency Injection
	// --------------------------------------------------
	postRepo := repository.NewPostRepository(db)

	postService := service.NewPostService(postRepo)

	postHandler := handlers.NewPostHandler(postService)

	// --------------------------------------------------
	// Gin Router
	// --------------------------------------------------
	r := gin.New()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Fix proxy warning
	r.SetTrustedProxies(nil)

	// --------------------------------------------------
	// Health Check
	// --------------------------------------------------
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "post-service",
		})
	})

	// --------------------------------------------------
	// API Routes
	// --------------------------------------------------
	api := r.Group("/api/v1")

	// Protected routes
	api.Use(middleware.Authenticate(cfg.JWTSecret))

	{
		api.POST("/posts", postHandler.Create)

		api.GET("/posts", postHandler.GetAll)

		api.GET("/posts/:id", postHandler.GetByID)

		api.PUT("/posts/:id", postHandler.Update)

		api.DELETE("/posts/:id", postHandler.Delete)
	}

	// --------------------------------------------------
	// HTTP Server
	// --------------------------------------------------
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("POST SERVICE RUNNING ON PORT:", cfg.Port)

	// --------------------------------------------------
	// Start Server
	// --------------------------------------------------
	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {

		log.Fatal("SERVER ERROR:", err)
	}
}