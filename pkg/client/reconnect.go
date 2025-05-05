package client

import (
	"context"
	"github.com/rs/zerolog"
	"math/rand"
	"net"
	"time"
)

func dialWithBackoff(ctx context.Context, addr string,
	log zerolog.Logger,
	p ReconnectPolicy) (net.Conn, error) {

	if p.MinDelay == 0 {
		p.MinDelay = time.Second
	}
	if p.MaxDelay == 0 {
		p.MaxDelay = 30 * time.Second
	}

	attempt := 0
	delay := p.MinDelay

	for {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			if attempt > 0 {
				log.Info().
					Int("attempt", attempt).
					Str("addr", addr).
					Msg("reconnected")
			}
			return conn, nil
		}

		if !p.Enable ||
			(p.MaxRetries > 0 && attempt >= p.MaxRetries) {
			return nil, err
		}

		// jittered exponential back‑off
		sleep := delay + time.Duration(rand.Int63n(int64(delay/2)))
		log.Warn().
			Err(err).
			Str("addr", addr).
			Dur("sleep", sleep).
			Int("attempt", attempt+1).
			Msg("dial failed – retrying")
		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		if delay *= 2; delay > p.MaxDelay {
			delay = p.MaxDelay
		}
		attempt++
	}
}
