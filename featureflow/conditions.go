package featureflow

import (
	"reflect"
	"strings"
	"regexp"
	"time"
)

// Every operator here is handed a user attribute as an interface{} — whatever the caller passed
// to WithAttribute — and a rule value decoded from JSON. Neither side can be trusted to be the
// type the operator wants, so nothing in this file asserts a type without checking. An operator
// that cannot make sense of its operands returns false: a mistargeted rule must fail to match,
// never panic into the host application's request.
func conditionsTest(op string, a interface{}, b []interface{}) bool {
	switch op {
	case "in":
		return in(a, b)
	case "notIn":
		return !in(a, b)
	}

	// Every remaining operator compares against the first condition value only. A condition that
	// arrives with no values is malformed rather than fatal, so fail closed instead of indexing
	// b[0] and panicking.
	if len(b) == 0 {
		return false
	}

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

// bothStrings reports whether both operands are strings, returning them ready to compare.
func bothStrings(a interface{}, b interface{}) (string, string, bool){
	aStr, aOk := a.(string)
	bStr, bOk := b.(string)
	return aStr, bStr, aOk && bOk
}

// toFloat64 widens any of Go's numeric types to float64.
//
// A rule value decoded from JSON is always a float64, but a user attribute is whatever the caller
// wrote: WithAttribute("age", 20) stores an int, which is the idiomatic Go — and what the README
// teaches. Comparing the two without widening either never matches or, if asserted, panics.
func toFloat64(value interface{}) (float64, bool){
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

// bothNumbers reports whether both operands are numeric, returning them widened to float64.
func bothNumbers(a interface{}, b interface{}) (float64, float64, bool){
	aNum, aOk := toFloat64(a)
	bNum, bOk := toFloat64(b)
	return aNum, bNum, aOk && bOk
}

func equals(a interface{}, b interface{}) bool{
	// Numbers are compared by value rather than by Go type. Interface equality is false for
	// int(20) against float64(20) however equal the numbers are, and the rule side is always a
	// float64 because it came from JSON.
	if aNum, bNum, ok := bothNumbers(a, b); ok {
		return aNum == bNum
	}
	// == on interfaces panics when the dynamic type is uncomparable — a slice, map or function,
	// any of which WithAttribute will happily accept — so check before comparing.
	if a == nil || b == nil {
		return a == b
	}
	if !reflect.TypeOf(a).Comparable() || !reflect.TypeOf(b).Comparable() {
		return false
	}
	return a == b
}

func contains(a interface{}, b interface{}) bool{
	if aStr, bStr, ok := bothStrings(a, b); ok {
		return strings.Contains(aStr, bStr)
	}
	return false
}

func startsWith(a interface{}, b interface{}) bool{
	if aStr, bStr, ok := bothStrings(a, b); ok {
		return strings.HasPrefix(aStr, bStr)
	}
	return false
}

func endsWith(a interface{}, b interface{}) bool{
	if aStr, bStr, ok := bothStrings(a, b); ok {
		return strings.HasSuffix(aStr, bStr)
	}
	return false
}

func matches(a interface{}, b interface{}) bool{
	aStr, bStr, ok := bothStrings(a, b)
	if !ok {
		return false
	}
	// An invalid pattern authored in the dashboard must not take the caller down either.
	if matched, err := regexp.MatchString(bStr, aStr); err == nil {
		return matched
	} else {
		return false
	}
}

func in(a interface{}, b []interface{}) bool{
	for _, bVar := range b {
		// equals, not ==, so that a numeric attribute matches a numeric list and an uncomparable
		// attribute returns false rather than panicking.
		if equals(a, bVar) {
			return true
		}
	}
	return false
}

//Numerics
func greaterThan(a interface{}, b interface{}) bool{
	if aNum, bNum, ok := bothNumbers(a, b); ok {
		return aNum > bNum
	}
	return false
}

func greaterThanOrEqual(a interface{}, b interface{}) bool{
	if aNum, bNum, ok := bothNumbers(a, b); ok {
		return aNum >= bNum
	}
	return false
}

func lessThan(a interface{}, b interface{}) bool{
	if aNum, bNum, ok := bothNumbers(a, b); ok {
		return aNum < bNum
	}
	return false
}

func lessThanOrEqual(a interface{}, b interface{}) bool{
	if aNum, bNum, ok := bothNumbers(a, b); ok {
		return aNum <= bNum
	}
	return false
}

//Dates
// dateLayouts are tried in order against either side of a date condition.
//
// The date-only form is the one that matters in practice: the dashboard's date picker emits
// "2026-07-03", and JavaScript's Date.parse — which the JS SDK and sdk-server both rely on —
// reads that as UTC midnight. Accepting only full timestamps here meant a scheduled rollout
// authored in the UI silently never fired for Go services while firing everywhere else.
//
// time.Parse applies UTC when the layout carries no zone, so the date-only layout gives UTC
// midnight without any further work.
var dateLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02",
}

func parseConditionDate(value interface{}) (time.Time, bool){
	str, ok := value.(string)
	if !ok {
		return time.Time{}, false
	}
	for _, layout := range dateLayouts {
		if parsed, err := time.Parse(layout, str); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

// Compared by instant (RFC3339), not lexicographically — a naive string comparison
// mishandles non-UTC offsets, e.g. "02:00-05:00" (07:00Z) sorts before "06:00Z" as a
// string even though it's the later instant.
func before(a interface{}, b interface{}) bool{
	aTime, aOk := parseConditionDate(a)
	bTime, bOk := parseConditionDate(b)
	if !aOk || !bOk {
		return false
	}
	return aTime.Before(bTime)
}

func after(a interface{}, b interface{}) bool{
	aTime, aOk := parseConditionDate(a)
	bTime, bOk := parseConditionDate(b)
	if !aOk || !bOk {
		return false
	}
	return aTime.After(bTime)
}
