package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	httpDelivery "github.com/bookmarketing/marketing-saas/api-go/internal/delivery/grpc"
	httpRouter "github.com/bookmarketing/marketing-saas/api-go/internal/delivery/http"
	"github.com/bookmarketing/marketing-saas/api-go/internal/infrastructure/database"
	"github.com/bookmarketing/marketing-saas/api-go/internal/infrastructure/repository"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Initialize repositories
	campaignRepo := repository.NewCampaignRepository(db.DB)
	dashboardRepo := repository.NewDashboardRepository(db.DB)

	// Initialize services for gRPC
	campaignService := service.NewCampaignService(campaignRepo)
	dashboardService := service.NewDashboardService(dashboardRepo)

	// Start gRPC server in a goroutine
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9090"
	}

	grpcServer := httpDelivery.NewGRPCServer(campaignService, dashboardService)
	go func() {
		log.Printf("Starting gRPC server on port %s", grpcPort)
		if err := grpcServer.Start(grpcPort); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Initialize and run HTTP server
	router := httpRouter.NewRouter(db.DB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting HTTP server on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
