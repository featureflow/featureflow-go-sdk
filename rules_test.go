package featureflow_go_sdk

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"fmt"
	"strconv"
	"encoding/json"
	"log"
	"github.com/davecgh/go-spew/spew"
)


type RulesTestContext struct{
	result bool
	rule Rule
	context_builder ContextBuilder
	variantValue float64
}

var rulesTestContext RulesTestContext

func theRuleIsADefaultRule() error {
	rulesTestContext.rule.DefaultRule = true
	return nil
}

func theRuleIsMatchedAgainstTheContext() error {
	rulesTestContext.result = RuleMatches(
		rulesTestContext.rule,
		rulesTestContext.context_builder.Build(),
	)

	return nil
}

func theResultFromTheMatchShouldBe(resultStr string) error {
	result := resultStr == "true"
	//context := rulesTestContext.context_builder.Build()
	if result != rulesTestContext.result{
		return fmt.Errorf("Expected %s to be %s",
			strconv.FormatBool(result),
			strconv.FormatBool(rulesTestContext.result),
		)
	}
	return nil
}

func theContextValuesAre(contextValuesTable *gherkin.DataTable) error {

	spew.Dump(contextValuesTable)

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
				log.Println("the fick", string(cell.Value))
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
			log.Println("---$$---jsaklf-------------", values)
			rulesTestContext.context_builder = rulesTestContext.context_builder.WithValues(key, values)
		} else {
			log.Println("!!!-------------------", value)
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
	return godog.ErrPending
}

func theResultingVariantShouldBe(arg1 string) error {
	return godog.ErrPending
}

func RulesFeatureContext(s *godog.Suite) {
	spew.Dump(s)
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

		rulesTestContext = RulesTestContext{
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
