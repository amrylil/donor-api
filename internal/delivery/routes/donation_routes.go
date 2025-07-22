package routes

import (
	"donor-api/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

func InitDonationRoutes(
	router *gin.RouterGroup,
	handler *handler.DonationHandler,
	authMiddleware gin.HandlerFunc,
) {
	donationsRoutes := router.Group("/donations")
	{
		donationsRoutes.Use(authMiddleware)
		donationsRoutes.POST("", handler.Create)
		donationsRoutes.GET("", handler.GetAll)
		donationsRoutes.GET("/:id", handler.GetByID)
		donationsRoutes.PUT("/:id", handler.Update)
		donationsRoutes.DELETE("/:id", handler.Delete)
	}
}
