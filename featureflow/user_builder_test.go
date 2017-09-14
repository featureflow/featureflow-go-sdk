package featureflow

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"fmt"
	"strings"
)

type userBuilderTestContextType struct {
	user_builder UserBuilder
	user         User
	error			error
}

var userBuilderTestContext userBuilderTestContextType

func thereIsAccessToTheUserBuilderModule() error {
	return nil
}

func theBuilderIsInitialisedWithTheId(id string) error {
	userBuilderTestContext.user_builder = NewUserBuilder(id)
	return nil
}

func theUserIsBuiltUsingTheBuilder() error {
	userBuilderTestContext.user, userBuilderTestContext.error = userBuilderTestContext.user_builder.Build()
	return nil
}

func theResultUserShouldHaveAnId(id string) error {
	if userBuilderTestContext.user.GetId() != id {
		return fmt.Errorf("Expected %s to be %s", userBuilderTestContext.user.GetId(), id)
	}
	return nil
}

func theResultUserShouldHaveNoAttributes() error {
	keys := userBuilderTestContext.user.GetAttributeKeys()
	filteredKeys := keys[:0]
	for _, key := range keys {
		if !strings.HasPrefix(key, "featureflow.") {
			filteredKeys = append(filteredKeys, key)
		}
	}

	if len(filteredKeys) > 0{
		return fmt.Errorf("Expected %d to be greater than 0", len(userBuilderTestContext.user.GetAttributeKeys()))
	}
	return nil
}

func theBuilderIsGivenTheFollowingAttributes(attributesTable *gherkin.DataTable) error {
	head := attributesTable.Rows[0].Cells

	for i := 1; i < len(attributesTable.Rows); i++ {
		key := ""
		value := ""
		for n, cell := range attributesTable.Rows[i].Cells {
			switch head[n].Value {
			case "key":
				key = cell.Value
			case "value":
				value = cell.Value
			default:
				return fmt.Errorf("unexpected column name: %s", head[n].Value)
			}
		}
		userBuilderTestContext.user_builder = userBuilderTestContext.user_builder.WithAttribute(key, value)
	}
	return nil
}

func theResultUserShouldHaveTheKeyWithAttribute(key, attribute string) error {
	userAttribute := userBuilderTestContext.user.GetAttributesForKey(key)[0]
	if userAttribute != attribute{
		return fmt.Errorf("Expected %s to be %s", userAttribute, attribute)
	}
	return nil
}

func theBuilderShouldThrowAnError() error {
	if _, err := userBuilderTestContext.user_builder.Build(); err == nil {
		return fmt.Errorf("Expected an error to have been thrown")
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^there is access to the Context Builder module$`, thereIsAccessToTheUserBuilderModule)
	s.Step(`^the builder is initialised with the id "([^"]*)"$`, theBuilderIsInitialisedWithTheId)
	s.Step(`^the user is built using the builder$`, theUserIsBuiltUsingTheBuilder)
	s.Step(`^the result user should have an id "([^"]*)"$`, theResultUserShouldHaveAnId)
	s.Step(`^the result user should have no attributes$`, theResultUserShouldHaveNoAttributes)
	s.Step(`^the builder is given the following attributes$`, theBuilderIsGivenTheFollowingAttributes)
	s.Step(`^the result user should have the key "([^"]*)" with attribute "([^"]*)"$`, theResultUserShouldHaveTheKeyWithAttribute)
	s.Step(`^the builder should throw an error$`, theBuilderShouldThrowAnError)

	s.BeforeScenario(func(interface{}){
		userBuilderTestContext = userBuilderTestContextType{}
	})
}