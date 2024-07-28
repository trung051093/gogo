package main

import (
	"context"
	"fmt"
	"gogo/modules/tracker/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// create new server
	server := server.New()

	// Start the server asynchronously.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
				cancel()
			}
		}()
		server.Start()
	}()

	// Listen for interrupt signal or until the server stops.
	<-ctx.Done()
}
