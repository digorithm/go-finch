package models

import (
	_ "github.com/lib/pq"
	"testing"
)

func newIngredientForTest(t *testing.T) *Ingredient {
	return NewIngredient(newDbForTest(t))
}

func createTestIngredient(t *testing.T, name string) (*IngredientRow, *Ingredient) {
	i := newIngredientForTest(t)

	test_ingredient := make(map[string]interface{})
	test_ingredient["name"] = name
	test_ingredient["carb_per_100g"] = 20.0
	test_ingredient["protein_per_100g"] = 30.0
	test_ingredient["fat_per_100g"] = 40.0
	test_ingredient["fiber_per_100g"] = 50.0
	test_ingredient["calories_per_100g"] = 60.0

	sqlResult, err := i.InsertIntoTable(nil, test_ingredient)

	if err != nil {
		t.Errorf("Creation of test ingredient should work. Error: %v", err)
	}

	test_ingredient_id, _ := sqlResult.LastInsertId()
	returned_ingredient, err := i.GetById(nil, test_ingredient_id)

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

	return returned_ingredient, i
}

func deleteTestIngredient(t *testing.T, i *Ingredient, id int64) {
	_, err := i.DeleteById(nil, id)
	if err != nil {
		t.Fatal("Something went wrong with test ingredient deletion. Error: %v", err)
	}
}

func TestGetIngredientById(t *testing.T) {

	returned_ingredient, ingredientObj := createTestIngredient(t, "beans")
	defer deleteTestIngredient(t, ingredientObj, returned_ingredient.ID)
}

func TestGetIngredientByName(t *testing.T) {
	//TODO: in the future, test with "bean", "BEANS", "BEAN" and implement this feature
	expected_normal_name := "beans"

	returned_ingredient, ingredientObj := createTestIngredient(t, expected_normal_name)
	defer deleteTestIngredient(t, ingredientObj, returned_ingredient.ID)

	// test normal case
	returned_ingredient_row, err := ingredientObj.GetByName(nil, expected_normal_name)

	if err != nil {
		t.Errorf("Get ingredient by name should work. Error: %v", err)
	}

	if returned_ingredient_row.Name != expected_normal_name {
		t.Errorf("Ingredients have different names")
		t.Errorf("Expected: %v", expected_normal_name)
		t.Errorf("Actual: %v", returned_ingredient_row.Name)
	}
}
