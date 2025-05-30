package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
)

func main() {
	cfg := config.New()
	repo := memdb.New()
	svc := service.New(repo)
	app := api.New(cfg.API, svc)
	app.Start()

	// TODO
	// Wait for interrupt signal (Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Block until signal is received
	<-stop

	// Gracefully shutdown the server
	fmt.Println("Shutting down server...")
	app.Close()
}
