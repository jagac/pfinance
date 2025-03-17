package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jagac/pfinance/internal/handlers"
	"github.com/jagac/pfinance/internal/jobs"
	"github.com/jagac/pfinance/internal/repositories"
	"github.com/jagac/pfinance/internal/routes"
	"github.com/jagac/pfinance/internal/services"
	"github.com/jagac/pfinance/pkg/cache"
	"github.com/jagac/pfinance/pkg/config"
	"github.com/jagac/pfinance/pkg/logger"
	"github.com/jagac/pfinance/pkg/worker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	cache := cache.NewCache[string, worker.TaskResult]()
	worker1 := worker.New(logger, cache)

	db, err := config.ConnectDB(ctx)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	go worker1.Run("Worker")

	repo := repositories.NewAssetRepository(db)
	assetService := services.NewAssetService(repo)
	returnCalc := services.NewReturnsCalculator(repo, cache)
	handler := handlers.NewAssetHandler(assetService, returnCalc)
	assetRouter := routes.NewAssetRouter(handler)
	assetRouter.RegisterRoutes(mux)

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	stockFetcher := services.NewStockFetcher()
	goldFetcher := services.NewGoldFetcher()

	goldTask := worker.Task{
		OriginContext: context.Background(),
		Name:          "goldPrice",
		Job:           jobs.FetchGoldJob(goldFetcher),
		TTL:           23 * time.Hour,
	}

	stockTask := worker.Task{
		OriginContext: context.Background(),
		Name:          "stockPrice",
		Job:           jobs.FetchStocksJob(repo, stockFetcher),
		TTL:           23 * time.Hour,
	}

	go func() {
		for range ticker.C {
			worker1.Enqueue(goldTask)
			worker1.Enqueue(stockTask)
		}
	}()

	srv := &http.Server{
		Addr:         ":3000",
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	var wg sync.WaitGroup
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Starting HTTP server on :3000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error:", "err", err.Error())
		}
	}()

	<-shutdownCh
	logger.Info("Shutting down gracefully...")

	if err := worker1.Shutdown(ctx); err != nil {
		logger.Error("Worker forced to shutdown:", "err", err.Error())
	}
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("HTTP server forced to shutdown:", "err", err.Error())
	}
	cancel()
	wg.Wait()
	logger.Info("Shutdown complete")

}
