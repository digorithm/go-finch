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
	fmt.Println("1")
	s := newScheduleForTest(t)
	h := newHouseForTest(t)
	fmt.Println("2")

	oldSchedule, err := h.GetHouseSchedule(nil, 2)
	fmt.Println("oldSchedule:")
	fmt.Println(oldSchedule)

	schedule, err := s.UpdateSchedule(nil, 2, 3, 3, 1)
	newSchedule, err := h.GetHouseSchedule(nil, 2)

	fmt.Println("newSchedule:")
	fmt.Println(newSchedule)

	if err != nil {
		t.Errorf("Updating schedule should work. Error: %v", err)
	}

	fmt.Println(schedule)

}
