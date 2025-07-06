package netcode

import (
    "net/http"
    "testing"
    "time"

    "github.com/gorilla/websocket"
)

func TestNewServer(t *testing.T) {
	addr := "localhost:8080"
	s := NewServer(addr)
	if s == nil {
		t.Fatal("expected server, got nil")
	}
	if s.Address != addr {
		t.Fatalf("expected address %s, got %s", addr, s.Address)
	}
}

func TestServerPeerConnections(t *testing.T) {
    s := NewServer("127.0.0.1:0")
    if err := s.Start(); err != nil {
        t.Fatalf("failed to start server: %v", err)
    }
    defer s.Stop()

    addr := s.listener.Addr().String()
    dial := func() (*websocket.Conn, *http.Response, error) {
        return websocket.DefaultDialer.Dial("ws://"+addr+"/ws", nil)
    }

    var conns []*websocket.Conn
    for i := 0; i < 3; i++ {
        c, _, err := dial()
        if err != nil {
            t.Fatalf("dial %d: %v", i, err)
        }
        conns = append(conns, c)
        if err := c.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
            t.Fatalf("write %d: %v", i, err)
        }
        _, msg, err := c.ReadMessage()
        if err != nil {
            t.Fatalf("read %d: %v", i, err)
        }
        if string(msg) != "hello" {
            t.Fatalf("expected echo hello, got %s", msg)
        }
    }

    if c, _, err := dial(); err == nil {
        c.Close()
        t.Fatal("expected connection limit reached")
    }

    // allow goroutines to clean up
    time.Sleep(50 * time.Millisecond)
    for _, c := range conns {
        c.Close()
    }
}
