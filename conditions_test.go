package featureflow_go_sdk

import (
	//"log"
	"github.com/DATA-DOG/godog"
)

type contextType struct {
	target string
	targetType  string
	value string
	values []string
	valueType string
	output bool
}

var context contextType

func theTargetIsAWithTheValueOf(targetType, target string) error {
	context.target = target
	context.targetType = targetType
	return nil
}

func theValueIsAWithTheValueOf(valueType, value string) error {
	context.value = value
	context.valueType = valueType
	return nil
}

func theOperatorTestIsRun(operator string) error {
	if context.targetType == "string" {
		context.output = TestString(operator, context.target, []string{context.value})
	} else{
		return godog.ErrPending
	}

	return nil
}

func theOutputShouldEqual(outputString string) error {
	var output = outputString == "true"
	if context.output != output{
		return godog.ErrUndefined
	}
	return godog.ErrUndefined
}

func theValueIsAnArrayOfValues(values string) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^the target is a "([^"]*)" with the value of "([^"]*)"$`, theTargetIsAWithTheValueOf)
	s.Step(`^the value is a "([^"]*)" with the value of "([^"]*)"$`, theValueIsAWithTheValueOf)
	s.Step(`^the operator test "([^"]*)" is run$`, theOperatorTestIsRun)
	s.Step(`^the output should equal "([^"]*)"$`, theOutputShouldEqual)
	s.Step(`^the value is an array of values "([^"]*)"$`, theValueIsAnArrayOfValues)

	s.BeforeScenario(func(interface{}) {
		context = contextType{}
	})
}
