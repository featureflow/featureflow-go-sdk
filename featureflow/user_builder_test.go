package featureflow

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

type userBuilderTestContextType struct {
	user_builder UserBuilder
	user         *User
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

func theBuilderIsGivenTheFollowingAttributes(attributesTable *godog.Table) error {
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

func theResultUserShouldHaveAAttributeWithKeyAndValue(key, value string) error {
	attrs := userBuilderTestContext.user.GetAttributesForKey(key)
	if len(attrs) == 0 {
		return fmt.Errorf("Expected user to have attribute %s, but it was not set", key)
	}
	userAttribute := fmt.Sprintf("%v", attrs[0])
	if userAttribute != value {
		return fmt.Errorf("Expected %s to be %s", userAttribute, value)
	}
	return nil
}

func theResultUserShouldHaveAAttributeWithKeyAndCurrentDatetimeInIso8601(key string) error {
	attrs := userBuilderTestContext.user.GetAttributesForKey(key)
	if len(attrs) == 0 {
		return fmt.Errorf("Expected user to have attribute %s, but it was not set", key)
	}
	value := fmt.Sprintf("%v", attrs[0])
	if _, err := time.Parse(time.RFC3339, value); err != nil {
		return fmt.Errorf("Expected %s to be a valid ISO8601 datetime: %v", value, err)
	}
	return nil
}

func theBuilderShouldThrowAnError() error {
	if _, err := userBuilderTestContext.user_builder.Build(); err == nil {
		return fmt.Errorf("Expected an error to have been thrown")
	}
	return nil
}

func UserBuilderFeatureContext(ctx *godog.ScenarioContext) {
	ctx.Step(`^there is access to the User Builder module$`, thereIsAccessToTheUserBuilderModule)
	ctx.Step(`^the builder is initialised with the id "([^"]*)"$`, theBuilderIsInitialisedWithTheId)
	ctx.Step(`^the user is built using the builder$`, theUserIsBuiltUsingTheBuilder)
	ctx.Step(`^the result user should have an id "([^"]*)"$`, theResultUserShouldHaveAnId)
	ctx.Step(`^the result user should have no attributes$`, theResultUserShouldHaveNoAttributes)
	ctx.Step(`^the builder is given the following attributes$`, theBuilderIsGivenTheFollowingAttributes)
	ctx.Step(`^the result user should have a attribute with key "([^"]*)" and value "([^"]*)"$`, theResultUserShouldHaveAAttributeWithKeyAndValue)
	ctx.Step(`^the result user should have a attribute with key "([^"]*)" and current datetime in iso8601$`, theResultUserShouldHaveAAttributeWithKeyAndCurrentDatetimeInIso8601)
	ctx.Step(`^the builder should throw an error$`, theBuilderShouldThrowAnError)

	ctx.Before(func(c context.Context, sc *godog.Scenario) (context.Context, error) {
		userBuilderTestContext = userBuilderTestContextType{}
		return c, nil
	})
}