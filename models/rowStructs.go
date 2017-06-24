package models

import (
	"reflect"
)

type HouseScheduleRow struct {
	Week   string `db:"day"`
	Type   string `db:"type"`
	Recipe string `db:"name"`
}

type ScheduleRow struct {
	HouseID  int64 `db:"house_id"`
	WeekID   int64 `db:"week_id"`
	TypeID   int64 `db:"type_id"`
	RecipeId int64 `db:"recipe_id"`
}

type HouseRow struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	GroceryDay string `db:"grocery_day"`
	HouseHold  int64  `db:"household_number"`
}

type HouseStorageRow struct {
	Name   string  `db:"name"`
	Amount float64 `db:"amount`
	Unit   string  `db:"name"`
}

type UserRow struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Username string `db:"username"`
}

type UserHouseRow struct {
	HouseRow
	OwnerRow
}

type OwnerRow struct {
	OwnType     int64  `db:"own_type"`
	Description string `db:"description"`
}

type UserOwnTypeRow struct {
	UserRow
	OwnerRow
}

type RecipeRow struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	Type      string `db:"type"`
	ServesFor int64  `db:"serves_for"`
}

type IngredientRow struct {
	ID              int64   `db:"id"`
	Name            string  `db:"name"`
	CarbPer100G     float64 `db:"carb_per_100g"`
	ProteinPer100G  float64 `db:"protein_per_100g"`
	FatPer100G      float64 `db:"fat_per_100g"`
	FiberPer100G    float64 `db:"fiber_per_100g"`
	CaloriesPer100G float64 `db:"calories_per_100g"`
}

type FullRecipeRow struct {
	RecipeRow
	StepID     int64   `db:"id"`
	Ingredient string  `db:"name"`
	Amount     float64 `db:"amount"`
	Unit       string  `db:"name"`
	Text       string  `db:"text"`
}

func createRecipeRows(data []interface{}) []RecipeRow {
	var recipes []RecipeRow

	var row RecipeRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.ID = v.Index(0).Interface().(int64)
		row.Name = v.Index(1).Interface().(string)
		row.Type = v.Index(2).Interface().(string)
		row.ServesFor = v.Index(3).Interface().(int64)

		recipes = append(recipes, row)
	}

	return recipes
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

		row.Name = v.Index(0).Interface().(string)
		row.Amount = v.Index(1).Interface().(float64)
		row.Unit = v.Index(2).Interface().(string)

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

func createScheduleRows(data []interface{}) []ScheduleRow {
	var schedule []ScheduleRow
	var s ScheduleRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		s.HouseID = v.Index(0).Interface().(int64)
		s.WeekID = v.Index(1).Interface().(int64)
		s.TypeID = v.Index(2).Interface().(int64)

		schedule = append(schedule, s)
	}

	return schedule
}
