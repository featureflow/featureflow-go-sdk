package featureflow

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"fmt"
)

type contextBuilderTestContextType struct {
	context_builder ContextBuilder
	context         Context
	error			error
}

var contextBuilderTestContext contextBuilderTestContextType

func thereIsAccessToTheContextBuilderModule() error {
	return nil
}

func theBuilderIsInitialisedWithTheKey(key string) error {
	contextBuilderTestContext.context_builder = NewContextBuilder(key)
	return nil
}

func theContextIsBuiltUsingTheBuilder() error {
	contextBuilderTestContext.context, contextBuilderTestContext.error = contextBuilderTestContext.context_builder.Build()
	return nil
}

func theResultContextShouldHaveAKey(key string) error {
	if contextBuilderTestContext.context.GetKey() != key {
		return fmt.Errorf("Expected %s to be %s", contextBuilderTestContext.context.GetKey(), key)
	}
	return nil
}

func theResultContextShouldHaveNoValues() error {
	if len(contextBuilderTestContext.context.GetValueKeys()) > 0{
		return fmt.Errorf("Expected %d to be greater than 0", len(contextBuilderTestContext.context.GetValueKeys()))
	}
	return nil
}

func theBuilderIsGivenTheFollowingValues(valuesTable *gherkin.DataTable) error {
	head := valuesTable.Rows[0].Cells

	for i := 1; i < len(valuesTable.Rows); i++ {
		key := ""
		value := ""
		for n, cell := range valuesTable.Rows[i].Cells {
			switch head[n].Value {
			case "key":
				key = cell.Value
			case "value":
				value = cell.Value
			default:
				return fmt.Errorf("unexpected column name: %s", head[n].Value)
			}
		}
		contextBuilderTestContext.context_builder = contextBuilderTestContext.context_builder.WithValue(key, value)
	}
	return nil
}

func theResultContextShouldHaveTheKeyWithValue(key, value string) error {
	contextValue := contextBuilderTestContext.context.GetValuesForKey(key)[0]
	if contextValue != value{
		return fmt.Errorf("Expected %s to be %s", contextValue, value)
	}
	return nil
}

func theBuilderShouldThrowAnError() error {
	if _, err := contextBuilderTestContext.context_builder.Build(); err == nil {
		return fmt.Errorf("Expected an error to have been thrown")
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^there is access to the Context Builder module$`, thereIsAccessToTheContextBuilderModule)
	s.Step(`^the builder is initialised with the key "([^"]*)"$`, theBuilderIsInitialisedWithTheKey)
	s.Step(`^the context is built using the builder$`, theContextIsBuiltUsingTheBuilder)
	s.Step(`^the result context should have a key "([^"]*)"$`, theResultContextShouldHaveAKey)
	s.Step(`^the result context should have no values$`, theResultContextShouldHaveNoValues)
	s.Step(`^the builder is given the following values$`, theBuilderIsGivenTheFollowingValues)
	s.Step(`^the result context should have the key "([^"]*)" with value "([^"]*)"$`, theResultContextShouldHaveTheKeyWithValue)
	s.Step(`^the builder should throw an error$`, theBuilderShouldThrowAnError)

	s.BeforeScenario(func(interface{}){
		contextBuilderTestContext = contextBuilderTestContextType{}
	})
}