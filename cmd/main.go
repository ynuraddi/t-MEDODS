package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ynuraddi/t-medods/config"
	"github.com/ynuraddi/t-medods/logger"
	"github.com/ynuraddi/t-medods/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gracefullShutdown(cancel)

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalln("failed to parse config:", err)
	}

	logger := logger.NewLogger(&config, os.Stdout)
	logger.Info("logger inited")

	clientMongo, err, close := mongoClient(&config, logger)
	if err != nil {
		log.Fatalln("failed init mongo client:", err)
	}
	defer close()

	repository := repository.New(config, logger, clientMongo)
}

func mongoClient(config *config.Config, logger logger.Logger) (client *mongo.Client, err error, close func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err, nil
	}

	return client, nil, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			logger.Error("failed disconnect mongo client", err)
			return
		}
		logger.Info("success disconnect mongo client")
	}
}

func gracefullShutdown(c context.CancelFunc) {
	osC := make(chan os.Signal, 1)
	signal.Notify(osC, os.Interrupt)

	go func() {
		log.Println(<-osC)
		c()
	}()
}
