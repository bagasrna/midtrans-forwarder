package main

import (
	"log"

	"midtrans-forwarder/config"
	"midtrans-forwarder/internal/delivery/http"
	"midtrans-forwarder/internal/repository"
	"midtrans-forwarder/internal/usecase"
	"midtrans-forwarder/pkg/database"
	"midtrans-forwarder/pkg/redis"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	cfg := config.Load()
	db :=database.NewDatabase(cfg.DBConfig)
	redisClient := redis.NewRedisClient(cfg.RedisConfig)

	userRepo := repository.NewUserRepository(db)
	resellerRepo := repository.NewResellerRepository(db, redisClient)

	userUseCase := usecase.NewUserUseCase(userRepo)
	resellerUseCase := usecase.NewResellerUseCase(resellerRepo)
	midtransUseCase := usecase.NewMidtransUseCase(resellerRepo, cfg.MidtransServerKey)

	handler.NewUserHandler(userUseCase, app)
	handler.NewResellerHandler(resellerUseCase, app)
	handler.NewMidtransHandler(midtransUseCase, app)

	log.Fatal(app.Listen(":8080"))
}