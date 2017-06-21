package models

import (
	"fmt"
	"reflect"
)

type HouseRow struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
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

type HouseStorageRow struct {
	HouseRow
	OwnerRow
}

type RecipeRow struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func createRecipeRows(recipes []RecipeRow, data []interface{}) []RecipeRow {

	var row RecipeRow
	var rows []RecipeRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.ID = v.Index(0).Interface().(int64)
		row.Name = v.Index(1).Interface().(string)

		rows = append(rows, row)
	}

	return rows
}

func createUserOwnTypeRows(users []UserOwnTypeRow, data []interface{}) []UserOwnTypeRow {

	var row UserOwnTypeRow
	var rows []UserOwnTypeRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.ID = v.Index(0).Interface().(int64)
		row.Email = v.Index(1).Interface().(string)
		row.Password = v.Index(2).Interface().(string)
		row.Username = v.Index(3).Interface().(string)
		row.OwnType = v.Index(4).Interface().(int64)
		row.Description = v.Index(5).Interface().(string)
	}

	fmt.Printf("%v", rows)
	return rows
}
