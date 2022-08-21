package router

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
)

type Router interface {
	Listen(ctx context.Context, wg *sync.WaitGroup) error
	Connect(ctx context.Context, wg *sync.WaitGroup) error
	Forward(Router)

	GetWriteChan() chan []byte
	GetReadChan() chan []byte
	SetWriteChan(ch chan []byte)
	SetReadChan(ch chan []byte)
}

type RouterRun struct {
	Router Router
	Run    func(ctx context.Context, wg *sync.WaitGroup) error
}

type Registry struct {
	mutex   sync.RWMutex
	reg     map[string]*RouterRun
	RunList []string
}

func NewRegistry() *Registry {
	return &Registry{
		reg: make(map[string]*RouterRun),
	}
}

func (r *Registry) Add(name string, routerRun *RouterRun) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.reg[name] = routerRun

	r.RunList = append(r.RunList, name)
}

func (r *Registry) Get(name string) *RouterRun {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.reg[name]
}

func (r *Registry) Run(ctx context.Context, wg *sync.WaitGroup) error {
	for _, name := range r.RunList {
		log.Info().Msgf("[Router] Starting %s", name)
		routerRun := r.Get(name)
		if routerRun == nil {
			log.Warn().Msgf("[Router] Cannot find router %s to run", name)
			continue
		}

		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			if err := routerRun.Run(ctx, wg); err != nil {
				log.Error().Err(err).Msgf("[Router] Error while running router %s", name)
			}
		}(name)
	}

	return nil
}
