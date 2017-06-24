package models

import (
	"errors"
	"fmt"

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

// CreateHouseSchedule can be called from anywhere, it checks for all possible cases
func (s *Schedule) CreateHouseSchedule(tx *sqlx.Tx, houseID int64) bool {

	query := "SELECT S.WEEK_ID, S.TYPE_ID, S.RECIPE_ID FROM SCHEDULE S WHERE HOUSE_ID = $1"
	data, err := s.GetCompoundModel(tx, query, houseID)
	if err != nil {
		fmt.Printf("%v", err)
	}

	res := createHouseScheduleRows(data)

	if len(res) == 0 {
		fmt.Println("In res == 0")
		s.CreateFullSchedule(tx, houseID)

	} else {
		s.InsertMissingSchedule(tx, houseID)
	}
	schedule := s.GetScheduleRow(tx, houseID)
	fmt.Println("schedule length:")
	fmt.Println(len(schedule))
	if err != nil {
		fmt.Printf("%v", err)
	}
	if len(schedule) == 21 {
		return true
	}
	return false
}

// CreateFullSchedule is called when there is no instance of house_id in schedule table
func (s *Schedule) CreateFullSchedule(tx *sqlx.Tx, houseID int64) {

	fmt.Println("In createfullschedule")

	i := 1
	data := make(map[string]interface{})

	for i < 8 {
		data["house_id"] = houseID
		data["week_id"] = i
		j := 1

		for j < 4 {
			data["type_id"] = j
			fmt.Println("Before insert:")
			fmt.Println(data)
			res, err := s.InsertIntoMultiKeyTable(tx, data)

			fmt.Println(res)
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			j++
		}
		i++
	}
}

func (s *Schedule) InsertMissingSchedule(tx *sqlx.Tx, houseID int64) {

	schedule := s.GetScheduleRow(tx, houseID)
	k := 0
	var i int64 = 1
	var j int64 = 1

	for i < 8 {
		for j < 4 {
			item := schedule[k]
			w := item.WeekID
			t := item.TypeID
			if (w != i) || (t != j) {
				data := make(map[string]interface{})
				data["house_id"] = houseID
				data["week_id"] = i
				data["type_id"] = j

				s.InsertIntoTable(tx, data)
			}
			j++
			k++
		}
		i++
	}

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

func (s *Schedule) GetHouseSchedule(tx *sqlx.Tx, houseID int64) ([]HouseScheduleRow, error) {

	query := "SELECT W.DAY, T.TYPE, R.NAME FROM RECIPE R, WEEKDAY W, MEAL_TYPE T, SCHEDULE S WHERE S.HOUSE_ID = $1 AND S.WEEK_ID = W.ID AND S.TYPE_ID = T.ID AND S.RECIPE_ID = R.ID ORDER BY WEEK_ID, TYPE_ID"

	data, err := s.GetCompoundModel(tx, query, houseID)

	schedule := createHouseScheduleRows(data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return schedule, err
}

func (s *Schedule) GetScheduleRow(tx *sqlx.Tx, houseID int64) []ScheduleRow {

	query := "SELECT * FROM SCHEDULE S WHERE S.HOUSE_ID = $1 ORDER BY WEEK_ID, TYPE_ID"

	data, err := s.GetCompoundModel(tx, query, houseID)
	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Println("this is the whole schedule:")
	fmt.Println(data)
	schedule := createScheduleRows(data)

	return schedule
}
