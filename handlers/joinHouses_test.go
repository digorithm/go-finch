package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"net/http/httptest"

	"encoding/json"

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

	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestGetUserInvitationsEndpoint(t *testing.T) {

	endpoint := "/invitations/users/6"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestInviteUserEndpoint(t *testing.T) {

	invitation := []byte(`{"house_id":3, "user_id":2}`)
	req, _ := http.NewRequest("POST", "/invitations/join", bytes.NewBuffer(invitation))

	request := SetTestDBEnv(req)
	response := httptest.NewRecorder()
	RouterForTest().ServeHTTP(response, request)
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Printf("%v", err)
	}
	var i []map[string]interface{}

	err = json.Unmarshal(body, &i)

	if err != nil {
		fmt.Println(err)
	}

	deleteInvite(t, i[0]["invite_id"].(float64))
	assert.Equal(t, 201, response.Code, "OK response is expected")
}

func deleteInvite(t *testing.T, inviteID float64) {

	endpoint := fmt.Sprintf("/invitations/%v", inviteID)
	method := "DELETE"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	fmt.Printf("response code for delete: %v", response.Code)
	assert.Equal(t, 204, response.Code, "OK response is expected")

}
