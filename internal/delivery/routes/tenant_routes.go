package routes

import (
  "donor-api/internal/delivery/http/handler"
  "github.com/gin-gonic/gin"
)

func InitTenantRoutes(
  router *gin.RouterGroup,
  handler *handler.TenantHandler,
) {
  tenantsRoutes := router.Group("/tenants")
  {
    tenantsRoutes.POST("", handler.Create)
    tenantsRoutes.GET("", handler.GetAll)
    tenantsRoutes.GET("/:id", handler.GetByID)
    tenantsRoutes.PUT("/:id", handler.Update)
    tenantsRoutes.DELETE("/:id", handler.Delete)
  }
}
