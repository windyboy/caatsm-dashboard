package main

import (
	"log"
	"os"

	"casstm-dashboard/handlers"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		// slog.Error("Error loading .env file", "err", err)
		log.Fatal(err)

	}
	// router := chi.NewRouter()
	// router.Handle("/*", public())
	// router.Get("/", handlers.Process(handlers.HandleHome))
	// fmt.Println("Hello, world!")
	listenAddr := os.Getenv("LISTEN_ADDR")
	e := echo.New()
	// e.Use(middleware.Logger())
	e.Static("/public", "public")
	handlers.SetupRoutes(e)
	// slog.Info("Starting server...", "listenAddr", listenAddr)
	// if err := http.ListenAndServe(listenAddr, router); err != nil {
	// 	slog.Error("Error starting server", "err", err)
	// 	log.Fatal(err)
	// }
	// nats := NatsHandler{}
	// go nats.Subscribe()

	e.Logger.Fatal(e.Start(listenAddr))
}
