package featureflow

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"fmt"
	"strconv"
	"encoding/json"
)


type rulesTestContextType struct{
	result                bool
	rule                  Rule
	context_builder       ContextBuilder
	variantValue          float64
	splitKeyVariantResult string
}

var rulesTestContext rulesTestContextType

func theRuleIsADefaultRule() error {
	rulesTestContext.rule.DefaultRule = true
	return nil
}

func theRuleIsMatchedAgainstTheContext() error {
	c, err := rulesTestContext.context_builder.Build()
	rulesTestContext.result = ruleMatches(
		rulesTestContext.rule,
		c,
	)

	return err
}

func theResultFromTheMatchShouldBe(resultStr string) error {
	result := resultStr == "true"
	if result != rulesTestContext.result{
		return fmt.Errorf("Expected %s to be %s",
			strconv.FormatBool(rulesTestContext.result),
			strconv.FormatBool(result),
		)
	}
	return nil
}

func theContextValuesAre(contextValuesTable *gherkin.DataTable) error {

	head := contextValuesTable.Rows[0].Cells

	for i := 1; i < len(contextValuesTable.Rows); i++ {
		key := ""
		var value Value
		var values []Value
		for n, cell := range contextValuesTable.Rows[i].Cells {
			switch head[n].Value {
			case "key":
				key = cell.Value
			case "value":
				if cell.Value[0] == '[' {
					json.Unmarshal([]byte(cell.Value), &values)
				} else{
					json.Unmarshal([]byte(cell.Value), &value)
				}
			default:
				return fmt.Errorf("unexpected column name: %s", head[n].Value)
			}
		}
		if values != nil {
			rulesTestContext.context_builder = rulesTestContext.context_builder.WithValues(key, values)
		} else {
			rulesTestContext.context_builder = rulesTestContext.context_builder.WithValue(key, value)
		}
	}
	return nil
}

func theRulesAudienceConditionsAre(audienceConditions *gherkin.DataTable) error {
	head := audienceConditions.Rows[0].Cells
	conditions := &rulesTestContext.rule.Audience.Conditions

	for i := 1; i < len(audienceConditions.Rows); i++ {
		condition := Condition{}
		for n, cell := range audienceConditions.Rows[i].Cells {
			switch head[n].Value {
			case "operator":
				condition.Operator = cell.Value
			case "target":
				condition.Target = cell.Value
			case "values":
				json.Unmarshal([]byte(cell.Value), &condition.Values)
			default:
				return fmt.Errorf("unexpected column name: %s", head[n].Value)
			}
		}
		rulesTestContext.rule.Audience.Conditions = append(*conditions, condition)
	}

	return nil
}


func theVariantValueOf(variantValue float64) error {
	rulesTestContext.variantValue = variantValue
	return nil
}

func theVariantSplitsAre(variantSplits *gherkin.DataTable) error {
	head := variantSplits.Rows[0].Cells
	splits := &rulesTestContext.rule.VariantSplits

	for i := 1; i < len(variantSplits.Rows); i++ {
		variantSplit := VariantSplit{}
		for n, cell := range variantSplits.Rows[i].Cells {
			switch head[n].Value {
			case "variantKey":
				variantSplit.VariantKey = cell.Value
			case "split":
				variantSplit.Split, _ = strconv.ParseFloat(cell.Value, 64)
			default:
				return fmt.Errorf("unexpected column name: %s", head[n].Value)
			}
		}
		rulesTestContext.rule.VariantSplits = append(*splits, variantSplit)
	}

	return nil
}

func theVariantSplitKeyIsCalculated() error {
	rulesTestContext.splitKeyVariantResult = getVariantSplitKey(
		rulesTestContext.rule.VariantSplits,
		rulesTestContext.variantValue,
	)
	return nil
}

func theResultingVariantShouldBe(variant string) error {
	if variant != rulesTestContext.splitKeyVariantResult{
		return fmt.Errorf("Expected %s to be %s", rulesTestContext.splitKeyVariantResult, variant)
	}
	return nil
}

func RulesFeatureContext(s *godog.Suite) {
	s.Step(`^the rule is a default rule$`, theRuleIsADefaultRule)
	s.Step(`^the rule is matched against the context$`, theRuleIsMatchedAgainstTheContext)
	s.Step(`^the result from the match should be (true|false)$`, theResultFromTheMatchShouldBe)
	s.Step(`^the context values are$`, theContextValuesAre)
	s.Step(`^the rule\'s audience conditions are$`, theRulesAudienceConditionsAre)
	s.Step(`^the variant value of (\d+)$`, theVariantValueOf)
	s.Step(`^the variant splits are$`, theVariantSplitsAre)
	s.Step(`^the variant split key is calculated$`, theVariantSplitKeyIsCalculated)
	s.Step(`^the resulting variant should be "([^"]*)"$`, theResultingVariantShouldBe)


	s.BeforeScenario(func(interface{}) {
		rulesTestContext = rulesTestContextType{
			rule: Rule{
				DefaultRule: false,
				VariantSplits:[]VariantSplit{},
				Audience: Audience{
					Conditions: []Condition{},
				},
			},
			context_builder: NewContextBuilder("anonymous"),
		}
	})
}
