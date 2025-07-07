package lobby

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
	"sync"
)

var ErrNotImplemented = errors.New("not implemented")

// LobbyService manages lobby lifecycle and queries.
type LobbyService interface {
	CreateLobby(ctx context.Context, hostID string, maxPlayers int, public bool) (string, error)
	ListLobbies(ctx context.Context, requesterID string) ([]*LobbySummary, error)
	SearchLobbies(ctx context.Context, requesterID, query string) ([]*LobbySummary, error)
	BlockUser(ctx context.Context, requesterID, targetID string) error
}

type Lobby struct {
	Code       string
	HostID     string
	MaxPlayers int
	Public     bool
	Players    int
}

type LobbyServer struct {
	mu      sync.Mutex
	lobbies map[string]*Lobby
	blocked map[string]map[string]bool
}

func NewLobbyServer() *LobbyServer {
	return &LobbyServer{
		lobbies: make(map[string]*Lobby),
		blocked: make(map[string]map[string]bool),
	}
}

var codeChars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomCode() (string, error) {
	b := make([]rune, 6)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeChars))))
		if err != nil {
			return "", err
		}
		b[i] = codeChars[n.Int64()]
	}
	return string(b), nil
}

func (ls *LobbyServer) CreateLobby(ctx context.Context, hostID string, maxPlayers int, public bool) (string, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	var code string
	var err error
	for {
		code, err = randomCode()
		if err != nil {
			return "", err
		}
		if _, exists := ls.lobbies[code]; !exists {
			break
		}
	}

	ls.lobbies[code] = &Lobby{
		Code:       code,
		HostID:     hostID,
		MaxPlayers: maxPlayers,
		Public:     public,
		Players:    1,
	}

	return code, nil
}

func (ls *LobbyServer) ListLobbies(ctx context.Context, requesterID string) ([]*LobbySummary, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	var res []*LobbySummary
	for _, l := range ls.lobbies {
		if !l.Public {
			continue
		}
		if blocked, ok := ls.blocked[l.HostID]; ok && blocked[requesterID] {
			continue
		}
		res = append(res, &LobbySummary{
			LobbyCode:  l.Code,
			HostId:     l.HostID,
			Players:    uint32(l.Players),
			MaxPlayers: uint32(l.MaxPlayers),
		})
	}
	return res, nil
}

func (ls *LobbyServer) SearchLobbies(ctx context.Context, requesterID, query string) ([]*LobbySummary, error) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	q := strings.ToLower(query)
	var res []*LobbySummary
	for _, l := range ls.lobbies {
		if !l.Public {
			continue
		}
		if blocked, ok := ls.blocked[l.HostID]; ok && blocked[requesterID] {
			continue
		}
		if strings.Contains(strings.ToLower(l.Code), q) || strings.Contains(strings.ToLower(l.HostID), q) {
			res = append(res, &LobbySummary{
				LobbyCode:  l.Code,
				HostId:     l.HostID,
				Players:    uint32(l.Players),
				MaxPlayers: uint32(l.MaxPlayers),
			})
		}
	}
	return res, nil
}

func (ls *LobbyServer) BlockUser(ctx context.Context, requesterID, targetID string) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	m, ok := ls.blocked[requesterID]
	if !ok {
		m = make(map[string]bool)
		ls.blocked[requesterID] = m
	}
	m[targetID] = true
	return nil
}
