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

	donationRepo := persistence.NewDonationRepository(db)
	donationUsecase := usecase.NewDonationUsecase(donationRepo)
	donationHandler := handler.NewDonationHandler(donationUsecase)

	eventRepo := persistence.NewEventRepository(db)
	eventUsecase := usecase.NewEventUsecase(eventRepo)
	eventHandler := handler.NewEventHandler(eventUsecase)

	locationRepo := persistence.NewLocationRepository(db)
	locationUsecase := usecase.NewLocationUsecase(locationRepo)
	locationHandler := handler.NewLocationHandler(locationUsecase)

	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		// Panggil fungsi inisialisasi untuk setiap fitur
		InitAuthRoutes(apiV1, authHandler, authMiddleware)
		InitProfileRoutes(apiV1, profileHanlder, authMiddleware)
		InitDonationRoutes(apiV1, donationHandler, authMiddleware)
		InitEventRoutes(apiV1, eventHandler)
		InitLocationRoutes(apiV1, locationHandler)

	}

	return router
}
