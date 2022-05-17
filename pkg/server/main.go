package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/e-wrobel/microservice/configuration"
	"github.com/e-wrobel/microservice/server"
)

const (
	configName = "config"
	configType = "yaml"
	configPath = "/etc"
)

func main() {
	config, err := configuration.New(configName, configType, configPath)
	if err != nil {
		log.Fatal(err)
	}

	s, err := server.New(config.SQLFile, config.Listen)
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGINT)
	signal.Notify(sigChan, syscall.SIGHUP)

	go func() {
		sig := <-sigChan
		log.Printf("Got signal: %s, gracefully stopping...", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), server.ShutDownTimeout)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Printf("Shutdown: %v", err)
		}
	}()

	s.Start()
}
