package featureflow

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
)

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

	suite := godog.TestSuite{
		Name:                "featureflow",
		ScenarioInitializer: InitializeScenario,
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
}
