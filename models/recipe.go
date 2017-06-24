package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func NewRecipe(db *sqlx.DB) *Recipe {
	recipe := &Recipe{}
	recipe.db = db
	recipe.table = "recipe"
	recipe.hasID = true

	return recipe
}

type Recipe struct {
	Base
}

func (r *Recipe) GetFullRecipe(tx *sqlx.Tx, recipeID int64) ([]FullRecipeRow, error) {

	query := "select r.id, r.name, r.type, r.serves_for, si.step_id, i.name, si.amount, u.name, s.text from recipe r inner join step_ingredient si on r.id = si.recipe_id inner join step s on s.id = si.step_id inner join ingredient i on i.id = si.ingredient_id inner join unit u on u.id = si.unit_id where r.id = $1"

	data, err := r.GetCompoundModel(tx, query, recipeID)

	FullRecipe := createFullRecipeRows(data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return FullRecipe, err

}

func (b *Base) GetRecipeForStruct(tx *sqlx.Tx, recipeQuery string, recipeID int64) ([]RecipeRow, error) {

	data, err := b.GetCompoundModel(tx, recipeQuery, recipeID)

	recipes := createRecipeRows(data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return recipes, err
}

// TODO: func (r *Recipe) GetNutritionalFacts(tx *sqlx.Tx, recipe_id int64) ()
