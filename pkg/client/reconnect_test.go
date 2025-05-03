package client

import (
    "context"
    "github.com/rs/zerolog"
    "io"
    "net"
    "testing"
    "time"
)

var sink = zerolog.New(io.Discard)

// ---------------------------------------------------------------------------
//  1. positive path— first two dials fail, third succeeds after we start
//     listening on the chosen port
//
// ---------------------------------------------------------------------------
func TestDialWithBackoff_SucceedsAfterServerComesUp(t *testing.T) {

    lh, err := net.Listen("tcp", "127.0.0.1:0")
    if err != nil {
        t.Fatalf("cannot grab ephemeral port: %v", err)
    }
    addr := lh.Addr().String()
    lh.Close() // nobody is listening now

    pol := ReconnectPolicy{
        Enable:     true,
        MinDelay:   10 * time.Millisecond,
        MaxDelay:   20 * time.Millisecond,
        MaxRetries: 0, // unlimited
    }

    ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel()

    done := make(chan struct{})
    var c net.Conn
    go func() {
        var e error
        c, e = dialWithBackoff(ctx, addr, sink, pol)
        if e != nil {
            t.Errorf("dialWithBackoff returned unexpected error: %v", e)
        }
        close(done)
    }()

    time.Sleep(40 * time.Millisecond)

    server, err := net.Listen("tcp", addr)
    if err != nil {
        t.Fatalf("failed to start listener: %v", err)
    }
    defer server.Close()

    go func() {
        conn, _ := server.Accept()
        if conn != nil {
            _ = conn.Close()
        }
    }()

    select {
    case <-done:
        if c == nil {
            t.Fatal("dialWithBackoff returned nil connection")
        }
        _ = c.Close()
    case <-ctx.Done():
        t.Fatal("dialWithBackoff never succeeded before context deadline")
    }
}

// ---------------------------------------------------------------------------
// 2) negative path — stop after MaxRetries and surface the final error
// ---------------------------------------------------------------------------
func TestDialWithBackoff_StopsAfterMaxRetries(t *testing.T) {
    pol := ReconnectPolicy{
        Enable:     true,
        MinDelay:   5 * time.Millisecond,
        MaxDelay:   10 * time.Millisecond,
        MaxRetries: 3, // small so the test is fast
    }
    ctx := context.Background()

    _, err := dialWithBackoff(ctx, "127.0.0.1:1", sink, pol) // nobody there
    if err == nil {
        t.Fatal("expected error after exhausting retries, got nil")
    }
}
