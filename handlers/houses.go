package handlers

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/digorithm/meal_planner/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func CreateHouseObj(r *http.Request) *models.House {
	db := r.Context().Value("db").(*sqlx.DB)

	houseObj := models.NewHouse(db)

	return houseObj
}

func GetHouseHandler(w http.ResponseWriter, r *http.Request) {

	HouseObj := CreateHouseObj(r)

	vars := mux.Vars(r)
	HouseID, err := strconv.Atoi(vars["house_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing house id"))
	}

	h, err := HouseObj.GetFullHouseInformation(nil, int64(HouseID))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("House not found"))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(h)
}

func PostHouseHandler(w http.ResponseWriter, r *http.Request) {

	HouseObj := CreateHouseObj(r)
	MemberOf := models.NewMember(r.Context().Value("db").(*sqlx.DB))

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing data in house request"))
		return
	}

	BodyStruct := make(map[string]interface{})

	_ = json.Unmarshal(body, &BodyStruct)

	HouseName := BodyStruct["name"].(string)
	HouseGroceryDay := BodyStruct["grocery_day"].(string)
	HouseHouseholdNumber := BodyStruct["household_number"].(float64)

	UserID := int64(BodyStruct["user_id"].(float64))

	h, err := HouseObj.CreateHouse(nil, HouseName, HouseGroceryDay, int64(HouseHouseholdNumber))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	MemberOf.AddOwner(nil, h.ID, UserID)

	ReturnedHouse, err := HouseObj.GetFullHouseInformation(nil, h.ID)

	w.WriteHeader(http.StatusOK)
	w.Write(ReturnedHouse)
}
