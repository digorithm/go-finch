package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"io/ioutil"

	"encoding/json"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetHouseRecipesEndpoint(t *testing.T) {

	endpoint := "/recipes/house/1"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestGetUserRecipesEndpoint(t *testing.T) {

	endpoint := "/recipes/user/1"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestGetRecipeByIDEndpoint(t *testing.T) {

	endpoint := "/recipes/1"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestGetAllRecipesEndpoint(t *testing.T) {

	endpoint := "/recipes"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
}
func TestGetRecipeBySearchStringEndpoint(t *testing.T) {

	endpoint := "/recipes?name=baked"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestAddRecipeEndpoint(t *testing.T) {

	// TODO: move this URL param to session user ID
	endpoint := "/recipes?author=1"
	method := "POST"

	recipe := []byte(`{
           "recipe_name":"feijoada",
						"type": ["Lunch", "Dinner"],
           "serves_for":"2",
           "steps":[
              {
                 "step_id":1,
                 "text":"description of the first step",
                 "step_ingredients":[
                    {
                       "name":"beans",
                       "amount":34.5,
                       "unit":10
                    },
                    {
                       "name":"rice",
                       "amount":14.5,
                       "unit":10
                    }
                 ]
              },
              {
                 "step_id":2,
                 "text":"description of the second step",
                 "step_ingredients":[
                    {
                       "name":"water",
                       "amount":4.5,
                       "unit":10
                    }
                 ]
              },
              {
                 "step_id":3,
                 "text":"description of the third step",
                 "step_ingredients":[
                    {
                       "name":"salt",
                       "amount":1.5,
                       "unit":10
                    }
                 ]
              }
           ]
        }`)

	request, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(recipe))

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)
	// TODO: check response
	assert.Equal(t, 201, response.Code, "OK response is expected")

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}

	var res []map[string]interface{}
	_ = json.Unmarshal(body, &res)

	// Remove inserted test recipe
	IDToDelete := res[0]["id"].(float64)

	DeleteTestRecipeEndpoint(int64(IDToDelete), t)
}

func DeleteTestRecipeEndpoint(RecipeToDelete int64, t *testing.T) {

	endpoint := fmt.Sprintf("/recipes/%v", RecipeToDelete)
	method := "DELETE"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
}
