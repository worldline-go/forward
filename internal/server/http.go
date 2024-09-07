package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/worldline-go/forward/internal/config"
	"github.com/worldline-go/forward/internal/handler"
	"golang.org/x/sync/errgroup"
)

var serverShutdownTimeout = 5 * time.Second

type Server struct {
	Server *http.Server
	Name   string
}

// ServeHTTP returns a new HTTP server.
func ServeHTTP() []Server {
	values := Parse(config.Application.Hosts, config.Application.Sockets)
	var servers []Server

	for _, value := range values {
		mux := http.NewServeMux()
		// Add handler functions here
		handler.SocketParser(value.Socket, mux.HandleFunc)

		s := &http.Server{
			Addr:    value.Host,
			Handler: mux,
		}

		servers = append(servers, Server{
			Server: s,
			Name:   value.Name,
		})
	}

	return servers
}

// StartHTTP starts the HTTP server.
func StartHTTP(server []Server) error {
	if server == nil {
		return nil
	}

	group := errgroup.Group{}

	for _, s := range server {
		group.Go(func() error {
			slog.Info(fmt.Sprintf("server %s on [%s]", s.Name, s.Server.Addr))
			if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

// StopHTTP stops the HTTP server.
func StopHTTP(server []Server) error {
	slog.Info("server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	wg := sync.WaitGroup{}

	for _, s := range server {
		if s.Server == nil {
			continue
		}

		wg.Add(1)
		go func(name string, s *http.Server) {
			defer wg.Done()

			if err := s.Shutdown(ctx); err != nil {
				slog.Error("server shutdown error", slog.String("error", err.Error()), slog.String("server", name))
			}
		}(s.Name, s.Server)
	}

	return nil
}
