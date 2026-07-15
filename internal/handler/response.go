package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type apiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func successResponse(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, apiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func errorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, apiResponse{
		Success: false,
		Message: message,
	})
}

func badRequest(c *gin.Context, message string) {
	errorResponse(c, http.StatusBadRequest, message)
}

func notFound(c *gin.Context, message string) {
	errorResponse(c, http.StatusNotFound, message)
}

func conflict(c *gin.Context, message string) {
	errorResponse(c, http.StatusConflict, message)
}

func internalError(c *gin.Context, message string) {
	errorResponse(c, http.StatusInternalServerError, message)
}
