package client

import (
	"fmt"
	"net"
)

type TakClient struct {
	Conn net.Conn
}

func NewTakClient(host string, port int) (*TakClient, error) {

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
		Conn: conn,
	}

	return client, nil
}

func (c *TakClient) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}

func (c *TakClient) SendMessage(message []byte) error {
	if c.Conn == nil {
		return fmt.Errorf("connection is not established")
	}

	// Send the message to the server
	_, err := c.Conn.Write(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}