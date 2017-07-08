package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"io/ioutil"

	"github.com/digorithm/meal_planner/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func AddRecipesHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing data in the recipe"))
		return
	}

	db := r.Context().Value("db").(*sqlx.DB)

	userObj := models.NewUser(db)
	recipeObj := models.NewRecipe(db)

	authorID := r.URL.Query()["author"]

	if len(authorID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID missing in the request to add a recipe"))
		return
	}

	authorIDString, _ := strconv.Atoi(authorID[0])

	returnedRecipe, err := userObj.AddRecipe(nil, body, int64(authorIDString))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Something went wrong when adding a recipe. Error: %v", err)))
		return
	}

	recipeTypes, _ := recipeObj.GetRecipeType(nil, returnedRecipe[0].ID)

	// Necessary to build the JSON response
	var fullRecipes [][]models.FullRecipeRow
	RecipesTypes := make(map[int64][]string)

	RecipesTypes[returnedRecipe[0].ID] = recipeTypes

	fullRecipes = append(fullRecipes, returnedRecipe)

	JSONResponse := buildFullRecipeJSONResponse(fullRecipes, RecipesTypes)

	w.WriteHeader(http.StatusCreated)
	w.Write(JSONResponse)
}

func DeleteRecipesHandler(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)

	recipeObj := models.NewRecipe(db)

	vars := mux.Vars(r)

	recipeID, err := strconv.Atoi(vars["recipe_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You must send an ID to delete a recipe"))
		return
	}

	_, err = recipeObj.DeleteById(nil, int64(recipeID))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not delete the recipe"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateRecipesHandler(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*sqlx.DB)

	recipeObj := models.NewRecipe(db)

	vars := mux.Vars(r)

	recipeID, err := strconv.Atoi(vars["recipe_id"])
	fieldToUpdate := vars["field"]

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ToUpdate := make(map[string]interface{})

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing data to update the recipe"))
		return
	}

	err = json.Unmarshal(body, &ToUpdate)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid json request"))
		return
	}

	switch fieldToUpdate {
	case "name", "serves_for":
		_, err = recipeObj.UpdateByID(nil, ToUpdate, int64(recipeID))
	case "type":
		fmt.Println("that a different case")
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not update"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("Recipe updated"))

}

func GetRecipesHandler(w http.ResponseWriter, r *http.Request) {

	db := r.Context().Value("db").(*sqlx.DB)

	recipeObj := models.NewRecipe(db)

	stringSearch := r.URL.Query()["name"]

	// Two cases here:
	// 1. Get all recipes regardless of name or ID
	// 2. Get recipes that match the search string
	if len(stringSearch) != 0 {
		fullRecipes, RecipesTypes, err := recipeObj.GetFullRecipesByStringSearch(nil, stringSearch[0])

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong"))
			return
		}

		JSONResponse := buildFullRecipeJSONResponse(fullRecipes, RecipesTypes)

		w.WriteHeader(http.StatusOK)
		w.Write(JSONResponse)

	} else {
		fullRecipes, RecipesTypes, err := recipeObj.GetFullRecipes(nil)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Something went wrong"))
			return
		}

		JSONResponse := buildFullRecipeJSONResponse(fullRecipes, RecipesTypes)

		w.WriteHeader(http.StatusOK)
		w.Write(JSONResponse)
	}

}

func GetRecipeByIDHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	recipeID, err := strconv.Atoi(vars["recipe_id"])

	db := r.Context().Value("db").(*sqlx.DB)

	recipeObj := models.NewRecipe(db)

	if err != nil {
		fmt.Println(err)
	}

	recipe, err := recipeObj.GetFullRecipe(nil, int64(recipeID))

	if err != nil {
		fmt.Printf("Something went wrong while fecthing the recipe. Error: %v", err)
	}

	if len(recipe) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Recipe not found"))
		return
	}

	var fullRecipes [][]models.FullRecipeRow

	RecipesTypes := make(map[int64][]string)
	recipeTypes, err := recipeObj.GetRecipeType(nil, recipe[0].ID)

	RecipesTypes[recipe[0].ID] = recipeTypes

	if err != nil {
		fmt.Printf("Error fecthing full recipe. Error: %v", err)
	}

	fullRecipes = append(fullRecipes, recipe)
	JSONResponse := buildFullRecipeJSONResponse(fullRecipes, RecipesTypes)

	w.WriteHeader(http.StatusOK)
	w.Write(JSONResponse)
}

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

	w.WriteHeader(http.StatusOK)
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

	w.WriteHeader(http.StatusOK)
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
