package models

import (
	"testing"

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

func TestInsertStorage(t *testing.T) {

	s := newStorageForTest(t)
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
