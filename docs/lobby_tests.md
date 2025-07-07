# Lobby Service Test Plan

This document captures the intended behavior of the lobby service as expressed in the unit tests. The tests are currently failing because the implementation is stubbed out.

## Goals

- Validate that creating a lobby generates a six character code.
- Ensure lobbies hosted by users who block the requester do not appear in list results.
- Support searching lobbies by code or host ID.
- Prevent blocked users from joining or seeing lobbies of the host who blocked them.

## Test Summaries

### `TestCreateLobbyGeneratesCode`

- Call `CreateLobby` with a host ID and options.
- Expect a 6 character lobby code in the result.

### `TestListLobbiesExcludesBlocked`

- Host creates a lobby then blocks another user.
- When the blocked user lists lobbies, the host's lobby should not appear.

### `TestSearchLobbiesMatchesCodeOrHost`

- Two lobbies are created with different host IDs.
- Searching with a partial string matching one lobby's code or host ID should return that lobby.

### `TestBlockUserPreventsJoin`

- Host creates a lobby and blocks a user.
- Listing lobbies as the blocked user must not return the host's lobby.

The server implementation currently returns `ErrNotImplemented` for all methods. The tests describe the expected behavior once the lobby service is fully implemented.
