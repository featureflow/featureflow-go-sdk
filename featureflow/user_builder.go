package featureflow

import (
	"errors"
	"time"
)

type Attribute interface{}


type UserBuilder interface {
	WithAttribute(string, Attribute) UserBuilder
	WithAttributes(string, []Attribute) UserBuilder
	Build() (*User, error)
}

type userBuilder struct {
	id string
	attributes map[string][]Attribute
}

type User struct{
	Id string `json:"id"`
	Attributes map[string][]Attribute `json:"attributes"`
}

func (cb *userBuilder) WithAttribute(key string, attribute Attribute) UserBuilder {
	cb.attributes[key] = []Attribute{attribute}
	return cb
}

func (cb *userBuilder) WithAttributes(Key string, Attributes []Attribute) UserBuilder {
	cb.attributes[Key] = Attributes
	return cb
}

func (cb *userBuilder) Build() (*User, error) {
	if len(cb.id) == 0 {
		return &User{}, errors.New("A user id is required")
	}
	return &User{
		Id: cb.id,
		Attributes: cb.attributes,
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

func (c *User) GetId() string {
	return c.Id
}

func (c *User) GetAttributes() map[string][]Attribute {
	return c.Attributes
}

func (c *User) GetAttributesForKey(key string) []Attribute {
	return c.Attributes[key]
}

func (c *User) GetAttributeKeys() []string {
	attributeKeys := []string{}
	for key, _ := range c.Attributes{
		attributeKeys = append(attributeKeys, key)
	}
	return attributeKeys
}