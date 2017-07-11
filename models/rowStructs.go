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
	ID          int64  `db:"id"`
	Email       string `db:"email"`
	Password    string `db:"password"`
	Username    string `db:"username"`
	OwnType     int64  `db:"own_type"`
	Description string `db:"description"`
}

type HouseUserOwnRow struct {
	HID         int64  `db:"house_id"`
	HouseNumber int64  `db:"household_number"`
	UID         int64  `db:"user_id"`
	Username    string `db:"username"`
	OwnType     string `db:"description"`
}

type FullRecipeRow struct {
	RecipeRow
	StepID     int64   `db:"step_id"`
	Ingredient string  `db:"ingredient"`
	Amount     float64 `db:"amount"`
	Unit       string  `db:"unit"`
	Text       string  `db:"text"`
}

type RecipeTypeNameRow struct {
	RecipeTypeRow
	TypeName string `db:"type"`
}

func createHouseUserOwnRows(data []interface{}) []HouseUserOwnRow {
	var users []HouseUserOwnRow
	var row HouseUserOwnRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.HID = v.Index(0).Interface().(int64)
		row.HouseNumber = v.Index(1).Interface().(int64)
		row.UID = v.Index(2).Interface().(int64)
		row.Username = v.Index(3).Interface().(string)
		row.OwnType = v.Index(4).Interface().(string)

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
		row.ServesFor = v.Index(2).Interface().(int64)
		row.StepID = v.Index(3).Interface().(int64)
		row.Ingredient = v.Index(4).Interface().(string)
		row.Amount = v.Index(5).Interface().(float64)
		row.Unit = v.Index(6).Interface().(string)
		row.Text = v.Index(7).Interface().(string)

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

func createUserHouseRows(data []interface{}) []UserHouseRow {

	var row UserHouseRow
	var houses []UserHouseRow

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

func createRecipeTypeNameRows(data []interface{}) []RecipeTypeNameRow {

	var row RecipeTypeNameRow
	var recipes []RecipeTypeNameRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.Recipe_id = v.Index(0).Interface().(int64)
		row.Type_id = v.Index(1).Interface().(int64)
		row.TypeName = v.Index(2).Interface().(string)
		recipes = append(recipes, row)
	}

	return recipes

}
