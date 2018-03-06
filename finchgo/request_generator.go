package finchgo

import (
	"encoding/json"
	"fmt"

	"github.com/Pallinder/go-randomdata"
)

func GenerateSignUp() []byte {
	Body := make(map[string]interface{})

	Body["name"] = randomdata.FirstName(randomdata.RandomGender)
	Body["email"] = randomdata.StringNumber(8, "") + randomdata.Email()

	profile := randomdata.GenerateProfile(randomdata.RandomGender)

	Body["password"] = profile.Login.Md5

	JSONBody, _ := json.Marshal(Body)

	return JSONBody
}

func GenerateRecipe() []byte {
	Body := make(map[string]interface{})

	Body["recipe_name"] = randomdata.SillyName()
	Body["type"] = []string{"Lunch", "Dinner", "Breakfast"}
	Body["serves_for"] = "2"

	Body["steps"] = make([]interface{}, 0, 0)

	NumberOfSteps := RandomIntRange(3, 10)

	for i := 0; i < NumberOfSteps; i++ {
		Step := make(map[string]interface{})
		Step["step_id"] = i + 1
		Step["text"] = randomdata.Paragraph()

		Step["step_ingredients"] = make([]interface{}, 0, 0)

		NumberOfIngredients := RandomIntRange(1, 5)

		for j := 0; j < NumberOfIngredients; j++ {

			StepIngredient := make(map[string]interface{})
			StepIngredient["name"] = randomdata.Noun() + randomdata.StringNumber(3, "")
			StepIngredient["amount"] = randomdata.Decimal(1, 10)
			StepIngredient["unit"] = 10

			Step["step_ingredients"] = append(Step["step_ingredients"].([]interface{}), StepIngredient)

		}

		Body["steps"] = append(Body["steps"].([]interface{}), Step)
	}

	JSONBody, err := json.Marshal(Body)

	if err != nil {
		fmt.Println(err)
	}

	return JSONBody
}

func GenerateIngredientList(NumberOfIngredients int) []byte {
	Body := make([]map[string]interface{}, 0, 0)

	for i := 0; i < NumberOfIngredients; i++ {
		IngredientBody := make(map[string]interface{})
		IngredientBody["name"] = randomdata.Noun() + randomdata.StringNumber(3, "")
		IngredientBody["amount"] = randomdata.Decimal(1, 10)
		IngredientBody["unit"] = 10
		Body = append(Body, IngredientBody)
	}

	JSONBody, err := json.Marshal(Body)

	if err != nil {
		fmt.Println(err)
	}

	return JSONBody

}
