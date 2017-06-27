package featureflow_go_sdk

type Value interface{}

type Context interface {
	GetKey() string
	GetValues(string) []Value
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

func NewContextBuilder(Key string) ContextBuilder {
	return &contextBuilder{key:Key, values: make(map[string][]Value)}
}

func (c *context) GetKey() string {
	return c.key
}

func (c *context) GetValues(key string) []Value {
	return c.values[key]
}