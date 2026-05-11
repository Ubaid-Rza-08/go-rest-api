// internal/handlers/auth.go
package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/Ubaid-Rza-08/go-rest-api/internal/models"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/service"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/utils"
)

// AuthHandler holds dependencies for authentication routes.
type AuthHandler struct {
	userSvc service.UserService
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(userSvc service.UserService) *AuthHandler {
	return &AuthHandler{userSvc: userSvc}
}

// Register godoc
// POST /api/v1/auth/register
// Creates a new user account.
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.userSvc.Register(&req)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			utils.Conflict(c, "email is already registered")
			return
		}
		utils.InternalError(c, "failed to create account")
		return
	}

	utils.Created(c, "account created successfully", user)
}

// Login godoc
// POST /api/v1/auth/login
// Authenticates a user and returns a JWT token.
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := h.userSvc.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCreds) {
			// Use 401 and a generic message to prevent user enumeration
			utils.Unauthorized(c, "invalid email or password")
			return
		}
		utils.InternalError(c, "login failed")
		return
	}

	utils.OK(c, "login successful", resp)
}