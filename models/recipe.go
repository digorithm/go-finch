package models

import (
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/jmoiron/sqlx"
)

func NewRecipe(db *sqlx.DB) *Recipe {
	recipe := &Recipe{}
	recipe.db = db
	recipe.table = "recipe"
	recipe.hasID = true

	return recipe
}

func NewRecipeType(db *sqlx.DB) *RecipeType {
	rType := &RecipeType{}
	rType.db = db
	rType.table = "recipe_type"
	rType.hasID = false

	return rType
}

type RecipeType struct {
	Base
}

type Recipe struct {
	Base
}

func (r *Recipe) GetRecipeType(tx *sqlx.Tx, recipeID int64) ([]string, error) {

	var types []string

	query := "select type from recipe_type inner join meal_type on meal_type.id = recipe_type.type_id where recipe_type.recipe_id = $1"

	rows, err := r.db.Queryx(query, recipeID)

	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		cols, err := rows.SliceScan()
		if err != nil {
			fmt.Printf("%v", err)
		}
		for _, col := range cols {
			types = append(types, col.(string))
		}
	}
	return types, err
}

// GetAllRecipes used by house and user
func (r *Recipe) GetAllRecipes(tx *sqlx.Tx) ([]RecipeRow, error) {

	query := "SELECT R.ID, R.NAME, R.SERVES_FOR FROM RECIPE R INNER JOIN USER_RECIPE U ON R.ID = U.RECIPE_ID"

	return r.GetAllRecipesForStruct(tx, query)
}

func (r *Recipe) GetAllRecipesByStringSearch(tx *sqlx.Tx, stringSearch string) ([]RecipeRow, error) {

	query := fmt.Sprintf("SELECT R.ID, R.NAME, R.SERVES_FOR FROM RECIPE R INNER JOIN USER_RECIPE U ON R.ID = U.RECIPE_ID WHERE R.NAME ILIKE '%%%v%%'", stringSearch)

	return r.GetAllRecipesForStruct(tx, query)
}

func (r *Recipe) GetFullRecipe(tx *sqlx.Tx, recipeID int64) ([]FullRecipeRow, error) {

	var FullRecipe []FullRecipeRow

	query := "select r.id, r.name, r.serves_for, si.step_id, i.name as ingredient, si.amount, u.name as unit, s.text from recipe r inner join step_ingredient si on r.id = si.recipe_id inner join step s on s.id = si.step_id and s.recipe_id = r.id inner join ingredient i on i.id = si.ingredient_id inner join unit u on u.id = si.unit_id where r.id = $1"

	data, err := r.GetCompoundModel(tx, query, recipeID)

	FullRecipe = createFullRecipeRows(data)

	if err != nil {
		fmt.Printf("Get Full Recipe failed: %v", err)
	}

	return FullRecipe, err

}
func (r *Recipe) GetFullRecipesByStringSearch(tx *sqlx.Tx, stringSearch string) ([][]FullRecipeRow, map[int64][]string, error) {

	recipes, err := r.GetAllRecipesByStringSearch(tx, stringSearch)

	if err != nil {
		fmt.Printf("Something went wrong while fecthing the recipes. Error: %v", err)
	}

	var fullRecipes [][]FullRecipeRow

	RecipesTypes := make(map[int64][]string)

	for _, recipe := range recipes {
		fullRecipe, err := r.GetFullRecipe(nil, recipe.ID)
		recipeTypes, err := r.GetRecipeType(nil, recipe.ID)

		RecipesTypes[recipe.ID] = recipeTypes

		if err != nil {
			fmt.Printf("Error fecthing full recipe. Error: %v", err)
		}
		fullRecipes = append(fullRecipes, fullRecipe)
	}

	return fullRecipes, RecipesTypes, err
}

func (r *Recipe) GetFullRecipes(tx *sqlx.Tx) ([][]FullRecipeRow, map[int64][]string, error) {

	recipes, err := r.GetAllRecipes(tx)

	if err != nil {
		fmt.Printf("Something went wrong while fecthing the recipes. Error: %v", err)
	}

	var fullRecipes [][]FullRecipeRow

	RecipesTypes := make(map[int64][]string)

	for _, recipe := range recipes {
		fullRecipe, err := r.GetFullRecipe(nil, recipe.ID)
		recipeTypes, err := r.GetRecipeType(nil, recipe.ID)

		RecipesTypes[recipe.ID] = recipeTypes

		if err != nil {
			fmt.Printf("Error fecthing full recipe. Error: %v", err)
		}
		fullRecipes = append(fullRecipes, fullRecipe)
	}

	return fullRecipes, RecipesTypes, err
}

func (b *Base) GetRecipeForStruct(tx *sqlx.Tx, recipeQuery string, recipeID int64) ([]RecipeRow, error) {

	data, err := b.GetCompoundModel(tx, recipeQuery, recipeID)

	recipes := createRecipeRows(data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return recipes, err
}

func (b *Base) GetAllRecipesForStruct(tx *sqlx.Tx, recipeQuery string) ([]RecipeRow, error) {

	data, err := b.GetCompoundModelWithoutID(tx, recipeQuery)

	recipes := createRecipeRows(data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return recipes, err
}

func (r *RecipeType) AddRecipeType(tx *sqlx.Tx, recipeID int64, typeID int64) error {

	data := make(map[string]interface{})
	data["recipe_id"] = recipeID
	data["type_id"] = typeID

	_, err := r.InsertIntoMultiKeyTable(tx, data)

	return err

}

func (r *RecipeType) GetRecipeType(tx *sqlx.Tx, recipeID int64) ([]RecipeTypeNameRow, error) {

	query := "SELECT R.RECIPE_ID, R.TYPE_ID, M.TYPE FROM RECIPE_TYPE R INNER JOIN MEAL_TYPE M ON R.TYPE_ID = M.ID WHERE R.RECIPE_ID = $1"
	data, err := r.GetCompoundModel(tx, query, recipeID)

	if err != nil {
		fmt.Printf("%v", err)
	}

	result := createRecipeTypeNameRows(data)

	return result, err
}

func (r *Recipe) UpdateRecipeType(tx *sqlx.Tx, recipeID int64, typeToUpdate []string) error {

	// Delete all types for this recipe

	query := "DELETE FROM recipe_type WHERE recipe_id = $1"
	_, err := r.db.Exec(query, recipeID)

	if err != nil {
		return err
	}

	for _, t := range typeToUpdate {

		TypeID, err := r.GetMealTypeByName(nil, t)

		if err != nil {
			return err
		}

		err = r.AddType(tx, recipeID, TypeID.ID)

		if err != nil {
			return err
		}

	}
	return nil
}

func (r *Recipe) AddType(tx *sqlx.Tx, recipeID, TypeID int64) error {
	query := fmt.Sprintf("INSERT INTO recipe_type (recipe_id, type_id) values (%v,%v)", recipeID, TypeID)
	_, err := r.db.Exec(query)

	if err != nil {
		fmt.Printf("Something went wrong. Error: %v", err)
		return err
	}
	return nil
}

func (r *Recipe) GetMealTypeByName(tx *sqlx.Tx, name string) (*MealRow, error) {
	mealType := &MealRow{}
	query := "SELECT * FROM meal_type WHERE type=$1"
	err := r.db.Get(mealType, query, name)

	return mealType, err
}

func (r *Recipe) UpdateRecipeStepIngredient(RecipeID, StepID int64, JSONRequest []byte) {

	// First we delete the step_id from the database

	query := "delete from step_ingredient where step_id = $1 and recipe_id = $2"
	_, err := r.db.Exec(query, StepID, RecipeID)

	query = "delete from step where id = $1 and recipe_id = $2"
	_, err = r.db.Exec(query, StepID, RecipeID)

	if err != nil {
		fmt.Println(err)
	}

	var step_db Base
	step_db.db = r.db
	step_db.table = "step"
	step_db.hasID = true

	var step_ingredient_db Base
	step_ingredient_db.db = r.db
	step_ingredient_db.table = "step_ingredient"
	step_ingredient_db.hasID = false

	ingredientObj := NewIngredient(r.db)

	stepData := make(map[string]interface{})
	stepData["recipe_id"] = RecipeID
	stepData["id"], _ = jsonparser.GetInt(JSONRequest, "step_id")
	stepData["text"], _ = jsonparser.GetString(JSONRequest, "text")

	_, err = step_db.InsertIntoTable(nil, stepData)

	if err != nil {
		fmt.Printf("Error while adding step into DB. Error: %v", err)
	}
	jsonparser.ArrayEach(JSONRequest, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

		stepIngredientData := make(map[string]interface{})
		stepIngredientData["recipe_id"] = stepData["recipe_id"]
		stepIngredientData["unit_id"], _ = jsonparser.GetInt(value, "unit")
		stepIngredientData["amount"], _ = jsonparser.GetFloat(value, "amount")
		stepIngredientData["step_id"] = stepData["id"]

		// Check if ingredient ID exists in the DB
		stepIngredientDataName, _ := jsonparser.GetString(value, "name")
		iRow, _ := ingredientObj.GetByName(nil, stepIngredientDataName)
		if iRow == nil {
			addedIRow, _ := ingredientObj.AddIngredient(nil, stepIngredientDataName)
			stepIngredientData["ingredient_id"] = addedIRow.ID
		} else {
			stepIngredientData["ingredient_id"] = iRow.ID
		}
		_, err = step_ingredient_db.InsertIntoTable(nil, stepIngredientData)
		if err != nil {
			fmt.Printf("Error while adding step ingredient into DB. Error: %v", err)
		}

	}, "step_ingredients")
}

// TODO: func (r *Recipe) GetNutritionalFacts(tx *sqlx.Tx, recipe_id int64) ()
