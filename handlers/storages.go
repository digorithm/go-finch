package handlers

import (
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/digorithm/meal_planner/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func CreateStorageObj(r *http.Request) *models.ItemInStorage {
	db := r.Context().Value("db").(*sqlx.DB)

	StorageObj := models.NewItemInStorage(db)

	return StorageObj
}

func GetStoragesHandler(w http.ResponseWriter, r *http.Request) {

	StorageObj := CreateStorageObj(r)

	vars := mux.Vars(r)
	HouseID, err := strconv.Atoi(vars["house_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing house id"))
		return
	}

	Storage, err := StorageObj.GetHouseStorage(nil, int64(HouseID))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("House not found"))
		return
	}

	StorageJSON, err := json.Marshal(Storage)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(StorageJSON)
}
