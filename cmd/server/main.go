package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nextlag/in_memory_db/internal/server/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := a.Run(ctx); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
