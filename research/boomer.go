package main

/*
import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"strings"

	"encoding/json"

	randomdata "github.com/Pallinder/go-randomdata"
	"github.com/myzhan/boomer"
)

var BaseURL = "http://localhost:8888"

func NewHTTPClient() *http.Client {
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	return &http.Client{Transport: tr}
}

type Ingredient struct {
	name   string
	amount float64
	unit   float64
}

func requestSignUp() ([]byte, error) {
	SignUpBody := GenerateSignUp()
	URL := strings.Join([]string{BaseURL, "/users"}, "")

	client := NewHTTPClient()

	resp, err := client.Post(URL, "application/json", bytes.NewBuffer(SignUpBody))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, err
}

func requestAddRecipe(AuthorID int64) ([]byte, error) {

	RecipeRequest := GenerateRecipe()

	ResourceURL := fmt.Sprintf("/recipes?author=%v", AuthorID)

	URL := strings.Join([]string{BaseURL, ResourceURL}, "")

	client := NewHTTPClient()
	resp, err := client.Post(URL, "application/json", bytes.NewBuffer(RecipeRequest))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, err
}

func requestSearchRecipeByName(RecipeName string) ([]byte, error) {

	ResourceURL := fmt.Sprintf("/recipes?name=%v", RecipeName)

	URL := strings.Join([]string{BaseURL, ResourceURL}, "")

	client := NewHTTPClient()
	resp, err := client.Get(URL)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, err
}

func requestCreateHouse(UserID int64) ([]byte, error) {
	ResourceURL := fmt.Sprintf("/houses")
	URL := strings.Join([]string{BaseURL, ResourceURL}, "")

	Request := make(map[string]interface{})
	Request["name"] = randomdata.SillyName()
	Request["user_id"] = float64(UserID)
	Request["grocery_day"] = "Friday"
	Request["household_number"] = 2.0

	JSONRequest, _ := json.Marshal(Request)

	client := NewHTTPClient()
	resp, err := client.Post(URL, "application/json", bytes.NewBuffer(JSONRequest))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, err
}

func requestCreateSchedule(HouseID int64) ([]byte, error) {
	ResourceURL := fmt.Sprintf("/schedules/create/%v", HouseID)
	URL := strings.Join([]string{BaseURL, ResourceURL}, "")

	client := NewHTTPClient()
	resp, err := client.Post(URL, "application/json", nil)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, err
}

func requestScheduleRecipes(HouseID int64, RecipesIDs []int64) {
	ResourceURL := fmt.Sprintf("/schedules/%v", HouseID)
	URL := strings.Join([]string{BaseURL, ResourceURL}, "")

	for _, id := range RecipesIDs {

		Request := make(map[string]interface{})
		Request["recipe_id"] = id
		Request["type"] = RandomIntRange(1, 4)
		Request["day"] = RandomIntRange(1, 7)

		JSONRequest, _ := json.Marshal(Request)

		client := NewHTTPClient()
		_, err := client.Post(URL, "application/json", bytes.NewBuffer(JSONRequest))

		if err != nil {
			fmt.Println(err)
		}
	}

}

func requestAddIngredientListToStorage(HouseID int64, NumberOFIngredients int) ([]byte, error) {

	ResourceURL := fmt.Sprintf("/storages/%v", HouseID)
	URL := strings.Join([]string{BaseURL, ResourceURL}, "")

	RequestBody := GenerateIngredientList(NumberOFIngredients)

	client := NewHTTPClient()
	resp, err := client.Post(URL, "application/json", bytes.NewBuffer(RequestBody))

	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, err
}

// requestUpdateStorage takes a [POST] /storages/{house_id}/ response and modify it to update the same storage
func requestUpdateStorage(HouseID int64, Storage []byte) ([]byte, error) {

	StorageStruct := JSONListToMap(Storage)

	for _, Ingredient := range StorageStruct {
		Ingredient["amount"] = Ingredient["amount"].(float64) / 2.0
		// Delete unecessary data
		delete(Ingredient, "unit")
		delete(Ingredient, "ingredient_id")
		delete(Ingredient, "house_id")

		// Change key name from ingredient_name to name and from unit_id to unit
		Ingredient["name"] = Ingredient["ingredient_name"]
		Ingredient["unit"] = Ingredient["unit_id"]
		delete(Ingredient, "ingredient_name")
		delete(Ingredient, "unit_id")
	}

	NewStorageUpdateRequest, err := json.Marshal(StorageStruct)

	if err != nil {
		fmt.Println(err)
	}
	ResourceURL := fmt.Sprintf("/storages/%v", HouseID)
	URL := strings.Join([]string{BaseURL, ResourceURL}, "")

	client := NewHTTPClient()
	resp, err := client.Post(URL, "application/json", bytes.NewBuffer(NewStorageUpdateRequest))

	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, err
}

// DefaultUserFlow executes a common user flow:
// - Signup, add a few recipes, modifiy a few recipes,
//	 add items to storage, remove items from storage, invite someone to house
//	 search for a few recipes by string and schedule a week
func DefaultUserFlow() {
	// Random number in the range 300 - 1500
	// This is the wait time between user flows, otherwise it would be concurrent
	// Thus super fast and would not be representative of a real user-flow behavior
	r := RandomIntRange(10, 20)
	time.Sleep(time.Duration(r) * time.Second)

	startTime := boomer.Now()

	// Make signup request
	SignUpResp, err := requestSignUp()
	SignUpStruct := JSONToMap(SignUpResp)
	UserID := int64(SignUpStruct["ID"].(float64))

	if err != nil {
		boomer.Events.Publish("request_failure", "http", "Signup", 0.0, err.Error())
	} else {

		boomer.Events.Publish("request_success", "http", "Signup", float64(boomer.Now()-startTime), int64(1))
	}

	// Make create house request
	CreateHouseResp, err := requestCreateHouse(UserID)
	HouseID := int64(JSONToMap(CreateHouseResp)["id"].(float64))

	if err != nil {
		boomer.Events.Publish("request_failure", "http", "CreateHouse", 0.0, err.Error())
	} else {

		boomer.Events.Publish("request_success", "http", "CreateHouse", float64(boomer.Now()-startTime), int64(1))
	}

	// Add 3-8 recipes
	NumberOfRecipes := RandomIntRange(3, 8)

	AddedRecipes := make([]int64, 0, 0)

	for i := 0; i < NumberOfRecipes; i++ {
		// Add a recipe every 10-30 seconds
		r = RandomIntRange(10, 30)
		time.Sleep(time.Duration(r) * time.Second)

		AddRecipeResp, err := requestAddRecipe(UserID)

		RecipeID := JSONListToMap(AddRecipeResp)[0]["id"]

		AddedRecipes = append(AddedRecipes, int64(RecipeID.(float64)))

		if err != nil {
			boomer.Events.Publish("request_failure", "http", "AddRecipe", 0.0, err.Error())
		} else {

			boomer.Events.Publish("request_success", "http", "AddRecipe", float64(boomer.Now()-startTime), int64(1))
		}

		_ = JSONListToMap(AddRecipeResp)
	}

	// Search for 5-20 recipes
	NumberOfSearches := RandomIntRange(5, 20)

	for i := 0; i < NumberOfSearches; i++ {
		// make search request every 5-10 seconds
		r = RandomIntRange(5, 10)
		time.Sleep(time.Duration(r) * time.Second)
		SearchString := randomdata.SillyName()
		_, err := requestSearchRecipeByName(SearchString)
		if err != nil {
			boomer.Events.Publish("request_failure", "http", "SearchRecipe", 0.0, err.Error())
		} else {

			boomer.Events.Publish("request_success", "http", "SearchRecipe", float64(boomer.Now()-startTime), int64(1))
		}
	}

	_, err = requestCreateSchedule(HouseID)
	requestScheduleRecipes(HouseID, AddedRecipes)

	endTime := boomer.Now()

	if err != nil {
		boomer.Events.Publish("request_failure", "http", "CreateSchedule", 0.0, err.Error())
	} else {

		boomer.Events.Publish("request_success", "http", "CreateSchedule", float64(endTime-startTime), int64(1))
	}

	r = RandomIntRange(10, 50)
	Storage, err := requestAddIngredientListToStorage(HouseID, r)

	requestUpdateStorage(HouseID, Storage)

	if err != nil {
		boomer.Events.Publish("request_failure", "http", "AddToStorage", 0.0, err.Error())
	} else {

		boomer.Events.Publish("request_success", "http", "AddToStorage", float64(endTime-startTime), int64(1))
	}

}

func main() {
	task := &boomer.Task{
		Name:   "Default",
		Weight: 1,
		Fn:     DefaultUserFlow,
	}

	boomer.Run(task)

}
*/
