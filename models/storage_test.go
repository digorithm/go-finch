package models

import (
	"testing"

	"fmt"

	_ "github.com/lib/pq"
)

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

func newStorageForTest(t *testing.T) *ItemInStorage {
	return NewItemInStorage(newDbForTest(t))
}

func TestAddIngToStorage(t *testing.T) {

	s := newStorageForTest(t)
	tearDown := TestSetup(t, s.db)
	defer tearDown(t, s.db)

	h := newHouseForTest(t)

	houseRow, err := h.CreateHouse(nil, "Test for Storage", "Wednesday", 10)

	if err != nil {
		t.Errorf("%v", err)
	}

	storage, err := s.AddIngToStorage(nil, houseRow.ID, 1, 0.4, 3)

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

func TestUpdateStorage(t *testing.T) {
	s := newStorageForTest(t)
	tearDown := TestSetup(t, s.db)
	defer tearDown(t, s.db)
	h := newHouseForTest(t)

	house, err := h.CreateHouse(nil, "Test For Storage 1", "Monday", 1)

	if err != nil {
		t.Errorf("%v", err)
	}

	_, err = s.AddIngToStorage(nil, house.ID, 2, 5, 2)

	if err != nil {
		t.Errorf("%v", err)
	}

	result, err := s.UpdateIngredient(nil, house.ID, 2, 9, 1)

	res, isEqual := checkStorageItem(result, house.ID, 2, 9, 1)

	if !isEqual {
		t.Errorf("Updating storage failed. Got: %v, Want: %v", err, res)
	}

	_, err = h.DeleteById(nil, house.ID)
	if err != nil {
		t.Fatalf("Deleting house by id should not fail. Error: %v", err)
	}

}

func TestGetHouseStorage(t *testing.T) {

	s := newStorageForTest(t)
	tearDown := TestSetup(t, s.db)
	defer tearDown(t, s.db)

	res, err := s.GetHouseStorage(nil, 3)

	if err != nil {
		t.Errorf("GetHouseStorage failed:%v", err)
	}

	var expected ItemInStorageRow

	expected.HouseID = 3
	expected.IngID = 2
	expected.IngName = "milk"
	expected.Amount = 0
	expected.UnitID = 2

	_, equal := checkStorageItem(expected, res[0].HouseID, res[0].IngID, res[0].Amount, res[0].UnitID)

	if !equal {
		t.Errorf("Updating storage failed. Got: %v, Want: %v", res[0], expected)
	}

}

func TestGetStorageIngredient(t *testing.T) {

	s := newStorageForTest(t)
	tearDown := TestSetup(t, s.db)
	defer tearDown(t, s.db)

	item, err := s.GetStorageIngredient(nil, 1, 1)

	if err != nil {
		t.Errorf("GetStorageIngredient failed:%v", err)
	}

	var expected ItemInStorageRow

	expected.HouseID = 1
	expected.IngID = 1
	expected.IngName = "potato"
	expected.Amount = 5
	expected.UnitID = 8

	_, equal := checkStorageItem(expected, item.HouseID, item.IngID, item.Amount, item.UnitID)

	if !equal {
		t.Errorf("Getting the item failed. Got: %v, Want: %v", item, expected)
	}

}

func TestNewIngAddList(t *testing.T) {
	s := newStorageForTest(t)
	tearDown := TestSetup(t, s.db)
	defer tearDown(t, s.db)
	i := newIngredientForTest(t)

	ingredient := []byte(`[{"name":"blueberry", "amount": 1	,"unit": 2 }]`)

	s.AddIngredientList(ingredient, 2)

	item, err := i.GetByName(nil, "blueberry")

	if err != nil {
		t.Errorf("GetByName failed:%v", err)
	}

	ing, err := s.GetStorageIngredient(nil, 2, item.ID)

	if err != nil {
		t.Errorf("GetStorageIngredient failed:%v", err)
	}

	expected, equal := checkStorageItem(ing, 2, item.ID, 1, 2)

	if !equal {
		t.Errorf("Adding new item failed. Got: %v, Want: %v", ing, expected)
	}

	i.DeleteById(nil, item.ID)

}

func TestExistingAddIngredient(t *testing.T) {

	s := newStorageForTest(t)
	tearDown := TestSetup(t, s.db)
	defer tearDown(t, s.db)
	i := newIngredientForTest(t)

	ingredient := []byte(`[{"name":"parmesan cheese", "amount":250, "unit":10 }]`)

	s.AddIngredientList(ingredient, 2)

	item, err := i.GetByName(nil, "parmesan cheese")

	if err != nil {
		t.Errorf("GetByName failed:%v", err)
	}

	ing, err := s.GetStorageIngredient(nil, 2, item.ID)

	if err != nil {
		t.Errorf("GetStorageIngredient failed:%v", err)
	}

	expected, equal := checkStorageItem(ing, 2, 3, 250, 10)

	if !equal {
		t.Errorf("Adding new item failed. Got: %v, Want: %v", ing, expected)
	}

	where := fmt.Sprintf("house_id = %v and ingredient_id = %v", 2, item.ID)

	s.DeleteFromTable(nil, where)

}
