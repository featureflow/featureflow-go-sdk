package featureflow

import (
	"context"
	"os"
	"sync/atomic"
	"testing"
	"time"
	// Embeds the IANA timezone database in the test binary so TestMain's
	// LoadLocation works on hosts without system tzdata (scratch/distroless CI images).
	_ "time/tzdata"

	"github.com/cucumber/godog"
)

// testTimezone is deliberately not UTC. A date-only condition value such as
// "2026-07-03" must mean UTC midnight in every SDK; an implementation that parses
// it in the host's local time passes those scenarios anyway whenever the runner
// happens to sit on UTC, which is exactly how the bug survived in three SDKs.
// Asia/Tokyo is +09:00 all year — a whole-hour offset with no DST, so the expected
// instants stay easy to reason about.
const testTimezone = "Asia/Tokyo"

// minScenarios is the floor below which the BDD suite is assumed to have stopped
// running rather than to have legitimately shrunk. 93 scenarios run today; raise this
// if the testbed grows substantially, but never lower it to match a sudden drop
// without first working out where the missing scenarios went.
const minScenarios = 90

// TestMain pins the process timezone for the whole test binary. Doing it here rather
// than in a per-test helper means a new test file cannot opt out by forgetting to call
// something, and removing it makes the date scenarios visibly stop being meaningful.
func TestMain(m *testing.M) {
	loc, err := time.LoadLocation(testTimezone)
	if err != nil {
		panic("test setup: cannot load timezone " + testTimezone + ": " + err.Error())
	}
	time.Local = loc
	os.Setenv("TZ", testTimezone)

	os.Exit(m.Run())
}

// InitializeScenario registers every *_test.go file's step definitions into a
// single suite. Each file still owns its own step regexes and Before hook.
func InitializeScenario(ctx *godog.ScenarioContext) {
	RulesFeatureContext(ctx)
	ConditionsFeatureContext(ctx)
	HashAlgorithmFeatureContext(ctx)
	UserBuilderFeatureContext(ctx)
	FeatureEvaluationFeatureContext(ctx)
	IntegrationFeatureContext(ctx)
}

// TestFeatures runs the shared scenarios from ../testbed/gherkin (a submodule of
// github.com/featureflow/featureflow-sdk-testbed) plus this repo's own
// features/integration.feature, via `go test ./...`.
//
// @builder-defers-implicit-attributes is excluded: NewUserBuilder injects implicit
// attributes (featureflow.user.id/date) immediately at construction time, not later
// at evaluate time, so this SDK runs @builder-injects-implicit-attributes instead.
// @json-value is excluded: not yet implemented in this SDK.
// @integration hits a real Featureflow server and is excluded unless
// FEATUREFLOW_TEST_API_KEY is set (see integration_test.go and features/integration.feature).
func TestFeatures(t *testing.T) {
	tags := "~@integration && ~@builder-defers-implicit-attributes && ~@json-value"
	if os.Getenv("FEATUREFLOW_TEST_API_KEY") != "" {
		tags = "~@builder-defers-implicit-attributes && ~@json-value"
	}

	// Strict catches steps that no longer match a definition. It does not catch a suite
	// that stopped loading its feature files altogether — "0 scenarios (0 passed)" is a
	// pass. The floor is a floor, not an exact count: scenarios are added upstream all
	// the time, but the total must never collapse.
	var scenariosRun int64
	initializer := func(ctx *godog.ScenarioContext) {
		InitializeScenario(ctx)
		ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
			atomic.AddInt64(&scenariosRun, 1)
			return ctx, nil
		})
	}

	suite := godog.TestSuite{
		Name:                "featureflow",
		ScenarioInitializer: initializer,
		Options: &godog.Options{
			Format: "pretty",
			Paths:  []string{"../testbed/gherkin", "features/integration.feature"},
			Tags:   tags,
			// Without Strict, a scenario whose steps have no matching definition is reported as
			// "undefined" and the suite still exits zero. The testbed is a submodule, so a step
			// renamed upstream silently stops being tested here — which is how the whole of
			// conditions.feature came to be skipped while `go test` stayed green.
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}

	if n := atomic.LoadInt64(&scenariosRun); n < minScenarios {
		t.Fatalf("only %d scenarios ran, expected at least %d — the testbed submodule is "+
			"probably uninitialised, or Paths/Tags no longer select the shared features", n, minScenarios)
	}
}
