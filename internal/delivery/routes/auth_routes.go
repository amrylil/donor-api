package routes

import (
	"donor-api/internal/delivery/http/handler"
	// "donor-api/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

// InitAuthRoutes mendaftarkan semua rute yang terkait dengan fitur autentikasi.
func InitAuthRoutes(
	router *gin.RouterGroup, // Menerima grup router utama (misal: /api/v1)
	authHandler *handler.AuthHandler,
	authMiddleware gin.HandlerFunc,
) {
	// Buat grup baru khusus untuk '/auth'
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// userRoutes := router.Group("/users")
	// {
	// 	userRoutes.Use(authMiddleware)
	// 	userRoutes.GET("/me", authHandler.GetProfile)
	// }
}
