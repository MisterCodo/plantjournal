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
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			// Create plant journal server
			s, err := server.NewServer(ctx, logger)
			if err != nil {
				cancel()
				logger.Error("failed to create server", "err", err)
			}

			// Listen for terminal signal and handle graceful shutdown of plant journal server
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-c
				s.Shutdown(ctx)
				cancel()
			}()

			// Start plant journal server
			err = s.Start(ctx)
			if err != nil && err != http.ErrServerClosed {
				logger.Error("failed to start server", "err", err)
				cancel()
			}

			// Wait
			<-ctx.Done()
		},
	}
)

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
