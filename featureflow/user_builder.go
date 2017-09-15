package featureflow

import (
	"errors"
	"time"
)

type Attribute interface{}

type User interface {
	GetId() string
	GetAttributes() map[string][]Attribute
	GetAttributesForKey(string) []Attribute
	GetAttributeKeys() []string
}

type UserBuilder interface {
	WithAttribute(string, Attribute) UserBuilder
	WithAttributes(string, []Attribute) UserBuilder
	Build() (User, error)
}

type userBuilder struct {
	id string
	attributes map[string][]Attribute
}

type user struct{
	id string `json:"id"`
	attributes map[string][]Attribute `json:"attributes"`
}

func (cb *userBuilder) WithAttribute(key string, attribute Attribute) UserBuilder {
	cb.attributes[key] = []Attribute{attribute}
	return cb
}

func (cb *userBuilder) WithAttributes(Key string, Attributes []Attribute) UserBuilder {
	cb.attributes[Key] = Attributes
	return cb
}

func (cb *userBuilder) Build() (User, error) {
	if len(cb.id) == 0 {
		return &user{}, errors.New("A user id is required")
	}
	return &user{
		id: cb.id,
		attributes: cb.attributes,
	}, nil
}

func NewUserBuilder(Id string) UserBuilder {
	attributes := make(map[string][]Attribute)
	attributes["featureflow.user.id"] = []Attribute{Id}
	attributes["featureflow.date"] = []Attribute{time.Now().Format(time.RFC3339)}
	return &userBuilder{
		id:Id,
		attributes: attributes,
	}
}

func (c *user) GetId() string {
	return c.id
}

func (c *user) GetAttributes() map[string][]Attribute {
	return c.attributes
}

func (c *user) GetAttributesForKey(key string) []Attribute {
	return c.attributes[key]
}

func (c *user) GetAttributeKeys() []string {
	attributeKeys := []string{}
	for key, _ := range c.attributes{
		attributeKeys = append(attributeKeys, key)
	}
	return attributeKeys
}