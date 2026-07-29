package featureflow

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

type conditionsContextType struct {
	target interface{}
	targetType  string
	values []interface{}
	valueType string
	output bool
}

var conditionsTestContext conditionsContextType

func theTargetIsAWithTheValueOf(targetType, target string) error {
	conditionsTestContext.targetType = targetType
	if targetType == "number"{

		if t, err := strconv.ParseFloat(target, 64); err != nil{
			return err
		} else {
			conditionsTestContext.target = t
		}
	} else{
		conditionsTestContext.target = target
	}
	return nil
}

func theAttributeIsAWithTheValueOf(valueType, value string) error {
	conditionsTestContext.valueType = valueType
	if valueType == "number"{
		if v, err := strconv.ParseFloat(value, 64); err != nil{
			return err
		} else {
			conditionsTestContext.values = []interface{}{v}
		}
	} else{
		conditionsTestContext.values = []interface{}{value}
	}
	return nil
}

func theOperatorTestIsRun(operator string) error {
	if conditionsTestContext.targetType == "string" || conditionsTestContext.targetType == "number" {
		conditionsTestContext.output = conditionsTest(operator, conditionsTestContext.target, conditionsTestContext.values)
	} else{
		return errors.New("Error, operator is not defined")
	}

	return nil
}

func theOutputShouldEqual(outputString string) error {
	var output = outputString == "true"

	if conditionsTestContext.output != output{
		return fmt.Errorf("Expected %s to equal %s", strconv.FormatBool(conditionsTestContext.output), strconv.FormatBool(output))
	}
	return nil
}

func theAttributeIsAnArrayOfValues(valuesString string) error {
	var values = strings.Split(valuesString, ", ")
	conditionsTestContext.values = make([]interface{}, len(values))
	for i := range values {
		conditionsTestContext.values[i] = values[i]
	}
	return nil
}

func ConditionsFeatureContext(ctx *godog.ScenarioContext) {
	ctx.Step(`^the target is a "([^"]*)" with the value of "([^"]*)"$`, theTargetIsAWithTheValueOf)
	ctx.Step(`^the attribute is a "([^"]*)" with the value of "([^"]*)"$`, theAttributeIsAWithTheValueOf)
	ctx.Step(`^the operator test "([^"]*)" is run$`, theOperatorTestIsRun)
	ctx.Step(`^the output should equal "([^"]*)"$`, theOutputShouldEqual)
	ctx.Step(`^the attribute is an array of values "([^"]*)"$`, theAttributeIsAnArrayOfValues)

	ctx.Before(func(c context.Context, sc *godog.Scenario) (context.Context, error) {
		conditionsTestContext = conditionsContextType{}
		return c, nil
	})
}
