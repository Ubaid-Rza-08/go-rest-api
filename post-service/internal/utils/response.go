package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSON(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, gin.H{
		"message": message,
		"data":    data,
	})
}

func OK(c *gin.Context, msg string, data interface{}) {
	JSON(c, http.StatusOK, msg, data)
}

func Created(c *gin.Context, msg string, data interface{}) {
	JSON(c, http.StatusCreated, msg, data)
}

func BadRequest(c *gin.Context, msg string) {
	JSON(c, http.StatusBadRequest, msg, nil)
}

func Unauthorized(c *gin.Context, msg string) {
	JSON(c, http.StatusUnauthorized, msg, nil)
}

func Forbidden(c *gin.Context, msg string) {
	JSON(c, http.StatusForbidden, msg, nil)
}

func NotFound(c *gin.Context, msg string) {
	JSON(c, http.StatusNotFound, msg, nil)
}

func Internal(c *gin.Context, msg string) {
	JSON(c, http.StatusInternalServerError, msg, nil)
}