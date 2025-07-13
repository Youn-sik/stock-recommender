package main

import (
	"log"
	"stock-recommender/backend/config"
	"stock-recommender/backend/database"
	"stock-recommender/backend/router"
	"stock-recommender/backend/services"
	"stock-recommender/backend/workers"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting Stock Recommender API on port %s", cfg.Port)

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize partition manager and create partitions
	partitionManager := services.NewPartitionManager(db)
	err = partitionManager.CreateMonthlyPartitions()
	if err != nil {
		log.Printf("Warning: Failed to create partitions: %v", err)
	}

	// Initialize services
	cacheService := services.NewCacheService(cfg)
	queueService, err := services.NewQueueService(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize queue service: %v", err)
		queueService = nil
	}

	// Initialize data collector service
	dataCollector := services.NewDataCollectorService(db, cfg)
	
	// Start scheduled data collection
	go dataCollector.StartScheduledCollection()

	aiClient := services.NewAIClient(cfg)
	indicatorService := services.NewIndicatorService()
	signalGenerator := services.NewSignalGeneratorService(db, indicatorService, aiClient, cacheService, queueService)

	// Start queue workers if queue service is available
	if queueService != nil {
		queueWorker := workers.NewQueueWorker(db, queueService, indicatorService, signalGenerator, aiClient, cacheService)
		err = queueWorker.StartWorkers()
		if err != nil {
			log.Printf("Warning: Failed to start queue workers: %v", err)
		} else {
			log.Println("Queue workers started successfully")
		}
	}

	// Setup router
	r := router.Setup(db, cfg)

	// Start server
	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}