package featureflow

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

type FeatureflowClient struct{
	ApiKey string
	Config Config

}

func (client* FeatureflowClient) getFeature(key string) (*Feature, error){
	var feature Feature

	return &feature, nil
}


type Config struct {
	Uri string
	FeatureStore FeatureStore
}

func Client(api_key string, config Config) (*FeatureflowClient, error){

	if config.FeatureStore == nil{
		featureStore, _ := NewInMemoryStore()
		config.FeatureStore = featureStore

		client := &http.Client{
			Timeout: time.Second * 10,
		}
		req, _ := http.NewRequest("GET", "https://app.featureflow.io/api/sdk/v1/features", nil)

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api_key))
		res, _ := client.Do(req)
		defer res.Body.Close()

		var features map[string]Feature

		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, features)

		log.Println(features)
	}

	return &FeatureflowClient{
		api_key,
		config,
	}, nil
}

func (client *FeatureflowClient) Evaluate(key string, context Context) Evaluate {
	feature, error := client.Config.FeatureStore.Get(key)

	if error != nil{
		return Evaluate{"off"}
	}

	for _, rule := range feature.Rules {
		if ruleMatches(rule, context){
			variant_value := getVariantValue(calculateHash("1", key, context.GetKey()))
			return Evaluate{getVariantSplitKey(rule.VariantSplits, variant_value)}
		}
	}

	return Evaluate{
		evaluated_variant: "off",
	}
}