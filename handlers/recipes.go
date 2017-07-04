package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/digorithm/meal_planner/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// GetHouseRecipesHandler will handle the API call to get all recipes of a given house,
// it will return a JSON for the endpoint /recipes/house/{house_id} as described in the docs
func GetHouseRecipesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	houseID, err := strconv.Atoi(vars["house_id"])

	db := r.Context().Value("db").(*sqlx.DB)

	houseObj := models.NewHouse(db)
	recipeObj := models.NewRecipe(db)

	if err != nil {
		fmt.Println(err)
	}

	recipes, err := houseObj.GetHouseRecipes(nil, int64(houseID))

	if err != nil {
		fmt.Printf("Something went wrong while fecthing the recipes. Error: %v", err)
	}

	var fullRecipes [][]models.FullRecipeRow

	RecipesTypes := make(map[int64][]string)

	for _, recipe := range recipes {
		fullRecipe, err := recipeObj.GetFullRecipe(nil, recipe.ID)
		recipeTypes, err := recipeObj.GetRecipeType(nil, recipe.ID)

		RecipesTypes[recipe.ID] = recipeTypes

		if err != nil {
			fmt.Printf("Error fecthing full recipe. Error: %v", err)
		}
		fullRecipes = append(fullRecipes, fullRecipe)
	}
	JSONResponse := buildFullRecipeJSONResponse(fullRecipes, RecipesTypes)

	w.Write(JSONResponse)
}

func GetUserRecipesHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["user_id"])

	db := r.Context().Value("db").(*sqlx.DB)

	userObj := models.NewUser(db)
	recipeObj := models.NewRecipe(db)

	if err != nil {
		fmt.Println(err)
	}

	recipes, err := userObj.GetUserRecipes(nil, int64(userID))

	if err != nil {
		fmt.Printf("Something went wrong while fecthing the recipes. Error: %v", err)
	}

	var fullRecipes [][]models.FullRecipeRow

	RecipesTypes := make(map[int64][]string)

	for _, recipe := range recipes {
		fullRecipe, err := recipeObj.GetFullRecipe(nil, recipe.ID)
		recipeTypes, err := recipeObj.GetRecipeType(nil, recipe.ID)

		RecipesTypes[recipe.ID] = recipeTypes

		if err != nil {
			fmt.Printf("Error fecthing full recipe. Error: %v", err)
		}
		fullRecipes = append(fullRecipes, fullRecipe)
	}
	JSONResponse := buildFullRecipeJSONResponse(fullRecipes, RecipesTypes)

	w.Write(JSONResponse)
}

func buildFullRecipeJSONResponse(recipes [][]models.FullRecipeRow, RecipesTypes map[int64][]string) []byte {

	finalRecipes := make([]map[string]interface{}, 0, 0)

	for _, recipe := range recipes {

		finalRecipe := make(map[string]interface{})

		if len(recipe) > 0 {
			recipeID := recipe[0].ID
			recipeName := recipe[0].Name
			servesFor := recipe[0].ServesFor
			recipeTypes := RecipesTypes[recipeID]

			steps := make([]map[string]interface{}, 0, 0)
			stepsIngredients := make(map[int64][]map[string]interface{})

			for _, row := range recipe {
				step := make(map[string]interface{})
				step["step_id"] = row.StepID
				step["text"] = row.Text

				singleStepIntegredient := make(map[string]interface{})
				singleStepIntegredient["name"] = row.Ingredient
				singleStepIntegredient["amount"] = row.Amount
				singleStepIntegredient["unit"] = row.Unit
				stepsIngredients[row.StepID] = append(stepsIngredients[row.StepID], singleStepIntegredient)

				step["step_ingredients"] = stepsIngredients[row.StepID]
				steps = append(steps, step)
			}

			steps = removeDuplicateStepID(steps)

			finalRecipe["id"] = recipeID
			finalRecipe["name"] = recipeName
			finalRecipe["type"] = recipeTypes
			finalRecipe["serves_for"] = servesFor
			finalRecipe["steps"] = steps

		}
		finalRecipes = append(finalRecipes, finalRecipe)

	}

	finalRecipesJSON, _ := json.MarshalIndent(finalRecipes, "", "    ")

	return finalRecipesJSON
}

// removeDuplicateStepID is the worst awful shameless workaround ever. Don't even try to understand why this was created.
func removeDuplicateStepID(steps []map[string]interface{}) []map[string]interface{} {
	for k, step := range steps {
		if k+1 < len(steps) {
			stepID := step["step_id"].(int64)
			nextStepID := steps[k+1]["step_id"].(int64)
			if stepID < nextStepID {
				// Continue to the next step
				continue
			} else if stepID == nextStepID {
				// same step, remove this current index
				steps = steps[:k+copy(steps[k:], steps[k+1:])]
			}
		}
	}
	return steps
}
