package featureflow

import (
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

type featureEvaluationTestContextType struct {
	feature feature
	result  string
}

var featureEvaluationTestContext featureEvaluationTestContextType

func theFeatureWithAnOffVariantKeyADefaultKeyOfIs(key, offVariantKey, defaultKey, enabledOrDisabled string) error {
	featureEvaluationTestContext.feature = feature{
		Key:           key,
		OffVariantKey: offVariantKey,
		Enabled:       enabledOrDisabled == "enabled",
		Rules: []rule{
			{DefaultRule: true, VariantSplits: []variantSplit{{VariantKey: defaultKey, Split: 100}}},
		},
	}
	return nil
}

// theFeatureIsEvaluatedWithAUser mirrors the evaluation logic in
// FeatureflowClient.Evaluate, without needing a real Client/polling — this tests the
// same enabled/disabled -> variant decision the client makes internally.
func theFeatureIsEvaluatedWithAUser(userId string) error {
	user, err := NewUserBuilder(userId).Build()
	if err != nil {
		return err
	}

	f := featureEvaluationTestContext.feature
	evaluatedVariant := f.OffVariantKey
	if f.Enabled {
		for _, r := range f.Rules {
			if ruleMatches(r, user) {
				variantValue := getVariantValue(calculateHash("1", f.Key, user.GetId()))
				evaluatedVariant = getVariantSplitKey(r.VariantSplits, variantValue)
				break
			}
		}
	}
	featureEvaluationTestContext.result = evaluatedVariant
	return nil
}

func theEvaluatedVariantShouldBe(variant string) error {
	if featureEvaluationTestContext.result != variant {
		return fmt.Errorf("Expected %s to be %s", featureEvaluationTestContext.result, variant)
	}
	return nil
}

func FeatureEvaluationFeatureContext(ctx *godog.ScenarioContext) {
	ctx.Step(`^the feature "([^"]*)" with an offVariantKey "([^"]*)", a default key of "([^"]*)" is (enabled|disabled)$`, theFeatureWithAnOffVariantKeyADefaultKeyOfIs)
	ctx.Step(`^the feature is evaluated with a user "([^"]*)"$`, theFeatureIsEvaluatedWithAUser)
	ctx.Step(`^the evaluated variant should be "([^"]*)"$`, theEvaluatedVariantShouldBe)

	ctx.Before(func(c context.Context, sc *godog.Scenario) (context.Context, error) {
		featureEvaluationTestContext = featureEvaluationTestContextType{}
		return c, nil
	})
}
