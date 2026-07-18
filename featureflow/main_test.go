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
	IntegrationFeatureContext(ctx)
}

// TestFeatures runs every .feature file under features/ via `go test ./...`.
// @integration scenarios hit a real Featureflow server and are excluded unless
// FEATUREFLOW_TEST_API_KEY is set (see integration_test.go and features/integration.feature).
func TestFeatures(t *testing.T) {
	tags := "~@integration"
	if os.Getenv("FEATUREFLOW_TEST_API_KEY") != "" {
		tags = ""
	}

	suite := godog.TestSuite{
		Name:                "featureflow",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			Tags:     tags,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
