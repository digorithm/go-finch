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
		fmt.Printf("%v", err)
		libhttp.HandleErrorJson(w, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetDaysHandler(w http.ResponseWriter, r *http.Request) {

	scheduleObj := CreateScheduleObj(r)

	res, err := scheduleObj.GetDays(nil)

	if err != nil {
		fmt.Printf("%v", err)
		libhttp.HandleErrorJson(w, err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateScheduleHandler(w http.ResponseWriter, r *http.Request) {

	Finch.ArtificialBlockingPoint("k4")

	scheduleObj := CreateScheduleObj(r)

	vars := mux.Vars(r)
	houseID, err := strconv.Atoi(vars["house_id"])

	if err != nil {
		fmt.Printf("%v", err)
		libhttp.HandleErrorJson(w, err)
	}

	created := scheduleObj.CreateHouseSchedule(nil, int64(houseID))

	var response []byte

	if created {
		response = []byte(`{"message":"Schedule for the house successfully created"}`)
		w.WriteHeader(http.StatusCreated)
	} else {
		response = []byte(`{"message":"Failed to create schedule"}`)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(response)

}

func DeleteScheduleHandler(w http.ResponseWriter, r *http.Request) {

	scheduleObj := CreateScheduleObj(r)

	vars := mux.Vars(r)
	houseID, err := strconv.Atoi(vars["house_id"])

	if err != nil {
		fmt.Printf("%v", err)
		libhttp.HandleErrorJson(w, err)
	}

	deleted := scheduleObj.DeleteSchedule(nil, int64(houseID))

	var response []byte

	if deleted {
		response = []byte(`{"message":"Schedule for the house successfully deleted"}`)
		w.WriteHeader(http.StatusOK)
	} else {
		response = []byte(`{"message":"Failed to delete schedule"}`)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(response)

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
