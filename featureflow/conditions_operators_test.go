package featureflow

import (
	"testing"
	"time"
)

// Unit tests for the condition operators. The shared cross-SDK scenarios in
// ../testbed/gherkin/conditions.feature pin the happy paths; these pin the edges that a rule
// author or a caller can reach by accident — wrong types, uncomparable types, malformed
// conditions and the date forms the dashboard produces. The overriding requirement is that none
// of them panic: an operator runs inside the caller's request goroutine, so a panic here takes
// down their request, not ours.

// allOperators is every operator conditionsTest recognises, plus one it does not.
var allOperators = []string{
	"equals",
	"contains",
	"startsWith",
	"endsWith",
	"matches",
	"in",
	"notIn",
	"greaterThan",
	"greaterThanOrEqual",
	"lessThan",
	"lessThanOrEqual",
	"before",
	"after",
	"noSuchOperator",
}

// runWithoutPanic runs fn, failing the test rather than unwinding if it panics.
func runWithoutPanic(t *testing.T, description string, fn func() bool) bool {
	t.Helper()

	result := false
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				t.Fatalf("%s panicked: %v", description, recovered)
			}
		}()
		result = fn()
	}()
	return result
}

// TestOperatorsWithMismatchedTypes covers every operator against operands it cannot make sense
// of. None may panic, and all must fail to match — except notIn, whose honest answer for an
// attribute that matches nothing in the list is true.
func TestOperatorsWithMismatchedTypes(t *testing.T) {
	operands := []struct {
		name      string
		attribute interface{}
		value     interface{}
	}{
		{"string attribute against a numeric rule value", "gold", 9000.0},
		{"numeric attribute against a string rule value", 9000.0, "gold"},
		{"int attribute against a string rule value", 20, "gold"},
		{"bool attribute", true, "gold"},
		{"nil attribute", nil, "gold"},
		{"nil rule value", "gold", nil},
		{"nil on both sides", nil, nil},
		{"slice attribute, which is uncomparable", []string{"admin"}, "gold"},
		{"map attribute, which is uncomparable", map[string]string{"tier": "gold"}, "gold"},
		{"struct attribute", struct{ A int }{1}, "gold"},
		{"timestamp attribute against a numeric rule value", "2026-07-03T00:00:00Z", 9000.0},
	}

	for _, op := range allOperators {
		for _, operand := range operands {
			t.Run(op+"/"+operand.name, func(t *testing.T) {
				got := runWithoutPanic(t, op, func() bool {
					return conditionsTest(op, operand.attribute, []interface{}{operand.value})
				})

				want := op == "notIn"
				// nil equals nil, and is therefore in a list containing nil — the
				// equality-based operators are the one exception to "no match".
				if operand.attribute == nil && operand.value == nil {
					switch op {
					case "equals", "in":
						want = true
					case "notIn":
						want = false
					}
				}

				if got != want {
					t.Errorf(
						"conditionsTest(%q, %#v, [%#v]) = %v, want %v",
						op, operand.attribute, operand.value, got, want,
					)
				}
			})
		}
	}
}

// TestOperatorsWithNoConditionValues covers a malformed condition that carries no values at all.
// Every operator but in/notIn reads the first value, so this used to be an index-out-of-range
// panic waiting on a bad rule.
func TestOperatorsWithNoConditionValues(t *testing.T) {
	for _, op := range allOperators {
		t.Run(op, func(t *testing.T) {
			got := runWithoutPanic(t, op, func() bool {
				return conditionsTest(op, "gold", []interface{}{})
			})

			// Nothing is in an empty list, so notIn is true and everything else is false.
			want := op == "notIn"
			if got != want {
				t.Errorf("conditionsTest(%q, \"gold\", []) = %v, want %v", op, got, want)
			}
		})
	}
}

// TestNumericOperatorsCoerceAttributeTypes covers the case the README teaches:
// WithAttribute("age", 20) stores an int, while the rule value arrives from JSON as a float64.
// Every Go numeric type must widen to float64 and compare on value.
func TestNumericOperatorsCoerceAttributeTypes(t *testing.T) {
	// The rule is "greaterThan 18", with 18 as a float64 exactly as the poller decodes it.
	ruleValue := []interface{}{float64(18)}

	tests := []struct {
		name      string
		attribute interface{}
		want      bool
	}{
		{"int", 20, true},
		{"int8", int8(20), true},
		{"int16", int16(20), true},
		{"int32", int32(20), true},
		{"int64", int64(20), true},
		{"uint", uint(20), true},
		{"uint8", uint8(20), true},
		{"uint16", uint16(20), true},
		{"uint32", uint32(20), true},
		{"uint64", uint64(20), true},
		{"float32", float32(20), true},
		{"float64", float64(20), true},
		{"int below the threshold", 17, false},
		{"int equal to the threshold", 18, false},
		{"int64 below the threshold", int64(17), false},
		{"float64 fractionally above the threshold", 18.5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runWithoutPanic(t, "greaterThan", func() bool {
				return conditionsTest("greaterThan", tt.attribute, ruleValue)
			})
			if got != tt.want {
				t.Errorf("greaterThan(%#v, 18.0) = %v, want %v", tt.attribute, got, tt.want)
			}
		})
	}
}

// TestValueOperatorsCompareNumbersAcrossTypes covers the remaining operators that have to see
// through the Go type of a numeric attribute, including the equality-based ones.
func TestValueOperatorsCompareNumbersAcrossTypes(t *testing.T) {
	tests := []struct {
		name      string
		op        string
		attribute interface{}
		values    []interface{}
		want      bool
	}{
		{"greaterThanOrEqual matches an equal int", "greaterThanOrEqual", 18, []interface{}{float64(18)}, true},
		{"greaterThanOrEqual rejects a lesser int", "greaterThanOrEqual", 17, []interface{}{float64(18)}, false},
		{"lessThan matches a lesser int64", "lessThan", int64(17), []interface{}{float64(18)}, true},
		{"lessThan rejects an equal int64", "lessThan", int64(18), []interface{}{float64(18)}, false},
		{"lessThanOrEqual matches an equal int64", "lessThanOrEqual", int64(18), []interface{}{float64(18)}, true},
		{"equals matches an int against a float64 rule value", "equals", 20, []interface{}{float64(20)}, true},
		{"equals rejects a different number", "equals", 20, []interface{}{float64(21)}, false},
		{"equals still compares strings exactly", "equals", "gold", []interface{}{"gold"}, true},
		{"equals does not conflate a number with its string form", "equals", 20, []interface{}{"20"}, false},
		{"in matches an int against numeric rule values", "in", 20, []interface{}{float64(18), float64(20)}, true},
		{"in rejects an absent number", "in", 19, []interface{}{float64(18), float64(20)}, false},
		{"notIn matches an absent number", "notIn", 19, []interface{}{float64(18), float64(20)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runWithoutPanic(t, tt.op, func() bool {
				return conditionsTest(tt.op, tt.attribute, tt.values)
			})
			if got != tt.want {
				t.Errorf("conditionsTest(%q, %#v, %#v) = %v, want %v", tt.op, tt.attribute, tt.values, got, tt.want)
			}
		})
	}
}

// TestStringOperators pins the string operators against well-typed operands, so that the
// defensive type checks added around them cannot quietly turn a real match into a non-match.
func TestStringOperators(t *testing.T) {
	tests := []struct {
		name      string
		op        string
		attribute interface{}
		value     interface{}
		want      bool
	}{
		{"contains matches a substring", "contains", "my-test-value", "-test-v", true},
		{"contains rejects an absent substring", "contains", "my-test-value", "nope", false},
		{"startsWith matches a prefix", "startsWith", "my-test-value", "my-test", true},
		{"startsWith rejects a non-prefix", "startsWith", "my-test-value", "-test-v", false},
		{"endsWith matches a suffix", "endsWith", "my-test-value", "test-value", true},
		{"endsWith rejects a non-suffix", "endsWith", "my-test-value", "-test-v", false},
		{"matches applies a pattern", "matches", "my-test-value", "my.test.+", true},
		{"matches rejects a non-matching pattern", "matches", "my-test-value", "wont-match", false},
		{"matches survives an invalid pattern", "matches", "my-test-value", "([unclosed", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runWithoutPanic(t, tt.op, func() bool {
				return conditionsTest(tt.op, tt.attribute, []interface{}{tt.value})
			})
			if got != tt.want {
				t.Errorf("conditionsTest(%q, %#v, %#v) = %v, want %v", tt.op, tt.attribute, tt.value, got, tt.want)
			}
		})
	}
}

// TestDateOnlyValueIsUTCMidnight pins the interpretation of a date-only value. The dashboard's
// date picker emits this form, and JavaScript's Date.parse — which the JS SDK and sdk-server rely
// on — reads it as UTC midnight. Any other reading shifts a scheduled rollout by up to a day
// relative to the rest of the platform.
func TestDateOnlyValueIsUTCMidnight(t *testing.T) {
	parsed, ok := parseConditionDate("2026-07-03")
	if !ok {
		t.Fatal(`parseConditionDate("2026-07-03") failed to parse`)
	}

	want := time.Date(2026, time.July, 3, 0, 0, 0, 0, time.UTC)
	if !parsed.Equal(want) {
		t.Errorf(`parseConditionDate("2026-07-03") = %s, want %s`, parsed, want)
	}
	if parsed.Location() != time.UTC {
		t.Errorf(`parseConditionDate("2026-07-03") is in %s, want UTC`, parsed.Location())
	}
}

// TestDateOperators covers before/after across every date form that reaches a rule, on either
// side of the comparison, plus the values that cannot be parsed at all.
func TestDateOperators(t *testing.T) {
	tests := []struct {
		name      string
		op        string
		attribute interface{}
		value     interface{}
		want      bool
	}{
		// A date-only rule value means UTC midnight, so the boundary falls between 23:59:59Z on
		// the previous day and 00:00:01Z on the day itself.
		{"a second past UTC midnight is after the date-only day", "after", "2026-07-03T00:00:01Z", "2026-07-03", true},
		{"a second before UTC midnight is not after the date-only day", "after", "2026-07-02T23:59:59Z", "2026-07-03", false},
		{"a second before UTC midnight is before the date-only day", "before", "2026-07-02T23:59:59Z", "2026-07-03", true},
		{"a second past UTC midnight is not before the date-only day", "before", "2026-07-03T00:00:01Z", "2026-07-03", false},
		{"UTC midnight itself is neither before nor after", "before", "2026-07-03T00:00:00Z", "2026-07-03", false},
		{"UTC midnight itself is not after either", "after", "2026-07-03T00:00:00Z", "2026-07-03", false},

		// Date-only on the attribute side, which is what a caller sets by hand.
		{"date-only attribute before a timestamp", "before", "2026-07-03", "2026-07-03T00:00:01Z", true},
		{"date-only attribute after a timestamp", "after", "2026-07-03", "2026-07-02T23:59:59Z", true},

		// Date-only on both sides.
		{"date-only both sides, earlier", "before", "2026-07-02", "2026-07-03", true},
		{"date-only both sides, equal", "before", "2026-07-03", "2026-07-03", false},
		{"date-only both sides, later", "after", "2026-07-04", "2026-07-03", true},

		// The full timestamp forms must keep working.
		{"fractional seconds", "before", "2017-03-09T02:39:46.182Z", "2017-04-09T02:39:46.182Z", true},
		{"no fractional seconds", "before", "2026-07-21T13:00:00Z", "2026-07-21T14:00:00Z", true},
		{"mixed fractional and whole seconds", "after", "2026-07-21T14:00:00.500Z", "2026-07-21T14:00:00Z", true},
		// Offsets compare by instant, not lexicographically: 02:00-05:00 is 07:00Z.
		{"an offset later than UTC compares by instant", "after", "2026-07-21T02:00:00-05:00", "2026-07-21T06:00:00Z", true},
		{"an offset earlier than UTC compares by instant", "before", "2026-07-21T23:00:00+10:00", "2026-07-21T14:00:00Z", true},

		// Unparseable values fail to match rather than erroring or panicking.
		{"unparseable attribute", "after", "not-a-date", "2026-07-03", false},
		{"unparseable attribute, before", "before", "not-a-date", "2026-07-03", false},
		{"unparseable rule value", "after", "2026-07-03T00:00:01Z", "the third of July", false},
		{"empty attribute", "before", "", "2026-07-03", false},
		{"month and day out of range", "before", "2026-13-45", "2026-07-03", false},
		{"a date with a time but no zone", "before", "2026-07-03T09:00:00", "2026-07-04", false},
		{"numeric attribute", "after", 1751500800000.0, "2026-07-03", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runWithoutPanic(t, tt.op, func() bool {
				return conditionsTest(tt.op, tt.attribute, []interface{}{tt.value})
			})
			if got != tt.want {
				t.Errorf("conditionsTest(%q, %#v, %#v) = %v, want %v", tt.op, tt.attribute, tt.value, got, tt.want)
			}
		})
	}
}

// TestRuleMatchesWithIntAttributeDoesNotPanic exercises the whole rule-matching path rather than
// a bare operator, since that is where a panic would actually escape into a caller's request.
func TestRuleMatchesWithIntAttributeDoesNotPanic(t *testing.T) {
	user, err := NewUserBuilder("user-123").
		WithAttribute("age", 20).
		WithAttribute("tier", "gold").
		Build()
	if err != nil {
		t.Fatalf("building user: %v", err)
	}

	overEighteen := rule{
		Audience: audience{
			Conditions: []condition{
				{Target: "age", Operator: "greaterThan", Values: []interface{}{float64(18)}},
			},
		},
	}
	scheduled := rule{
		Audience: audience{
			Conditions: []condition{
				{Target: "featureflow.date", Operator: "after", Values: []interface{}{"2020-01-01"}},
			},
		},
	}
	mistyped := rule{
		Audience: audience{
			Conditions: []condition{
				{Target: "tier", Operator: "greaterThan", Values: []interface{}{float64(18)}},
			},
		},
	}

	tests := []struct {
		name string
		rule rule
		want bool
	}{
		{"an int attribute matches a numeric rule", overEighteen, true},
		{"a date-only rule value fires against featureflow.date", scheduled, true},
		{"a string attribute under a numeric operator fails to match", mistyped, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runWithoutPanic(t, "ruleMatches", func() bool {
				return ruleMatches(tt.rule, user)
			})
			if got != tt.want {
				t.Errorf("ruleMatches(%s) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
