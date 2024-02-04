package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"plantjournal/server"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "plant",
		Short: "Self-hosted plant journal to track plant care.",
		Run: func(_cmd *cobra.Command, _args []string) {
			ctx, cancel := context.WithCancel(context.Background())

			// Create plant journal server
			s, err := server.NewServer(ctx)
			if err != nil {
				cancel()
				slog.Error("failed to create server", "err", err)
			}

			// Listen for terminal signal and handle graceful shutdown of plant journal server
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-c
				slog.Info("closing server")
				s.Shutdown(ctx)
				cancel()
			}()

			// Start plant journal server
			err = s.Start(ctx)
			if err != nil && err != http.ErrServerClosed {
				slog.Error("failed to start server")
				cancel()
			}

			// Wait
			<-ctx.Done()
		},
	}
)

func main() {
	slog.Info("Plant Journal")

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
