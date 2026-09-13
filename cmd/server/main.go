package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"industrial-platform-BE/internal/middleware"
	"industrial-platform-BE/internal/platform"
)

func main() {
	cfg := platform.Load()

	platform.InitLogger(cfg.AppEnv)
	defer platform.Sync()

	dbPool, err := platform.NewDBPool(cfg.DatabaseURL)
	if err != nil {
		platform.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	platform.Log.Info("Database connected",
		zap.String("env", cfg.AppEnv),
		zap.String("port", cfg.AppPort),
	)

	router := gin.New()
	router.Use(middleware.Recovery())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.RequestLogger())

	router.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := dbPool.Ping(ctx); err != nil {
			c.JSON(503, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		platform.Log.Info("Starting server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			platform.Log.Fatal("Server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	platform.Log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		platform.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	platform.Log.Info("Server exited")
}
