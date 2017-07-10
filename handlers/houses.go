package handlers

import (
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
