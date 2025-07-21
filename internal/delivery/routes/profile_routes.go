package routes

import (
	"donor-api/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

func InitProfileRoutes(
	router *gin.RouterGroup, // Menerima grup router utama (misal: /api/v1)
	profileHandler *handler.ProfileHandler,
	authMiddleware gin.HandlerFunc,
) {
	userRoutes := router.Group("/profile")

	{
		userRoutes.Use(authMiddleware)
		userRoutes.GET("/", profileHandler.GetProfile)
		userRoutes.POST("/create", profileHandler.CreateMyDetail)
		userRoutes.PUT("/update", profileHandler.UpdateProfile)
		userRoutes.PUT("/detail/update", profileHandler.UpdateMyDetail)
	}
}
