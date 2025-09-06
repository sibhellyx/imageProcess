package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sibhellyx/imageProccesor/config"
	"github.com/sibhellyx/imageProccesor/internal/server"
)

func main() {

	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	server := server.NewServer(ctx, cfg)
	go func() {
		<-sigChan
		server.Shutdown()
		cancel()
	}()

	server.Serve()
}
