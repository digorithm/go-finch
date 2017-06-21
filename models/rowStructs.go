package models

import (
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
