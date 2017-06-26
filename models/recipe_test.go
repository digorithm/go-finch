package models

import (
	"testing"

	_ "github.com/lib/pq"
)

func createVarsForFullRecipe(id int64, name string, t string, s int64, si int64, in string, a float64, un string, te string) FullRecipeRow {
	var fRecipe FullRecipeRow

	fRecipe.ID = id
	fRecipe.Name = name
	fRecipe.Type = t
	fRecipe.ServesFor = s
	fRecipe.StepID = si
	fRecipe.Ingredient = in
	fRecipe.Amount = a
	fRecipe.Unit = un
	fRecipe.Text = te

	return fRecipe
}

func newRecipeForTest(t *testing.T) *Recipe {
	return NewRecipe(newDbForTest(t))
}

func TestGetFullRecipe(t *testing.T) {

	r := newRecipeForTest(t)
	var f1 = createVarsForFullRecipe(1, "Baked Potato", "Lunch/Dinner", 4, 1, "potato", 4, "litre", "peel and cut the potatoes into an inch thick disks")
	var f2 = createVarsForFullRecipe(1, "Baked Potato", "Lunch/Dinner", 4, 2, "milk", 0.25, "cup", "mix the milk and parmesan together")
	var f3 = createVarsForFullRecipe(1, "Baked Potato", "Lunch/Dinner", 4, 2, "parmesan cheese", 1, "grams", "mix the milk and parmesan together")
	var result []FullRecipeRow

	fullRecipe, err := r.GetFullRecipe(nil, 1)
	if err != nil {
		t.Errorf("Getting full recipe should work. Error: %v", err)
	}

	result = append(result, f1, f2, f3)
	i := 0
	for i < len(fullRecipe) {
		if result[i] != fullRecipe[i] {
			t.Errorf("Get FullRecipe failed, got: %d, want: %d", fullRecipe[i], result[i])
		}
		i++
	}
}
