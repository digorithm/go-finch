package models

import (
	"reflect"
)

type HouseScheduleRow struct {
	Week   string `db:"day"`
	Type   string `db:"type"`
	Recipe string `db:"name"`
}

type HouseStorageRow struct {
	ItemInStorageRow
	IngName  string `db:"name"`
	UnitName string `db:"name"`
}

type UserHouseRow struct {
	HouseRow
	OwnerRow
}

type UserOwnTypeRow struct {
	UserRow
	OwnerRow
}

type FullRecipeRow struct {
	RecipeRow
	StepID     int64   `db:"id"`
	Ingredient string  `db:"name"`
	Amount     float64 `db:"amount"`
	Unit       string  `db:"name"`
	Text       string  `db:"text"`
}

func createUserOwnTypeRows(data []interface{}) []UserOwnTypeRow {
	var users []UserOwnTypeRow
	var row UserOwnTypeRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.ID = v.Index(0).Interface().(int64)
		row.Email = v.Index(1).Interface().(string)
		row.Password = v.Index(2).Interface().(string)
		row.Username = v.Index(3).Interface().(string)
		row.OwnType = v.Index(4).Interface().(int64)
		row.Description = v.Index(5).Interface().(string)

		users = append(users, row)
	}

	return users
}

func createHouseStorageRows(data []interface{}) []HouseStorageRow {
	var storage []HouseStorageRow
	var row HouseStorageRow

	for i := 0; i < len(data); i++ {
		v := reflect.ValueOf(data[i])

		row.HouseID = v.Index(0).Interface().(int64)
		row.IngID = v.Index(1).Interface().(int64)
		row.Amount = v.Index(2).Interface().(float64)
		row.UnitID = v.Index(3).Interface().(int64)
		row.IngName = v.Index(4).Interface().(string)
		row.UnitName = v.Index(5).Interface().(string)

		storage = append(storage, row)
	}

	return storage
}

func createFullRecipeRows(data []interface{}) []FullRecipeRow {
	var fullRecipe []FullRecipeRow
	var row FullRecipeRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.ID = v.Index(0).Interface().(int64)
		row.Name = v.Index(1).Interface().(string)
		row.Type = v.Index(2).Interface().(string)
		row.ServesFor = v.Index(3).Interface().(int64)
		row.StepID = v.Index(4).Interface().(int64)
		row.Ingredient = v.Index(5).Interface().(string)
		row.Amount = v.Index(6).Interface().(float64)
		row.Unit = v.Index(7).Interface().(string)
		row.Text = v.Index(8).Interface().(string)

		fullRecipe = append(fullRecipe, row)
	}

	return fullRecipe
}

func createHouseScheduleRows(data []interface{}) []HouseScheduleRow {
	var schedule []HouseScheduleRow
	var row HouseScheduleRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.Week = v.Index(0).Interface().(string)
		row.Type = v.Index(1).Interface().(string)
		row.Recipe = v.Index(2).Interface().(string)

		schedule = append(schedule, row)
	}

	return schedule
}

func createUserHouseRows(houses []UserHouseRow, data []interface{}) []UserHouseRow {

	var row UserHouseRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.ID = v.Index(0).Interface().(int64)
		row.Name = v.Index(1).Interface().(string)
		row.GroceryDay = v.Index(2).Interface().(string)
		row.HouseHold = v.Index(3).Interface().(int64)
		row.OwnType = v.Index(4).Interface().(int64)
		row.Description = v.Index(5).Interface().(string)

		houses = append(houses, row)
	}

	return houses
}
