package models

import (
	"errors"
	"fmt"
	"reflect"

	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Schedule struct {
	Base
}

func NewSchedule(db *sqlx.DB) *Schedule {

	schedule := &Schedule{}
	schedule.db = db
	schedule.table = "schedule"
	schedule.hasID = false

	return schedule
}

func (s *Schedule) GetHouseSchedule(tx *sqlx.Tx, houseID int64) ([]HouseScheduleRow, error) {

	query := "SELECT W.DAY, T.TYPE, R.NAME FROM RECIPE R, WEEKDAY W, MEAL_TYPE T, SCHEDULE S WHERE S.HOUSE_ID = $1 AND S.WEEK_ID = W.ID AND S.TYPE_ID = T.ID AND S.RECIPE_ID = R.ID ORDER BY WEEK_ID, TYPE_ID"

	data, err := s.GetCompoundModel(tx, query, houseID)

	schedule := createHouseScheduleRows(data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return schedule, err
}

func (s *Schedule) UpdateSchedule(tx *sqlx.Tx, hID int64, wID int64, tID int64, rID int64) (sql.Result, error) {

	query := fmt.Sprintf("UPDATE %v SET RECIPE_ID = %v WHERE HOUSE_ID = %v AND WEEK_ID = %v AND TYPE_ID = %v", s.table, rID, hID, wID, tID)

	tx, wrapInSingleTransaction, err := s.newTransactionIfNeeded(tx)

	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty")
	}
	if err != nil {
		return nil, err
	}

	sqlResult, err := tx.Exec(query)

	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	if err != nil {
		fmt.Printf("Error is: %v", err)
	}

	return sqlResult, err

}

// CreateHouseSchedule can be called from anywhere, it checks for all possible cases
func (s *Schedule) CreateHouseSchedule(tx *sqlx.Tx, houseID int64) bool {

	res := s.GetCurrentScheduleRows(tx, houseID)

	if len(res) == 0 {

		s.CreateFullSchedule(tx, houseID)

	} else {

		s.InsertMissingSchedule(tx, houseID, res)
	}

	schedule := s.GetCurrentScheduleRows(tx, houseID)

	if len(schedule) == 21 {
		return true
	}
	return false
}

// CreateFullSchedule is called when there is no instance of house_id in schedule table
func (s *Schedule) CreateFullSchedule(tx *sqlx.Tx, houseID int64) {

	i := 1
	data := make(map[string]interface{})

	for i < 8 {
		data["house_id"] = houseID
		data["week_id"] = i
		j := 1

		for j < 4 {
			data["type_id"] = j
			_, err := s.InsertIntoMultiKeyTable(tx, data)

			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			j++
		}
		i++
	}
}

func (s *Schedule) InsertMissingSchedule(tx *sqlx.Tx, houseID int64, schedule []ScheduleRow) {

	var k int
	var i int64 = 1
	var j int64 = 1
	data := make(map[string]interface{})
	data["house_id"] = houseID

	for i < 8 {
		j = 1
		for j < 4 {

			item := schedule[k]
			w := item.WeekID
			t := item.TypeID

			if (w != i) || (t != j) {
				data["week_id"] = i
				data["type_id"] = j
				s.InsertIntoMultiKeyTable(tx, data)
			} else {
				if k < len(schedule)-1 {
					k++
				}
			}
			j++
		}
		i++
	}

}

func (s *Schedule) GetCurrentScheduleRows(tx *sqlx.Tx, houseID int64) []ScheduleRow {

	query := "SELECT S.WEEK_ID, S.TYPE_ID FROM SCHEDULE S WHERE HOUSE_ID = $1 ORDER BY WEEK_ID, TYPE_ID"

	data, err := s.GetCompoundModel(tx, query, houseID)
	if err != nil {
		fmt.Printf("%v", err)
	}

	schedule := ExistingScheduleRows(data)

	return schedule
}

func ExistingScheduleRows(data []interface{}) []ScheduleRow {
	var schedule []ScheduleRow
	var s ScheduleRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		s.WeekID = v.Index(0).Interface().(int64)
		s.TypeID = v.Index(1).Interface().(int64)

		schedule = append(schedule, s)
	}

	return schedule
}
