package main

import (
	"essay/src/internal/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting application")

	appInstance := app.NewApp()
	defer appInstance.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("Received signal: %s. Shutting down...", sig)
		appInstance.Close()
		os.Exit(0)
	}()

	handler := appInstance.ServeMux()

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
