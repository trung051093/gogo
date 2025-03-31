package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	server "gogo/modules/snake_battle/server"
)

// Inspired by https://github.com/raeperd/kickstart.go/blob/main/main.go
func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Create a new context to listen for SIGINT and SIGTERM signals for graceful shutdown
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(log)

	srv := server.New()
	// Use a channel to capture any server errors
	errChan := make(chan error, 1)

	// Start the server in a goroutine to allow the main goroutine to listen for signals
	go func() {
		srv.Start()
	}()

	select {
	// Server exited with an error
	case err := <-errChan:
		log.ErrorContext(ctx, "server exited with error",
			slog.Any("error", err),
		)
		return err
	// Received a signal, shutdown the server
	case <-ctx.Done():
		slog.InfoContext(ctx, "shutting down server")

		// Create a new context with timeout for the shutdown process
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		// Shutdown the HTTP server first
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}

		// After server is shut down, cancel the main context to close all resources
		cancel()

		// Add post-shutdown cleanup here if necessary, for example:
		// - Close queue listeners
		// - Close database connection

		return nil
	}
}
