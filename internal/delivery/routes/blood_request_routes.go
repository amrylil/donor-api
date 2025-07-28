package routes

import (
	"donor-api/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

func InitBloodRequestRoutes(
	router *gin.RouterGroup,
	handler *handler.BloodRequestHandler,
) {
	blood_requestsRoutes := router.Group("/blood-requests")
	{
		blood_requestsRoutes.POST("", handler.Create)
		blood_requestsRoutes.GET("", handler.GetAll)
		blood_requestsRoutes.GET("/:id", handler.GetByID)
		blood_requestsRoutes.PUT("/:id", handler.Update)
		blood_requestsRoutes.DELETE("/:id", handler.Delete)
	}
}
