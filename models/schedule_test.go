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
	fmt.Println("2")
	schedule, err := s.UpdateSchedule(nil, 2, 1, 3, 2)

	if err != nil {
		t.Errorf("Updating schedule should work. Error: %v", err)
	}

	fmt.Println(schedule)

}
