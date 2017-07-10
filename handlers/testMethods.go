package handlers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

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
	router.HandleFunc("/recipes/{recipe_id}", DeleteRecipesHandler).Methods("DELETE")
	router.HandleFunc("/recipes/{recipe_id}/{field}", UpdateRecipesHandler).Methods("PUT")
	router.HandleFunc("/recipes", GetRecipesHandler).Methods("GET")
	router.HandleFunc("/recipes", AddRecipesHandler).Methods("POST")

	router.HandleFunc("/users", GetUsersHandler).Methods("GET")
	router.HandleFunc("/users/{user_id}", GetUserByIDHandler).Methods("GET")
	router.HandleFunc("/users", PostSignup).Methods("POST")
	router.HandleFunc("/users/{user_id}", DeleteUser).Methods("DELETE")

	router.HandleFunc("/invitations/users/{user_id}", GetUserInvitations).Methods("GET")
	router.HandleFunc("/invitations/houses/{house_id}", GetHouseInvitations).Methods("GET")
	router.HandleFunc("/invitations/join", InviteUser).Methods("POST")
	router.HandleFunc("/invitations/{invite_id}", DeleteInvitation).Methods("DELETE")

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

func GetTestDB() *sqlx.DB {
	u, err := libunix.CurrentUser()
	if err != nil {
		fmt.Println(err)
	}
	db, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://%v@localhost:5432/meal_planner?sslmode=disable", u))

	if err != nil {
		fmt.Println(err)
	}
	return db
}

func ExtractInterfaceSliceOfStrings(t interface{}) []string {
	var str []string

	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)

		for i := 0; i < s.Len(); i++ {
			str = append(str, s.Index(i).Interface().(string))
		}
		return str
	}
	return str
}
