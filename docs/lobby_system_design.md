# Lobby System Design

This document outlines the design for a lobby management system built on top of
the existing websocket server. The goal is to allow clients to create lobbies,
list and search open lobbies, and block unwanted peers.

## Overview

- **Lobby Codes**: Each lobby has a short alphanumeric code (e.g., `AB12CD`).
- **Open Lobbies**: Public lobbies are discoverable via a list and search API.
- **Blocking**: Clients can block other users by ID to prevent them from joining
  their lobby or sending invites. Blocked clients should not see the blocking
  user's lobby when listing or searching for open lobbies.

The server maintains lobby state in memory and communicates with clients using
Protobuf messages over WebSocket.

## Dataflow

```
Client ---> Server: CreateLobbyRequest
Server ---> Client: CreateLobbyResponse (contains lobby code)

Client ---> Server: ListLobbiesRequest
Server ---> Client: ListLobbiesResponse (open lobbies summary)

Client ---> Server: SearchLobbiesRequest (query string)
Server ---> Client: SearchLobbiesResponse

Client ---> Server: BlockUserRequest (target user ID)
Server ---> Client: BlockUserResponse
```

1. **Create Lobby**: Client sends `CreateLobbyRequest` with lobby parameters.
   Server generates a unique code and stores the lobby record. The response
   includes the lobby code so others can join.
2. **List Open Lobbies**: Client requests a list of public lobbies. Server
   returns summaries for each lobby (code, host name, player count) excluding
   lobbies hosted by users who have blocked the requester.
3. **Search Lobbies**: Client supplies a query string. Server performs a
   case-insensitive search over lobby codes and host names and returns matches,
   again omitting lobbies where the host has blocked the requester.
4. **Block Users**: The server records blocked user IDs on a per-user basis.
   Blocked users are prevented from joining the blocking client's lobby and will
   not see that lobby in search or list responses.

## Message Types (Protobuf)

```protobuf
message CreateLobbyRequest {
  string host_id = 1;
  uint32 max_players = 2;
  bool public = 3;
}

message CreateLobbyResponse {
  string lobby_code = 1;
  bool   success = 2;
}

message ListLobbiesRequest {}

message LobbySummary {
  string lobby_code = 1;
  string host_id    = 2;
  uint32 players    = 3;
  uint32 max_players = 4;
}

message ListLobbiesResponse {
  repeated LobbySummary lobbies = 1;
}

message SearchLobbiesRequest {
  string query = 1; // code or host partial
}

message SearchLobbiesResponse {
  repeated LobbySummary results = 1;
}

message BlockUserRequest {
  string user_id = 1; // user to block
}

message BlockUserResponse {
  bool success = 1;
}
```

## Server Interfaces

The lobby subsystem exposes the following Go interfaces:

```go
// LobbyService manages lobby lifecycle and queries.
type LobbyService interface {
    CreateLobby(ctx context.Context, hostID string, maxPlayers int, public bool) (string, error)
    ListLobbies(ctx context.Context) ([]LobbySummary, error)
    SearchLobbies(ctx context.Context, query string) ([]LobbySummary, error)
    BlockUser(ctx context.Context, requesterID, targetID string) error
}
```

### Implementation Notes

- Lobby data is stored in memory using a `map[string]*Lobby` keyed by lobby
  code. Blocking rules are stored per user in a `map[string]map[string]bool`.
- Lobby codes are six-character uppercase strings generated using a secure
  random source.
- Searches run over open lobbies only and match if the code or host ID contains
  the query substring.
- List and search results exclude lobbies whose host has blocked the requesting
  user so blocked clients remain unaware of those lobbies.
- The server should periodically prune empty lobbies to free resources.

