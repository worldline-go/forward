package router

import (
	"fmt"

	"github.com/worldline-go/forward/internal/config"
)

func Parse(r *config.Router, registry *Registry) error {
	for _, i := range r.TCPListen {
		router := NewTCP(i.Host)
		routerRun := &RouterRun{
			Router: router,
			Run:    router.Listen,
		}
		registry.Add(i.Name, routerRun)
	}

	for _, i := range r.TCPConnect {
		router := NewTCP(i.URL)
		routerRun := RouterRun{
			Router: router,
			Run:    router.Connect,
		}
		registry.Add(i.Name, &routerRun)
	}

	for _, i := range r.Forward {
		var fromRouter Router
		var toRouter Router

		if from := registry.Get(i.From); from != nil {
			fromRouter = from.Router
		}

		if to := registry.Get(i.To).Router; to != nil {
			toRouter = to
		}

		if fromRouter == nil || toRouter == nil {
			return fmt.Errorf("[Router] Cannot Forwarding from %s to %s", i.From, i.To)
		}

		fromRouter.Forward(toRouter)
	}

	return nil
}
