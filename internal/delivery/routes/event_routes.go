package routes

import (
	"donor-api/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func InitEventRoutes(
	router *gin.RouterGroup,
	handler *handler.EventHandler,
) {
	eventsRoutes := router.Group("/events")
	{
		eventsRoutes.POST("", handler.Create)
		eventsRoutes.GET("", handler.GetAll)
		eventsRoutes.GET("/:id", handler.GetByID)
		eventsRoutes.PUT("/:id", handler.Update)
		eventsRoutes.DELETE("/:id", handler.Delete)
	}
}
