package routes

import (
	"donor-api/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func InitStockRoutes(
	router *gin.RouterGroup,
	handler *handler.StockHandler,
) {
	stocksRoutes := router.Group("/stocks")
	{
		stocksRoutes.POST("", handler.Create)
		stocksRoutes.GET("", handler.GetAll)
		stocksRoutes.GET("/:id", handler.GetByID)
		stocksRoutes.PUT("/:id", handler.Update)
		stocksRoutes.DELETE("/:id", handler.Delete)
	}
}
