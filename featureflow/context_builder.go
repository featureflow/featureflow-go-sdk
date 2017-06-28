package featureflow

import "errors"

type Value interface{}

type contextInterface interface {
	GetKey() string
	GetValuesForKey(string) []Value
	GetValueKeys() []string
}

type contextBuilderInterface interface {
	WithValue(string, Value) contextBuilderInterface
	WithValues(string, []Value) contextBuilderInterface
	Build() contextInterface
}

type contextBuilder struct {
	key string
	values map[string][]Value
}

type context struct{
	key string
	values map[string][]Value
}

func (cb *contextBuilder) WithValue(key string, value Value) contextBuilderInterface {
	cb.values[key] = []Value{value}
	return cb
}

func (cb *contextBuilder) WithValues(Key string, Values []Value) contextBuilderInterface {
	cb.values[Key] = Values
	return cb
}

func (cb *contextBuilder) Build() contextInterface {
	return &context{
		key: cb.key,
		values: cb.values,
	}
}

func NewContextBuilder(Key string) (contextBuilderInterface, error) {
	if len(Key) == 0 {
		return nil, errors.New("A key must be defined")
	}
	return &contextBuilder{key:Key, values: make(map[string][]Value)}, nil
}

func (c *context) GetKey() string {
	return c.key
}

func (c *context) GetValuesForKey(key string) []Value {
	return c.values[key]
}

func (c *context) GetValueKeys() []string {
	valueKeys := []string{}
	for key, _ := range c.values{
		valueKeys = append(valueKeys, key)
	}
	return valueKeys
}