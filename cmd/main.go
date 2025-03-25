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
	"github.com/jagac/pfinance/internal/middleware"
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

	loggingConfig := middleware.LoggingConfig{Logger: logger}
	corsConfig := middleware.CORSConfig{}
	logMiddleware := loggingConfig.Middleware
	corsMiddleware := corsConfig.Middleware
	stockFetcher := services.NewStockFetcher()
	goldFetcher := services.NewGoldFetcher()

	repo := repositories.NewAssetRepository(db)
	newRepo := repositories.NewAssetReturnHistoryRepository(db)
	assetService := services.NewAssetService(repo)
	returnService := services.NewHistoricReturns(repo, newRepo, goldFetcher, stockFetcher)
	returnCalc := services.NewReturnsCalculator(repo, cache, newRepo)
	handler := handlers.NewAssetHandler(assetService, returnCalc, returnService)
	assetRouter := routes.NewAssetRouter(handler, logMiddleware, corsMiddleware)
	assetRouter.RegisterRoutes(mux)

	hourlyTicker := time.NewTicker(31 * time.Minute)
	defer hourlyTicker.Stop()
	dailyTicker := time.NewTicker(24 * time.Hour)
	defer dailyTicker.Stop()

	goldTask := worker.Task{
		OriginContext: context.Background(),
		Name:          "goldPrice",
		Job:           jobs.FetchGoldJob(goldFetcher),
		TTL:           29 * time.Minute,
	}

	stockTask := worker.Task{
		OriginContext: context.Background(),
		Name:          "stockPrice",
		Job:           jobs.FetchStocksJob(repo, stockFetcher),
		TTL:           29 * time.Minute,
	}

	dailyReturnTask := worker.Task{
		OriginContext: context.Background(),
		Name:          "dailyReturn",
		Job:           jobs.TotalReturnsJob(returnService),
		TTL:           10 * time.Minute,
	}

	go func() {
		for range hourlyTicker.C {
			worker1.Enqueue(goldTask)
			worker1.Enqueue(stockTask)
		}
	}()
	go func() {
		for range dailyTicker.C {
			worker1.Enqueue(dailyReturnTask)
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
