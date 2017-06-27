package featureflow_go_sdk

import "log"

func RuleMatches(rule Rule, context Context) bool{
	log.Println(context.GetValues("role"))
	if rule.DefaultRule == true{
		return true
	} else {
		for _, condition := range rule.Audience.Conditions{
			pass := false
			for _, value := range context.GetValues(condition.Target){
				if Test(condition.Operator, value, condition.Values){
					pass = true
					break
				}
			}
			if !pass {
				return false
			}
		}
		return true
	}
}
