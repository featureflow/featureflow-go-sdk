package featureflow_go_sdk

import (
	"log"
	"crypto/sha1"
	"fmt"
	"strconv"
)

func RuleMatches(rule Rule, context Context) bool{
	log.Println(context.GetValuesForKey("role"))
	if rule.DefaultRule == true{
		return true
	} else {
		for _, condition := range rule.Audience.Conditions{
			pass := false
			for _, value := range context.GetValuesForKey(condition.Target){
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

func GetVariantSplitKey(variant_splits []VariantSplit, variant_value float64) string{
	percent := 0.0

	for _, variant_split := range variant_splits {
		percent += variant_split.Split
		if percent >= variant_value {
			return variant_split.VariantKey
		}
	}

	return "off"
}

func CalculateHash(salt, feature, contextKey string) string{
	toHash := fmt.Sprintf("%s:%s:%s",salt, feature, contextKey)
	hasher := sha1.New()
	hasher.Write([]byte(toHash))
	return fmt.Sprintf("%x", hasher.Sum(nil))[0:15]
}

func GetVariantValue(hash string) float64{
	hashInt, _ := strconv.ParseInt(hash, 16, 64)
	return float64(hashInt % 100 + 1)
}