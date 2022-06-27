package handler

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/rs/zerolog/log"
)

// SocketDirector is a reverse proxy director that rewrites the URL path to
func SocketDirector(req *http.Request) {
	req.Header.Add("X-Forwarded-Host", req.Host)
	req.URL.Scheme = "http"
	req.URL.Host = "localhost"
}

// SocketHandler returns a handler that proxies the request to the given path.
func SocketHandler(socketURLPath string, socketPath string, methods *FilterMethods) func(http.ResponseWriter, *http.Request) {
	sockClient := &httputil.ReverseProxy{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
				con, err := new(net.Dialer).DialContext(ctx, "unix", socketPath)
				if err != nil {
					return nil, err
				}

				return con, nil
			},
		},
		Director: SocketDirector,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if methods.Match(r.Method) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, socketURLPath)
			if r.URL.Path[0] != '/' {
				r.URL.Path = "/" + r.URL.Path
			}

			sockClient.ServeHTTP(w, r)

			return
		}

		// TODO: return allowed methods
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed")) //nolint:errcheck
	}
}

// SocketParser parse socketconfig and record in the mux.
func SocketParser(socketConfigs []string, handlerFunc func(string, func(http.ResponseWriter, *http.Request))) {
	// /var/run/docker.sock:/docker/:*,-POST,-PUT,-DELETE
	for _, socketConfig := range socketConfigs {
		socketList := strings.Split(socketConfig, ":")

		if len(socketList) < 2 {
			continue
		}

		socketPath := socketList[0]
		socketURLPath := socketList[1]

		socketMethods := NewFilterMethods()
		if len(socketList) > 2 {
			socketMethods.Parse(strings.Split(socketList[2], ","))
		}

		log.Info().Msgf("socket route [%s] to [%s]; %s", socketURLPath, socketPath, socketMethods)

		handlerFunc(socketURLPath, SocketHandler(socketURLPath, socketPath, socketMethods))
	}
}
