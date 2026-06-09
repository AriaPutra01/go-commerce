package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/AriaPutra01/go-commerce/internal/cache"
	"github.com/AriaPutra01/go-commerce/internal/config"
	"github.com/AriaPutra01/go-commerce/internal/module/auth"
	"github.com/AriaPutra01/go-commerce/internal/token"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	RDB      *redis.Client
	App      *gin.Engine
	Log      *slog.Logger
	Config   *config.Config
	JWTMaker *token.JWTMaker
}

func Bootstrap(config *BootstrapConfig) {
	cache := cache.NewRedisStore(config.RDB)
	authRepository := auth.NewRepository(config.DB)
	authService := auth.NewService(config.JWTMaker, authRepository, cache)
	authHandler := auth.NewHandler(authService)
	authRoute := auth.NewRoute(config.App, authHandler)
	authRoute.RegisterRoute()
}

func GracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
