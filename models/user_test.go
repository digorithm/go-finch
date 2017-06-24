package models

import (
	"testing"

	"fmt"

	_ "github.com/lib/pq"
)

func newUserForTest(t *testing.T) *User {
	return NewUser(newDbForTest(t))
}

func deleteTestUser(t *testing.T, u *User, id int64) {
	_, err := u.DeleteById(nil, id)
	if err != nil {
		t.Fatal("Something went wrong with test user deletion. Error: %v", err)
	}
}

func TestUserSignup(t *testing.T) {
	u := newUserForTest(t)

	// Signup
	userRow, err := u.Signup(nil, newEmailForTest(), "username", "abc123", "abc123")
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

	userRow, err := u.Signup(nil, newEmailForTest(), "username", "abc123", "abc123")
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

	userRow, err := u.Signup(nil, newEmailForTest(), "username", "abc123", "abc123")
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

	userRow, err := u.Signup(nil, newEmailForTest(), "username", "abc123", "abc123")
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
	u := newUserForTest(t)

	// Define a test recipe.
	// It will come from the client request as a JSON.
	// The handler will extract the maps[] from the JSON just like we are doing
	// down here and pass them to User.AddRecipe(...)

	recipe := make(map[string]interface{})
	recipe["name"] = "feijoada"
	recipe["type"] = "Lunch/Dinner"
	recipe["serves_for"] = 3

	steps := make([]map[string]interface{}, 0, 0)
	s1 := make(map[string]interface{})
	s2 := make(map[string]interface{})
	s3 := make(map[string]interface{})

	// TODO: Fix this representation. This way the amount/unit can't describe many ingredients
	s1["step_id"] = 1
	s1["text"] = "Text to describe first step"
	s1["step_ingredients"] = []string{"beans", "rice"}
	s1["amount"] = 4.5
	s1["unit"] = 10

	s2["step_id"] = 2
	s2["text"] = "Text to describe second step"
	s2["step_ingredients"] = []string{"water"}
	s2["amount"] = 34.5
	s2["unit"] = 10

	s3["step_id"] = 3
	s3["text"] = "Text to describe third step"
	s3["step_ingredients"] = []string{"salt"}
	s3["amount"] = 42.5
	s3["unit"] = 10
	steps = append(steps, s1, s2, s3)

	_, err := u.AddRecipe(nil, recipe, steps)

	if err != nil {
		t.Errorf("Add recipe should work. Err: %v", err)
	}
}

func TestGetUserRecipes(t *testing.T) {

	u := newUserForTest(t)
	// Insert new recipe with this user, return id
	// getUserRecipes passing this id

	recipes, err := u.GetUserRecipes(nil, 2)

	if err != nil {
		t.Errorf("Generic get recipe should work")
	}

	fmt.Println(recipes)

}
