package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "github.com/tamaqazaq/subscription-service/docs"

	"github.com/tamaqazaq/subscription-service/config"
	"github.com/tamaqazaq/subscription-service/internal/adapters/http"
	"github.com/tamaqazaq/subscription-service/internal/infrastructure/postgres"
	"github.com/tamaqazaq/subscription-service/internal/usecase"
)

// @title Subscription Service API
// @version 1.0
// @description REST API for managing subscriptions
// @host localhost:8080
// @BasePath /

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", "host="+cfg.DBHost+
		" port="+cfg.DBPort+
		" user="+cfg.DBUser+
		" password="+cfg.DBPassword+
		" dbname="+cfg.DBName+
		" sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	repo := postgres.NewSubscriptionRepository(db)
	uc := usecase.NewSubscriptionUsecase(repo)
	handler := http.NewSubscriptionHandler(uc)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	sub := router.Group("/subscriptions")
	{
		sub.POST("", handler.CreateSubscription)
		sub.GET("", handler.GetAll)
		sub.GET("/:id", handler.GetByID)
		sub.PUT("/:id", handler.Update)
		sub.DELETE("/:id", handler.Delete)
		sub.GET("/total", handler.GetTotal)
	}

	log.Println("Server running on port " + cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
