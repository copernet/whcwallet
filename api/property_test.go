package api_test

import (
	"net/url"
	"testing"
)

func TestGetProperties(t *testing.T) {
	uri := "/property/list"

	res := checkGet(uri, t)
	r := res.Result.([]interface{})
	for _, v := range r {
		t.Log(v.(map[string]interface{}))
	}
}

func TestListByOwner(t *testing.T) {
	uri := "/property/listbyowner"
	param := url.Values{
		"address": {
			"bitcoincash:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd",
			"bitcoincash:qr6pnv8vf9y0zjuceqtsynwp7zqfqcgkmsqza54ef5",
		},
	}

	res := checkPost(uri, param, t)
	r := res.Result.([]interface{})
	for _, v := range r {
		t.Log(v.(map[string]interface{}))
	}
}

func TestListOwners(t *testing.T) {
	uri := "/property/listowners/459"
	res := checkPost(uri, nil, t)
	r := res.Result.([]interface{})
	for _, v := range r {
		t.Log(v.(map[string]interface{}))
	}
}
