package main

import (
	"f4allgo-restaurant/internal/adapter/primary/cli"
	"f4allgo-restaurant/internal/adapter/secondary/eventpublisher"
	"f4allgo-restaurant/internal/adapter/secondary/storage"
	"f4allgo-restaurant/internal/boot"
	"f4allgo-restaurant/internal/core/service"
	"fmt"
	"os"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
)

func main() {
	// Load application configuration from the filesystem (if present).
	boot.LoadConfig()

	// Get the database connection and transaction manager.
	db := boot.GetDatabaseConnection()
	trManager := boot.GetTransactionManager(db)

	// Secondary adapters
	restaurantRepository := storage.NewRestaurantPostgresRepository(db, trmgorm.DefaultCtxGetter, nil)
	outboxPublisher := eventpublisher.NewDomainEventOutboxPublisher(db, trmgorm.DefaultCtxGetter, boot.GetLogger(), boot.GetConfig(), nil)

	// Core service
	restaurantService := service.NewDefaultRestaurantService(restaurantRepository, outboxPublisher, trManager)

	// Primary adapters
	restaurantCli := cli.NewRestaurantCli(restaurantService)

	if err := restaurantCli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
