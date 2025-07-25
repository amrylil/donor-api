package routes

import (
	"donor-api/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

func InitLocationRoutes(
	router *gin.RouterGroup,
	handler *handler.LocationHandler,
) {
	locationsRoutes := router.Group("/locations")
	{
		locationsRoutes.POST("", handler.Create)
		locationsRoutes.GET("", handler.GetAll)
		locationsRoutes.GET("/:id", handler.GetByID)
		locationsRoutes.PUT("/:id", handler.Update)
		locationsRoutes.DELETE("/:id", handler.Delete)
		locationsRoutes.GET("/by-user-location", handler.GetAllByUserLocation)
	}
}
