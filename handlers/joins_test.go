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

type AcceptResponse struct {
	ID        int64 `json:"id"`
	Household int64 `json:"household_number"`
	Users     []struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Ownership string `json:"ownership"`
	} `json:"users"`
}

func TestGetHouseInvitationsEndpoint(t *testing.T) {

	endpoint := "/invitations/houses/3"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestGetHouseRequestsEndpoint(t *testing.T) {

	endpoint := "/requests/houses/5"
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

func TestGetUserJoinsEndpoint(t *testing.T) {

	endpoint := "/requests/users/1"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestRequestJoinEndpoint(t *testing.T) {

	invitation := []byte(`{"house_id":6, "user_id":3}`)
	req, _ := http.NewRequest("POST", "/requests/join", bytes.NewBuffer(invitation))

	request := SetTestDBEnv(req)
	response := httptest.NewRecorder()
	RouterForTest().ServeHTTP(response, request)
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Printf("%v", err)
	}
	var i map[string]interface{}

	err = json.Unmarshal(body, &i)

	if err != nil {
		fmt.Println(err)
	}

	deleteInvite(t, i["invite_id"].(float64))
	assert.Equal(t, 201, response.Code, "OK response is expected")
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
	var i map[string]interface{}

	err = json.Unmarshal(body, &i)

	if err != nil {
		fmt.Println(err)
	}

	deleteInvite(t, i["invite_id"].(float64))
	assert.Equal(t, 201, response.Code, "OK response is expected")
}

func TestRespondInvite(t *testing.T) {

	inviteID := int64(mockInvite(t))

	inv := fmt.Sprintf(`{"invite_id":%v, "accepts":true}`, inviteID)
	invitation := []byte(inv)
	req, _ := http.NewRequest("POST", "/invitations/respond", bytes.NewBuffer(invitation))

	request := SetTestDBEnv(req)
	response := httptest.NewRecorder()
	RouterForTest().ServeHTTP(response, request)

	deleteMember(t)

}

func mockInvite(t *testing.T) float64 {

	invitation := []byte(`{"house_id":4, "user_id":1}`)
	req, _ := http.NewRequest("POST", "/invitations/join", bytes.NewBuffer(invitation))

	request := SetTestDBEnv(req)
	response := httptest.NewRecorder()
	RouterForTest().ServeHTTP(response, request)
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Printf("%v", err)
	}
	var i map[string]interface{}

	err = json.Unmarshal(body, &i)

	if err != nil {
		fmt.Println(err)
	}
	return i["invite_id"].(float64)

}

func deleteInvite(t *testing.T, inviteID float64) {

	endpoint := fmt.Sprintf("/invitations/%v", inviteID)
	method := "DELETE"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 204, response.Code, "OK response is expected")

}

func deleteMember(t *testing.T) {

	endpoint := "/houses/4/users/1"
	method := "DELETE"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 204, response.Code, "OK response is expected")
}
