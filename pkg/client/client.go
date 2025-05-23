package client

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/coreywagehoft/go-tak/pkg/cot"
	"github.com/coreywagehoft/go-tak/pkg/cotproto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	idleTimeout = 5 * time.Minute
	pingTimeout = time.Second * 15

	// SAMulticast UDP
	saMulticastAddr string = "239.2.3.1:6969"
	// ChatMultiCast UDP
	chatMultiCastAddr string = "224.10.10.1:17012"

	maxDatagramSize int = 8192
)

type Config struct {
	Host                 string
	Port                 int
	Logger               zerolog.Logger
	Reconnect            ReconnectPolicy
	SAMultiCastHandler   func(msg *cot.CotMessage)
	SAMultiCast          bool
	ChatMultiCastHandler func(msg *cot.CotMessage)
	ChatMultiCast        bool
}

type ReconnectPolicy struct {
	Enable     bool
	MinDelay   time.Duration
	MaxDelay   time.Duration
	MaxRetries int
}

type TakClient struct {
	Conn                 net.Conn
	SAMultiCastConn      net.Conn
	ChatMultiCastConn    net.Conn
	Logger               zerolog.Logger
	sendChan             chan []byte
	cancel               context.CancelFunc
	SAMultiCastHandler   func(msg *cot.CotMessage)
	ChatMultiCastHandler func(msg *cot.CotMessage)
	takVersion           int32
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
		Conn:                 conn,
		sendChan:             make(chan []byte, 50),
		SAMultiCastHandler:   config.SAMultiCastHandler,
		ChatMultiCastHandler: config.ChatMultiCastHandler,
		takVersion:           0,
	}

	if reflect.ValueOf(config.Logger).IsZero() {
		client.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().Timestamp().Logger()
	} else {
		client.Logger = config.Logger
	}

	ctx, client.cancel = context.WithCancel(ctx)

	go client.handleWrite()
	if config.SAMultiCast {
		go client.handleSARead(ctx)
	}

	if config.ChatMultiCast {
		go client.handleChatRead(ctx)
	}

	go client.pinger(ctx)
	return &client, nil
}

/* func (c *TakClient) NewSAMultiCastConn(ctx context.Context) error {
	if c.SAMultiCastConn != nil {
		return fmt.Errorf("multicast connection already exists")
	}

	addr, err := net.ResolveUDPAddr("udp", saMulticastAddr)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to resolve multicast address")
	}

	conn, err := net.ListenPacket("udp", addr.String())
	if err != nil {
		c.Logger.Error().Err(err).Msg("Error listening on UDP")

		return err
	}
	defer conn.Close()

	p := ipv4.NewPacketConn(conn)
	if err := p.JoinGroup(nil, addr); err != nil {
		c.Logger.Error().Err(err).Msg("Error joining multicast group")

		return err
	}
	if err := p.SetMulticastLoopback(false); err != nil {
		c.Logger.Error().Err(err).Msg("Error setting multicast loopback")
	}

	c.SAMultiCastConn = conn

		go func() {
			defer conn.Close()

			for {
				buf := make([]byte, maxDatagramSize)
				n, _, err := conn.ReadFrom(buf)
				if err != nil {
					c.Logger.Error().Err(err).Msg("error reading from multicast connection")
					break
				}

				c.Logger.Debug().Msg(fmt.Sprintf("received %d bytes from multicast connection", n))
			}
		}()

	return nil
} */

func (c *TakClient) NewSAMultiCastConn() error {
	if c.SAMultiCastConn != nil {
		return fmt.Errorf("multicast connection already exists")
	}

	addr, err := net.ResolveUDPAddr("udp", saMulticastAddr)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to resolve multicast address")
	}

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		c.Logger.Error().Err(err).Msg("Error listening on UDP")

		return err
	}

	conn.SetReadBuffer(maxDatagramSize)

	c.SAMultiCastConn = conn

	return nil
}

func (c *TakClient) NewChatMultiCastConn() error {
	if c.ChatMultiCastConn != nil {
		return fmt.Errorf("multicast connection already exists")
	}

	addr, err := net.ResolveUDPAddr("udp", chatMultiCastAddr)
	if err != nil {
		c.Logger.Error().Err(err).Msg("failed to resolve multicast address")
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		c.Logger.Error().Err(err).Msg("Error listening on UDP")

		return err
	}

	conn.SetReadBuffer(maxDatagramSize)

	c.ChatMultiCastConn = conn

	return nil
}

func (c *TakClient) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}

	if c.SAMultiCastConn != nil {
		return c.SAMultiCastConn.Close()
	}

	return nil
}

func (c *TakClient) SetVersion(n int32) {
	atomic.StoreInt32(&c.takVersion, n)
}

func (c *TakClient) GetVersion() int32 {
	return atomic.LoadInt32(&c.takVersion)
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

func (c *TakClient) handleSARead(ctx context.Context) {
	defer c.Stop()

	er := cot.NewTagReader(c.SAMultiCastConn)
	pr := cot.NewProtoReader(c.SAMultiCastConn)

	for ctx.Err() == nil {
		var msg *cot.CotMessage

		var err error

		switch c.GetVersion() {
		case 0:
			msg, err = c.processXMLRead(er)
		case 1:
			msg, err = c.processProtoRead(pr)
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				c.Logger.Info().Err(err).Msg("EOF")

				break
			}

			c.Logger.Warn().Err(err).Msg("error")

			break
		}

		if msg == nil {
			continue
		}

		// ping
		if msg.GetType() == "t-x-c-t" {
			c.Logger.Info().Str("ping.from", msg.GetUID()).Msg("ping")

			if err := c.SendCot(cot.MakePong()); err != nil {
				c.Logger.Error().Err(err).Msg("SendMsg error")
			}

			continue
		}

		// pong
		if msg.GetType() == "t-x-c-t-r" {
			continue
		}

		c.Logger.Info().Msg(fmt.Sprintf("msg: %s", msg))

		c.SAMultiCastHandler(msg)
	}
}

func (c *TakClient) handleChatRead(ctx context.Context) {
	defer c.Stop()

	er := cot.NewTagReader(c.ChatMultiCastConn)
	pr := cot.NewProtoReader(c.ChatMultiCastConn)

	for ctx.Err() == nil {
		var msg *cot.CotMessage

		var err error

		switch c.GetVersion() {
		case 0:
			msg, err = c.processXMLRead(er)
		case 1:
			msg, err = c.processProtoRead(pr)
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				c.Logger.Info().Err(err).Msg("EOF")

				break
			}

			c.Logger.Warn().Err(err).Msg("error")

			break
		}

		if msg == nil {
			continue
		}

		// ping
		if msg.GetType() == "t-x-c-t" {
			c.Logger.Info().Str("ping.from", msg.GetUID()).Msg("ping")

			if err := c.SendCot(cot.MakePong()); err != nil {
				c.Logger.Error().Err(err).Msg("SendMsg error")
			}

			continue
		}

		// pong
		if msg.GetType() == "t-x-c-t-r" {
			continue
		}

		c.Logger.Info().Msg(fmt.Sprintf("msg: %s", msg))

		c.ChatMultiCastHandler(msg)
	}
}

//nolint:nilnil
func (c *TakClient) processXMLRead(er *cot.TagReader) (*cot.CotMessage, error) {
	tag, dat, err := er.ReadTag()
	if err != nil {
		return nil, err
	}

	if tag == "?xml" {
		return nil, nil
	}

	if tag == "auth" {
		// <auth><cot username=\"test\" password=\"111111\" uid=\"ANDROID-xxxx\ callsign=\"zzz\""/></auth>
		return nil, nil
	}

	/* if tag != "event" {
		log.Error().Str("tag", tag).Str("data", string(dat[:])).Msg("bad tag")
		return nil, fmt.Errorf("bad tag: %s", dat)
	} */

	log.Error().Str("tag", tag).Str("data", string(dat[:])).Msg("bad tag")

	ev := new(cot.Event)
	if err := xml.Unmarshal(dat, ev); err != nil {
		return nil, fmt.Errorf("xml decode error: %w, client: %s", err, string(dat))
	}

	// h.setActivity()

	c.Logger.Debug().Str("xml.event", string(dat)).Msg("decoded")

	if ev.Type == "t-x-takp-q" {
		ver := ev.Detail.GetFirst("TakControl").GetFirst("TakRequest").GetAttr("version")
		if ver == "1" {
			if err := c.sendEvent(cot.ProtoChangeOkMsg()); err == nil {
				c.Logger.Info().Msg(fmt.Sprintf("client switch to v.1"))
				c.SetVersion(1)

				return nil, nil
			}

			return nil, fmt.Errorf("error on send ok: %w", err)
		}
	}

	if ev.Type == "t-x-takp-v" {
		if ps := ev.Detail.GetFirst("TakControl").GetFirst("TakProtocolSupport"); ps != nil {
			v := ps.GetAttr("version")
			c.Logger.Info().Msg("server supports protocol v" + v)

			if v == "1" {
				c.Logger.Debug().Msg("sending v1 req")
				_ = c.sendEvent(cot.VersionReqMsg(1))
			}
		} else {
			c.Logger.Warn().Msg("invalid protocol support message: " + string(dat))
		}

		return nil, nil
	}

	if ev.Type == "t-x-takp-r" {
		if n := ev.Detail.GetFirst("TakControl").GetFirst("TakResponse"); n != nil {
			status := n.GetAttr("status")
			c.Logger.Info().Msg(fmt.Sprintf("server switches to v1: %s", status))

			if status == "true" {
				c.SetVersion(1)
			} else {
				c.Logger.Error().Msg(fmt.Sprintf("got TakResponce with status %s: %s", status, ev.Detail))
			}
		}

		return nil, nil
	}

	return cot.EventToProto(ev)
}

func (c *TakClient) processProtoRead(r *cot.ProtoReader) (*cot.CotMessage, error) {
	msg, err := r.ReadProtoBuf()
	if err != nil {
		return nil, err
	}

	// h.setActivity()

	var d *cot.Node
	d, err = cot.DetailsFromString(msg.GetCotEvent().GetDetail().GetXmlDetail())

	c.Logger.Debug().Msg(fmt.Sprintf("proto msg: %s", msg))

	return &cot.CotMessage{TakMessage: msg, Detail: d}, err
}

func (c *TakClient) sendEvent(evt *cot.Event) error {
	if c.GetVersion() != 0 {
		return fmt.Errorf("bad client version")
	}

	msg, err := xml.Marshal(evt)
	if err != nil {
		return err
	}

	c.Logger.Debug().Msg("sending" + string(msg))

	// Send the message to the server
	c.sendChan <- msg

	return nil
}

func (c *TakClient) SendCot(msg *cotproto.TakMessage) error {
	if c.Conn == nil {
		return fmt.Errorf("connection is not established")
	}

	if msg == nil {
		return fmt.Errorf("message cannot be nil")
	}

	switch c.GetVersion() {
	case 0:
		c.Logger.Debug().Interface("TakMessage", &msg).Msg("Sending message")

		// Convert the message to bytes
		buf, err := xml.Marshal(cot.ProtoToEvent(msg))
		if err != nil {
			return err
		}

		// Send the message to the server
		c.sendChan <- buf

		return nil
	case 1:
		c.Logger.Debug().Interface("TakMessage", &msg).Msg("Sending message")

		buf, err := cot.MakeProtoPacket(msg)
		if err != nil {
			return err
		}

		// Send the message to the server
		c.sendChan <- buf
	}

	return fmt.Errorf("unknown version: %d", c.GetVersion())
}

func (c *TakClient) Stop() {

	if c.Conn != nil {
		c.Conn.Close()
	}

	if c.SAMultiCastConn != nil {
		c.SAMultiCastConn.Close()
	}

	close(c.sendChan)
}
