package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/digorithm/meal_planner/libhttp"
	"github.com/digorithm/meal_planner/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func CreateScheduleObj(r *http.Request) *models.Schedule {
	db := r.Context().Value("db").(*sqlx.DB)

	scheduleObj := models.NewSchedule(db)

	return scheduleObj
}

func GetMealTypesHandler(w http.ResponseWriter, r *http.Request) {

	scheduleObj := CreateScheduleObj(r)

	res, err := scheduleObj.GetMeals(nil)

	if err != nil {
		fmt.Printf("getAllUnitsHandler failed: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetDaysHandler(w http.ResponseWriter, r *http.Request) {

	scheduleObj := CreateScheduleObj(r)

	res, err := scheduleObj.GetDays(nil)

	if err != nil {
		fmt.Printf("getAllUnitsHandler failed: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetScheduleHandler(w http.ResponseWriter, r *http.Request) {

	scheduleObj := CreateScheduleObj(r)

	vars := mux.Vars(r)
	houseID, err := strconv.Atoi(vars["house_id"])

	res, err := scheduleObj.GetSchedule(nil, int64(houseID))

	if err != nil {
		fmt.Printf("%v", err)
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func ModifyScheduleHandler(w http.ResponseWriter, r *http.Request) {

	scheduleObj := CreateScheduleObj(r)

	vars := mux.Vars(r)
	houseID, err := strconv.Atoi(vars["house_id"])

	newRecipe, err := ioutil.ReadAll(r.Body)

	res, err := scheduleObj.ModifySchedule(nil, int64(houseID), newRecipe)

	if err != nil {
		fmt.Printf("%v", err)
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func NewFullScheduleHandler(w http.ResponseWriter, r *http.Request) {

}
