package api_test

import (
	"net/url"
	"testing"
)

func TestGetBalanceForAddress(t *testing.T) {
	uri := "/balance/addresses"
	data := url.Values{
		"address": {
			"bchtest:qqvk2qjqudp687azp2ln0nrdh7afehq0zgkp0kph25",
			"bchtest:qpz8py2yqyp7x2aaqfvrwdk4jf2c23ypzczr6weclj",
		},
	}

	res := checkPost(uri, data, t)

	t.Log(res.Result)
}
