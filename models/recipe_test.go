package models

import (
	"testing"

	_ "github.com/lib/pq"
)

func newRecipeForTest(t *testing.T) *Recipe {
	return NewRecipe(newDbForTest(t))
}

/*func TestGetFullRecipe(t *testing.T) {
	r := newRecipeForTest(t)

	FullRecipe, err := r.GetFullRecipe(nil, 1)

	if err != nil {
		t.Errorf("Getting full recipe should work. Error: %v", err)
	}
	fmt.Println(FullRecipe)
}*/
