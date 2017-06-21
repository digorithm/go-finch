package models

import (
	"reflect"
)

type HouseScheduleRow struct {
	Week   string `db:"day"`
	Type   string `db:"type"`
	Recipe string `db:"name"`
}

type UserHouseRow struct {
	HouseRow
	OwnerRow
}

type HouseRow struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
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

type OwnerRow struct {
	OwnType     int64  `db:"own_type"`
	Description string `db:"description"`
}

type UserOwnTypeRow struct {
	UserRow
	OwnerRow
}

type RecipeRow struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	image string `db:recipe`
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
	RecipeID   int64   `db:"id"`
	RecipeName string  `db:"name"`
	StepID     int64   `db:"id"`
	Ingredient string  `db:"name"`
	Amount     float64 `db:"amount"`
	Unit       string  `db:"name"`
	Text       string  `db:"text"`
}

func createRecipeRows(recipes []RecipeRow, data []interface{}) []RecipeRow {

	var row RecipeRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.ID = v.Index(0).Interface().(int64)
		row.Name = v.Index(1).Interface().(string)

		recipes = append(recipes, row)
	}

	return recipes
}

func createUserOwnTypeRows(users []UserOwnTypeRow, data []interface{}) []UserOwnTypeRow {

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

func createHouseStorageRows(storage []HouseStorageRow, data []interface{}) []HouseStorageRow {

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

func createFullRecipeRows(fullRecipe []FullRecipeRow, data []interface{}) []FullRecipeRow {

	var row FullRecipeRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.RecipeID = v.Index(0).Interface().(int64)
		row.RecipeName = v.Index(1).Interface().(string)
		row.StepID = v.Index(2).Interface().(int64)
		row.Ingredient = v.Index(3).Interface().(string)
		row.Amount = v.Index(4).Interface().(float64)
		row.Unit = v.Index(5).Interface().(string)
		row.Text = v.Index(6).Interface().(string)

		fullRecipe = append(fullRecipe, row)
	}

	return fullRecipe
}

func createHouseScheduleRows(schedule []HouseScheduleRow, data []interface{}) []HouseScheduleRow {

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
		row.OwnType = v.Index(2).Interface().(int64)
		row.Description = v.Index(3).Interface().(string)

		houses = append(houses, row)
	}

	return houses
}
