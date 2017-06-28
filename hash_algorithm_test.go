package featureflow_go_sdk

import (
	"github.com/DATA-DOG/godog"
	"fmt"
)

type hashAlgorithmContextType struct {
	salt string
	feature string
	contextKey string
	hashResult string
	result float64
}

var hashAlgorithmContext hashAlgorithmContextType

func theSaltIsTheFeatureIsAndTheKeyIs(salt, feature, contextKey string) error {
	hashAlgorithmContext = hashAlgorithmContextType{
		salt: salt,
		feature: feature,
		contextKey: contextKey,
	}
	return nil
}

func theVariantValueIsCalculated() error {
	hashAlgorithmContext.hashResult = CalculateHash(
		hashAlgorithmContext.salt,
		hashAlgorithmContext.feature,
		hashAlgorithmContext.contextKey,
	)

	hashAlgorithmContext.result = GetVariantValue(hashAlgorithmContext.hashResult)
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

func HashAlgorithmFeatureContext(s *godog.Suite) {
	s.Step(`^the salt is "([^"]*)", the feature is "([^"]*)" and the key is "([^"]*)"$`, theSaltIsTheFeatureIsAndTheKeyIs)
	s.Step(`^the variant value is calculated$`, theVariantValueIsCalculated)
	s.Step(`^the hash value calculated should equal "([^"]*)"$`, theHashValueCalculatedShouldEqual)
	s.Step(`^the result from the variant calculation should be (\d+)$`, theResultFromTheVariantCalculationShouldBe)
}