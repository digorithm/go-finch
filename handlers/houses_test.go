package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetHouseEndpoint(t *testing.T) {

	endpoint := "/houses/1"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	ResponseStruct := make(map[string]interface{})
	ResponseBody, _ := ioutil.ReadAll(response.Body)

	err := json.Unmarshal(ResponseBody, &ResponseStruct)

	if err != nil {
		t.Fatal("Could not parse json response")
	}

	assert.Equal(t, ResponseStruct["name"], "My Lovely Home")
	assert.Equal(t, ResponseStruct["grocery_day"], "Friday")
	assert.Equal(t, ResponseStruct["household_number"], float64(2))

	assert.Equal(t, ResponseStruct["residents"].([]interface{})[0].(map[string]interface{})["name"], "guli")
	assert.Equal(t, ResponseStruct["residents"].([]interface{})[0].(map[string]interface{})["ownership"], "owner")

	assert.Equal(t, ResponseStruct["residents"].([]interface{})[1].(map[string]interface{})["name"], "digo")
	assert.Equal(t, ResponseStruct["residents"].([]interface{})[1].(map[string]interface{})["ownership"], "resident")

	assert.Equal(t, ResponseStruct["residents"].([]interface{})[2].(map[string]interface{})["name"], "joe")
	assert.Equal(t, ResponseStruct["residents"].([]interface{})[2].(map[string]interface{})["ownership"], "blocked")
}
func TestPostHouseEndpoint(t *testing.T) {

	endpoint := "/houses"
	method := "POST"
	RequestBody := []byte(`{
		"name": "My other lovely home",
		"user_id": 1,
		"grocery_day": "Friday",
		"household_number": 2
	}`)
	request, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(RequestBody))

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	ResponseStruct := make(map[string]interface{})
	ResponseBody, _ := ioutil.ReadAll(response.Body)

	err := json.Unmarshal(ResponseBody, &ResponseStruct)

	if err != nil {
		t.Fatal("Could not parse json response")
	}

	assert.Equal(t, ResponseStruct["name"], "My other lovely home")
	assert.Equal(t, ResponseStruct["grocery_day"], "Friday")
	assert.Equal(t, ResponseStruct["household_number"], float64(2))

	assert.Equal(t, ResponseStruct["residents"].([]interface{})[0].(map[string]interface{})["name"], "guli")
	assert.Equal(t, ResponseStruct["residents"].([]interface{})[0].(map[string]interface{})["ownership"], "owner")

	IDToDelete := ResponseStruct["id"].(float64)

	DeleteTestHouse(t, int64(IDToDelete))
}

func TestDeleteHouseEndpoint(t *testing.T) {
	endpoint := "/houses"
	method := "POST"
	RequestBody := []byte(`{
		"name": "My other lovely home",
		"user_id": 1,
		"grocery_day": "Friday",
		"household_number": 2
	}`)
	request, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(RequestBody))

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	ResponseStruct := make(map[string]interface{})
	ResponseBody, _ := ioutil.ReadAll(response.Body)

	err := json.Unmarshal(ResponseBody, &ResponseStruct)

	if err != nil {
		t.Fatal("Could not parse json response")
	}

	IDToDelete := ResponseStruct["id"].(float64)

	DeleteTestHouse(t, int64(IDToDelete))

}
func TestUpdateHouseEndpoint(t *testing.T) {
	endpoint := "/houses"
	method := "POST"
	RequestBody := []byte(`{
		"name": "My other lovely home",
		"user_id": 1,
		"grocery_day": "Friday",
		"household_number": 2
	}`)
	request, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(RequestBody))

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	ResponseStruct := make(map[string]interface{})
	ResponseBody, _ := ioutil.ReadAll(response.Body)

	err := json.Unmarshal(ResponseBody, &ResponseStruct)

	if err != nil {
		t.Fatal("Could not parse json response")
	}

	IDToUpdate := ResponseStruct["id"].(float64)

	endpoint = fmt.Sprintf("/houses/%v", IDToUpdate)
	method = "PUT"
	RequestBody = []byte(`{
		"name": "My very lovely home",
		"grocery_day": "Saturday",
		"household_number": 3
	}`)
	request, _ = http.NewRequest(method, endpoint, bytes.NewBuffer(RequestBody))

	request = SetTestDBEnv(request)

	response = httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	ResponseBody, _ = ioutil.ReadAll(response.Body)

	err = json.Unmarshal(ResponseBody, &ResponseStruct)

	if err != nil {
		t.Fatal("Could not parse json response")
	}

	assert.Equal(t, ResponseStruct["name"], "My very lovely home")
	assert.Equal(t, ResponseStruct["grocery_day"], "Saturday")
	assert.Equal(t, ResponseStruct["household_number"], float64(3))

	// After everything
	DeleteTestHouse(t, int64(IDToUpdate))

}

func DeleteTestHouse(t *testing.T, IDToDelete int64) {

	endpoint := fmt.Sprintf("/houses/%v", IDToDelete)
	method := "DELETE"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

}
