package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Ubaid-Rza-08/api-gateway/internal/config"
	"github.com/Ubaid-Rza-08/api-gateway/internal/gateway"
	"github.com/Ubaid-Rza-08/api-gateway/internal/middleware"
)

func main() {

	cfg := config.Load()

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	r.SetTrustedProxies(nil)

	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	// ------------------------------------------------
	// Health
	// ------------------------------------------------

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "gateway running",
		})
	})

	// ------------------------------------------------
	// Proxies
	// ------------------------------------------------

	authProxy := gateway.ReverseProxy(
		cfg.AuthServiceURL,
	)

	postProxy := gateway.ReverseProxy(
		cfg.PostServiceURL,
	)

	// ------------------------------------------------
	// AUTH SERVICE
	// ------------------------------------------------

	r.Any("/api/v1/auth/login", authProxy)

	r.Any("/api/v1/auth/register", authProxy)

	r.Any("/api/v1/profile", authProxy)

	r.Any("/api/v1/users", authProxy)

	r.Any("/api/v1/users/:id", authProxy)

	// ------------------------------------------------
	// POST SERVICE
	// ------------------------------------------------

	// EXACT ROUTES
	r.Any("/api/v1/posts", postProxy)

	r.Any("/api/v1/posts/:id", postProxy)

	// ------------------------------------------------
	// HTTP SERVER
	// ------------------------------------------------

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("API Gateway running on port:", cfg.Port)

	if err := srv.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {

		log.Fatal(err)
	}
}