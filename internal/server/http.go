package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/forward/internal/config"
	"github.com/worldline-go/forward/internal/handler"
	"github.com/worldline-go/forward/internal/server/mux"
)

var serverShutdownTimeout = 5 * time.Second

// ServeHTTP returns a new HTTP server.
func ServeHTTP() *http.Server {
	mux := &mux.Mux{}

	// Add handler functions here
	handler.SocketParser(config.Application.Serve.Socket, mux.HandleFunc)

	// Not found
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found")) //nolint:errcheck
	})

	s := &http.Server{
		Addr:           config.Application.Host,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
	}

	return s
}

// StartHTTP starts the HTTP server.
func StartHTTP(server *http.Server) error {
	if server == nil {
		return nil
	}

	log.Info().Msgf("server on [%s]", config.Application.Host)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	log.Info().Msg("server stopped gracefully")

	return nil
}

// StopHTTP stops the HTTP server.
func StopHTTP(server *http.Server) error {
	if server == nil {
		return nil
	}

	log.Info().Msg("server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
