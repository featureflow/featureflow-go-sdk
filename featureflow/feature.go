package featureflow

type feature struct {
	Key string `json:"key"`
	Rules []rule `json:"rules"`
	OffVariantKey string `json:"offVariantKey"`
	Enabled bool `json:"enabled"`
}

type rule struct {
	DefaultRule   bool `json:"defaultRule,omitempty"`
	Audience      audience `json:"audience,omitempty"`
	VariantSplits []variantSplit `json:"variantSplits"`
}

type audience struct {
	Conditions []condition `json:"conditions,omitempty"`
}

type condition struct {
	Target string `json:"target"`
	Operator string `json:"operator"`
	Values []interface{} `json:"values"`
}

type variantSplit struct {
	VariantKey string `json:"variantKey"`
	Split float64 `json:"split"`
}