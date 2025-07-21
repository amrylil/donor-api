package routes

import (
	"donor-api/internal/delivery/http/handler"
	"donor-api/internal/delivery/http/middleware"
	"donor-api/internal/infrastructure/persistence"
	"donor-api/internal/infrastructure/security"
	"donor-api/internal/usecase"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewAPIRoutes(db *gorm.DB) *gin.Engine {
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	jwtExpHoursStr := os.Getenv("JWT_EXPIRATION_IN_HOURS")
	jwtExpHours, _ := strconv.ParseInt(jwtExpHoursStr, 10, 64)

	userRepo := persistence.NewUserRepository(db)
	jwtService := security.NewJWTService(jwtSecret, jwtExpHours)
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtService)
	authHandler := handler.NewAuthHandler(authUsecase)
	authMiddleware := middleware.AuthMiddleware(jwtService)

	userUsecase := usecase.NewUserUsecase(userRepo)
	profileHanlder := handler.NewProfileHandler(userUsecase)

	router := gin.Default()

	// Buat grup utama untuk API
	apiV1 := router.Group("/api/v1")
	{
		// Panggil fungsi inisialisasi untuk setiap fitur
		InitAuthRoutes(apiV1, authHandler, authMiddleware)
		InitProfileRoutes(apiV1, profileHanlder, authMiddleware)
	}

	return router
}
