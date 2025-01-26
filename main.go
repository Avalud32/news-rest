package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"news-rest/config"
	"news-rest/db"
	"news-rest/handlers"
	"news-rest/repository"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	dbConn := db.GetDB(cfg)

	err = db.PingDB(cfg)
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	newsRepo := repository.NewNewsRepository(dbConn)
	newsHandler := handlers.NewNewsHandler(newsRepo)

	app := fiber.New()
	api := app.Group("/api")
	api.Get("/list", newsHandler.GetNewsList)
	api.Post("/edit/:Id", newsHandler.UpdateNews)

	go func() {
		log.Printf("Server started on port: %s", cfg.ServerPort)
		if err := app.Listen(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	log.Println("Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server Shutdown complete.")
}
