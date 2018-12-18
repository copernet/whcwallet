package api_test

import (
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"strings"

	"encoding/json"
	"testing"

	"github.com/copernet/whcwallet/routers"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

var router *gin.Engine

func newRouter() *gin.Engine {
	if router != nil {
		return router
	}

	router = routers.InitRouter()
	return router
}

func Get(uri string) ([]byte, error) {
	router := newRouter()

	req := httptest.NewRequest("GET", uri, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	ret := w.Result()
	defer ret.Body.Close()

	body, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func PostForm(uri string, param url.Values) ([]byte, error) {
	router := newRouter()

	req := httptest.NewRequest("POST", uri, strings.NewReader(param.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	ret := w.Result()
	defer ret.Body.Close()

	body, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// convenient function
func isStatusOk(res *Response) bool {
	if res.Code != 0 || res.Message != "" {
		return false
	}

	return true
}

func checkGet(uri string, t *testing.T) *Response {
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
		t.Error("Request result with error:", res.Message)
	}

	return &res
}

func checkPost(uri string, param url.Values, t *testing.T) *Response {
	ret, err := PostForm(uri, param)

	if err != nil {
		t.Error(err)
	}

	var res Response
	err = json.Unmarshal(ret, &res)
	if err != nil {
		t.Error(err)
	}

	if !isStatusOk(&res) {
		t.Error("Request result with error:", res.Message)
	}

	return &res
}
