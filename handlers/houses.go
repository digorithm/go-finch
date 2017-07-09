package handlers

import (
	"net/http"

	"github.com/digorithm/meal_planner/models"
	"github.com/jmoiron/sqlx"
)

func CreateHouseObj(r *http.Request) *models.House {
	db := r.Context().Value("db").(*sqlx.DB)

	houseObj := models.NewHouse(db)

	return houseObj
}
