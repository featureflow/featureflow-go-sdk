package featureflow

import "fmt"

type FeatureflowClient struct{
	ApiKey string
	Config Config
	EventsClient EventsClient
	FailoverVariants map[string]string
}

type Config struct {
	Url string
	FeatureStore FeatureStore
	FeatureRegistrations []FeatureRegistration
	DisableEvents bool
}


func Client(api_key string, config Config) (*FeatureflowClient, error){
	//TODO LOG> Featureflow initializing

	if len(api_key) == 0{
		return nil, fmt.Errorf("Api_Key must exist")
	}

	if config.FeatureStore == nil{
		featureStore, _ := NewInMemoryStore()
		config.FeatureStore = featureStore
	}

	eventsClient := NewEventsClient(api_key, config.DisableEvents)

	url := "https://app.featureflow.io/api/sdk/v1/features"
	go newPollingClient(api_key, url, config.FeatureStore)

	failoverVariants := make(map[string]string)

	if config.FeatureRegistrations != nil{
		//TODO LOG> Featureflow registering features with featureflow
		eventsClient.registerFeaturesEvent(config.FeatureRegistrations)
		for _, registration := range config.FeatureRegistrations{
			failoverVariants[registration.Key] = registration.FailoverVariant
		}
	}
	//TODO LOG> Featureflow initialized
	return &FeatureflowClient{
		api_key,
		config,
		eventsClient,
		failoverVariants,
	}, nil
}


func (client *FeatureflowClient) Evaluate(key string, context Context) Evaluate {
	feature, error := client.Config.FeatureStore.Get(key)

	var evaluatedVariant string = "off"

	if error != nil{
		failover := client.FailoverVariants[key]
		if len(failover) > 0 {
			//TODO WARN> Using failover variant of {failover}
			evaluatedVariant = failover
		} else{
			//TODO WARN> Using default failover variant of "off"
		}
	} else{
		for _, rule := range feature.Rules {
			if ruleMatches(rule, context){
				variant_value := getVariantValue(calculateHash("1", key, context.GetKey()))
				evaluatedVariant = getVariantSplitKey(rule.VariantSplits, variant_value)
				break
			}
		}
	}

	return Evaluate{
		feature_key: key,
		evaluated_variant: evaluatedVariant,
		context: context,
		eventsClient: client.EventsClient,
	}
}