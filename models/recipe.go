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

func (r *Recipe) GetFullRecipe(tx *sqlx.Tx, recipe_id int64) ([]FullRecipeRow, error) {

	var FullRecipe []FullRecipeRow

	query := "select r.id, r.name, si.step_id, i.name, si.amount, u.name, s.text from recipe r inner join step_ingredient si on r.id = si.recipe_id inner join step s on s.id = si.step_id inner join ingredient i on i.id = si.ingredient_id inner join unit u on u.id = si.unit_id where r.id = $1"

	data, err := r.GetCompoundModel(nil, query, recipe_id)

	FullRecipe = createFullRecipeRows(FullRecipe, data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return FullRecipe, err

}
