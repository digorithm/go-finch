package handlers

import (
	"fmt"
	"net/http"

	"github.com/digorithm/meal_planner/models"
	"github.com/jmoiron/sqlx"
)

func CreateUnitObj(r *http.Request) *models.Unit {
	db := r.Context().Value("db").(*sqlx.DB)

	unitObj := models.NewUnit(db)

	return unitObj
}

func GetAllUnitsHandler(w http.ResponseWriter, r *http.Request) {
	unitObj := CreateUnitObj(r)

	res, err := unitObj.GetAllUnits(nil)

	if err != nil {
		fmt.Printf("getAllUnitsHandler failed: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
