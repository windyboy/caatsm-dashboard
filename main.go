package main

import (
	"log"
	"os"

	"casstm-dashboard/handlers"
	"casstm-dashboard/internal/config"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env file not loaded: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	listenAddr := cfg.Server.ListenAddr
	if envAddr := os.Getenv("LISTEN_ADDR"); envAddr != "" {
		listenAddr = envAddr
	}
	if listenAddr == "" {
		listenAddr = ":3000"
	}

	e := echo.New()
	e.Static("/public", "public")
	handlers.SetupRoutes(e)

	e.Logger.Fatal(e.Start(listenAddr))
}
