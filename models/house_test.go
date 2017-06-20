package models

import (
	"fmt"
	_ "github.com/lib/pq"
	"testing"
)

var userString string = "SELECT U.ID, U.EMAIL, U.PASSWORD, U.USERNAME, O.OWN_TYPE, O.DESCRIPTION FROM USER_INFO U INNER JOIN MEMBER_OF M ON M.USER_ID = U.ID INNER JOIN OWNERSHIP O ON O.OWN_TYPE = M.OWN_TYPE WHERE M.HOUSE_ID = $1"
var uRow []UserOwnTypeRow

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

	fmt.Println(h)
	recipes, err := h.GetHouseRecipes(nil, 1)
	if err != nil {
		t.Errorf("Getting users should work. Error: %v", err)
	}

	fmt.Println(recipes)
}

func TestRowGetters(t *testing.T) {
	h := newHouseForTest(t)

	fmt.Println("1")

	fmt.Println(uRow)
	fmt.Println(userString)
	fmt.Println(h)
	users, err := h.rowGetter(nil, uRow, userString, 1)
	if err != nil {
		t.Errorf("Getting users should work. Error: %v", err)
	}

	fmt.Println(users)

}
