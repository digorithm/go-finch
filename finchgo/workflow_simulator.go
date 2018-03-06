package finchgo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"strings"

	"encoding/json"

	randomdata "github.com/Pallinder/go-randomdata"
)

func NewWorkflowSimulator(BaseURL string) *WorkflowSimulator {
	WorkflowSimulator := &WorkflowSimulator{}
	WorkflowSimulator.BaseURL = BaseURL
	WorkflowSimulator.StopChan = make(chan bool)

	return WorkflowSimulator
}

type WorkflowSimulator struct {
	BaseURL  string
	StopChan chan bool
}

func NewHTTPClient() http.Client {
	tr := &http.Transport{
		DisableKeepAlives: true,
	}
	return http.Client{Transport: tr}
}

type Ingredient struct {
	name   string
	amount float64
	unit   float64
}

type ScheduleSettings struct {
	Low, Med, High int
	Days           map[string][]map[string]string
}

func (w *WorkflowSimulator) requestSignUp() ([]byte, error) {
	SignUpBody := GenerateSignUp()

	URL := strings.Join([]string{w.BaseURL, "/users"}, "")

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

func (w *WorkflowSimulator) requestAddRecipe(AuthorID int64) ([]byte, error) {

	RecipeRequest := GenerateRecipe()

	ResourceURL := fmt.Sprintf("/recipes?author=%v", AuthorID)

	URL := strings.Join([]string{w.BaseURL, ResourceURL}, "")

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

func (w *WorkflowSimulator) requestSearchRecipeByName(RecipeName string) ([]byte, error) {

	ResourceURL := fmt.Sprintf("/recipes?name=%v", RecipeName)

	URL := strings.Join([]string{w.BaseURL, ResourceURL}, "")

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

func (w *WorkflowSimulator) requestCreateHouse(UserID int64) ([]byte, error) {
	ResourceURL := fmt.Sprintf("/houses")
	URL := strings.Join([]string{w.BaseURL, ResourceURL}, "")

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

func (w *WorkflowSimulator) requestCreateSchedule(HouseID int64) ([]byte, error) {
	ResourceURL := fmt.Sprintf("/schedules/create/%v", HouseID)
	URL := strings.Join([]string{w.BaseURL, ResourceURL}, "")

	client := NewHTTPClient()
	resp, err := client.Post(URL, "application/json", nil)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, err
}

func (w *WorkflowSimulator) requestScheduleRecipes(HouseID int64, RecipesIDs []int64) {
	ResourceURL := fmt.Sprintf("/schedules/%v", HouseID)
	URL := strings.Join([]string{w.BaseURL, ResourceURL}, "")

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

func (w *WorkflowSimulator) requestAddIngredientListToStorage(HouseID int64, NumberOFIngredients int) ([]byte, error) {

	ResourceURL := fmt.Sprintf("/storages/%v", HouseID)
	URL := strings.Join([]string{w.BaseURL, ResourceURL}, "")

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
func (w *WorkflowSimulator) requestUpdateStorage(HouseID int64, Storage []byte) ([]byte, error) {

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
	URL := strings.Join([]string{w.BaseURL, ResourceURL}, "")

	client := NewHTTPClient()
	resp, err := client.Post(URL, "application/json", bytes.NewBuffer(NewStorageUpdateRequest))

	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return body, err
}

func (w *WorkflowSimulator) SignupTask(semaphore chan bool, ConcurrentUsers int) {

	// Random number in the range 300 - 1500
	// This is the wait time between user flows, otherwise it would be concurrent
	// Thus super fast and would not be representative of a real user-flow behavior
	r := RandomIntRange(1, 10)
	time.Sleep(time.Duration(r) * time.Second)

	// Make a few signup request
	SignUpResp, err := w.requestSignUp()
	SignUpStruct := JSONToMap(SignUpResp)
	UserID := int64(SignUpStruct["ID"].(float64))
	if err != nil {
		fmt.Println(err)
	}

	// Make create house request
	CreateHouseResp, err := w.requestCreateHouse(UserID)
	HouseID := int64(JSONToMap(CreateHouseResp)["id"].(float64))

	if err != nil {

		fmt.Println(err)
	}
	// Add 3-8 recipes
	NumberOfRecipes := RandomIntRange(1, 2)

	AddedRecipes := make([]int64, 0, 0)

	for i := 0; i < NumberOfRecipes; i++ {
		// Add a recipe every 10-30 seconds
		r = RandomIntRange(10, 30)
		time.Sleep(time.Duration(r) * time.Second)

		AddRecipeResp, err := w.requestAddRecipe(UserID)

		RecipeID := JSONListToMap(AddRecipeResp)[0]["id"]

		AddedRecipes = append(AddedRecipes, int64(RecipeID.(float64)))

		if err != nil {

			fmt.Println(err)
		}

		_ = JSONListToMap(AddRecipeResp)
	}

	// Search for 50-1000 recipes
	NumberOfSearches := RandomIntRange(5, 20)

	for i := 0; i < NumberOfSearches; i++ {
		// make search request every 5-10 seconds
		r = RandomIntRange(5, 10)
		time.Sleep(time.Duration(r) * time.Second)
		SearchString := randomdata.SillyName()
		_, err := w.requestSearchRecipeByName(SearchString)
		if err != nil {

			fmt.Println(err)
		}
	}

	_, err = w.requestCreateSchedule(HouseID)
	w.requestScheduleRecipes(HouseID, AddedRecipes)

	if err != nil {

		fmt.Println(err)
	}

	r = RandomIntRange(10, 50)
	Storage, err := w.requestAddIngredientListToStorage(HouseID, r)

	w.requestUpdateStorage(HouseID, Storage)

	if err != nil {
		fmt.Println(err)
	}

	semaphore <- true
}

func (w *WorkflowSimulator) getWeekday() string {
	now := time.Now()

	return now.Weekday().String()
}

func (w *WorkflowSimulator) getDayPeriod() string {

	hour := time.Now().Hour()

	if hour < 12 {
		return "Morning"
	} else if hour >= 12 && hour < 17 {
		return "Afternoon"
	}
	return "Night"

}

func (w *WorkflowSimulator) readSettings(settingsPath string) ScheduleSettings {
	file, _ := os.Open(settingsPath)
	decoder := json.NewDecoder(file)
	ScheduleSettings := ScheduleSettings{}
	err := decoder.Decode(&ScheduleSettings)
	if err != nil {
		fmt.Println("error:", err)
	}
	return ScheduleSettings
}

func (w *WorkflowSimulator) getConcurrentUsers(schedule ScheduleSettings) int {

	// Default value just in case
	ConcurrentUsers := RandomIntRange(10, 50)
	weekday := w.getWeekday()
	currentPeriod := w.getDayPeriod()

	for day := range schedule.Days {
		if day == weekday {
			for _, period := range schedule.Days[day] {

				for j, load := range period {

					if j == currentPeriod {
						switch load {
						case "low":
							return schedule.Low
						case "med":
							return schedule.Med
						case "high":
							return schedule.High
						}

					}
				}
			}
		}
	}

	return ConcurrentUsers
}

func (w *WorkflowSimulator) monitorSemaphore(ticker *time.Ticker, semaphore chan bool, concurrentUsers int) {
	for {
		select {
		case <-ticker.C:
			fmt.Printf("%v users working \n", concurrentUsers-len(semaphore))
		case <-w.StopChan:
			fmt.Println("Shutting down the simulation")
			w.StopChan <- true
			return
		}
	}
}

func (w *WorkflowSimulator) userSpawner(semaphore chan bool, concurrentUsers int) {
	fmt.Printf("Spawning %v concurrent users \n", concurrentUsers)

	SpawnerClock := time.NewTicker(300 * time.Millisecond)
	for {
		select {
		case <-w.StopChan:
			fmt.Println("No users being spawned")
			return
		case <-SpawnerClock.C:
			<-semaphore
			go w.SignupTask(semaphore, concurrentUsers)
		}
	}
}

func (w *WorkflowSimulator) monitorUpdate(wg *sync.WaitGroup, settings ScheduleSettings, concurrentUsers int) {
	ticker := time.NewTicker(45 * time.Second)
	for {
		select {
		case <-ticker.C:
			newConcurrentUsers := w.getConcurrentUsers(settings)
			// if they are different, run wg.Done() to kill previous goroutine
			if newConcurrentUsers != concurrentUsers {
				fmt.Printf("Killing previous goroutine. New number of concurrent users:: %v", newConcurrentUsers)
				wg.Done()
				return
			}
		}
	}
}

func (w *WorkflowSimulator) Run(ConcurrentUsers int) {

	for {

		// We are controlling concurrency by using the semaphore pattern with channels
		// This means we will only spawn <ConcurrentUsers> users at once.
		// Once a user is finished with the workflow, we spawn another one
		Semaphore := make(chan bool, ConcurrentUsers)

		// Populate the semaphore channel with the right amount of max concurrent users
		for i := 0; i < ConcurrentUsers; i++ {
			Semaphore <- true
		}

		var wg sync.WaitGroup

		wg.Add(1)

		// This is to monitor the Semaphore channel, ticker is the frequency to monitor it
		ticker := time.NewTicker(5 * time.Second)
		go w.monitorSemaphore(ticker, Semaphore, ConcurrentUsers)

		// Spin the user spawner
		go w.userSpawner(Semaphore, ConcurrentUsers)

		wg.Wait()
		ticker.Stop()

		//FirstTime = false
	}
}

func (w *WorkflowSimulator) Stop() {
	w.StopChan <- true
}
