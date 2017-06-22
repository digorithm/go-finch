package models

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func newHouseForTest(t *testing.T) *House {
	return NewHouse(newDbForTest(t))
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
		t.Errorf("Getting house recipes should work. Error: %v", err)
	}

	var r1 RecipeRow
	var r2 RecipeRow
	r1.ID = 1
	r1.Name = "Baked Potato"
	r1.Type = "Lunch/Dinner"
	r1.ServesFor = 4

	r2.ID = 4
	r2.Name = "Roast Chicken"
	r2.Type = "Lunch/Dinner"
	r2.ServesFor = 4
	var result []RecipeRow

	result = append(result, r1, r2)
	i := 0
	for i < len(recipes) {
		if result[i] != recipes[i] {
			t.Errorf("House Recipes, got: %d, want: %d", recipes[i], result[i])
		}

		i++
	}

	fmt.Println("recipes")
}

func TestGetHouseSchedule(t *testing.T) {
	h := newHouseForTest(t)

	schedule, err := h.GetHouseSchedule(nil, 2)
	if err != nil {
		t.Errorf("Getting house schedule should work. Error: %v", err)
	}

	var s1 HouseScheduleRow
	var s2 HouseScheduleRow

	s1.Week = "Saturday"
	s1.Type = "Breakfast"
	s1.Recipe = "No Flour Pancake"
	s2.Week = "Tuesday"
	s2.Type = "Breakfast"
	s2.Recipe = "No Flour Pancake"

	var result []HouseScheduleRow

	result = append(result, s1, s2)

	i := 0
	for i < len(schedule) {
		if result[i] != schedule[i] {
			t.Errorf("House Schedule, got: %d, want: %d", schedule[i], result[i])
		}

		i++
	}

	fmt.Println("schedule")

}

func TestUpdateHouseHold(t *testing.T) {
	h := newHouseForTest(t)

	house, err := h.UpdateHouseHold(nil, 3, 2)
	if err != nil {
		t.Errorf("Updating household should work. Error: %v", err)
	}

	fmt.Println(house)

}
