package model

import (
	"testing"
)

func TestQueryTransactions2(t *testing.T) {
	var property_id int64 = 1
	var start int64 = 0
	var count int64 = 1
	var address string = "bitcoincash:qp0uasws2tz2pclmxjpujx4066vc5k6lvvmcq7jzdq"
	ret, err := QueryTransactions(property_id, address,start, count)
	if err != nil {
		t.Error(err)
	}

	t.Log(ret)
}

func TestQueryTransactions(t *testing.T) {
	var property_id int64 = 1
	var start int64 = 0
	var count int64 = 1
	var address string = "bitcoincash:qp0uasws2tz2pclmxjpujx4066vc5k6lvvmcq7jzdq"
	ret, err := QueryTransactions(property_id,address, start, count)
	if err != nil {
		t.Error(err)
	}

	t.Log(ret)
}

func TestQueryTotal(t *testing.T) {
	var property_id int64 = 1
	var address string = "bitcoincash:qp0uasws2tz2pclmxjpujx4066vc5k6lvvmcq7jzdq"
	ret, err := QueryTotal(property_id,address)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}

func TestGetHistoryDetail(t *testing.T) {
	result, err := GetHistorySpDetail("78dcea0da8914a88e7d03b64bbd7c01096491b76c1007f78cbd0b1817594af7b")
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}

func TestGetHistoryList(t *testing.T) {

	addrs := []string{
		"bitcoincash:qr6pnv8vf9y0zjuceqtsynwp7zqfqcgkmsqza54ef5",
		"bitcoincash:qp0uasws2tz2pclmxjpujx4066vc5k6lvvmcq7jzdq",
		"bitcoincash:qp0uasws2tz2pclmxjpujx4066vc5k6lvvmcq7jzdq",
	}
	var pageNo = 1
	list, err := GetHistoryList(addrs, 50, 1,pageNo)
	if err != nil {
		t.Error(err)
	}

	if len(list) != 2 {
		t.Error("the number of result is incorrect")
	}
}
