package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetContextValue(c *gin.Context, key string) (*uuid.UUID, error) {
	valueContext, exists := c.Get(key)
	if !exists {
		return nil, errors.New("tenant ID not found in context")
	}
	value, ok := valueContext.(uuid.UUID)
	if !ok {
		return nil, errors.New("invalid tenant ID format")
	}
	return &value, nil
}
