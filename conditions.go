package featureflow_go_sdk

import (
	"strings"
	"regexp"
	"strconv"
)

func Test(op string, a interface{}, b []interface{}) bool {
	switch op {
	case "equals":
		return equals(a, b[0])
	case "contains":
		return contains(a, b[0])
	case "startsWith":
		return startsWith(a, b[0])
	case "endsWith":
		return endsWith(a, b[0])
	case "matches":
		return matches(a, b[0])
	case "in":
		return in(a, b)
	case "notIn":
		return !in(a,b)
	case "greaterThan":
		return greaterThan(a, b[0])
	case "greaterThanOrEqual":
		return greaterThanOrEqual(a, b[0])
	case "lessThan":
		return lessThan(a, b[0])
	case "lessThanOrEqual":
		return lessThanOrEqual(a, b[0])
	case "before":
		return before(a, b[0])
	case "after":
		return after(a, b[0])
	default:
		return false
	}
}

func equals(a interface{}, b interface{}) bool{
	return a == b
}

func contains(a interface{}, b interface{}) bool{
	return strings.Contains(a.(string), b.(string))
}

func startsWith(a interface{}, b interface{}) bool{
	return strings.HasPrefix(a.(string), b.(string))
}

func endsWith(a interface{}, b interface{}) bool{
	return strings.HasSuffix(a.(string), b.(string))
}

func matches(a interface{}, b interface{}) bool{
	if matched, err := regexp.MatchString(b.(string), a.(string)); err == nil {
		return matched
	} else {
		return false
	}
}

func in(a interface{}, b []interface{}) bool{
	for _, bVar := range b {
		if bVar == a {
			return true
		}
	}
	return false
}

//Numerics
func greaterThan(aStr interface{}, bStr interface{}) bool{
	a, _ := strconv.ParseFloat(aStr.(string), 64)
	b, _ := strconv.ParseFloat(bStr.(string), 64)
	return a > b
}

func greaterThanOrEqual(aStr interface{}, bStr interface{}) bool{
	a, _ := strconv.ParseFloat(aStr.(string), 64)
	b, _ := strconv.ParseFloat(bStr.(string), 64)
	return a >= b
}

func lessThan(aStr interface{}, bStr interface{}) bool{
	a, _ := strconv.ParseFloat(aStr.(string), 64)
	b, _ := strconv.ParseFloat(bStr.(string), 64)
	return a < b
}

func lessThanOrEqual(aStr interface{}, bStr interface{}) bool{
	a, _ := strconv.ParseFloat(aStr.(string), 64)
	b, _ := strconv.ParseFloat(bStr.(string), 64)
	return a <= b
}

//Dates
func before(a interface{}, b interface{}) bool{
	return a.(string) < b.(string)
}

func after(a interface{}, b interface{}) bool{
	return a.(string) > b.(string)
}



