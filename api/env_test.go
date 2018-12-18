package api_test

import (
	"encoding/json"
	"testing"
)

func TestGetEnv(t *testing.T) {
	uri := "/env"
	ret, err := Get(uri)
	if err != nil {
		t.Error(err)
	}

	var res Response
	err = json.Unmarshal(ret, &res)
	if err != nil {
		t.Error(err)
	}

	if !isStatusOk(&res) {
		t.Error(res.Message)
	}

	result, ok := res.Result.(map[string]interface{})
	if !ok {
		t.Error("result format error")
	}

	_, ok = result["testnet"].(bool)
	if !ok {
		t.Error("the wanted type is boolean")
	}
}
