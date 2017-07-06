package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"net/http/httptest"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsersEndpoint(t *testing.T) {

	endpoint := "/users"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestGetUserByIDEndpoint(t *testing.T) {

	endpoint := "/users/1"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestPostUserEndpoint(t *testing.T) {

	user := []byte(`{"name":"aName", "email":"anEmail@email.com", "password":"aPassword"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(user))

	request := SetTestDBEnv(req)
	response := httptest.NewRecorder()
	RouterForTest().ServeHTTP(response, request)
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Printf("%v", err)
	}
	var i map[string]interface{}
	_ = json.Unmarshal(body, &i)

	deleteUser(t, i["ID"].(float64))
	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func deleteUser(t *testing.T, userID float64) {

	endpoint := fmt.Sprintf("/users/%v", userID)
	method := "DELETE"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	fmt.Printf("response code for delete: %v", response.Code)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}
