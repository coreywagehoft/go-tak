package client

import (
	"context"
	"encoding/xml"
	"fmt"
	"net"
	"time"

	"github.com/coreywagehoft/go-tak/pkg/cot"
	"github.com/coreywagehoft/go-tak/pkg/cotproto"
	"github.com/rs/zerolog"
)

const (
	idleTimeout = 5 * time.Minute
	pingTimeout = time.Second * 15
)

type TakClient struct {
	Conn     net.Conn
	Logger   zerolog.Logger
	sendChan chan []byte
	cancel   context.CancelFunc
}

func NewTakClient(ctx context.Context, host string, port int) (*TakClient, error) {

	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	if port <= 0 || port > 65535 {
		return nil, fmt.Errorf("port must be between 1 and 65535")
	}

	// Connect to the server
	conn, err := net.Dial("tcp", net.JoinHostPort(host, fmt.Sprint(port)))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	client := &TakClient{
		Conn:     conn,
		sendChan: make(chan []byte, 50),
	}

	ctx, client.cancel = context.WithCancel(ctx)

	go client.handleWrite()
	go client.pinger(ctx)

	return client, nil
}

func (c *TakClient) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}

func (c *TakClient) pinger(ctx context.Context) {
	ticker := time.NewTicker(pingTimeout)
	defer ticker.Stop()

	for ctx.Err() == nil {
		select {
		case <-ticker.C:
			// TODO ADD LOGGER
			c.Logger.Debug().Msg("Sending ping")

			if err := c.SendCot(cot.MakePing("go-tak-client")); err != nil {
				c.Logger.Error().Err(err).Msg("sendMsg error")
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *TakClient) handleWrite() {
	for msg := range c.sendChan {
		if _, err := c.Conn.Write(msg); err != nil {
			// TODO LOGGER
			c.Stop()

			break
		}
	}
}

func (c *TakClient) SendCot(msg *cotproto.TakMessage) error {
	if c.Conn == nil {
		return fmt.Errorf("connection is not established")
	}

	if msg == nil {
		return fmt.Errorf("message cannot be nil")
	}

	c.Logger.Debug().Interface("TakMessage", &msg).Msg("Sending message")

	// Convert the message to bytes
	buf, err := xml.Marshal(cot.ProtoToEvent(msg))
	if err != nil {
		return err
	}

	// Send the message to the server
	c.sendChan <- buf

	return nil
}

func (c *TakClient) Stop() {

	if c.Conn != nil {
		c.Conn.Close()
	}

	close(c.sendChan)
}
