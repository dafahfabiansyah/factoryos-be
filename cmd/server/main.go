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
	appmiddleware "industrial-platform-BE/internal/middleware"
	"industrial-platform-BE/internal/module/auth/handler"
	authmiddleware "industrial-platform-BE/internal/module/auth/middleware"
	"industrial-platform-BE/internal/module/auth/repository"
	"industrial-platform-BE/internal/module/auth/service"
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

	// Initialize auth dependencies
	userRepo := repository.NewUserRepository(dbPool)
	orgRepo := repository.NewOrganizationRepository(dbPool)
	refreshTokenRepo := repository.NewRefreshTokenRepository(dbPool)
	passwordResetRepo := repository.NewPasswordResetTokenRepository(dbPool)
	emailVerifyRepo := repository.NewEmailVerificationTokenRepository(dbPool)

	passwordService := service.NewPasswordService()
	tokenService := service.NewTokenService(cfg.JWTSecret, cfg.JWTExpiryHour)

	authService := service.NewAuthService(
		userRepo,
		orgRepo,
		refreshTokenRepo,
		passwordResetRepo,
		emailVerifyRepo,
		passwordService,
		tokenService,
		cfg.JWTExpiryHour,           // access token TTL
		24*30*time.Hour,             // refresh token TTL (30 days)
		1*time.Hour,                 // password reset TTL (1 hour)
		24*time.Hour,                // email verify TTL (24 hours)
	)

	authHandler := handler.NewAuthHandler(authService)
	authMiddleware := authmiddleware.NewAuthMiddleware(tokenService)

	router := gin.New()
	router.Use(appmiddleware.Recovery())
	router.Use(appmiddleware.ErrorHandler())
	router.Use(appmiddleware.RequestLogger())

	// Health check (no auth)
	router.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := dbPool.Ping(ctx); err != nil {
			c.JSON(503, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Auth routes (public)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.Refresh)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.POST("/forgot-password", authHandler.ForgotPassword)
		authGroup.POST("/reset-password", authHandler.ResetPassword)
		authGroup.POST("/verify-email", authHandler.VerifyEmail)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.GET("/auth/me", authHandler.GetProfile)
		// Admin only routes
		admin := protected.Group("/")
		admin.Use(authMiddleware.RequireAdmin())
		{
			// admin routes here
		}
	}

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
