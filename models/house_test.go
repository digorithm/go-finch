package models

import (
	"testing"

	_ "github.com/lib/pq"
)

func createVarsForGetUsers(id int64, email string, pWord string, uName string, ownT int64, desc string) UserOwnTypeRow {
	var user UserOwnTypeRow

	user.ID = id
	user.Email = email
	user.Password = pWord
	user.Username = uName
	user.OwnType = ownT
	user.Description = desc

	return user
}

func createVarsForGetRecipes(id int64, name string, typem string, serves int64) RecipeRow {
	var recipe RecipeRow

	recipe.ID = id
	recipe.Name = name
	recipe.Type = typem
	recipe.ServesFor = serves

	return recipe
}

func createVarsForGetSchedule(day string, typem string, name string) HouseScheduleRow {
	var schedule HouseScheduleRow

	schedule.Week = day
	schedule.Type = typem
	schedule.Recipe = name

	return schedule
}

func checkStorageItem(storage ItemInStorageRow, hID int64, ingID int64, amt float64, uID int64) (ItemInStorageRow, bool) {

	var s ItemInStorageRow
	s.HouseID = hID
	s.IngID = ingID
	s.Amount = amt
	s.UnitID = uID

	a := (storage.HouseID == s.HouseID)
	b := (storage.IngID == s.IngID)
	c := (storage.Amount == s.Amount)
	d := (storage.UnitID == s.UnitID)

	isEqual := a && b && c && d

	return s, isEqual

}

func newHouseForTest(t *testing.T) *House {
	return NewHouse(newDbForTest(t))
}

func newStorageForTest(t *testing.T) *ItemInStorage {
	return NewItemInStorage(newDbForTest(t))
}

func TestHouseCRUD(t *testing.T) {
	h := newHouseForTest(t)

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
	var u1 = createVarsForGetUsers(1, "gulipek5@gmail.com", "password", "guli", 1, "owner")
	var u2 = createVarsForGetUsers(2, "rod.dearaujo@gmail.com", "password1", "digo", 2, "resident")
	var u3 = createVarsForGetUsers(4, "iamjoe@gmail.com", "password3", "joe", 3, "not allowed")
	var result []UserOwnTypeRow

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
	var r1 = createVarsForGetRecipes(2, "Beans with rice", "Lunch/Dinner", 6)
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

	house, err := h.UpdateHouseHold(nil, 3, 2)
	if err != nil {
		t.Errorf("Updating house schedule should work. Error: %v", err)
	}

	if house != 1 {
		t.Errorf("Update House Schedule failed, got: %d, want: %d", house, 1)
	}

	h.UpdateHouseHold(nil, 4, 2)

}

func TestInsertStorage(t *testing.T) {

	s := newStorageForTest(t)
	h := newHouseForTest(t)

	houseRow, err := h.CreateHouse(nil, "Test for Storage", "Wednesday", 10)

	if err != nil {
		t.Errorf("%v", err)
	}

	storage, err := s.InsertToStorage(nil, houseRow.ID, 1, 0.4, 3)

	if err != nil {
		t.Errorf("Inserting ingredient to storage should work. Error: %v", err)
	}

	res, isEqual := checkStorageItem(storage, houseRow.ID, 1, 0.4, 3)

	if !isEqual {
		t.Errorf("Inserting to storage failed. Got: %v, Want: %v", err, res)
	}

	_, err = h.DeleteById(nil, houseRow.ID)
	if err != nil {
		t.Fatalf("Deleting house by id should not fail. Error: %v", err)
	}

}
