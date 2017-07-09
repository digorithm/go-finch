package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/digorithm/meal_planner/libhttp"
	"github.com/digorithm/meal_planner/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func CreateJoinObj(r *http.Request) *models.Join {
	db := r.Context().Value("db").(*sqlx.DB)

	joinObj := models.NewJoin(db)

	return joinObj
}

func GetHouseInvitations(w http.ResponseWriter, r *http.Request) {

	joinObj := CreateJoinObj(r)

	vars := mux.Vars(r)
	houseID, err := strconv.Atoi(vars["house_id"])

	if err != nil {
		fmt.Printf("%v", err)
		libhttp.HandleErrorJson(w, err)
	}

	inviteJSON, err := joinObj.GetHouseInvitations(nil, int64(houseID))

	w.WriteHeader(http.StatusOK)
	w.Write(inviteJSON)

}
