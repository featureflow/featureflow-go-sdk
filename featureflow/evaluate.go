package featureflow


type Evaluate struct {
	feature_key string
	evaluated_variant string
	user User
	eventsClient EventsClient
}

func (e Evaluate) Is(variant string) bool{
	e.eventsClient.evaluateEvent(e.feature_key, e.evaluated_variant, variant, e.user)
	return e.evaluated_variant == variant
}

func (e Evaluate) IsOn() bool{
	return e.Is("on")
}

func (e Evaluate) IsOff() bool{
	return e.Is("off")
}

func (e Evaluate) Value() string{
	return e.evaluated_variant
}