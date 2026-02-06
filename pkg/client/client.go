package client

import (
	"context"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"reflect"
	"time"

	"github.com/coreywagehoft/go-tak/pkg/cot"
	"github.com/coreywagehoft/go-tak/pkg/cotproto"
	"github.com/rs/zerolog"
)

const (
	idleTimeout = 5 * time.Minute
	pingTimeout = time.Second * 15
)

type Config struct {
	Host      string
	Port      int
	Logger    zerolog.Logger
	Reconnect ReconnectPolicy
}

type ReconnectPolicy struct {
	Enable     bool
	MinDelay   time.Duration
	MaxDelay   time.Duration
	MaxRetries int
}

type TakClient struct {
	Conn     net.Conn
	Logger   zerolog.Logger
	sendChan chan []byte
	cancel   context.CancelFunc
}

func NewTakClient(ctx context.Context, config Config) (*TakClient, error) {

	if config.Host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	if config.Port <= 0 || config.Port > 65535 {
		return nil, fmt.Errorf("port must be between 1 and 65535")
	}

	addr := net.JoinHostPort(config.Host, fmt.Sprint(config.Port))
	conn, err := dialWithBackoff(ctx, addr, config.Logger, config.Reconnect)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	client := TakClient{
		Conn:     conn,
		sendChan: make(chan []byte, 50),
	}

	if reflect.ValueOf(config.Logger).IsZero() {
		client.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().Timestamp().Logger()
	} else {
		client.Logger = config.Logger
	}

	ctx, client.cancel = context.WithCancel(ctx)

	go client.handleWrite()
	go client.pinger(ctx)
	return &client, nil
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
			c.Logger.Debug().
				Str("remote", c.Conn.RemoteAddr().String()).
				Dur("interval", pingTimeout).
				Msg("sending ping")

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

			c.Logger.Error().
				Err(err).
				Int("bytes", len(msg)).
				Str("remote", c.Conn.RemoteAddr().String()).
				Msg("TCP write failed, stopping client")

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

	// Format according to Traditional Protocol spec:
	// <?xml version='1.0' encoding='UTF-8' standalone='yes'?>\n<event>...</event>
	fullMsg := []byte(xml.Header)
	fullMsg = append(fullMsg, buf...)

	// Send the message to the server
	c.sendChan <- fullMsg

	return nil
}

func (c *TakClient) Stop() {

	if c.Conn != nil {
		c.Conn.Close()
	}

	close(c.sendChan)
}
