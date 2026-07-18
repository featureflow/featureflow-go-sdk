package featureflow

import (
	"context"
	"fmt"
	"os"

	"github.com/cucumber/godog"
)

type integrationTestContextType struct{
	client *FeatureflowClient
	error error
	result bool
}

var integrationTestContext integrationTestContextType

func thereIsAccessToTheFeatureflowLibrary() error {
	//it wil be available, strongly typed
	return nil
}

// theFeatureflowClientIsInitializedWithTheConfiguredApiKey reads the key and base URL
// from the environment instead of a literal in the .feature file, so no real SDK key
// is checked into source. Scenarios using this step are tagged @integration and
// excluded by default (see TestFeatures) — they only run when a harness sets
// FEATUREFLOW_TEST_API_KEY.
func theFeatureflowClientIsInitializedWithTheConfiguredApiKey() error {
	config := Config{BaseURL: os.Getenv("FEATUREFLOW_TEST_BASE_URL")}
	integrationTestContext.client, integrationTestContext.error = Client(os.Getenv("FEATUREFLOW_TEST_API_KEY"), config)
	return nil
}

func theFeatureflowClientIsInitializedWithTheApiKey(api_key string) error {
	integrationTestContext.client, integrationTestContext.error = Client(api_key, Config{BaseURL: os.Getenv("FEATUREFLOW_TEST_BASE_URL")})
	return nil
}

func theFeatureWithUserIdIsEvaluatedWithTheVariantValue(featureKey, userId, variantValue string) error {
	user, _ := NewUserBuilder(userId).Build()
	integrationTestContext.result = integrationTestContext.client.Evaluate(featureKey, user).Is(variantValue)
	return nil
}

func theResultOfTheEvaluationShouldEqual(value string) error {
	expected := value == "true"
	if integrationTestContext.result != expected {
		return fmt.Errorf("Expected %t to be %t", integrationTestContext.result, expected)
	}
	return nil
}

func theFeatureflowClientIsInitializedWithNoApiKey() error {
	integrationTestContext.client, integrationTestContext.error = Client("", Config{})
	return nil
}

func theFeatureflowClientShouldThrowAnError() error {
	if integrationTestContext.error == nil{
		return fmt.Errorf("Expected the featureflow client to throw an error")
	}
	return nil
}

func IntegrationFeatureContext(ctx *godog.ScenarioContext) {
	ctx.Step(`^there is access to the Featureflow library$`, thereIsAccessToTheFeatureflowLibrary)
	ctx.Step(`^the FeatureflowClient is initialized with the configured apiKey$`, theFeatureflowClientIsInitializedWithTheConfiguredApiKey)
	ctx.Step(`^the FeatureflowClient is initialized with the apiKey "([^"]*)"$`, theFeatureflowClientIsInitializedWithTheApiKey)
	ctx.Step(`^the feature "([^"]*)" with user id "([^"]*)" is evaluated with the variant value "([^"]*)"$`, theFeatureWithUserIdIsEvaluatedWithTheVariantValue)
	ctx.Step(`^the result of the evaluation should equal (true|false)$`, theResultOfTheEvaluationShouldEqual)
	ctx.Step(`^the FeatureflowClient is initialized with no apiKey$`, theFeatureflowClientIsInitializedWithNoApiKey)
	ctx.Step(`^the featureflow client should throw an error$`, theFeatureflowClientShouldThrowAnError)

	ctx.Before(func(c context.Context, sc *godog.Scenario) (context.Context, error) {
		integrationTestContext = integrationTestContextType{}
		return c, nil
	})
}
