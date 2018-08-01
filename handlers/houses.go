package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/digorithm/meal_planner/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func CreateHouseObj(r *http.Request) *models.House {
	db := r.Context().Value("db").(*sqlx.DB)

	houseObj := models.NewHouse(db)

	return houseObj
}

func CreateMemberObj(r *http.Request) *models.Member {
	db := r.Context().Value("db").(*sqlx.DB)

	memberObj := models.NewMember(db)

	return memberObj
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

	Finch.ArtificialBlockingPoint("k3")

	HouseObj := CreateHouseObj(r)
	MemberOf := CreateMemberObj(r)

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

func DeleteHouseHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	HouseID, err := strconv.Atoi(vars["house_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing house id"))
		return
	}

	HouseObj := CreateHouseObj(r)

	_, err = HouseObj.DeleteById(nil, int64(HouseID))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateHouseHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	HouseID, err := strconv.Atoi(vars["house_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing house id"))
		return
	}

	HouseObj := CreateHouseObj(r)

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing data in house request"))
		return
	}

	BodyStruct := make(map[string]interface{})

	_ = json.Unmarshal(body, &BodyStruct)

	ToUpdate := make(map[string]interface{})

	if BodyStruct["name"] != nil {
		ToUpdate["name"] = BodyStruct["name"]
	}

	if BodyStruct["grocery_day"] != nil {
		ToUpdate["grocery_day"] = BodyStruct["grocery_day"]
	}

	if BodyStruct["household_number"] != nil {
		ToUpdate["household_number"] = BodyStruct["household_number"]
	}

	_, err = HouseObj.UpdateByID(nil, ToUpdate, int64(HouseID))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	UpdatedHouse, err := HouseObj.GetFullHouseInformation(nil, int64(HouseID))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(UpdatedHouse)
}
