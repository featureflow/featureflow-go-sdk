package featureflow

import (
	"time"
	"log"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func newPollingClient(api_key string, url string, featureStore FeatureStore){
	var etag string = ""
	go getFeatures(api_key, url, &etag, featureStore)

	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <- ticker.C:
				getFeatures(api_key, url, &etag, featureStore)
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func getFeatures(api_key string, url string, etag *string, featureStore FeatureStore){
	featureClient := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api_key))
	req.Header.Set("If-None-Match", *etag)

	res, getErr := featureClient.Do(req)
	if getErr != nil {
		log.Println(getErr)
	}

	if res.StatusCode == 200 {
		*etag = res.Header.Get("ETag")

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Println(readErr)
		}
		res.Body.Close()

		var features map[string]Feature

		json.Unmarshal(body, &features)

		featuresMap := make(map[string]*Feature)

		for key, _ := range features{
			var f = features[key]
			featuresMap[key] = &f
		}

		featureStore.SetAll(featuresMap)
	}
}

