package helper

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return parts[1], nil
}
