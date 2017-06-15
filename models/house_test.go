package models

import (
	_ "github.com/lib/pq"
	"testing"
)

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
