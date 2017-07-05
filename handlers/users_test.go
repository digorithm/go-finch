package handlers

import (
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
