package models

import (
	"reflect"
)

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
	ID   int64  `db:"id"`
	Name string `db:"name"`
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
