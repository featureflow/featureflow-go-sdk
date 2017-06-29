package featureflow

import (
	"fmt"
	"log"
	"os"
)

const LOG_INFO = "[info]"
const LOG_ERROR = "[error]"

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
	Logger *log.Logger
}


func Client(api_key string, config Config) (*FeatureflowClient, error){
	if (config.Logger == nil){
		config.Logger = log.New(os.Stderr, "Featureflow:", log.LstdFlags)
	}
	//TODO LOG> Featureflow initializing
	config.Logger.Println(LOG_INFO, "initializing client")

	if len(api_key) == 0{
		return nil, fmt.Errorf("Api_Key must exist")
	}

	if config.FeatureStore == nil{
		featureStore, _ := NewInMemoryStore()
		config.FeatureStore = featureStore
	}

	eventsClient := NewEventsClient(api_key, &config)

	url := "https://app.featureflow.io/api/sdk/v1/features"
	go newPollingClient(api_key, url, &config)

	failoverVariants := make(map[string]string)

	if config.FeatureRegistrations != nil{
		eventsClient.registerFeaturesEvent(config.FeatureRegistrations)
		for _, registration := range config.FeatureRegistrations{
			config.Logger.Println(LOG_INFO, "Registering feature with key " + registration.Key)
			failoverVariants[registration.Key] = registration.FailoverVariant
		}
	}
	config.Logger.Println(LOG_INFO, "client initialized")
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
		providedOrDefaultString := "default"
		if len(failover) > 0 {
			evaluatedVariant = failover
			providedOrDefaultString = "provided"
		}
		client.Config.Logger.Println(
			LOG_INFO,
			fmt.Sprintf(
				"Evaluating nil feature '%s' using the %s failover '%s'",
				key,
				providedOrDefaultString,
				evaluatedVariant,
			),
		)
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