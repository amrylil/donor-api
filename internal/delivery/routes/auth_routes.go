package routes

import (
	"donor-api/internal/delivery/http/handler"
	"donor-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(
	router *gin.RouterGroup,
	authHandler *handler.AuthHandler,
	authMiddleware gin.HandlerFunc,
) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/register/admin",
			authMiddleware,
			middleware.RequireRoles("superadmin"),
			authHandler.RegisterAdmin,
		)

		authRoutes.POST("/register/super-admin", authHandler.RegisterSuperAdmin)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/google", authHandler.GoogleAuth)
	}

}
