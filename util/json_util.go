package util

import (
	"encoding/json"
)

type Result map[string]interface{}

func JsonStringArrayToMap(src []string) []Result {
	var result []Result
	for _, value := range src {
		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(value), &dat); err == nil {
			result = append(result, dat)
		}
	}
	return result
}

func JsonStringToMap(src string) Result {
	var result Result
	if err := json.Unmarshal([]byte(src), &result); err == nil {
		return result
	}
	return nil
}
