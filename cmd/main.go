package main

import (
	"fmt"
	"log"
	"scorer/internal/client"
	"scorer/internal/config"
	"scorer/internal/handlers"
	"scorer/internal/repositories"
	"scorer/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New()

	AppConfig := config.LoadConfig()

	//get database connection to pass into repos
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", AppConfig.DBParams.Host, AppConfig.DBParams.User, AppConfig.DBParams.Password, AppConfig.DBParams.Name, AppConfig.DBParams.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("couldn't connect to database", err)
	}
	defer close(db)

	//initialize vars
	productRepo := repositories.NewProductRepository(db)
	merchantRepo := repositories.NewMerchantRepository(db)
	merchantProductRepo := repositories.NewMerchantProductRepository(db)
	queServiceClient := client.NewQueServiceClient("")
	service := services.NewServicesFuncs(productRepo, merchantRepo, merchantProductRepo, queServiceClient)
	handler := handlers.NewHandler(service)

	app.Post("/score", handler.HandleScore)

	app.Listen(AppConfig.ServerParams.ListenPort)
}

func close(db *gorm.DB) {
	db_, err := db.DB()
	if err != nil {
		panic(err)
	}
	db_.Close()
}