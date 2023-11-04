package main

import (
	"f4allgo-restaurant/internal/adapter/primary/api/rest"
	"f4allgo-restaurant/internal/adapter/secondary/eventpublisher"
	"f4allgo-restaurant/internal/adapter/secondary/storage"
	"f4allgo-restaurant/internal/boot"
	"f4allgo-restaurant/internal/core/service"
	"net/http"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load application configuration from the filesystem (if present).
	boot.LoadConfig()

	// Prints the banner with the application name (if configured)
	boot.PrintBanner()

	// Get the database connection and transaction manager.
	gormDb := boot.GetDatabaseConnection()
	trManager := boot.GetTransactionManager(gormDb)
	sqlDb, _ := gormDb.DB()

	// Inits tally scope and gets the reporter.
	r := boot.InitTallyReporter(sqlDb)

	// Inits health checks and gets the handler.
	h := boot.GetHealthHandler(sqlDb)

	// Secondary adapter for RestaurantRepository port.
	restaurantRepository := storage.NewRestaurantPostgresRepository(gormDb, trmgorm.DefaultCtxGetter, boot.GetTallyScope())

	// Secondary adapter for DomainEventPublisher port.
	outboxPublisher := eventpublisher.NewDomainEventOutboxPublisher(gormDb, trmgorm.DefaultCtxGetter, boot.GetLogger(), boot.GetTallyScope())

	// Core service
	restaurantService := service.NewDefaultRestaurantService(restaurantRepository, outboxPublisher, trManager)

	// Primary adapters
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

	err := router.Run(":8080")
	if err != nil {
		panic("failed to listen on port 8080")
	}
}
