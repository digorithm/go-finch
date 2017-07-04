package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/digorithm/meal_planner/libunix"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// Function used to test HTTP endpoints in the REST API.
// How to use:
// 1. Add the route to be tested in the component
// 2. Add the handler that will handle a route
// 3. Write the test to call that route
func RouterForTest() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/recipes/house/{house_id}/", GetHouseRecipesHandler).Methods("GET")
	router.HandleFunc("/recipes/user/{user_id}/", GetUserRecipesHandler).Methods("GET")
	return router
}

func SetTestDBEnv(request *http.Request) *http.Request {
	u, err := libunix.CurrentUser()
	if err != nil {
		fmt.Println(err)
	}
	db, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://%v@localhost:5432/meal_planner?sslmode=disable", u))

	if err != nil {
		fmt.Println(err)
	}

	ctx := request.Context()
	ctx = context.WithValue(ctx, "db", db)

	request = request.WithContext(ctx)

	return request
}
func TestGetHouseRecipesEndpoint(t *testing.T) {

	endpoint := "/recipes/house/1/"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
}
func TestGetUserRecipesEndpoint(t *testing.T) {

	endpoint := "/recipes/user/1/"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
}
