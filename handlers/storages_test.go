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

	assert.Equal(t, res[0]["ingredient_name"], "potato")
	assert.Equal(t, res[0]["amount"].(float64), 5.0)
}

func TestPostStorageEndpoint(t *testing.T) {

	endpoint := "/storages/2"
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

	assert.Equal(t, 201, response.Code, "Created response is expected")

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}

	var res []map[string]interface{}
	_ = json.Unmarshal(body, &res)

	assert.Equal(t, res[1]["ingredient_name"], "apple")
	assert.Equal(t, res[1]["amount"].(float64), 800.0)

	assert.Equal(t, res[0]["ingredient_name"], "pasta")
	assert.Equal(t, res[0]["amount"].(float64), 500.0)

	DeleteStorage2(t)

	GetStorage2(t)

}

func TestDeleteStorage(t *testing.T) {

	CreateStorage2(t)

	endpoint := "/storages/2"
	method := "DELETE"

	RequestBody := []byte(`[
									{
										"name": "apple"
									},
									{
										"name": "pasta"
									}
								]`)

	request, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(RequestBody))

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "Created response is expected")

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}

	var res []map[string]interface{}
	_ = json.Unmarshal(body, &res)

	assert.Equal(t, res[0]["ingredient_name"], "cassava")
	assert.Equal(t, res[0]["amount"].(float64), 12.0)

	assert.Equal(t, res[1]["ingredient_name"], "cherry")
	assert.Equal(t, res[1]["amount"].(float64), 3.0)

	DeleteStorage2(t)

}

func DeleteStorage2(t *testing.T) {

	endpoint := "/storages/all/2"
	method := "DELETE"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func GetStorage2(t *testing.T) {

	endpoint := "/storages/2"
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

	var res1 []map[string]interface{}
	_ = json.Unmarshal(body, &res1)

	assert.Empty(t, res1, "The storage of house should be empty")
}

func CreateStorage2(t *testing.T) {

	endpoint := "/storages/2"
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
									}, 
									{
										"name": "cassava",
										"amount": 12,
										"unit": 4
									},
									{
										"name": "cherry",
										"amount": 3,
										"unit": 7
									}
								]`)

	request, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(RequestBody))

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

}
