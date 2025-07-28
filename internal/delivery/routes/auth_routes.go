package routes

import (
	"donor-api/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(
	router *gin.RouterGroup,
	authHandler *handler.AuthHandler,
	authMiddleware gin.HandlerFunc,
) {
	// Buat grup baru khusus untuk '/auth'
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/google", authHandler.GoogleAuth)
	}

}
