# SDK harness

Small `net/http` app in `main.go` that uses the SDK source from this repo directly (same Go module — no build or install step needed).

```bash
# from the repo root
FEATUREFLOW_SERVER_KEY=sdk-srv-env-<your_key> go run ./examples/harness

# Staging / custom app host
FEATUREFLOW_BASE_URL=https://beta.featureflow-staging.com FEATUREFLOW_SERVER_KEY=sdk-srv-env-<your_key> go run ./examples/harness
```

Then open <http://127.0.0.1:3456/> to evaluate features from the browser, or use the JSON API:

| Endpoint | Description |
|---|---|
| `GET /health` | liveness check |
| `GET /api/config` | effective `baseUrl` and events setting |
| `GET /api/evaluate?feature=<key>&userId=<id>&role=<r>&tier=<t>` | evaluate one feature for a user |

Environment variables:

- `FEATUREFLOW_SERVER_KEY` (required) — server SDK key (`sdk-srv-env-…`).
- `FEATUREFLOW_BASE_URL` — override the default `https://app.featureflow.io`.
- `FEATUREFLOW_DISABLE_EVENTS` — the harness disables the events API by default; set to `false` to send evaluation events like a production client.
- `PORT` — listen port (default `3456`).

The harness registers a `harness-demo` feature (failover `off`) on startup, so that key exists in your environment to experiment with. See [CONTRIBUTING.md](../../CONTRIBUTING.md) for the full development guide.
