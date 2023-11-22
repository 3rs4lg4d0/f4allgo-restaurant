package main

import (
	"net"
	"net/http"

	pb "f4allgo-restaurant/internal/adapter/primary/grpc"
	"f4allgo-restaurant/internal/adapter/secondary/eventpublisher"
	"f4allgo-restaurant/internal/adapter/secondary/storage"
	"f4allgo-restaurant/internal/boot"
	"f4allgo-restaurant/internal/core/service"

	trmgorm "github.com/avito-tech/go-transaction-manager/gorm"
	"google.golang.org/grpc"
)

func main() {
	// Load application configuration from the filesystem (if present).
	boot.LoadConfig()

	// Prints the banner with the application name (if configured)
	boot.PrintBanner("Service type: gRPC")

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
	rsServer := pb.NewRestaurantServiceServer(restaurantService)

	startServers(rsServer, r.HTTPHandler(), h)
}

func startServers(server pb.RestaurantServiceServer, metricsHandler http.Handler, healthHandler http.Handler) {
	go func() {
		lis, err := net.Listen("tcp", ":8081")
		if err != nil {
			panic("failed to listen on port 8081")
		}

		grpcServer := grpc.NewServer()
		pb.RegisterRestaurantServiceServer(grpcServer, server)
		err = grpcServer.Serve(lis)
		if err != nil {
			panic("failed to start the gRPC service")
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/metrics", metricsHandler)
	mux.Handle("/health", healthHandler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic("failed to listen on port 8080")
	}
}
