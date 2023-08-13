package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ynuraddi/t-medods/config"
	"github.com/ynuraddi/t-medods/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gracefullShutdown(cancel)

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("failed to parse config", err)
	}

	logger := logger.NewLogger(&config, os.Stdout)
	logger.Info("logger inited")
}

func gracefullShutdown(c context.CancelFunc) {
	osC := make(chan os.Signal, 1)
	signal.Notify(osC, os.Interrupt)

	go func() {
		log.Println(<-osC)
		c()
	}()
}
