package models

import (
	"errors"
	"fmt"
	"reflect"

	"database/sql"

	"encoding/json"

	"github.com/jmoiron/sqlx"
)

type Schedule struct {
	Base
}

func NewSchedule(db *sqlx.DB) *Schedule {

	schedule := &Schedule{}
	schedule.db = db
	schedule.table = "schedule"
	schedule.hasID = false

	return schedule
}

func (s *Schedule) GetSchedule(tx *sqlx.Tx, houseID int64) ([]byte, error) {

	rawSchedule, err := s.GetHouseSchedule(tx, houseID)

	if err != nil {
		fmt.Printf("GetSchedule failed: %v", err)
	}

	return s.GetScheduleFullRecipes(tx, rawSchedule)
}

func (s *Schedule) GetDays(tx *sqlx.Tx) ([]byte, error) {

	query := "SELECT W.ID, W.DAY FROM WEEKDAY W"

	data, err := s.GetCompoundModelWithoutID(tx, query)

	if err != nil {
		fmt.Printf("getAllUnits failed: %v", err)
	}

	return AllFixedTablesJSON(data), err
}

func (s *Schedule) GetMeals(tx *sqlx.Tx) ([]byte, error) {

	query := "SELECT M.ID, M.TYPE FROM MEAL_TYPE M"

	data, err := s.GetCompoundModelWithoutID(tx, query)

	if err != nil {
		fmt.Printf("getAllUnits failed: %v", err)
	}

	return AllFixedTablesJSON(data), err
}

func (s *Schedule) GetHouseSchedule(tx *sqlx.Tx, houseID int64) ([]HouseScheduleRow, error) {

	query := "SELECT W.DAY, T.TYPE, R.NAME, R.ID FROM RECIPE R, WEEKDAY W, MEAL_TYPE T, SCHEDULE S WHERE S.HOUSE_ID = $1 AND S.WEEK_ID = W.ID AND S.TYPE_ID = T.ID AND S.RECIPE_ID = R.ID ORDER BY WEEK_ID, TYPE_ID"

	data, err := s.GetCompoundModel(tx, query, houseID)

	schedule := createHouseScheduleRows(data)

	if err != nil {
		fmt.Printf("GetHouseSchedule failed: %v", err)
	}

	return schedule, err
}

func (s *Schedule) ModifySchedule(tx *sqlx.Tx, houseID int64, newRecipe []byte) ([]byte, error) {

	info := make(map[string]interface{})
	err := json.Unmarshal(newRecipe, &info)

	if err != nil {
		fmt.Printf("ModifySchedule failed: %v", err)
	}

	wDay := info["day"].(float64)
	tOf := info["type"].(float64)
	rID := info["recipe_id"].(float64)

	_, err = s.UpdateSchedule(tx, houseID, int64(wDay), int64(tOf), int64(rID))
	if err != nil {
		fmt.Printf("GetNewSchedule failed: %v", err)
	}

	houseScheduleRaw, _ := s.GetHouseSchedule(tx, houseID)

	return s.GetScheduleFullRecipes(tx, houseScheduleRaw)

}

func (s *Schedule) GetScheduleFullRecipes(tx *sqlx.Tx, HouseScheduleRaw []HouseScheduleRow) ([]byte, error) {

	r := NewRecipe(s.db)
	finalSchedule := make(map[string]interface{})

	for _, schedule := range HouseScheduleRaw {

		recipeID := schedule.RID
		rType, _ := r.GetRecipeType(tx, recipeID)
		recipe, _ := r.GetFullRecipe(tx, recipeID)

		finalRecipe := BuildScheduleRecipeJSON(rType, recipe)

		if len(finalRecipe) > 0 {
			key := fmt.Sprintf("%v_%v", schedule.Week, schedule.Type)
			finalSchedule[key] = finalRecipe
		}
	}

	return json.MarshalIndent(finalSchedule, "", "    ")

}

func BuildScheduleRecipeJSON(rType []string, recipe []FullRecipeRow) map[string]interface{} {

	finalRecipe := make(map[string]interface{})

	if len(recipe) == 0 {
		return finalRecipe
	}

	recipeID := recipe[0].ID
	recipeName := recipe[0].Name
	servesFor := recipe[0].ServesFor
	recipeType := rType

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
	finalRecipe["type"] = recipeType
	finalRecipe["serves_for"] = servesFor
	finalRecipe["steps"] = steps

	return finalRecipe

}

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

func (s *Schedule) UpdateSchedule(tx *sqlx.Tx, hID, wDay, tOf, rID int64) (sql.Result, error) {

	query := fmt.Sprintf("UPDATE SCHEDULE SET RECIPE_ID = %v WHERE HOUSE_ID = %v AND WEEK_ID = %v AND TYPE_ID = %v", rID, hID, wDay, tOf)

	tx, wrapInSingleTransaction, err := s.newTransactionIfNeeded(tx)

	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty")
	}
	if err != nil {
		return nil, err
	}

	sqlResult, err := tx.Exec(query)

	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	if err != nil {
		fmt.Printf("Error in UpdateSchedule: %v", err)
	}

	return sqlResult, err

}

// CreateHouseSchedule can be called from anywhere, it checks for all possible cases
func (s *Schedule) CreateHouseSchedule(tx *sqlx.Tx, houseID int64) bool {

	res := s.GetCurrentScheduleRows(tx, houseID)

	if len(res) == 0 {

		s.CreateFullSchedule(tx, houseID)

	} else {

		s.InsertMissingSchedule(tx, houseID, res)
	}

	schedule := s.GetCurrentScheduleRows(tx, houseID)

	if len(schedule) == 28 {
		return true
	}
	return false
}

// CreateFullSchedule is called when there is no instance of house_id in schedule table
func (s *Schedule) CreateFullSchedule(tx *sqlx.Tx, houseID int64) {

	i := 1
	data := make(map[string]interface{})

	for i < 8 {
		data["house_id"] = houseID
		data["week_id"] = i
		j := 1

		for j < 5 {
			data["type_id"] = j
			_, err := s.InsertIntoMultiKeyTable(tx, data)

			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			j++
		}
		i++
	}
}

func (s *Schedule) DeleteSchedule(tx *sqlx.Tx, houseID int64) bool {

	where := fmt.Sprintf("house_id = %v", houseID)
	sRes, err := s.DeleteFromTable(tx, where)

	if err != nil {
		fmt.Printf("Error in DeleteSchedule: %v", err)
	}

	res, _ := sRes.RowsAffected()
	if res == 28 {
		return true
	}

	return false

}

func (s *Schedule) InsertMissingSchedule(tx *sqlx.Tx, houseID int64, schedule []ScheduleRow) {

	var k int
	var i int64 = 1
	var j int64 = 1
	data := make(map[string]interface{})
	data["house_id"] = houseID

	for i < 8 {
		j = 1
		for j < 5 {

			item := schedule[k]
			w := item.WeekID
			t := item.TypeID

			if (w != i) || (t != j) {
				data["week_id"] = i
				data["type_id"] = j
				s.InsertIntoMultiKeyTable(tx, data)
			} else {
				if k < len(schedule)-1 {
					k++
				}
			}
			j++
		}
		i++
	}

}

func (s *Schedule) GetCurrentScheduleRows(tx *sqlx.Tx, houseID int64) []ScheduleRow {

	query := "SELECT S.WEEK_ID, S.TYPE_ID FROM SCHEDULE S WHERE HOUSE_ID = $1 ORDER BY WEEK_ID, TYPE_ID"

	data, err := s.GetCompoundModel(tx, query, houseID)
	if err != nil {
		fmt.Printf("%v", err)
	}

	schedule := ExistingScheduleRows(data)

	return schedule
}

func ExistingScheduleRows(data []interface{}) []ScheduleRow {
	var schedule []ScheduleRow
	var s ScheduleRow

	for i := 0; i < len(data); i++ {

		v := reflect.ValueOf(data[i])

		s.WeekID = v.Index(0).Interface().(int64)
		s.TypeID = v.Index(1).Interface().(int64)

		schedule = append(schedule, s)
	}

	return schedule
}
