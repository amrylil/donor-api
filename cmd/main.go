package main

import (
	"donor-api/internal/delivery/routes"
	"donor-api/internal/entity"
	"donor-api/internal/infrastructure/database"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

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

	log.Printf("🚀 Server berjalan di http://localhost:8000")
	router.Run(":8080")
}
