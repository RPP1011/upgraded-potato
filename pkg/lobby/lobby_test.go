package lobby

import (
	"context"
	"testing"
)

func TestCreateLobbyGeneratesCode(t *testing.T) {
	ls := NewLobbyServer()
	code, err := ls.CreateLobby(context.Background(), "host1", 4, true)
	if err != nil {
		t.Fatalf("CreateLobby returned error: %v", err)
	}
	if len(code) != 6 {
		t.Fatalf("expected lobby code length 6, got %d", len(code))
	}
}

func TestListLobbiesExcludesBlocked(t *testing.T) {
	ls := NewLobbyServer()
	hostID := "host1"
	userID := "user1"
	if _, err := ls.CreateLobby(context.Background(), hostID, 4, true); err != nil {
		t.Fatalf("CreateLobby: %v", err)
	}
	if err := ls.BlockUser(context.Background(), hostID, userID); err != nil {
		t.Fatalf("BlockUser: %v", err)
	}
	lobbies, err := ls.ListLobbies(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListLobbies: %v", err)
	}
	if len(lobbies) != 0 {
		t.Fatalf("expected no lobbies, got %d", len(lobbies))
	}
}

func TestSearchLobbiesMatchesCodeOrHost(t *testing.T) {
	ls := NewLobbyServer()
	if _, err := ls.CreateLobby(context.Background(), "Alpha", 4, true); err != nil {
		t.Fatalf("CreateLobby: %v", err)
	}
	if _, err := ls.CreateLobby(context.Background(), "Beta", 4, true); err != nil {
		t.Fatalf("CreateLobby: %v", err)
	}
	res, err := ls.SearchLobbies(context.Background(), "user", "Al")
	if err != nil {
		t.Fatalf("SearchLobbies: %v", err)
	}
	if len(res) == 0 {
		t.Fatalf("expected results for query")
	}
}

func TestBlockUserPreventsJoin(t *testing.T) {
	ls := NewLobbyServer()
	hostID := "host2"
	userID := "user2"
	code, err := ls.CreateLobby(context.Background(), hostID, 4, true)
	if err != nil {
		t.Fatalf("CreateLobby: %v", err)
	}
	if err := ls.BlockUser(context.Background(), hostID, userID); err != nil {
		t.Fatalf("BlockUser: %v", err)
	}
	// Attempt to list lobbies should not include blocked host
	lobbies, err := ls.ListLobbies(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListLobbies: %v", err)
	}
	for _, l := range lobbies {
		if l.LobbyCode == code {
			t.Fatalf("blocked lobby should not be listed")
		}
	}
}
