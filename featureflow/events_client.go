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
	if e.Config.DisableEvents{
		return
	}
	body, _ := json.Marshal(features)
	go e.sendEvent(
		"Register Features",
		http.MethodPut,
		e.Config.BaseURL+"/api/sdk/v1/register",
		body,
	)
}

type evaluateEventType struct{
	FeatureKey string 			`json:"featureKey"`
	EventType string 			`json:"type"`
	EvaluatedVariant string 	`json:"evaluatedVariant"`
	ExpectedVariant string 		`json:"expectedVariant"`
	Timestamp time.Time 		`json:"timestamp""`
	User *User 					`json:"user"`
}

func (e*EventsClient) evaluateEvent(key, evaluatedVariant, expectedVariant string, user *User){
	if e.Config.DisableEvents{
		return
	}
	body, _ := json.Marshal(
		[]evaluateEventType{
			{key, "evaluate", evaluatedVariant, expectedVariant, time.Now(),user},
		},
	)
	go e.sendEvent("evaluate", http.MethodPost, e.Config.BaseURL+"/api/sdk/v1/events", body)
}


func (e*EventsClient) sendEvent(event_type, method, url string, body []byte){
	client := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.ApiKey))
	req.Header.Set("X-Featureflow-Client", fmt.Sprintf("GoClient/1.0.0"))
	res, _ := client.Do(req)

	if res != nil {
		if res.StatusCode >= 400{
			e.Config.Logger.Println(
				LOG_ERROR,
				fmt.Sprintf("unable to send event %s to %s. Failed with response status %d", event_type, url, res.StatusCode),
			)
		}
		res.Body.Close()
	} else{
		e.Config.Logger.Println(LOG_ERROR, "unable to send event %s to %s. Internal SDK error", event_type, url)
	}

}

type EventsClient struct{
	ApiKey string
	Config *Config
}

func NewEventsClient(api_key string, config *Config) EventsClient {
	return EventsClient{api_key, config}
}

