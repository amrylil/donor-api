package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetContextValue(c *gin.Context, key string) (*uuid.UUID, error) {
	valueContext, exists := c.Get(key)
	if !exists {
		nilUUID := uuid.Nil
		return &nilUUID, nil
	}

	value, ok := valueContext.(uuid.UUID)
	if !ok {
		return nil, errors.New("invalid UUID format")
	}

	return &value, nil
}

func GetContextAnyValue(c *gin.Context, key string) (any, error) {
	valueContext, exists := c.Get(key)
	if !exists {
		return nil, errors.New("value context not found")

	}

	return &valueContext, nil
}
