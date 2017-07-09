package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"encoding/json"

	"io/ioutil"

	"github.com/digorithm/meal_planner/libhttp"
	"github.com/digorithm/meal_planner/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func CreateUserObj(r *http.Request) *models.User {
	db := r.Context().Value("db").(*sqlx.DB)

	userObj := models.NewUser(db)

	return userObj
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {

	userObj := CreateUserObj(r)

	users, err := userObj.AllUsers(nil)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}

	JSONResponse := buildAllUsersJSONResponse(users)

	w.WriteHeader(http.StatusOK)
	w.Write(JSONResponse)
}

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {

	var user []*models.UserRow
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["user_id"])

	userObj := CreateUserObj(r)

	if err != nil {
		fmt.Println(err)
	}

	u, err := userObj.GetById(nil, int64(userID))

	if err != nil {
		fmt.Printf("Something went wrong while fecthing the user. Error: %v", err)
	}

	user = append(user, u)

	if len(user) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	JSONResponse := buildAllUsersJSONResponse(user)

	w.WriteHeader(http.StatusOK)
	w.Write(JSONResponse)
}

func buildAllUsersJSONResponse(users []*models.UserRow) []byte {

	finalUsers := make([]map[string]interface{}, 0, 0)

	for _, user := range users {

		finalUser := make(map[string]interface{})
		finalUser["id"] = user.ID
		finalUser["name"] = user.Username
		finalUser["email"] = user.Email

		finalUsers = append(finalUsers, finalUser)
	}

	finalUsersJSON, _ := json.MarshalIndent(finalUsers, "", "	")

	return finalUsersJSON
}

/*func GetSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/users/users-external.html.tmpl", "templates/users/signup.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, nil)
}*/

func PostSignup(w http.ResponseWriter, r *http.Request) {

	userObj := CreateUserObj(r)

	body, err := ioutil.ReadAll(r.Body)

	userJSON, err := userObj.Signup(nil, body)

	if err != nil {
		fmt.Printf("%v", err)
		libhttp.HandleErrorJson(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(userJSON)

	// Perform login
	//PostLogin(w, r)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	userObj := CreateUserObj(r)

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])

	_, err = userObj.DeleteById(nil, int64(userID))

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong in Delete user with email"))
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

/*

func GetLoginWithoutSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/users/users-external.html.tmpl", "templates/users/login.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, nil)
}

// GetLogin get login page.
func GetLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "meal_planner-session")

	currentUserInterface := session.Values["user"]
	if currentUserInterface != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	GetLoginWithoutSession(w, r)
}

// PostLogin performs login.
func PostLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	email := r.FormValue("email")
	password := r.FormValue("password")

	u := createUserObj(r)

	user, err := u.GetUserByEmailAndPassword(nil, email, password)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	session, _ := sessionStore.Get(r, "meal_planner-session")
	session.Values["user"] = user

	err = session.Save(r, w)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	http.Redirect(w, r, "/", 302)
}


func GetLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "meal_planner-session")

	delete(session.Values, "user")
	session.Save(r, w)

	http.Redirect(w, r, "/login", 302)
}

func PostPutDeleteUsersID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	method := r.FormValue("_method")
	if method == "" || strings.ToLower(method) == "post" || strings.ToLower(method) == "put" {
		PutUsersID(w, r)
	} else if strings.ToLower(method) == "delete" {
		DeleteUsersID(w, r)
	}
}

func PutUsersID(w http.ResponseWriter, r *http.Request) {
	userId, err := getIDFromPath(w, r)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	db := r.Context().Value("db").(*sqlx.DB)

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "meal_planner-session")

	currentUser := session.Values["user"].(*models.UserRow)

	if currentUser.ID != userId {
		err := errors.New("Modifying other user is not allowed.")
		libhttp.HandleErrorJson(w, err)
		return
	}

	email := r.FormValue("Email")
	password := r.FormValue("Password")
	passwordAgain := r.FormValue("PasswordAgain")

	u := models.NewUser(db)

	currentUser, err = u.UpdateEmailAndPasswordById(nil, currentUser.ID, email, password, passwordAgain)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	// Update currentUser stored in session.
	session.Values["user"] = currentUser
	err = session.Save(r, w)

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	http.Redirect(w, r, "/", 302)
}

*/
