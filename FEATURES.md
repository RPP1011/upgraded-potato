# Feature List and Acceptance Criteria

This document outlines the desired features for the `upgraded-potato` websocket netcode library and the requirements for accepting each feature. The project is focused on a peer-to-peer (P2P) architecture where each node typically connects to only one or two peers at a time.

## 1. WebSocket Server Initialization
- **Description:** Ability to create and configure a WebSocket server.
- **Acceptance:**
  - `NewServer(address)` returns a non-nil server instance.
  - Server stores the provided `address` value.

## 2. Peer Connections
- **Description:** Manage a small number of simultaneous peer connections.
- **Acceptance:**
  - Server maintains stable connections with up to three peers.
  - Each peer can send and receive messages.

## 3. Protobuf Message Serialization
- **Description:** Use Protocol Buffers for encoding network messages.
- **Acceptance:**
  - Messages are serialized to Protobuf format before transmission.
  - Server correctly deserializes received Protobuf messages.

## 4. .NET 3.5 Client Compatibility
- **Description:** Provide compatibility for .NET 3.5 clients.
- **Acceptance:**
  - Example client code in .NET 3.5 can establish a connection and exchange messages.
  - Cross-platform handshake demonstrates interoperability between Go server and .NET client.

## 5. Lobby Management
- **Description:** Clients can create and join lobbies for games.
- **Acceptance:**
  - API exposes commands to create, join, and leave lobbies.
  - Lobby size limits are enforced.

## 6. Matchmaking
- **Description:** Automatic matching of players into lobbies or games.
- **Acceptance:**
  - Matchmaking algorithm groups clients based on configurable criteria (e.g., skill, region).
  - Client receives notification when a match is found.

## 7. Secure Communication
- **Description:** Optional TLS support for encrypted connections.
- **Acceptance:**
  - Server can be configured with TLS certificates.
  - Clients can establish secure (wss://) connections.

## 8. Real-Time Game State Sync
- **Description:** Broadcast game state updates at a regular tick rate.
- **Acceptance:**
  - Server sends state updates to all connected clients at the configured interval.
  - Clients receive and apply state updates without data loss.

