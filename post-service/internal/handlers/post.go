// internal/handlers/post.go
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Ubaid-Rza-08/post-service/internal/models"
	"github.com/Ubaid-Rza-08/post-service/internal/service"
	"github.com/Ubaid-Rza-08/post-service/internal/utils"
)

type PostHandler struct {
	service *service.PostService
}

func NewPostHandler(s *service.PostService) *PostHandler {
	return &PostHandler{
		service: s,
	}
}

// -----------------------------------------------------
// CREATE POST
// -----------------------------------------------------

func (h *PostHandler) Create(c *gin.Context) {

	var req models.CreatePostRequest

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("BIND ERROR:", err)

		utils.BadRequest(c, err.Error())
		return
	}

	// Extract userID from JWT middleware
	userID, exists := c.Get("userID")
	if !exists {
		fmt.Println("JWT USER ID NOT FOUND")

		utils.Unauthorized(c, "unauthorized")
		return
	}

	// Create post
	post, err := h.service.Create(
		&req,
		userID.(string),
	)

	if err != nil {
		fmt.Println("CREATE POST ERROR:", err)

		utils.Internal(c, err.Error())
		return
	}

	utils.Created(c, "post created successfully", post)
}

// -----------------------------------------------------
// GET ALL POSTS
// -----------------------------------------------------

func (h *PostHandler) GetAll(c *gin.Context) {

	posts, err := h.service.GetAll()
	if err != nil {

		fmt.Println("GET ALL POSTS ERROR:", err)

		utils.Internal(c, "failed to fetch posts")
		return
	}

	utils.OK(c, "posts fetched successfully", posts)
}

// -----------------------------------------------------
// GET POST BY ID
// -----------------------------------------------------

func (h *PostHandler) GetByID(c *gin.Context) {

	id := c.Param("id")

	post, err := h.service.GetByID(id)
	if err != nil {

		fmt.Println("GET POST ERROR:", err)

		utils.NotFound(c, "post not found")
		return
	}

	utils.OK(c, "post fetched successfully", post)
}

// -----------------------------------------------------
// UPDATE POST
// -----------------------------------------------------

func (h *PostHandler) Update(c *gin.Context) {

	id := c.Param("id")

	var req models.UpdatePostRequest

	// Validate body
	if err := c.ShouldBindJSON(&req); err != nil {

		fmt.Println("UPDATE BIND ERROR:", err)

		utils.BadRequest(c, err.Error())
		return
	}

	// JWT user
	userID, exists := c.Get("userID")
	if !exists {

		fmt.Println("JWT USER ID NOT FOUND")

		utils.Unauthorized(c, "unauthorized")
		return
	}

	post, err := h.service.Update(
		id,
		&req,
		userID.(string),
	)

	if err != nil {

		fmt.Println("UPDATE POST ERROR:", err)

		switch {
		case errors.Is(err, service.ErrForbidden):
			utils.Forbidden(c, "you cannot update this post")

		default:
			utils.Internal(c, "failed to update post")
		}

		return
	}

	utils.OK(c, "post updated successfully", post)
}

// -----------------------------------------------------
// DELETE POST
// -----------------------------------------------------

func (h *PostHandler) Delete(c *gin.Context) {

	id := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {

		fmt.Println("JWT USER ID NOT FOUND")

		utils.Unauthorized(c, "unauthorized")
		return
	}

	err := h.service.Delete(
		id,
		userID.(string),
	)

	if err != nil {

		fmt.Println("DELETE POST ERROR:", err)

		switch {
		case errors.Is(err, service.ErrForbidden):
			utils.Forbidden(c, "you cannot delete this post")

		default:
			utils.Internal(c, "failed to delete post")
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "post deleted successfully",
	})
}