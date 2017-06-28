package featureflow

import (
	"crypto/sha1"
	"fmt"
	"strconv"
)

func ruleMatches(rule rule, context contextInterface) bool{
	if rule.DefaultRule == true{
		return true
	} else {
		for _, condition := range rule.Audience.Conditions{
			pass := false
			for _, value := range context.GetValuesForKey(condition.Target){
				if conditionsTest(condition.Operator, value, condition.Values){
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

func getVariantSplitKey(variant_splits []variantSplit, variant_value float64) string{
	percent := 0.0

	for _, variant_split := range variant_splits {
		percent += variant_split.Split
		if percent >= variant_value {
			return variant_split.VariantKey
		}
	}

	return "off"
}

func calculateHash(salt, feature, contextKey string) string{
	toHash := fmt.Sprintf("%s:%s:%s",salt, feature, contextKey)
	hasher := sha1.New()
	hasher.Write([]byte(toHash))
	return fmt.Sprintf("%x", hasher.Sum(nil))[0:15]
}

func getVariantValue(hash string) float64{
	hashInt, _ := strconv.ParseInt(hash, 16, 64)
	return float64(hashInt % 100 + 1)
}