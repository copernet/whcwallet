package model

import (
	"fmt"
	"testing"
)

func TestGetAllActiveCrowdSale(t *testing.T) {

	ret, err := GetAllActiveCrowdSale("",20, 1)
	if err != nil {
		t.Error(err)
	}

	for _, v := range ret {
		fmt.Println(v)
	}

	t.Log(ret)
}

func TestGetAllActiveCrowdSaleCount(t *testing.T) {

	ret, err := GetAllActiveCrowdSaleCount("3")
	if err != nil {
		t.Error(err)
	}

	t.Log(ret)
}
