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

func GetUserInvitations(w http.ResponseWriter, r *http.Request) {
	joinObj := CreateJoinObj(r)

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])

	if err != nil {
		fmt.Printf("%v", err)
		libhttp.HandleErrorJson(w, err)
	}

	inviteJSON, err := joinObj.GetUserInvitations(nil, int64(userID))

	w.WriteHeader(http.StatusOK)
	w.Write(inviteJSON)
}

func InviteUser(w http.ResponseWriter, r *http.Request) {

	joinObj := CreateJoinObj(r)

	invitation, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Printf("%v", err)
		libhttp.HandleErrorJson(w, err)
	}

	responseJSON, err := joinObj.AddInvitation(nil, invitation)

	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

func DeleteInvitation(w http.ResponseWriter, r *http.Request) {

	joinObj := CreateJoinObj(r)

	vars := mux.Vars(r)
	inviteID, err := strconv.Atoi(vars["invite_id"])

	if err != nil {
		fmt.Printf("%v", err)
		libhttp.HandleErrorJson(w, err)
	}

	err = joinObj.DeleteInvitation(nil, int64(inviteID))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
