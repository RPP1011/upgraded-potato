package netcode

import "testing"

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
