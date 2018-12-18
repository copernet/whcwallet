package api_test

import (
	"reflect"
	"testing"

	"github.com/copernet/whcwallet/api/data"
)

func TestGetCategories(t *testing.T) {
	uri := "/category/"

	res := checkGet(uri, t)

	result, ok := res.Result.(map[string]interface{})
	if !ok {
		t.Error("result format error")
	}

	if len(result) != len(data.Categories) {
		t.Error("the number of the category items is error")
	}
}

func TestGetSubCategories(t *testing.T) {
	uri := "/category/subcategories?category=Information%20and%20communication"

	res := checkGet(uri, t)
	result, ok := res.Result.([]interface{})
	if !ok {
		t.Error("result format error")
	}

	expected := data.Categories["Information and communication"]
	for idx, item := range result {
		str, ok := item.(string)
		if !ok {
			t.Errorf("expected: string, but got: %v", reflect.TypeOf(item))
		}

		if str != expected[idx] {
			t.Errorf("the subcategory result expected: %s, but got: %s", expected[idx], item)
		}
	}
}
