package main

import (
	"f4allgo-restaurant/internal/adapter/primary/rest"
	"f4allgo-restaurant/internal/adapter/secondary/eventpublisher"
	"f4allgo-restaurant/internal/adapter/secondary/storage"
	"f4allgo-restaurant/internal/boot"
	"f4allgo-restaurant/internal/core/service"
	"fmt"
	"net/http"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load application configuration from the filesystem (if present).
	boot.LoadConfig()

	// Prints the banner with the application name (if configured)
	boot.PrintBanner("Service type: REST", fmt.Sprintf("Listening on port %d", boot.GetConfig().AppPort))

	// Get the database connection and transaction manager.
	gormDB, sqlDB := boot.GetDatabaseConnection()
	trManager := boot.GetTransactionManager(gormDB)

	// Inits tally scope and gets the reporter.
	r := boot.InitTallyReporter(sqlDB)

	// Inits health checks and gets the handler.
	h := boot.GetHealthHandler(sqlDB)

	// Secondary adapter for RestaurantRepository port.
	restaurantRepository := storage.NewRestaurantPostgresRepository(gormDB, trmgorm.DefaultCtxGetter, boot.GetTallyScope())

	// Secondary adapter for DomainEventPublisher port.
	outboxPublisher := eventpublisher.NewDomainEventOutboxPublisher(gormDB, trmgorm.DefaultCtxGetter, boot.GetLogger(), boot.GetConfig(), boot.GetTallyScope())

	// Core service
	restaurantService := service.NewDefaultRestaurantService(restaurantRepository, outboxPublisher, trManager)

	// Primary adapter
	restaurantHandler := rest.NewRestaurantHandler(restaurantService)

	startGinServer(restaurantHandler, r.HTTPHandler(), h)
}

func startGinServer(restaurantHandler *rest.RestaurantHandler, metricsHandler http.Handler, healthHandler http.Handler) {
	gin.SetMode(boot.GetConfig().GinMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/metrics", gin.WrapH(metricsHandler))
	router.GET("/health", gin.WrapH(healthHandler))

	api := router.Group("/api/v1")
	api.GET("/restaurants", restaurantHandler.GetRestaurants)
	api.POST("/restaurants", restaurantHandler.CreateRestaurant)
	api.DELETE("/restaurants/:restaurantId", restaurantHandler.DeleteRestaurant)
	api.GET("/restaurants/:restaurantId", restaurantHandler.GetRestaurant)
	api.PUT("/restaurants/:restaurantId/menu", restaurantHandler.UpdateMenu)

	err := router.Run(fmt.Sprintf(":%d", boot.GetConfig().AppPort))
	if err != nil {
		panic(fmt.Sprintf("failed to listen on port %d", boot.GetConfig().AppPort))
	}
}
