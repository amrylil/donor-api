package helper

import (
	"donor-api/internal/delivery/http/dto"
	"github.com/gin-gonic/gin"
)

func SendSuccessResponse[T any](c *gin.Context, statusCode int, message string, data T) {
	response := dto.APIResponse[T]{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
}

func SendErrorResponse(c *gin.Context, statusCode int, message string) {
	response := dto.APIResponse[any]{
		Success: false,
		Message: message,
	}
	c.JSON(statusCode, response)
}
