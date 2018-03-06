// Package handlers provides request handlers.
package handlers

import (
	"encoding/json"
	"errors"
	"github.com/digorithm/meal_planner/finchgo"
	"github.com/digorithm/meal_planner/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
)

var WorkflowSimulator = finchgo.NewWorkflowSimulator("http://localhost:8888")
var Finch *finchgo.Finch

func getCurrentUser(w http.ResponseWriter, r *http.Request) *models.UserRow {
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "meal_planner-session")
	return session.Values["user"].(*models.UserRow)
}

func getIDFromPath(w http.ResponseWriter, r *http.Request) (int64, error) {
	userIDString := mux.Vars(r)["id"]
	if userIDString == "" {
		return -1, errors.New("user id cannot be empty")
	}

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

// RespondWithError sends a JSON formatted error
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON sends a JSON formatted response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func StartWorkflowSimulatorHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	ConcurrentUsers, err := strconv.Atoi(vars["concurrent_users"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing number of concurrent users"))
	}

	go WorkflowSimulator.Run(ConcurrentUsers)

	w.WriteHeader(http.StatusOK)
}

func StopWorkflowSimulatorHandler(w http.ResponseWriter, r *http.Request) {

	go WorkflowSimulator.Stop()

	w.WriteHeader(http.StatusOK)
}
