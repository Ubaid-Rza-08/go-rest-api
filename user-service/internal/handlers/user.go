// internal/handlers/user.go
package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Ubaid-Rza-08/go-rest-api/internal/middleware"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/models"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/service"
	"github.com/Ubaid-Rza-08/go-rest-api/internal/utils"
)

// UserHandler holds dependencies for user CRUD routes.
type UserHandler struct {
	userSvc service.UserService
}

// NewUserHandler creates a UserHandler.
func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// GetProfile godoc
// GET /api/v1/profile
// Returns the authenticated user's own profile.
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)

	user, err := h.userSvc.GetByID(userID.(string))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			utils.NotFound(c, "user not found")
			return
		}
		utils.InternalError(c, "failed to fetch profile")
		return
	}

	utils.OK(c, "profile fetched", user)
}

// GetAll godoc
// GET /api/v1/users?page=1&page_size=20
// Lists all users — admin only.
func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := h.userSvc.GetAll(page, pageSize)
	if err != nil {
		utils.InternalError(c, "failed to fetch users")
		return
	}

	utils.OK(c, "users fetched", gin.H{
		"users":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetByID godoc
// GET /api/v1/users/:id
// Returns a single user by ID.
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userSvc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			utils.NotFound(c, "user not found")
			return
		}
		utils.InternalError(c, "failed to fetch user")
		return
	}

	utils.OK(c, "user fetched", user)
}

// Update godoc
// PUT /api/v1/users/:id
// Updates a user's profile. Users can update themselves; admins can update anyone.
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	callerID, _ := c.Get(middleware.ContextUserID)
	callerRole, _ := c.Get(middleware.ContextUserRole)

	user, err := h.userSvc.Update(id, &req, callerID.(string), callerRole.(string))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			utils.NotFound(c, "user not found")
		case errors.Is(err, service.ErrEmailTaken):
			utils.Conflict(c, "email is already taken")
		case errors.Is(err, service.ErrForbidden):
			utils.Forbidden(c, "you cannot update this user")
		default:
			utils.InternalError(c, "failed to update user")
		}
		return
	}

	utils.OK(c, "user updated", user)
}

// Delete godoc
// DELETE /api/v1/users/:id
// Soft-deletes a user. Admin or the user themselves.
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	callerID, _ := c.Get(middleware.ContextUserID)
	callerRole, _ := c.Get(middleware.ContextUserRole)

	if err := h.userSvc.Delete(id, callerID.(string), callerRole.(string)); err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			utils.NotFound(c, "user not found")
		case errors.Is(err, service.ErrForbidden):
			utils.Forbidden(c, "you cannot delete this user")
		default:
			utils.InternalError(c, "failed to delete user")
		}
		return
	}

	utils.OK(c, "user deleted successfully", nil)
}