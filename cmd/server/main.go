package main

import (
	"log"

	"samasalaire-backend/internal/config"
	"samasalaire-backend/internal/database"
	"samasalaire-backend/internal/router"
)

func main() {
	cfg := config.Load()

	client := database.Connect(cfg)
	db := database.GetDB(client, cfg.DBName)

	database.CreateIndexes(db)

	r := router.New(db)

	log.Printf("HTTP server listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
