package models

import (
	"testing"

	_ "github.com/lib/pq"
)

func createVarsForGetUsers(hID, hNum, uID int64, uName, oType string) HouseUserOwnRow {
	var user HouseUserOwnRow

	user.HID = hID
	user.HouseNumber = hNum
	user.UID = uID
	user.Username = uName
	user.OwnType = oType

	return user
}

func createVarsForGetRecipes(id int64, name string, serves int64) RecipeRow {
	var recipe RecipeRow

	recipe.ID = id
	recipe.Name = name
	recipe.ServesFor = serves

	return recipe
}

func createVarsForGetSchedule(day string, typem string, name string, rid int64) HouseScheduleRow {
	var schedule HouseScheduleRow

	schedule.Week = day
	schedule.Type = typem
	schedule.Recipe = name
	schedule.RID = rid

	return schedule
}

func newHouseForTest(t *testing.T) *House {
	return NewHouse(newDbForTest(t))
}

func TestHouseCRUD(t *testing.T) {
	h := newHouseForTest(t)
	tearDown := TestSetup(t, h.db)
	defer tearDown(t, h.db)

	// Create house
	houseRow, err := h.CreateHouse(nil, "my lovely home", "Monday", 5)

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
	tearDown := TestSetup(t, h.db)
	defer tearDown(t, h.db)

	var u1 = createVarsForGetUsers(1, 2, 1, "guli", "owner")
	var u2 = createVarsForGetUsers(1, 2, 2, "digo", "resident")
	var u3 = createVarsForGetUsers(1, 2, 4, "joe", "blocked")
	var result []HouseUserOwnRow

	users, err := h.GetHouseUsers(nil, 1)
	if err != nil {
		t.Errorf("Getting users should work. Error: %v", err)
	}

	result = append(result, u1, u2, u3)
	i := 0
	for i < len(users) {
		if result[i] != users[i] {
			t.Errorf("House Users, got: %d, want: %d", users[i], result[i])
		}

		i++
	}

}

func TestGetRecipes(t *testing.T) {

	h := newHouseForTest(t)
	tearDown := TestSetup(t, h.db)
	defer tearDown(t, h.db)

	var r1 = createVarsForGetRecipes(2, "Beans with rice", 6)
	var result []RecipeRow

	recipes, err := h.GetHouseRecipes(nil, 3)
	if err != nil {
		t.Errorf("Getting house recipes should work. Error: %v", err)
	}

	result = append(result, r1)
	i := 0
	for i < len(recipes) {
		if result[i] != recipes[i] {
			t.Errorf("House Recipes, got: %d, want: %d", recipes[i], result[i])
		}

		i++
	}
}

func TestUpdateHouseHold(t *testing.T) {
	h := newHouseForTest(t)
	tearDown := TestSetup(t, h.db)
	defer tearDown(t, h.db)

	res, err := h.UpdateHouseHold(nil, 3, 2)
	if err != nil {
		t.Errorf("Updating house schedule should work. Error: %v", err)
	}

	if res.HouseHold != 3 {
		t.Errorf("Update House Schedule failed, got: %d, want: %d", res.HouseHold, 3)
	}

	h.UpdateHouseHold(nil, 4, 2)

}
