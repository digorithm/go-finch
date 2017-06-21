package models

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func newIngredientForTest(t *testing.T) *Ingredient {
	return NewIngredient(newDbForTest(t))
}

func deleteTestIngredient(t *testing.T, i *Ingredient, id int64) {
	_, err := i.DeleteById(nil, id)
	if err != nil {
		t.Fatal("Something went wrong with test ingredient deletion. Error: %v", err)
	}
}

func TestGetIngredientById(t *testing.T) {
	i := newIngredientForTest(t)

	test_ingredient := make(map[string]interface{})
	test_ingredient["name"] = "beans"
	test_ingredient["carb_per_100g"] = 20.0
	test_ingredient["protein_per_100g"] = 30.0
	test_ingredient["fat_per_100g"] = 40.0
	test_ingredient["fiber_per_100g"] = 50.0
	test_ingredient["calories_per_100g"] = 60.0

	sqlResult, err := i.InsertIntoTable(nil, test_ingredient)

	if err != nil {
		fmt.Errorf("%v", err)
	}

	test_ingredient_id, _ := sqlResult.LastInsertId()
	returned_ingredient, err := i.GetById(nil, test_ingredient_id)

	defer deleteTestIngredient(t, i, test_ingredient_id)

	if err != nil {
		t.Errorf("Getting ingredient by id should work. Error: %v", err)
	}
	if returned_ingredient.Name != test_ingredient["name"] ||
		returned_ingredient.CarbPer100G != test_ingredient["carb_per_100g"] ||
		returned_ingredient.ProteinPer100G != test_ingredient["protein_per_100g"] ||
		returned_ingredient.FatPer100G != test_ingredient["fat_per_100g"] ||
		returned_ingredient.FiberPer100G != test_ingredient["fiber_per_100g"] ||
		returned_ingredient.CaloriesPer100G != test_ingredient["calories_per_100g"] {

		t.Errorf("Returned ingredient and test test ingredient are different")
	}
}
