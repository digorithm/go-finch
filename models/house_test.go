package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

var userString string = "SELECT U.ID, U.EMAIL, U.PASSWORD, U.USERNAME, O.OWN_TYPE, O.DESCRIPTION FROM USER_INFO U INNER JOIN MEMBER_OF M ON M.USER_ID = U.ID INNER JOIN OWNERSHIP O ON O.OWN_TYPE = M.OWN_TYPE WHERE M.HOUSE_ID = $1"
var uRow UserOwnTypeRow

func newHouseForTest(t *testing.T) *House {
	return NewHouse(newDbForTest(t))
}

func TestHouseCRUD(t *testing.T) {
	h := newHouseForTest(t)

	// Create house
	houseRow, err := h.CreateHouse(nil, "my lovely home")

	if err != nil {
		t.Errorf("Creating house should work. Error: %v", err)
	}

	// Test deletion
	_, err = h.DeleteById(nil, houseRow.ID)
	if err != nil {
		t.Fatalf("Deleting house by id should not fail. Error: %v", err)
	}
}

func TestGetUsers(t *testing.T) {
	h := newHouseForTest(t)

	users, err := h.GetHouseUsers(nil, 1)
	if err != nil {
		t.Errorf("Getting users should work. Error: %v", err)
	}

	fmt.Println(users)
}

func TestGetRecipes(t *testing.T) {
	h := newHouseForTest(t)

	recipes, err := h.GetHouseRecipes(nil, 1)
	if err != nil {
		t.Errorf("Getting users should work. Error: %v", err)
	}

	var r1 RecipeRow
	var r2 RecipeRow
	r1.ID = 1
	r1.Name = "Baked Potato"
	r2.ID = 4
	r2.Name = "Roast Chicken"
	var result []RecipeRow

	result = append(result, r1, r2)

	assert.Equal(t, result, recipes, "Two objects should be the same")

	fmt.Println(recipes)
}
