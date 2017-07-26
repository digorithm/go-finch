package models

import (
	"testing"

	"fmt"

	"encoding/json"

	_ "github.com/lib/pq"
)

func newUserForTest(t *testing.T) *User {
	return NewUser(newDbForTest(t))
}

func newUserRowForTest(t *testing.T) (*User, []byte) {
	u := newUserForTest(t)
	email := newEmailForTest()
	//var userRow *UserRow

	user := fmt.Sprintf(`{"name":"username", "password":"abc123", "email":"%v"}`, email)

	newSignup := []byte(user)
	row, _ := u.Signup(nil, newSignup)

	//_ = json.Unmarshal(row, userRow)

	return u, row
}

func deleteTestUser(t *testing.T, u *User, id int64) {
	_, err := u.DeleteById(nil, id)
	if err != nil {
		t.Fatal("Something went wrong with test user deletion. Error: %v", err)
	}
}

func deleteTestRecipe(t *testing.T, id int64) {
	r := newRecipeForTest(t)
	_, err := r.DeleteById(nil, id)
	if err != nil {
		t.Fatal("Something went wrong with test recipe deletion. Error: %v", err)
	}
}

func TestUserSignup(t *testing.T) {
	u := newUserForTest(t)

	email := newEmailForTest()
	user := fmt.Sprintf(`{"name":"username", "password":"abc123", "email":"%v"}`, email)

	newSignup := []byte(user)
	usr, err := u.Signup(nil, newSignup)

	var userRow *UserRow

	_ = json.Unmarshal(usr, &userRow)

	defer deleteTestUser(t, u, userRow.ID)

	if err != nil {
		t.Errorf("Signing up user should work. Error: %v", err)
	}
	if userRow == nil {
		t.Fatal("Signing up user should work.")
	}
	if userRow.ID <= 0 {
		t.Fatal("Signing up user should work.")
	}
}

func TestGetUserById(t *testing.T) {
	u := newUserForTest(t)

	email := newEmailForTest()
	user := fmt.Sprintf(`{"name":"username", "password":"abc123", "email":"%v"}`, email)

	newSignup := []byte(user)
	usr, err := u.Signup(nil, newSignup)

	var userRow *UserRow

	_ = json.Unmarshal(usr, &userRow)

	defer deleteTestUser(t, u, userRow.ID)

	returningUserRow, err := u.GetById(nil, userRow.ID)

	if err != nil {
		t.Errorf("Find user by ID should work")
	}

	if userRow.ID != returningUserRow.ID {
		t.Errorf("IDs did not match!")
	}

}

func TestGetUserByEmail(t *testing.T) {
	u := newUserForTest(t)
	email := newEmailForTest()
	user := fmt.Sprintf(`{"name":"username", "password":"abc123", "email":"%v"}`, email)

	newSignup := []byte(user)
	usr, err := u.Signup(nil, newSignup)

	var userRow *UserRow

	_ = json.Unmarshal(usr, &userRow)
	defer deleteTestUser(t, u, userRow.ID)

	returningUserRow, err := u.GetByEmail(nil, userRow.Email)

	if err != nil {
		t.Errorf("Find user by Email should work")
	}

	if userRow.Email != returningUserRow.Email {
		t.Errorf("Emails did not match!")
	}

}

func TestGetUserByUsername(t *testing.T) {
	u := newUserForTest(t)
	email := newEmailForTest()
	user := fmt.Sprintf(`{"name":"username", "password":"abc123", "email":"%v"}`, email)

	newSignup := []byte(user)
	usr, err := u.Signup(nil, newSignup)

	var userRow *UserRow

	_ = json.Unmarshal(usr, &userRow)
	defer deleteTestUser(t, u, userRow.ID)

	returningUserRow, err := u.GetByUsername(nil, userRow.Username)

	if err != nil {
		t.Errorf("Find user by Username should work")
	}

	if userRow.Username != returningUserRow.Username {
		t.Errorf("Usernames did not match!")
	}

}

func TestAddRecipe(t *testing.T) {

	u, row := newUserRowForTest(t)

	var uRow UserRow

	err := json.Unmarshal(row, &uRow)

	fmt.Printf("u:: %v\n", uRow)

	if err != nil {
		fmt.Println(err)
	}

	defer deleteTestUser(t, u, uRow.ID)

	// Define a test recipe.
	// It will come from the client request as a JSON.
	// The handler will extract the maps[] from the JSON just like we are doing
	// down here and pass them to User.AddRecipe(...)

	test_recipe := []byte(`{
		"recipe_name": "feijoada",
		"type": ["Lunch", "Dinner"],
		"serves_for": "2",
		"steps": [
			{
				"step_id": 1,
				"text": "description of the first step",
				"step_ingredients": [
					{"name": "beans", "amount": 34.5, "unit": 10},
					{"name": "rice", "amount": 14.5, "unit": 10}
				]
			},
			{
				"step_id": 2,
				"text": "description of the second step",
				"step_ingredients": [
					{"name": "water", "amount": 4.5, "unit": 10}
				]
			},
			{
				"step_id": 3,
				"text": "description of the third step",
				"step_ingredients": [
					{"name": "salt", "amount": 1.5, "unit": 10}
				]
			}
		]
	}`)

	returnedRecipe, err := u.AddRecipe(nil, test_recipe, uRow.ID)

	if err != nil {
		t.Errorf("Add recipe should work. Err: %v", err)
	}

	if len(returnedRecipe) == 0 {
		t.Errorf("Empty returnedRecipe: %v", returnedRecipe)
	}

	if returnedRecipe[0].Name != "feijoada" {
		t.Errorf("Recipes have different names.")
		t.Errorf("Expected: feijoada")
		t.Errorf("Actual: %v", returnedRecipe[0].Name)
	}

	if returnedRecipe[0].ServesFor != 2 {
		t.Errorf("Recipes have different ServesFor.")
		t.Errorf("Expected: 2")
		t.Errorf("Actual: %v", returnedRecipe[0].ServesFor)
	}

	if returnedRecipe[0].Ingredient != "beans" {
		t.Errorf("Wrong ingredient.")
		t.Errorf("Expected: beans")
		t.Errorf("Actual: %v", returnedRecipe[0].Ingredient)
	}

	if returnedRecipe[1].Ingredient != "rice" {
		t.Errorf("Wrong ingredient.")
		t.Errorf("Expected: rice")
		t.Errorf("Actual: %v", returnedRecipe[1].Ingredient)
	}

	if returnedRecipe[2].Ingredient != "water" {
		t.Errorf("Wrong ingredient.")
		t.Errorf("Expected: water")
		t.Errorf("Actual: %v", returnedRecipe[2].Ingredient)
	}

	if returnedRecipe[3].Ingredient != "salt" {
		t.Errorf("Wrong ingredient.")
		t.Errorf("Expected: salt")
		t.Errorf("Actual: %v", returnedRecipe[3].Ingredient)
	}

	deleteTestRecipe(t, returnedRecipe[0].ID)

}

func TestGetUserRecipes(t *testing.T) {

	u := newUserForTest(t)
	var r1 = createVarsForGetRecipes(2, "Beans with rice", 6)
	var r2 = createVarsForGetRecipes(3, "No Flour Pancake", 2)
	var result []RecipeRow

	recipes, err := u.GetUserRecipes(nil, 2)
	if err != nil {
		t.Errorf("Generic get recipe should work")
	}

	result = append(result, r1, r2)
	i := 0
	for i < len(recipes) {
		if result[i] != recipes[i] {
			t.Errorf("House Users, got: %d, want: %d", recipes[i], result[i])
		}
		i++
	}

}

func TestGetMembership(t *testing.T) {

	u := newUserForTest(t)

	own, err := u.getMembership(nil, 2, 1)

	if err != nil {
		fmt.Printf("%v", err)
	}

	if own.OwnType != 2 {
		t.Errorf("Get Membership, got: %d, want: %d", own.OwnType, 2)
	}

}
