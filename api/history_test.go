package api_test

import (
	"net/url"
	"testing"
)

func TestGetHistory(t *testing.T) {
	uri := "/history/id/548"
	data := url.Values{"start": {"0"}, "count": {"15"}}

	res := checkPost(uri, data, t)

	result, ok := res.Result.(map[string]interface{})
	if !ok {
		t.Error("result format error")
	}

	for _, value := range result {
		t.Log(value ,"--------------xx")
	}
}

func TestGetHistoryDetail(t *testing.T) {
	uri := "/history/detail?tx_hash=b5246ea270a395ace66eb834469b2e96144d634d7ea9c74ff9fec6b82584dcab"

	res := checkGet(uri, t)
	result, ok := res.Result.(map[string]interface{})
	if !ok {
		t.Error("result format error")
	}

	t.Log(result)
	for _, value := range result {
		t.Log(value)
	}

}

func TestGetHistoryList(t *testing.T) {
	uri := "/history/list"
	data := url.Values{
		"address": {
			"bchtest:qpz8py2yqyp7x2aaqfvrwdk4jf2c23ypzczr6weclj",
			"bitcoincash:qr7tq66epzg06mjaq9m6f6sk6ldxstpvevy0yrdhhd",
			//"bitcoincash:qq6qag6mv2fzuq73qanm6k60wppy23djnv7ddk3lpk",
		},
	}

	res := checkPost(uri, data, t)
	result, ok := res.Result.([]interface{})
	if !ok {
		t.Error("result formate error")
	}
	for _, value := range result {
		t.Log(value)
	}
}