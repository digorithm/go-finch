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
	HouseID int64   `db:"house_id" json:"house_id"`
	IngID   int64   `db:"ingredient_id" json:"ingredient_id"`
	Amount  float64 `db:"amount" json:"amount"`
	UnitID  int64   `db:"unit_id" json:"unit_id"`
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

type RecipeTypeRow struct {
	Recipe_id int64 `db:"recipe_id"`
	Type_id   int64 `db:"type_id"`
}

func createRecipeRows(data []interface{}) []RecipeRow {
	var recipes []RecipeRow

	var row RecipeRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		row.ID = v.Index(0).Interface().(int64)
		row.Name = v.Index(1).Interface().(string)
		row.ServesFor = v.Index(2).Interface().(int64)

		recipes = append(recipes, row)
	}

	return recipes
}

func createItemInStorage(data []interface{}) ItemInStorageRow {

	var item ItemInStorageRow

	v := reflect.ValueOf(data[0])

	item.HouseID = v.Index(0).Interface().(int64)
	item.IngID = v.Index(1).Interface().(int64)
	item.Amount = v.Index(2).Interface().(float64)
	item.UnitID = v.Index(3).Interface().(int64)

	return item

}

func createOwnerRow(data []interface{}) *OwnerRow {

	own := &OwnerRow{}

	v := reflect.ValueOf(data[0])

	own.OwnType = v.Index(0).Interface().(int64)
	own.Description = v.Index(1).Interface().(string)

	return own
}

func createMemberOfRow(data []interface{}) MemberOfRow {

	var member MemberOfRow

	v := reflect.ValueOf(data[0])

	member.HouseID = v.Index(0).Interface().(int64)
	member.UserID = v.Index(1).Interface().(int64)
	member.OwnType = v.Index(2).Interface().(int64)

	return member
}
