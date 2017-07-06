package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/digorithm/meal_planner/libunix"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Function used to test HTTP endpoints in the REST API.
// How to use:
// 1. Add the route to be tested in the component
// 2. Add the handler that will handle a route
// 3. Write the test to call that route
func RouterForTest() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/recipes/house/{house_id}", GetHouseRecipesHandler).Methods("GET")
	router.HandleFunc("/recipes/user/{user_id}", GetUserRecipesHandler).Methods("GET")
	router.HandleFunc("/recipes/{recipe_id}", GetRecipeByIDHandler).Methods("GET")
	router.HandleFunc("/recipes", GetRecipesHandler).Methods("GET")
	router.HandleFunc("/users", GetUsersHandler).Methods("GET")
	router.HandleFunc("/users/{user_id}", GetUserByIDHandler).Methods("GET")
	router.HandleFunc("/users", PostSignup).Methods("POST")
	router.HandleFunc("/users/{user_id}", DeleteUser).Methods("DELETE")

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
