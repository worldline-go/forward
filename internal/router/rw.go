package router

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/rs/zerolog/log"
)

func ReadConn(ctx context.Context, wg *sync.WaitGroup, c io.ReadWriter, ch *chan []byte) error {
	defer wg.Done()

	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			n, err := c.Read(buf)
			if err != nil {
				if err != io.EOF {
					return fmt.Errorf("unable to read from connection: %w", err)
				}
				log.Debug().Msgf("Returning r")
				return nil
			}
			log.Debug().Msgf("Read con %s", string(buf[:n]))

			// send to another area
			select {
			case *ch <- buf[:n]:
			default:
				log.Debug().Msgf("Returning rd")
				return fmt.Errorf("unable to send to channel")
			}
		}
	}
}

func WriteConn(ctx context.Context, wg *sync.WaitGroup, c io.ReadWriter, ch *chan []byte) error {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return nil
		case buf := <-*ch:
			log.Debug().Msgf("Write con %s", string(buf))
			_, err := c.Write(buf)
			if err != nil {
				log.Debug().Err(err).Msgf("Returning w")
				if err != io.EOF {
					// return fmt.Errorf("unable to read from connection: %w", err)
				}
				// return nil
			}
		}
	}
}
