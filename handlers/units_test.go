package handlers

import (
	"net/http"
	"testing"

	"net/http/httptest"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUnitsEndpoint(t *testing.T) {

	endpoint := "/units"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

}
