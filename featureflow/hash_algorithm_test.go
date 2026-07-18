package featureflow

import (
	"fmt"

	"github.com/cucumber/godog"
)

type hashAlgorithmContextType struct {
	salt string
	feature string
	userId string
	hashResult string
	result float64
}

var hashAlgorithmContext hashAlgorithmContextType

func theSaltIsTheFeatureIsAndTheUserIdIs(salt, feature, userId string) error {
	hashAlgorithmContext = hashAlgorithmContextType{
		salt: salt,
		feature: feature,
		userId: userId,
	}
	return nil
}

func theVariantValueIsCalculated() error {
	hashAlgorithmContext.hashResult = calculateHash(
		hashAlgorithmContext.salt,
		hashAlgorithmContext.feature,
		hashAlgorithmContext.userId,
	)

	hashAlgorithmContext.result = getVariantValue(hashAlgorithmContext.hashResult)
	return nil
}

func theHashValueCalculatedShouldEqual(hash string) error {
	if hash != hashAlgorithmContext.hashResult{
		return fmt.Errorf("Expected %s to be %s", hashAlgorithmContext.hashResult, hash)
	}
	return nil
}

func theResultFromTheVariantCalculationShouldBe(result float64) error {
	if result != hashAlgorithmContext.result{
		return fmt.Errorf("Expected %f to be %f", hashAlgorithmContext.result, result)
	}
	return nil
}

func HashAlgorithmFeatureContext(ctx *godog.ScenarioContext) {
	ctx.Step(`^the salt is "([^"]*)", the feature is "([^"]*)" and the user id is "([^"]*)"$`, theSaltIsTheFeatureIsAndTheUserIdIs)
	ctx.Step(`^the variant value is calculated$`, theVariantValueIsCalculated)
	ctx.Step(`^the hash value calculated should equal "([^"]*)"$`, theHashValueCalculatedShouldEqual)
	ctx.Step(`^the result from the variant calculation should be (\d+)$`, theResultFromTheVariantCalculationShouldBe)
}