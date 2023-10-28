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
	tally "github.com/uber-go/tally/v4"
)

func main() {
	// Load application configuration from the filesystem (if present).
	boot.LoadConfig()

	// Prints the banner with the application name (if configured)
	boot.PrintBanner()

	// Get the database connection and transaction manager.
	db := boot.GetDatabaseConnection()
	trManager := boot.GetTransactionManager(db)

	// Inits tally scope and gets the reporter.
	r := boot.GetTallyReporter()

	// Inits health checks and gets the handler.
	sqlDb, _ := db.DB()
	h := boot.GetHealthHandler(sqlDb)

	// Secondary adapter for RestaurantRepository port.
	timer := boot.GetTallyScope().Tagged(map[string]string{"repository": "restaurant"}).Timer("database_durations")
	restaurantRepository := storage.NewRestaurantPostgresRepository(db, trmgorm.DefaultCtxGetter, timer)

	// Secondary adapter for DomainEventPublisher port.
	restaurantCreated := boot.GetTallyScope().Tagged(map[string]string{"event_type": "RestaurantCreated"}).Counter("outgoing_events")
	restaurantDeleted := boot.GetTallyScope().Tagged(map[string]string{"event_type": "RestaurantDeleted"}).Counter("outgoing_events")
	restaurantMenuUpdated := boot.GetTallyScope().Tagged(map[string]string{"event_type": "RestaurantMenuUpdated"}).Counter("outgoing_events")
	eventCounters := map[string]tally.Counter{
		"RestaurantCreated":     restaurantCreated,
		"RestaurantDeleted":     restaurantDeleted,
		"RestaurantMenuUpdated": restaurantMenuUpdated,
	}
	outboxPublisher := eventpublisher.NewDomainEventOutboxPublisher(db, trmgorm.DefaultCtxGetter, boot.GetLogger(), eventCounters)

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
