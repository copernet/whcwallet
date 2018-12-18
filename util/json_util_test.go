package util

import (
	"fmt"
	"testing"
)

func TestJsonStringArryToMap(t *testing.T) {
	var a []string
	a = append(a, "")
	b := JsonStringArrayToMap(a)
	fmt.Println(b)
}
