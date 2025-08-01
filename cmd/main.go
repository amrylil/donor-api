package main

import (
	"donor-api/internal/delivery/routes"
	"donor-api/internal/entity"
	"donor-api/internal/infrastructure/database"
	"fmt"
	"log"

	_ "donor-api/docs"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Donor App API
// @version         1.0
// @description     API untuk mengelola donor darah, stok darah, dan informasi pendonor dengan autentikasi OAuth.
// @termsOfService  https://donor-darah.duckdns.org/terms

// @contact.name   Ulil Amry
// @contact.url    https://donor-darah.duckdns.org
// @contact.email  support@donor-darah.duckdns.org

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      donor-darah.duckdns.org
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke database: %v", err)
	}
	fmt.Println("✅ Berhasil terhubung ke database!")

	err = db.AutoMigrate(&entity.User{}, &entity.UserDetail{}, &entity.Donation{}, &entity.Event{}, &entity.Stock{}, &entity.Location{}, &entity.BloodRequest{})
	if err != nil {
		log.Fatalf("❌ Gagal melakukan migrasi database: %v", err)
	}
	fmt.Println("✅ Migrasi database berhasil!")

	router := routes.NewAPIRoutes(db)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("🚀 Server berjalan di http://localhost:8000")
	router.Run(":8080")
}
