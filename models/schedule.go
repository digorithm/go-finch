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
	schedule.hasID = true

	return schedule
}

func (s *Schedule) UpdateSchedule(tx *sqlx.Tx, hID int64, wID int64, tID int64, rID int64) (sql.Result, error) {

	query := fmt.Sprintf("UPDATE %v SET RECIPE_ID = %v WHERE HOUSE_ID = %v AND WEEK_ID = %v AND TYPE_ID = %v", s.table, rID, hID, wID, tID)

	tx, wrapInSingleTransaction, err := s.newTransactionIfNeeded(tx)

	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
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
