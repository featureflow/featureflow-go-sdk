package featureflow_go_sdk

import "errors"

type Value interface{}

type Context interface {
	GetKey() string
	GetValuesForKey(string) []Value
	GetValueKeys() []string
}

type ContextBuilder interface {
	WithValue(string, Value) ContextBuilder
	WithValues(string, []Value) ContextBuilder
	Build() Context
}

type contextBuilder struct {
	key string
	values map[string][]Value
}

type context struct{
	key string
	values map[string][]Value
}

func (cb *contextBuilder) WithValue(key string, value Value) ContextBuilder {
	cb.values[key] = []Value{value}
	return cb
}

func (cb *contextBuilder) WithValues(Key string, Values []Value) ContextBuilder {
	cb.values[Key] = Values
	return cb
}

func (cb *contextBuilder) Build() Context {
	return &context{
		key: cb.key,
		values: cb.values,
	}
}

func NewContextBuilder(Key string) (ContextBuilder, error) {
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