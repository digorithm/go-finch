package handlers

import (
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

/*func TestInviteUserEndpoint(t *testing.T) {

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
	_ = json.Unmarshal(body, &i)

	//deleteUser(t, i["ID"].(float64))
	assert.Equal(t, 201, response.Code, "OK response is expected")
}*/
