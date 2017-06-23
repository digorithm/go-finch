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
	h := newHouseForTest(t)

	oldSchedule, err := h.GetHouseSchedule(nil, 2)
	fmt.Println("oldSchedule:")
	fmt.Println(oldSchedule)

	schedule, err := s.UpdateSchedule(nil, 2, 3, 3, 1)
	newSchedule, err := h.GetHouseSchedule(nil, 2)

	fmt.Println("newSchedule:")
	fmt.Println(newSchedule)

	s.UpdateSchedule(nil, 2, 3, 3, 4)

	row, e := schedule.RowsAffected()

	if row != 1 {
		t.Errorf("Update Schedule failed, got: %d, want: %d, with error: %d", row, 1, e)
	}
	if err != nil {
		t.Errorf("Updating schedule should work. Error: %v", err)
	}

}
