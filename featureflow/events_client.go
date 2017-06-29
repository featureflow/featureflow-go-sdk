package featureflow

import (
	"net/http"
	"log"
	"fmt"
	"time"
	"bytes"
	"encoding/json"
)

func (e*EventsClient) registerFeaturesEvent(features []FeatureRegistration){
	if e.Disabled{
		return
	}
	body, _ := json.Marshal(features)
	go sendEvent(
		"Register Features",
		http.MethodPut,
		"https://app.featureflow.io/api/sdk/v1/register",
		e.ApiKey,
		body,
	)
}

func (e*EventsClient) evaluateEvent(key, evaluatedVariant, expectedVariant string, context Context){
	if e.Disabled{
		return
	}
	body, _ := json.Marshal(
		struct{
			FeatureKey string `json:"featureKey"`
			EvaluatedVariant string `json:"evaluatedVariant"`
			ExpectedVaraint string `json:"expectedVariant"`
			Context Context `json:"context"`
		}{key, evaluatedVariant, expectedVariant, context},
	)
	go sendEvent("evaluate", http.MethodPost, "https://app.featureflow.io/api/sdk/v1/events", e.ApiKey, body)
}


func sendEvent(event_name, method, url, api_key string, body []byte){
	client := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api_key))
	res, _ := client.Do(req)
	res.Body.Close()
}

type EventsClient struct{
	ApiKey string
	Disabled bool
}

func NewEventsClient(api_key string, disabled bool) EventsClient {
	return EventsClient{api_key, disabled}
}

