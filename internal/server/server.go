package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"restaurant-inventory-api/internal/config"
)

func Run(cfg config.Config, handler http.Handler) error {
	httpServer := &http.Server{
		Addr:         cfg.HTTPAddr(),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("starting %s on %s", cfg.AppName, cfg.HTTPAddr())
		serverErrors <- httpServer.ListenAndServe()
	}()

	shutdownSignals := make(chan os.Signal, 1)
	signal.Notify(shutdownSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case signal := <-shutdownSignals:
		log.Printf("shutdown signal received: %s", signal)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			return err
		}

		log.Printf("%s stopped gracefully", cfg.AppName)
		return nil
	}
}
