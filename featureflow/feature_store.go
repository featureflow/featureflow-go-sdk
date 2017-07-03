package featureflow

import (
	"fmt"
	"sync"
)

type FeatureStore interface {
	Get(string) (*feature, error)
	Set(string, *feature) error
	SetAll(map[string]*feature) error
}

type inMemoryStore struct{
	features map[string]*feature
	sync.RWMutex
}

func (store *inMemoryStore) Get(key string) (*feature, error){
	store.RLock()
	defer store.RUnlock()

	feature := store.features[key]

	if feature != nil {
		return feature, nil
	} else {
		return nil, fmt.Errorf("feature %s was not found", key)
	}
}

func (store *inMemoryStore) Set(key string, feature *feature) error{
	store.Lock()
	defer store.Unlock()

	store.features[key] = feature
	return nil
}

func (store *inMemoryStore) SetAll(features map[string]*feature) error{
	store.Lock()
	defer store.Unlock()

	store.features = features
	return nil
}

func NewInMemoryStore() (*inMemoryStore, error){
	return &inMemoryStore{
		features: make(map[string]*feature),
	}, nil
}