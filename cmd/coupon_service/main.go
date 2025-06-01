package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	repo := memdb.NewRepository()
	svc := service.New(repo)
	app := api.New(cfg, svc)

	appErr := make(chan error, 1)
	go func() {
		appErr <- app.Start()
	}()

	// this api should stay alive maximum one year once it starts
	preconfiguredShutDown := 365 * 24 * time.Hour
	oneYearCtx, cancel := context.WithTimeout(context.Background(), preconfiguredShutDown)
	defer cancel()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// block the server alive until one of these conditions are met
	select {
	case stopSignal := <-quit:
		log.Printf("Quit signal received: %s\n", stopSignal)
	case err := <-appErr:
		log.Printf("Error initializing the server: %s\n", err)
	case <-oneYearCtx.Done():
		log.Println("Preconfigured time to quit the server arrived")
	}

	//shutdown the service with 20s delay
	log.Println("Shutting down server...")
	app.Shutdown()
}
