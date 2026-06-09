package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AriaPutra01/go-commerce/internal/app"
	"github.com/AriaPutra01/go-commerce/internal/cache"
	"github.com/AriaPutra01/go-commerce/internal/config"
	"github.com/AriaPutra01/go-commerce/internal/database"
	"github.com/AriaPutra01/go-commerce/internal/logger"
	"github.com/AriaPutra01/go-commerce/internal/token"
)

func main() {
	cfg := config.Load()
	log := logger.NewLogger(cfg)
	engine := app.NewGin(cfg, log)
	jwtMaker := token.NewJWTMaker(cfg.SecretKey)
	db := database.NewDatabase(cfg, log)
	rdb := cache.NewRedis(cfg)

	app.Bootstrap(&app.BootstrapConfig{
		DB:       db,
		RDB:      rdb,
		App:      engine,
		Log:      log,
		Config:   cfg,
		JWTMaker: jwtMaker,
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      engine,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go app.GracefulShutdown(server, done)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Info("Graceful shutdown complete.")
}
