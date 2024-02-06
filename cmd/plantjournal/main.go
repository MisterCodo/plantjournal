package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MisterCodo/plantjournal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config *server.Config
	addr   string
	port   int

	rootCmd = &cobra.Command{
		Use:   "plant",
		Short: "Self-hosted plant journal to track plant care.",
		Run: func(_cmd *cobra.Command, _args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			// Create plant journal server.
			s, err := server.NewServer(ctx, logger, config)
			if err != nil {
				cancel()
				logger.Error("failed to create server", "err", err)
			}

			// Listen for terminal signal and handle graceful shutdown of plant journal server.
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-c
				s.Shutdown(ctx)
				cancel()
			}()

			// Start plant journal server.
			err = s.Start(ctx)
			if err != nil && err != http.ErrServerClosed {
				logger.Error("failed to start server", "err", err)
				cancel()
			}

			// Wait.
			<-ctx.Done()
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&addr, "addr", "a", "", "address of server")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "port of server")

	err := viper.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	if err != nil {
		panic(err)
	}

	viper.SetDefault("addr", "")
	viper.SetDefault("port", 8080)
	viper.SetEnvPrefix("plant")
}

func initConfig() {
	viper.AutomaticEnv()
	config = &server.Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
