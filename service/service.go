package service

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Start(serviceName string, port, signingKey string, registerHandlersFunc func(*fiber.App, string), gracefulShutdownFunc func() error) {

	service := fiber.New(fiber.Config{
		AppName: serviceName,
	})
	service.Use(logger.New())

	registerHandlersFunc(service, signingKey)

	log.Printf("Starting: %s\n", serviceName)

	errs := make(chan error)
	go func() {
		if err := service.Listen(":" + port); err != nil {
			errs <- err
		}
	}()
	go func() { errs <- gracefulShutdown(service, gracefulShutdownFunc) }()

	if err := <-errs; err != nil {
		log.Fatal(err)
	}
}

func gracefulShutdown(app *fiber.App, gracefulShutdownFunc func() error) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	log.Println("Gracefully shutting down...")
	if err := app.Shutdown(); err != nil {
		return err
	}

	log.Println("Running cleanup tasks...")
	if err := gracefulShutdownFunc(); err != nil {
		return err
	}

	log.Println("Service was successful shutdown.")
	return nil
}
