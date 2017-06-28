package featureflow


type Evaluate struct {
	evaluated_variant string
}

func (e Evaluate) Is(variant string) bool{
	return e.evaluated_variant == variant
}

func (e Evaluate) IsOn() bool{
	return e.Is("on")
}

func (e Evaluate) IsOff() bool{
	return e.Is("off")
}