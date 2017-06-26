package models

import "reflect"

type MealRow struct {
	ID   int64  `db:"id"`
	Type string `db:"type"`
}

type WeekRow struct {
	ID  int64  `db:"id"`
	Day string `db:"day"`
}

type HouseRow struct {
	ID         int64  `db:"id"`
	Name       string `db:"name"`
	GroceryDay string `db:"grocery_day"`
	HouseHold  int64  `db:"household_number"`
}

type UserRow struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password"`
	Username string `db:"username"`
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

type UnitRow struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type StepRow struct {
	RecipeID int64  `db:"recipe_id`
	ID       int64  `db:"id"`
	Text     string `db:"text"`
}

type StepIngredientRow struct {
	RecipeID int64   `db:"recipe_id"`
	StepID   int64   `db:"step_id"`
	IngID    int64   `db:"ingredient_id"`
	UnitID   int64   `db:"unit_id"`
	Amount   float64 `db:"amount"`
}

type ItemInStorageRow struct {
	HouseID int64   `db:"house_id"`
	IngID   int64   `db:"ingredient_id"`
	Amount  float64 `db:"amount`
	UnitID  int64   `db:"unit_id"`
}

type UserRecipeRow struct {
	UserID   int64 `db:"user_id"`
	RecipeID int64 `db:"recipe_id"`
}

type ScheduleRow struct {
	HouseID  int64 `db:"house_id"`
	WeekID   int64 `db:"week_id"`
	TypeID   int64 `db:"type_id"`
	RecipeId int64 `db:"recipe_id"`
}

type OwnerRow struct {
	OwnType     int64  `db:"own_type"`
	Description string `db:"description"`
}

type MemberOfRow struct {
	UserID  int64 `db:"user_id"`
	HouseID int64 `db:"house_id"`
	OwnType int64 `db:"own_type"`
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
