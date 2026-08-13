package featureflow

import (
	"context"
	"os"
	"strings"
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
// running rather than to have legitimately shrunk. 112 scenarios run today; raise this
// if the testbed grows substantially, but never lower it to match a sudden drop
// without first working out where the missing scenarios went.
const minScenarios = 98

// The shared testbed carries scenarios for behaviour this SDK does not implement, so the
// feature files are listed rather than globbed. A list on its own trades one silent failure
// for another — a file added upstream would simply never run — so every file in the testbed
// must appear in exactly one of these two sets, and TestTestbedFeatureFilesAreAccountedFor
// fails on any that does not. Moving a file from unsupported to supported is then a
// deliberate act with a visible diff.
var (
	supportedFeatures = []string{
		"application_tag.feature",
		"bucketing.feature",
		"conditions.feature",
		"feature_evaluation.feature",
		"json_value.feature", // present but tag-excluded; see TestFeatures
		"rules.feature",
		"user_builder.feature",
	}

	// Each entry names the backlog item that would close the gap, so an unimplemented
	// suite cannot quietly become permanent.
	unsupportedFeatures = map[string]string{
		"events.feature":            "client-side event summarisation is not implemented in this SDK",
		"goals.feature":             "goal tracking is not implemented in this SDK (SDK-BACKLOG item 7)",
		"sdk_config.feature":        "server-driven SDK config is not implemented in this SDK",
		"tracked_exposures.feature": "per-flag exposure fidelity is not implemented in this SDK",
	}
)

const testbedGherkin = "../testbed/gherkin"

// TestTestbedFeatureFilesAreAccountedFor makes an upstream addition to the testbed a
// build failure here rather than a silent omission. It is the counterpart to Strict:
// Strict catches a step that stopped matching, this catches a whole file that was never
// wired up.
func TestTestbedFeatureFilesAreAccountedFor(t *testing.T) {
	entries, err := os.ReadDir(testbedGherkin)
	if err != nil {
		t.Fatalf("cannot read %s — is the testbed submodule initialised? %v", testbedGherkin, err)
	}

	known := make(map[string]bool, len(supportedFeatures)+len(unsupportedFeatures))
	for _, name := range supportedFeatures {
		known[name] = true
	}
	for name := range unsupportedFeatures {
		known[name] = true
	}

	seen := make(map[string]bool, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".feature") {
			continue
		}
		seen[name] = true
		if !known[name] {
			t.Errorf("testbed feature file %q is in neither supportedFeatures nor "+
				"unsupportedFeatures — add step definitions and list it as supported, or "+
				"record why it is not implemented here", name)
		}
	}

	for name := range known {
		if !seen[name] {
			t.Errorf("feature file %q is listed here but no longer exists in the testbed — "+
				"it was probably renamed upstream", name)
		}
	}
}

// featurePaths returns the absolute set of feature files godog should run.
func featurePaths() []string {
	paths := make([]string, 0, len(supportedFeatures)+1)
	for _, name := range supportedFeatures {
		paths = append(paths, testbedGherkin+"/"+name)
	}
	return append(paths, "features/integration.feature")
}

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
	ApplicationTagFeatureContext(ctx)
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
			Paths:  featurePaths(),
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
