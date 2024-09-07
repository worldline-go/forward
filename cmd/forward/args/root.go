package args

import (
	"context"
	"log/slog"

	"github.com/rakunlabs/into"
	"github.com/spf13/cobra"
	"github.com/worldline-go/forward/internal/config"
	"github.com/worldline-go/forward/internal/info"
	"github.com/worldline-go/forward/internal/server"
)

var rootCmd = &cobra.Command{
	Use:           "forward",
	Short:         "export connections",
	Long:          "export socket to HTTP",
	SilenceErrors: true,
	SilenceUsage:  true,
	Example: `  Multiple hosts and sockets:
    forward -H x@0.0.0.0:8080 -s x@/test.sock -H y@0.0.0.0:8081 -s y@/test2.sock
  Share docker socket with only GET method:
    forward -s /var/run/docker.sock:/:*,-POST,-PUT,-DELETE`,
	Version: "v0.0.0",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runForward(cmd.Context())
	},
}

// Execute is the entry point for the application.
func Execute(ctx context.Context) error {
	rootCmd.Version = info.AppInfo.Version
	rootCmd.Long += "\n" + longInfo()

	return rootCmd.ExecuteContext(ctx)
}

func longInfo() string {
	return "version: " + info.AppInfo.Version + " commit: " + info.AppInfo.BuildCommit + " buildDate:" + info.AppInfo.BuildDate
}

func init() {
	rootCmd.Flags().StringArrayVarP(&config.Application.Hosts, "host", "H", config.Application.Hosts, "Host to listen on")
	rootCmd.Flags().StringArrayVarP(&config.Application.Sockets, "socket", "s", config.Application.Sockets, "Socket to export: /var/run/docker.sock:/:*,-POST,-PUT,-DELETE")
}

func runForward(_ context.Context) error {
	slog.Info("forward " + longInfo())

	httpServer := server.ServeHTTP()

	into.ShutdownAdd(func() error {
		return server.StopHTTP(httpServer)
	}, into.WithShutdownName("http server"))

	return server.StartHTTP(httpServer)
}
