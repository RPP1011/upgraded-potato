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

## Test Reporter Integration

The project uses the [dorny/test-reporter](https://github.com/dorny/test-reporter) action
to surface Go test results in pull requests. Below is the recommended setup for
public repositories.

**.github/workflows/ci.yml**

```yaml
name: CI
on:
  pull_request:
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: go test -json ./... > report.json
      - uses: actions/upload-artifact@v4
        if: ${{ !cancelled() }}
        with:
          name: test-results
          path: report.json
```

**.github/workflows/test-report.yml**

```yaml
name: Test Report
on:
  workflow_run:
    workflows: ["CI"]
    types:
      - completed
permissions:
  contents: read
  actions: read
  checks: write
jobs:
  report:
    runs-on: ubuntu-latest
    steps:
      - uses: dorny/test-reporter@v2
        with:
          artifact: test-results
          name: Go Tests
          path: report.json
          reporter: golang-json
```

When you push new commits, the CI workflow uploads results as an artifact, and
the Test Report workflow converts them into annotated reports on GitHub.

## TestReporter CLI

If you prefer uploading results directly from your own environment, you can use the TestReporter command-line tool.

1. Install the CLI from <https://testreporter.com>.
2. Generate a JSON test report:
   ```sh
   go test -json ./... > report.json
   ```
3. Upload the report:
   ```sh
   testreporter upload --file report.json --format golang-json --token <your-token>
   ```
   Replace `<your-token>` with the API token from your TestReporter account.
   The CLI sends the test results to your project so they appear in the TestReporter dashboard.

## Design Documents

- [Lobby System Design](docs/lobby_system_design.md)

