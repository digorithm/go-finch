package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
)

func TestHouseStorageEndpoint(t *testing.T) {

	endpoint := "/storages/1"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}

	var res []map[string]interface{}
	_ = json.Unmarshal(body, &res)

	assert.Equal(t, res[0]["ingredient"], "potato")
	assert.Equal(t, res[0]["amount"].(float64), 5.0)
}

func TestPostStorageEndpoint(t *testing.T) {

	endpoint := "/storages/1"
	method := "POST"

	RequestBody := []byte(`[
									{
										"name": "apple",
										"amount": 800,
										"unit": 10
									},
									{
										"name": "pasta",
										"amount": 500,
										"unit": 10
									}
								]`)

	request, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(RequestBody))

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	endpoint = "/storages/1"
	method = "GET"

	request, _ = http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response = httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}

	var res []map[string]interface{}
	_ = json.Unmarshal(body, &res)

	assert.Equal(t, res[1]["ingredient"], "pasta")
	assert.Equal(t, res[1]["amount"].(float64), 500.0)

	assert.Equal(t, res[2]["ingredient"], "apple")
	assert.Equal(t, res[2]["amount"].(float64), 800.0)

}
