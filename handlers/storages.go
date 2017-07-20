package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/digorithm/meal_planner/libhttp"
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

func PostStoragesHandler(w http.ResponseWriter, r *http.Request) {

	StorageObj := CreateStorageObj(r)

	vars := mux.Vars(r)
	HouseID, err := strconv.Atoi(vars["house_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing house id"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = StorageObj.AddIngredientList(body, int64(HouseID))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

	w.WriteHeader(http.StatusCreated)
	w.Write(StorageJSON)
}

func DeleteHouseStorage(w http.ResponseWriter, r *http.Request) {

	storageObj := CreateStorageObj(r)

	vars := mux.Vars(r)
	houseID, err := strconv.Atoi(vars["house_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing house id"))
		return
	}

	err = storageObj.DeleteStorage(nil, int64(houseID))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed To Delete House Storage"))
		return
	}

	w.WriteHeader(http.StatusOK)

}

func DeleteFromStorage(w http.ResponseWriter, r *http.Request) {

	storageObj := CreateStorageObj(r)

	vars := mux.Vars(r)
	houseID, err := strconv.Atoi(vars["house_id"])

	items, err := ioutil.ReadAll(r.Body)

	resJSON, err := storageObj.DeleteStorageItems(nil, items, int64(houseID))

	if err != nil {
		fmt.Printf("%v", err)
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)

}
