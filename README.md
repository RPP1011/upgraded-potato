# upgraded-potato

A websocket-based multiplayer netcode library made in Go with Protobuf that supports .NET 3.5 clients.

## Development

This project follows a test-driven workflow. Run tests using the `Makefile`:

```sh
make test
```

Tests will also run automatically in GitHub Actions on each push or pull request.

## Test-driven Development

All features should be implemented using a test-driven workflow:

1. **Write a failing test** that describes the new behavior or regression fix.
2. **Add just enough code** to make the failing test pass.
3. **Run** `make test` or `go test ./...` to verify all tests succeed.
4. **Refactor** while keeping tests green.

This cycle keeps the library reliable and ensures new changes are covered by automated tests.

## Design Documents

- [Lobby System Design](docs/lobby_system_design.md)

