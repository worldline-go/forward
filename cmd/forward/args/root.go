package args

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/worldline-go/forward/internal/config"
	"github.com/worldline-go/forward/internal/router"
	"github.com/worldline-go/igconfig"
	"github.com/worldline-go/igconfig/loader"
	"github.com/worldline-go/logz"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "forward",
	Short:   "export connections",
	Long:    "export socket to HTTP",
	Version: config.AppInfo.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := logz.SetLogLevel(config.Application.LogLevel); err != nil {
			log.Warn().Err(err).Msg("unable to set log level")
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		logConfig := log.With().Str("component", "config").Logger()
		ctxConfig := logConfig.WithContext(cmd.Context())

		loaders := []loader.Loader{
			&loader.Consul{},
			&loader.Vault{},
			&loader.File{},
			&loader.Env{},
		}

		if err := igconfig.LoadWithLoadersWithContext(ctxConfig, "", &config.Application, loaders[3]); err != nil {
			return fmt.Errorf("unable to load prefix settings: %v", err)
		}

		loader.ConsulConfigPathPrefix = config.LoadConfig.Prefix.Consul
		loader.VaultSecretBasePath = config.LoadConfig.Prefix.Vault

		if err := igconfig.LoadWithLoadersWithContext(ctxConfig, config.LoadConfig.AppName, &config.Application, loaders...); err != nil {
			return fmt.Errorf("unable to load configuration settings: %v", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runForward(cmd.Context())
	},
}

// Execute is the entry point for the application.
func Execute(ctx context.Context) error {
	rootCmd.Version = config.AppInfo.Version
	rootCmd.Long += "\nversion: " + config.AppInfo.Version + " commit: " + config.AppInfo.BuildCommit + " buildDate:" + config.AppInfo.BuildDate
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	rootCmd.Flags().StringVarP(&config.Application.LogLevel, "log-level", "l", config.Application.LogLevel, "Log level: debug, info, warn, error")
}

func runForward(ctxParent context.Context) error {
	log.WithLevel(zerolog.NoLevel).Msgf("forward [%s]", config.AppInfo.Version)

	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

	ctx, cancel := context.WithCancel(ctxParent)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			log.Error().Msg("unable continue process, shutting down...")
		case <-sig:
			log.Info().Msg("gracefully shutting down...")
			cancel()
		}
	}()

	// set registry
	registry := router.NewRegistry()

	// parse config
	if err := router.Parse(&config.Application.Router, registry); err != nil {
		return err
	}

	// start routers
	if err := registry.Run(ctx, wg); err != nil {
		return err
	}

	// routerTCP1 := router.NewTCP("localhost:8080")
	// if err := routerTCP1.Connect(ctx, wg); err != nil {
	// 	return err
	// }

	// routerTCP2 := router.NewTCP("localhost:9090")
	// if err := routerTCP2.Listen(ctx, wg); err != nil {
	// 	return err
	// }

	// routerTCP1.Forward(routerTCP2)

	wg.Wait()

	return nil
}
