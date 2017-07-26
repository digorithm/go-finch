package models

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewIngredient(db *sqlx.DB) *Ingredient {
	ingredient := &Ingredient{}
	ingredient.db = db
	ingredient.table = "ingredient"
	ingredient.hasID = true

	return ingredient
}

type Ingredient struct {
	Base
}

func (i *Ingredient) ingredientRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*IngredientRow, error) {
	ingredientId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}
	return i.GetById(tx, ingredientId)
}

func (i *Ingredient) GetById(tx *sqlx.Tx, ingredient_id int64) (*IngredientRow, error) {
	ingredient := &IngredientRow{}
	query := fmt.Sprintf("select * from %v where id=$1", i.table)
	err := i.db.Get(ingredient, query, ingredient_id)

	return ingredient, err
}

func (i *Ingredient) GetByName(tx *sqlx.Tx, name string) (*IngredientRow, error) {

	ingredient := &IngredientRow{}
	query := fmt.Sprintf("select * from %v where name=$1", i.table)
	err := i.db.Get(ingredient, query, name)

	if err != nil {
		return nil, err
	}

	return ingredient, err
}

func (i *Ingredient) AddIngredients(tx *sqlx.Tx, ingredients []string) ([]int64, error) {

	var insertedIDs []int64
	var err error

	for _, ing := range ingredients {
		// Check if ingredient already exists in the db
		iRow, err := i.GetByName(tx, ing)
		if iRow == nil && err != nil {
			// Add ingredient to db
			iRow, err, _ = i.AddIngredient(tx, ing)
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
		}
		insertedIDs = append(insertedIDs, iRow.ID)
	}
	return insertedIDs, err
}

func (i *Ingredient) AddIngredient(tx *sqlx.Tx, name string) (*IngredientRow, error, int64) {
	data := make(map[string]interface{})
	data["name"] = name
	// TODO: get ingredient nutrients externally, for we just insert any number
	data["carb_per_100g"] = 0.0
	data["protein_per_100g"] = 0.0
	data["fiber_per_100g"] = 0.0
	data["fat_per_100g"] = 0.0
	data["calories_per_100g"] = 0.0

	result, err := i.InsertIntoTable(tx, data)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	InsertedID, err := result.LastInsertId()
	IRow, err := i.ingredientRowFromSqlResult(tx, result)
	return IRow, err, InsertedID
}
