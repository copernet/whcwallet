package api_test

import (
	"testing"
)

func TestListActiveCrowdSales(t *testing.T) {
	uri := "/crowdsale/list/active"

	res := checkGet(uri, t)

	result, ok := res.Result.([]interface{})
	if !ok {
		t.Error("result format error")
	}

	for _, value := range result {
		r := value.(map[string]interface{})
		if r["active"] != true {
			t.Error("recive data error")
		}
	}

	t.Log(result)
}
