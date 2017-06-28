package featureflow

import "fmt"

type FeatureStore interface {
	Get(string) (*Feature, error)
	Add(string, *Feature) error
	Clear() error
}

type inMemoryStore struct{
	features map[string]*Feature
}

func (store *inMemoryStore) Get(key string) (*Feature, error){
	if feature, ok := store.features[key]; ok {
		return feature, nil
	} else {
		return feature, fmt.Errorf("feature %s was not found", key)
	}
}

func (store *inMemoryStore) Add(key string, feature *Feature) error{
	store.features[key] = feature
	return nil
}

func (store *inMemoryStore) Clear() error{
	store.features = make(map[string]*Feature)
	return nil
}

func NewInMemoryStore() (FeatureStore, error){
	return &inMemoryStore{
		features: make(map[string]*Feature),
	}, nil
}