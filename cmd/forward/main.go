package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/forward/cmd/forward/args"
	"github.com/worldline-go/forward/internal/info"
)

var (
	// Populated by goreleaser during build
	version = "v0.0.0"
	commit  = "?"
	date    = ""
)

func main() {
	// change information
	info.AppInfo.Version = version
	info.AppInfo.BuildCommit = commit
	info.AppInfo.BuildDate = date

	if err := args.Execute(context.Background()); err != nil {
		log.Error().Err(err).Msg("failed to execute")
		os.Exit(1)
	}
}
