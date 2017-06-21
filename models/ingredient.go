package models

import (
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

func (i *Ingredient) GetById(tx *sqlx.Tx, ingredient_id int64) (*IngredientRow, error) {
	ingredient := &IngredientRow{}
	query := fmt.Sprintf("select * from %v where id=$1", i.table)
	err := i.db.Get(ingredient, query, ingredient_id)

	return ingredient, err
}
