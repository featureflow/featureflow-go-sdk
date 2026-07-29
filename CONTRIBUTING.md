# featureflow-go-sdk

> Featureflow Go SDK

Get your Featureflow account at [featureflow.io](http://www.featureflow.io)

# Contributor Guidelines

Featureflow SDKs are open source and we welcome contributions from all developers.

## Local setup

You need **Go 1.21+** (see `go.mod`). No other tooling is required — dependencies are managed with Go modules and resolved from `go.sum`.

```bash
git clone https://github.com/featureflow/featureflow-go-sdk.git
cd featureflow-go-sdk
go build ./...
```

## Building and testing

```bash
# build
go build ./...

# static checks
go vet ./...

# format (run before committing)
go fmt ./...

# run the test suite (unit + BDD scenarios)
go test ./...
```

The test suite is BDD-style, driven by [godog](https://github.com/cucumber/godog) (Cucumber for Go):

- Scenarios live in `featureflow/features/*.feature`.
- Step definitions live in the matching `featureflow/*_test.go` files (e.g. `rules_test.go` implements `features/rules.feature`).
- `featureflow/main_test.go` wires every step context into a single `TestFeatures` suite, so `go test ./...` runs everything.

When adding or changing behavior, update **both** the `.feature` scenario and its Go step definitions — an unmatched step will fail the suite as "undefined".

For verbose scenario output:

```bash
go test ./featureflow -v
```

## Integration tests

Scenarios tagged `@integration` (see `featureflow/features/integration.feature`) hit a **real Featureflow environment** and are excluded from `go test ./...` by default. To run them, set an SDK key for an environment that has a `test-integration` feature configured to match the scenario examples:

```bash
FEATUREFLOW_TEST_API_KEY=sdk-srv-env-<your_key> go test ./...
```

- **Staging / custom app host:** set `FEATUREFLOW_TEST_BASE_URL` to point the client somewhere other than the default `https://app.featureflow.io`, e.g.
  `FEATUREFLOW_TEST_BASE_URL=https://beta.featureflow-staging.com FEATUREFLOW_TEST_API_KEY=sdk-srv-env-<key> go test ./...`
- Never hard-code real SDK keys in `.feature` files or step definitions — keys are always read from the environment.

## Manual test harness (example server)

A small `net/http` app in [`examples/harness/`](examples/harness/) uses the SDK source from this repo directly — no build step needed:

```bash
FEATUREFLOW_SERVER_KEY=sdk-srv-env-<your_key> go run ./examples/harness
```

- Defaults to `http://127.0.0.1:3456/` (override with `PORT`). The home page can trigger evaluations in the browser; JSON endpoints: `GET /health`, `GET /api/config`, `GET /api/evaluate?feature=<key>&userId=<id>&role=<r>&tier=<t>`.
- The harness turns off the events client by default (no calls to the events API). Set `FEATUREFLOW_DISABLE_EVENTS=false` to send evaluation events like a full production client.
- **Staging / custom app host:** set `FEATUREFLOW_BASE_URL` to point somewhere other than the default `https://app.featureflow.io`. Example:
  `FEATUREFLOW_BASE_URL=https://beta.featureflow-staging.com FEATUREFLOW_SERVER_KEY=sdk-srv-env-<key> go run ./examples/harness`
- On startup the harness registers a `harness-demo` feature (failover `off`), so that key exists in your environment to experiment with.

See [`examples/harness/README.md`](examples/harness/README.md) for details.

**SDK log output:** the client logs to **stderr** with the prefix `Featureflow:` by default. Pass your own `*log.Logger` via `Config.Logger` to redirect or silence it.

## License

Apache-2.0
