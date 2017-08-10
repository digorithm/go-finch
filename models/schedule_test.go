package models

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func newScheduleForTest(t *testing.T) *Schedule {
	return NewSchedule(newDbForTest(t))
}

func TestUpdateSchedule(t *testing.T) {

	s := newScheduleForTest(t)
	tearDown := TestSetup(t, s.db)
	defer tearDown(t, s.db)

	_, err := s.GetHouseSchedule(nil, 2)

	schedule, err := s.UpdateSchedule(nil, 2, 3, 3, 1)
	_, err = s.GetHouseSchedule(nil, 2)

	s.UpdateSchedule(nil, 2, 3, 3, 4)

	row, e := schedule.RowsAffected()

	if row != 1 {
		t.Errorf("Update Schedule failed, got: %d, want: %d, with error: %d", row, 1, e)
	}
	if err != nil {
		t.Errorf("Updating schedule should work. Error: %v", err)
	}

}

func TestGetHouseSchedule(t *testing.T) {

	s := newScheduleForTest(t)
	tearDown := TestSetup(t, s.db)
	defer tearDown(t, s.db)

	var s1 = createVarsForGetSchedule("Tuesday", "Breakfast", "No Flour Pancake", 3)
	var s2 = createVarsForGetSchedule("Wednesday", "Lunch", "Roast Chicken", 4)
	var s3 = createVarsForGetSchedule("Saturday", "Breakfast", "No Flour Pancake", 3)
	var result []HouseScheduleRow

	schedule, err := s.GetHouseSchedule(nil, 2)
	if err != nil {
		t.Errorf("Getting house schedule should work. Error: %v", err)
	}

	result = append(result, s1, s2, s3)
	i := 0
	for i < len(schedule) {
		if result[i] != schedule[i] {
			t.Errorf("Get House Schedule failed, got: %d, want: %d", schedule[i], result[i])
		}

		i++
	}

}

func TestCreateFullSchedule(t *testing.T) {

	s := newScheduleForTest(t)
	tearDown := TestSetup(t, s.db)
	defer tearDown(t, s.db)

	h := newHouseForTest(t)

	house, err := h.CreateHouse(nil, "Our Home", "Friday", 2)
	if err != nil {
		fmt.Printf("%v", err)
	}

	houseID := house.ID
	inserted := s.CreateHouseSchedule(nil, houseID)

	if inserted != true {
		t.Errorf("Full schedule insertion failed, got: %d, want: true", inserted)
	}

	_, er := h.DeleteById(nil, houseID)
	if er != nil {
		t.Fatalf("Deleting house by id should not fail. Error: %v", err)
	}

}

func TestInsertMissingSchedule(t *testing.T) {

	s := newScheduleForTest(t)
	tearDown := TestSetup(t, s.db)
	defer tearDown(t, s.db)

	inserted := s.CreateHouseSchedule(nil, 1)

	if inserted != true {
		t.Errorf("Full schedule insertion failed, got: %d, want: true", inserted)
	}
}
