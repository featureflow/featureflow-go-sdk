package featureflow

import (
	"time"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func newPollingClient(api_key string, url string, config *Config){
	var etag string = ""
	getFeatures(api_key, url, &etag, config)

	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <- ticker.C:
				getFeatures(api_key, url, &etag, config)
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func getFeatures(api_key string, url string, etag *string, config *Config){
	featureClient := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		config.Logger.Println(LOG_ERROR, err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api_key))
	req.Header.Set("If-None-Match", *etag)

	res, getErr := featureClient.Do(req)
	if getErr != nil {
		config.Logger.Println(LOG_ERROR, getErr)
		return
	}

	defer res.Body.Close()
	if res.StatusCode == 200 {
		*etag = res.Header.Get("ETag")

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			config.Logger.Println(LOG_ERROR, readErr)
		}

		var features map[string]feature

		json.Unmarshal(body, &features)

		featuresMap := make(map[string]*feature)

		for key, _ := range features{
			var f = features[key]
			featuresMap[key] = &f
		}

		config.Logger.Println(LOG_INFO, "updating features")
		config.FeatureStore.SetAll(featuresMap)
	} else if res.StatusCode >= 400{
		config.Logger.Println(LOG_ERROR, fmt.Sprintf("request for features failed with response status %d", res.StatusCode))
	}
}

