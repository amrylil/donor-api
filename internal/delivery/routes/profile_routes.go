package routes

import (
	"donor-api/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

func InitProfileRoutes(
	router *gin.RouterGroup,
	profileHandler *handler.ProfileHandler,
	authMiddleware gin.HandlerFunc,
) {
	profilRoutes := router.Group("/profile")

	{
		profilRoutes.Use(authMiddleware)
		profilRoutes.GET("/", profileHandler.GetProfile)
		profilRoutes.POST("/detail/create", profileHandler.CreateMyDetail)
		profilRoutes.PUT("/update", profileHandler.UpdateProfile)
		profilRoutes.PUT("/detail/update", profileHandler.UpdateMyDetail)
	}

	userRoutes := router.Group("/users").Use(authMiddleware)

	userRoutes.GET("/", profileHandler.GetAll)
	userRoutes.POST("/", profileHandler.CreateAllUserData)
}
