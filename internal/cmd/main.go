package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VolkHackVH/todo-list.git/internal/db"
	"github.com/VolkHackVH/todo-list.git/internal/router"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	//todo: init logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	//todo: load .env for file
	if err := godotenv.Load(); err != nil {
		logger.Warn(".env file not found")
	}

	//todo: init database manager
	dbManager := db.NewDBManager(os.Getenv("DATABASE_URL"), logger)

	//? database connect and retry
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dbManager.Connect(ctx); err != nil {
		logger.Fatal("DB connection failed: %w", err)
	}

	healtCtx, healtCancel := context.WithCancel(context.Background())
	defer healtCancel()
	dbManager.StartHealthCheck(healtCtx, 60*time.Second)

	//todo: init router
	r := router.InitRouter(dbManager.Queries)

	//todo: server settings
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//todo: startup server
	go func() {
		logger.Infof("Start server on %v", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server startup failed: %w", err)
		}
	}()

	//todo: graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server is shutdowning...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server shutdown error: %w", err)
	}

	logger.Info("Server Shutdown")
}
