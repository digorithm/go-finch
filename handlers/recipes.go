package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/digorithm/meal_planner/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

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

	for _, recipe := range recipes {
		fullRecipe, err := recipeObj.GetFullRecipe(nil, recipe.ID)

		if err != nil {
			fmt.Printf("Error fecthing full recipe. Error: %v", err)
		}
		fullRecipes = append(fullRecipes, fullRecipe)
	}

	JSONResponse := buildFullRecipeJSONResponse(fullRecipes)

	fmt.Println(string(JSONResponse))

	w.Write([]byte("hello"))
}

func buildFullRecipeJSONResponse(recipes [][]models.FullRecipeRow) []byte {
	for _, recipe := range recipes {
		fmt.Println(recipe)
	}

	return []byte("In development")
}
