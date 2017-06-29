package featureflow

import (
	"fmt"
	"sync"
)

type FeatureStore interface {
	Get(string) (*Feature, error)
	Set(string, *Feature) error
	SetAll(map[string]*Feature) error
}

type inMemoryStore struct{
	features map[string]*Feature
	sync.RWMutex
}

func (store *inMemoryStore) Get(key string) (*Feature, error){
	store.RLock()
	defer store.RUnlock()

	feature := store.features[key]

	if feature != nil {
		return feature, nil
	} else {
		return nil, fmt.Errorf("feature %s was not found", key)
	}
}

func (store *inMemoryStore) Set(key string, feature *Feature) error{
	store.Lock()
	defer store.Unlock()

	store.features[key] = feature
	return nil
}

func (store *inMemoryStore) SetAll(features map[string]*Feature) error{
	store.Lock()
	defer store.Unlock()

	store.features = features
	return nil
}

func NewInMemoryStore() (*inMemoryStore, error){
	return &inMemoryStore{
		features: make(map[string]*Feature),
	}, nil
}