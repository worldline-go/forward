package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/forward/cmd/forward/args"
	"github.com/worldline-go/logz"
)

func main() {
	logz.InitializeLog(nil)

	if err := args.Execute(context.Background()); err != nil {
		log.Error().Err(err).Msg("failed to execute")
		os.Exit(1)
	}
}
