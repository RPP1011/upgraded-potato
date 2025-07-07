package lobby

import (
	"context"
	"errors"
)

var ErrNotImplemented = errors.New("not implemented")

// LobbyService manages lobby lifecycle and queries.
type LobbyService interface {
	CreateLobby(ctx context.Context, hostID string, maxPlayers int, public bool) (string, error)
	ListLobbies(ctx context.Context, requesterID string) ([]*LobbySummary, error)
	SearchLobbies(ctx context.Context, requesterID, query string) ([]*LobbySummary, error)
	BlockUser(ctx context.Context, requesterID, targetID string) error
}

type LobbyServer struct{}

func NewLobbyServer() *LobbyServer {
	return &LobbyServer{}
}

func (ls *LobbyServer) CreateLobby(ctx context.Context, hostID string, maxPlayers int, public bool) (string, error) {
	return "", ErrNotImplemented
}

func (ls *LobbyServer) ListLobbies(ctx context.Context, requesterID string) ([]*LobbySummary, error) {
	return nil, ErrNotImplemented
}

func (ls *LobbyServer) SearchLobbies(ctx context.Context, requesterID, query string) ([]*LobbySummary, error) {
	return nil, ErrNotImplemented
}

func (ls *LobbyServer) BlockUser(ctx context.Context, requesterID, targetID string) error {
	return ErrNotImplemented
}
