package router

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/rs/zerolog/log"
)

type TCP struct {
	host       string
	fWriteChan chan []byte
	fReadChan  chan []byte
}

func NewTCP(host string) *TCP {
	return &TCP{
		host:       host,
		fWriteChan: make(chan []byte, 256),
		fReadChan:  make(chan []byte, 256),
	}
}

func (r *TCP) GetWriteChan() chan []byte {
	return r.fWriteChan
}

func (r *TCP) GetReadChan() chan []byte {
	return r.fReadChan
}

func (r *TCP) SetWriteChan(ch chan []byte) {
	r.fWriteChan = ch
}

func (r *TCP) SetReadChan(ch chan []byte) {
	r.fReadChan = ch
}

func (r *TCP) Listen(ctx context.Context, wg *sync.WaitGroup) error {
	listener, err := net.Listen("tcp", r.host)
	if err != nil {
		return fmt.Errorf("unable to tcp listen on %s: %w", r.host, err)
	}

	log.Info().Msgf("[TCP] Listening on %s", r.host)

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		listener.Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			conn, err := listener.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					log.Error().Err(err).Msgf("unable to accept connection")
				}
				return
			}

			log.Debug().Msgf("Accepted LocalAddr: %v", conn.LocalAddr())

			wg.Add(1)
			go func(c net.Conn) {
				defer func() {
					c.Close()
					wg.Done()
				}()

				wgRW := sync.WaitGroup{}

				wgRW.Add(2)
				go ReadConn(ctx, &wgRW, c, &r.fWriteChan)
				go WriteConn(ctx, &wgRW, c, &r.fReadChan)

				wgRW.Wait()
			}(conn)
		}
	}()

	return nil
}

func (r *TCP) Connect(ctx context.Context, wg *sync.WaitGroup) error {
	var d net.Dialer
	conn, err := d.Dial("tcp", r.host)
	if err != nil {
		return fmt.Errorf("unable to tcp connect on %s: %w", r.host, err)
	}

	log.Info().Msgf("[TCP] Connected to %s", r.host)

	wg.Add(1)
	go func(c net.Conn) {
		defer func() {
			c.Close()
			wg.Done()
		}()

		wgRW := sync.WaitGroup{}

		wgRW.Add(2)
		go ReadConn(ctx, &wgRW, c, &r.fWriteChan)
		go WriteConn(ctx, &wgRW, c, &r.fReadChan)

		wgRW.Wait()
	}(conn)

	return nil
}

func (r *TCP) Forward(r2 Router) {
	r.SetWriteChan(r2.GetReadChan())
	r2.SetWriteChan(r.GetReadChan())
}
