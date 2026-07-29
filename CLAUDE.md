# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Server-side Go SDK for Featureflow (`github.com/featureflow/featureflow-go-sdk`), a feature-flag evaluation client. It polls the Featureflow API for flag config, caches it in memory, and evaluates flags locally against a `User` — no network round-trip per evaluation. See the root workspace `CLAUDE.md` (one directory up) for how this SDK fits into the broader Featureflow system; this repo has no CLAUDE.md of its own beyond this file.

## Commands

```bash
go build ./...
go vet ./...
go test ./...              # unit + BDD scenarios, excludes @integration
go fmt ./...
```

There's no single-test flag for the BDD suite (it's one `TestFeatures` entry point over all `.feature` files) — to narrow scope, temporarily comment out unwanted `.feature` files or `ctx.Step`/`FeatureContext` registrations in [main_test.go](featureflow/main_test.go), or run `go test -run TestFeatures ./...` (which still runs the whole suite, just isolates it from any future non-BDD `Test*` funcs).

Integration scenarios in [features/integration.feature](featureflow/features/integration.feature) hit a real Featureflow server and are tagged `@integration`. They're skipped unless `FEATUREFLOW_TEST_API_KEY` (and optionally `FEATUREFLOW_TEST_BASE_URL`) is set — see the tag logic in [main_test.go](featureflow/main_test.go).

Manual test harness (browser + JSON API against a real environment): `FEATUREFLOW_SERVER_KEY=sdk-srv-env-<key> go run ./examples/harness` — see [examples/harness/README.md](examples/harness/README.md).

## Architecture

**Evaluation is entirely local.** `Client()` in [featureflow_client.go](featureflow/featureflow_client.go) starts a background poller ([polling_client.go](featureflow/polling_client.go)) that fetches `{BaseURL}/api/sdk/v1/features` every 30s (ETag-conditional) and writes the full flag set into a `FeatureStore`. `client.Evaluate(key, user)` reads from that store synchronously and computes the variant in-process — it never blocks on network I/O.

Evaluation flow (`FeatureflowClient.Evaluate` → returns an `Evaluate` struct with `.Is()`/`.IsOn()`/`.IsOff()`/`.Value()`):
1. Look up the feature in `Config.FeatureStore`. If missing, fall back to `FailoverVariants[key]` (set via `Config.WithFeatures`) or `"off"`.
2. If found and disabled, return `feature.OffVariantKey`.
3. If enabled, walk `feature.Rules` in order; the first rule whose `audience.Conditions` all match the user (via [conditions.go](featureflow/conditions.go)'s operators — `equals`, `contains`, `in`, `greaterThan`, `before`, etc.) wins, or a `DefaultRule: true` rule always matches.
4. The winning rule's `VariantSplits` are bucketed deterministically by SHA-1 hash of `salt:featureKey:userId` (see `calculateHash`/`getVariantValue`/`getVariantSplitKey` in [evaluate_helpers.go](featureflow/evaluate_helpers.go)) — same user + feature always gets the same variant. This must stay in sync with the bucketing algorithm used by the other SDKs and backend (see root `CLAUDE.md`: SHA-1 of `salt:featureKey:bucketKey` mod 100).
5. Calling `.Is(variant)` on the result fires an async `evaluate` event to `{BaseURL}/api/sdk/v1/events` (fire-and-forget goroutine in [events_client.go](featureflow/events_client.go)), unless `Config.DisableEvents` is set.

**Key files:**
- [featureflow_client.go](featureflow/featureflow_client.go) — `Client()` constructor, `Config` struct, `Evaluate()` entry point.
- [feature.go](featureflow/feature.go) — wire types (`feature`, `rule`, `audience`, `condition`, `variantSplit`) deserialized from the polling response.
- [feature_store.go](featureflow/feature_store.go) — `FeatureStore` interface + the default `inMemoryStore` (mutex-guarded map). Swappable via `Config.FeatureStore`.
- [user_builder.go](featureflow/user_builder.go) — `NewUserBuilder(id)` builder for `User`; auto-injects `featureflow.user.id` and `featureflow.date` attributes (the latter is what `before`/`after` date-rule conditions and time-of-day rules match against).
- [feature_registration.go](featureflow/feature_registration.go) — `WithFeature(key, failover).AddVariant(...).Build()`, used for `Config.WithFeatures` to pre-register flags/failovers with the server and set local failover variants.
- [conditions.go](featureflow/conditions.go) / [evaluate_helpers.go](featureflow/evaluate_helpers.go) — rule-matching and bucketing logic; this is the part that must match other SDKs' evaluation semantics.

**BDD test wiring:** each `*_test.go` file owns step definitions for one concern (`conditions_test.go`, `rules_test.go`, `hash_algorithm_test.go`, `user_builder_test.go`, `integration_test.go`) and exposes a `XxxFeatureContext(ctx *godog.ScenarioContext)` function. [main_test.go](featureflow/main_test.go) registers all of them into one `godog.TestSuite` that runs every `.feature` file under `featureflow/features/`. When changing behavior, keep the `.feature` scenario and its Go step definitions in sync — an orphaned step regex or scenario fails silently as "undefined."

## Conventions

- Public API surface (`Client`, `UserBuilder`, `Evaluate`, `FeatureRegistrationBuilder`) uses the builder pattern returning interfaces; internal wire types (`feature`, `rule`, `condition`, etc.) are unexported structs.
- `Config.Logger` defaults to `log.New(os.Stderr, "Featureflow:", log.LstdFlags)` if not supplied; log lines are prefixed `LOG_INFO`/`LOG_ERROR`.
- No external HTTP mocking library is used for the non-integration tests — the polling/events clients are exercised indirectly via the BDD scenarios and the real network calls are only reached in `@integration` tests.
