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
	BaseURL string
	FeatureStore FeatureStore
	WithFeatures []FeatureRegistration
	DisableEvents bool
	Logger *log.Logger
}


func Client(api_key string, config Config) (*FeatureflowClient, error){
	if config.Logger == nil{
		config.Logger = log.New(os.Stderr, "Featureflow:", log.LstdFlags)
	}
	config.Logger.Println(LOG_INFO, "initializing client")

	if len(api_key) == 0{
		return nil, fmt.Errorf("Api_Key must exist")
	}

	if config.FeatureStore == nil{
		featureStore, _ := NewInMemoryStore()
		config.FeatureStore = featureStore
	}

	if config.BaseURL == ""{
		config.BaseURL = "https://app.featureflow.io"
	}

	config.Logger.Println(LOG_INFO, fmt.Sprintf("Connecting to %s", config.BaseURL))

	eventsClient := NewEventsClient(api_key, &config)

	newPollingClient(api_key, config.BaseURL+"/api/sdk/v1/features", &config)

	failoverVariants := make(map[string]string)

	if config.WithFeatures != nil{
		eventsClient.registerFeaturesEvent(config.WithFeatures)
		for _, registration := range config.WithFeatures{
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


func (client *FeatureflowClient) EvaluateBasic(key, userId string) Evaluate{
	user, _ := NewUserBuilder(userId).Build()
	return client.Evaluate(key, user)
}

func (client *FeatureflowClient) Evaluate(key string, user *User) Evaluate {
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
		if feature.Enabled {
			for _, rule := range feature.Rules {
				if ruleMatches(rule, user){
					variant_value := getVariantValue(calculateHash("1", key, user.GetId()))
					evaluatedVariant = getVariantSplitKey(rule.VariantSplits, variant_value)
					break
				}
			}
		} else {
			evaluatedVariant = feature.OffVariantKey
		}

	}

	return Evaluate{
		feature_key: key,
		evaluated_variant: evaluatedVariant,
		user: user,
		eventsClient: client.EventsClient,
	}
}