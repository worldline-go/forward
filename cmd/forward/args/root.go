package args

import (
	"context"
	"os"
	"os/signal"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/forward/internal/config"
	"github.com/worldline-go/forward/internal/info"
	"github.com/worldline-go/forward/internal/server"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "forward",
	Short:   "export connections",
	Long:    "export socket to HTTP\n" + "version: " + info.AppInfo.Version + " commit: " + info.AppInfo.BuildCommit + " buildDate:" + info.AppInfo.BuildDate,
	Version: info.AppInfo.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.InitializeLog()
		config.SetLogLevel(config.Application.LogLevel)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runForward(cmd.Context())
	},
}

// Execute is the entry point for the application.
func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	rootCmd.Flags().StringVarP(&config.Application.Host, "host", "H", config.Application.Host, "Host to listen on, default: 0.0.0.0:8080")
	rootCmd.Flags().StringArrayVarP(&config.Application.Serve.Socket, "socket", "s", config.Application.Serve.Socket, "Socket to export: /var/run/docker.sock:/docker/:*,-POST,-PUT,-DELETE")
}

func runForward(ctx context.Context) error {
	chNotify := make(chan os.Signal, 1)
	signal.Notify(chNotify, os.Interrupt)
	defer func() {
		signal.Stop(chNotify)
		close(chNotify)
	}()

	httpServer := server.ServeHTTP()

	go func() {
		select {
		case <-ctx.Done():
			log.Error().Msg("unable continue process, shutting down...")
		case <-chNotify:
			log.Info().Msg("gracefully shutting down...")
		}

		server.StopHTTP(httpServer) //nolint:errcheck
	}()

	if err := server.StartHTTP(httpServer); err != nil {
		return err
	}

	return nil
}
