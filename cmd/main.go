package main

import (
	"fmt"
	"log"
	"scorer/internal/client"
	"scorer/internal/config"
	"scorer/internal/handlers"
	"scorer/internal/metrics"
	"scorer/internal/redis_client"
	"scorer/internal/repositories"
	"scorer/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	//metrics
	reg := prometheus.NewRegistry()
	metric := metric.NewMetric(reg)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	//initialize vars
	productRepo := repositories.NewProductRepository(db)
	merchantRepo := repositories.NewMerchantRepository(db)
	merchantProductRepo := repositories.NewMerchantProductRepository(db)
	redisClient := redis_client.NewRedisClient(AppConfig.QueueParams.Address, AppConfig.QueueParams.Password, AppConfig.QueueParams.Number, AppConfig.QueueParams.Protocol)
	queServiceClient := client.NewQueServiceClient(AppConfig.MicroserviceParams.QueueServerAddress, redisClient)
	service := services.NewServicesFuncs(productRepo, merchantRepo, merchantProductRepo, queServiceClient, metric)
	handler := handlers.NewHandler(service)

	app.Post("/score", handler.HandleScore)

	// expose /metrics on the same Fiber app/port
	app.Get("/metrics", adaptor.HTTPHandler(promHandler))

	app.Listen(AppConfig.ServerParams.ListenPort)
}

func close(db *gorm.DB) {
	db_, err := db.DB()
	if err != nil {
		panic(err)
	}
	db_.Close()
}
