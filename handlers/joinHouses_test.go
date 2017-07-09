package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"net/http/httptest"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetHouseInvitationsEndpoint(t *testing.T) {

	endpoint := "/invitations/houses/3"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	fmt.Printf("response body: %v", response.Body)
	assert.Equal(t, 200, response.Code, "OK response is expected")

}
