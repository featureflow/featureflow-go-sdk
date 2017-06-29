package featureflow

import (
	"github.com/DATA-DOG/godog"
	"fmt"
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

func theFeatureflowClientIsInitializedWithTheApiKey(api_key string) error {
	integrationTestContext.client, integrationTestContext.error = Client(api_key, Config{})
	return nil
}

func theFeatureWithContextKeyIsEvaluatedWithTheValue(featureKey, contextKey, variantValue string) error {
	context, _ := NewContextBuilder(contextKey).Build()
	integrationTestContext.result = integrationTestContext.client.Evaluate(featureKey, context).Is(variantValue)
	return nil
}

func theResultOfTheEvaluationShouldEqual(value string) error {
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

func IntegrationFeatureContext(s *godog.Suite) {
	s.Step(`^there is access to the Featureflow library$`, thereIsAccessToTheFeatureflowLibrary)
	s.Step(`^the FeatureflowClient is initialized with the apiKey "([^"]*)"$`, theFeatureflowClientIsInitializedWithTheApiKey)
	s.Step(`^the feature "([^"]*)" with context key "([^"]*)" is evaluated with the value "([^"]*)"$`, theFeatureWithContextKeyIsEvaluatedWithTheValue)
	s.Step(`^the result of the evaluation should equal (true|false)$`, theResultOfTheEvaluationShouldEqual)
	s.Step(`^the FeatureflowClient is initialized with no apiKey$`, theFeatureflowClientIsInitializedWithNoApiKey)
	s.Step(`^the featureflow client should throw an error$`, theFeatureflowClientShouldThrowAnError)

	s.BeforeScenario(func(interface{}) {
		integrationTestContext = integrationTestContextType{}
	})
}
