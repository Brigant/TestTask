package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Brigant/TestTask/app/api"
	"github.com/Brigant/TestTask/app/api/handler"
	"github.com/Brigant/TestTask/app/config"
	"github.com/Brigant/TestTask/app/service"
	"github.com/Brigant/TestTask/app/storage"
	"github.com/go-playground/validator/v10"
)

func main() {
	const timoutLimit = 5

	cfg, err := config.InitConfig()
	if err != nil {
		log.Printf("failed to initialize config: %s", err)

		return
	}

	storage, err := storage.New(cfg)
	if err != nil {
		log.Printf("failed to initialize storage: %s", err)
		return
	}

	service := service.New(storage, cfg.Server)

	validator := validator.New()

	handler := handler.NewHandler(service, cfg, validator)
	server := api.NewServer(cfg, handler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.HTTPServer.Listen(cfg.Server.AppAddress); err != nil {
			log.Printf("Start and Listen error: %s", err)
		}
	}()

	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), timoutLimit*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown server error: %s", err)

		return
	}

	if err := storage.Close(); err != nil {
		log.Printf("Close repository error: %s", err)

		return
	}

	log.Println("server stopped")
}
