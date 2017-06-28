package featureflow

type Feature struct {
	Key string `json:"key"`
	Rules []Rule `json:"rules"`
	OffVariantKey string `json:"offVariantKey"`
}

type Rule struct {
	DefaultRule   bool `json:"defaultRule,omitempty"`
	Audience      Audience `json:"audience,omitempty"`
	VariantSplits []VariantSplit `json:"variantSplits"`
}

type Audience struct {
	Conditions []Condition `json:"conditions,omitempty"`
}

type Condition struct {
	Target string `json:"target"`
	Operator string `json:"operator"`
	Values []interface{} `json:"values"`
}

type VariantSplit struct {
	VariantKey string `json:"variantKey"`
	Split float64 `json:"split"`
}