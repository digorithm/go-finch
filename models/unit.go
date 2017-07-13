package models

import (
	"fmt"
	"reflect"

	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type Unit struct {
	Base
}

func NewUnit(db *sqlx.DB) *Unit {

	unit := &Unit{}
	unit.db = db
	unit.table = "unit"
	unit.hasID = true

	return unit
}

func (u *Unit) GetAllUnits(tx *sqlx.Tx) ([]byte, error) {

	query := "SELECT U.ID, U.NAME FROM UNIT U "

	data, err := u.GetCompoundModelWithoutID(tx, query)

	if err != nil {
		fmt.Printf("getAllUnits failed: %v", err)
	}

	return AllFixedTablesJSON(data), err
}

func AllFixedTablesJSON(data []interface{}) []byte {

	allUnits := make(map[string]interface{})

	for _, d := range data {

		v := reflect.ValueOf(d)

		var key string

		key = v.Index(1).Interface().(string)
		allUnits[key] = v.Index(0).Interface().(int64)

	}

	allUnitsJSON, _ := json.Marshal(allUnits)
	return allUnitsJSON

}
