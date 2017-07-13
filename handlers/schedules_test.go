package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
)

func TestGetHouseScheduleEndpoint(t *testing.T) {

	endpoint := "/schedules/1"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

}
func TestModifyHouseScheduleEndpoint(t *testing.T) {

	endpoint := "/schedules/1"
	method := "POST"

	newRecipe := []byte(`{"recipe_id":1, "type":4, "day":4}`)
	request, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(newRecipe))

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	/*body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("result: %v", string(body))*/

}

func TestGetMealTypesEndpoint(t *testing.T) {

	endpoint := "/meals"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

}

func TestGetAllDaysEndpoint(t *testing.T) {

	endpoint := "/days"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

}
