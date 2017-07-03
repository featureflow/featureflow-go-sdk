package featureflow

type FeatureRegistration struct{
	Key string `json:"key"`
	FailoverVariant string `json:"failoverVariant"`
	Variants []FeatureRegistrationVariant `json:"variants,omitempty"`
}

type FeatureRegistrationVariant struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type FeatureRegistrationBuilder interface{
	AddVariant(string, string) FeatureRegistrationBuilder
	Build() FeatureRegistration
}

type featureRegistrationBuilder struct{
	featureRegistration FeatureRegistration
}

func (f *featureRegistrationBuilder) AddVariant(key, value string) FeatureRegistrationBuilder {
	f.featureRegistration.Variants = append(f.featureRegistration.Variants, FeatureRegistrationVariant{key, value})
	return f
}

func (f* featureRegistrationBuilder) Build() FeatureRegistration {
	if len(f.featureRegistration.Variants) < 2{
		f.featureRegistration.Variants = []FeatureRegistrationVariant{
			{"on","On"},
			{"off","Off"},
		}
	}
	return f.featureRegistration
}

func WithFeature(key, failover string) FeatureRegistrationBuilder{
	return &featureRegistrationBuilder{
		FeatureRegistration{
			Key: key,
			FailoverVariant: failover,
			Variants: []FeatureRegistrationVariant{},
		},
	}
}