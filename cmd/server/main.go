package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"music_library/internal/config"
	"music_library/internal/handlers"
	"music_library/internal/models"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "music_library/docs"
)

func main() {
	log := logrus.New()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Auto migrations
	if err := db.AutoMigrate(&models.Song{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	router := gin.Default()
	handler := &handlers.Handler{
		DB:     db,
		Config: cfg,
		Log:    log,
	}

	// Routes
	router.GET("/songs", handler.GetSongs)
	router.GET("/songs/:id/text", handler.GetSongText)
	router.DELETE("/songs/:id", handler.DeleteSong)
	router.PUT("/songs/:id", handler.UpdateSong)
	router.POST("/songs", handler.AddSong)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := fmt.Sprintf(":%s", cfg.Port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
