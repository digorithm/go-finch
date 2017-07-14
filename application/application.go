package application

import (
	"net/http"

	"github.com/carbocation/interpose"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"github.com/digorithm/meal_planner/handlers"
	"github.com/digorithm/meal_planner/middlewares"
)

// New is the constructor for Application struct.
func New(config *viper.Viper) (*Application, error) {
	dsn := config.Get("dsn").(string)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	cookieStoreSecret := config.Get("cookie_secret").(string)

	app := &Application{}
	app.config = config
	app.dsn = dsn
	app.db = db
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, err
}

// Application is the application object that runs HTTP server.
type Application struct {
	config       *viper.Viper
	dsn          string
	db           *sqlx.DB
	sessionStore sessions.Store
}

func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetDB(app.db))
	middle.Use(middlewares.SetSessionStore(app.sessionStore))
	middle.Use(middlewares.Log())

	middle.UseHandler(app.Mux())

	return middle, nil
}

func (app *Application) Mux() *gorilla_mux.Router {
	//MustLogin := middlewares.MustLogin

	router := gorilla_mux.NewRouter()

	//router.Handle("/", MustLogin(http.HandlerFunc(handlers.GetHome))).Methods("GET")
	/*router.HandleFunc("/", handlers.GetMain).Methods("GET")

	router.HandleFunc("/signup", handlers.GetSignup).Methods("GET")
	router.HandleFunc("/signup", handlers.PostSignup).Methods("POST")
	router.HandleFunc("/login", handlers.GetLogin).Methods("GET")
	router.HandleFunc("/login", handlers.PostLogin).Methods("POST")
	router.HandleFunc("/logout", handlers.GetLogout).Methods("GET")

	router.Handle("/users/{id:[0-9]+}", MustLogin(http.HandlerFunc(handlers.PostPutDeleteUsersID))).Methods("POST", "PUT", "DELETE")*/

	router.HandleFunc("/recipes/house/{house_id}", handlers.GetHouseRecipesHandler).Methods("GET")
	router.HandleFunc("/recipes/user/{user_id}", handlers.GetUserRecipesHandler).Methods("GET")
	router.HandleFunc("/recipes/{recipe_id}", handlers.GetRecipeByIDHandler).Methods("GET")
	router.HandleFunc("/recipes/{recipe_id}", handlers.DeleteRecipesHandler).Methods("DELETE")
	router.HandleFunc("/recipes/{recipe_id}/{field}", handlers.UpdateRecipesHandler).Methods("PUT")
	router.HandleFunc("/recipes/{recipe_id}/step/{step_id}", handlers.UpdateRecipeStepIngredientHandler).Methods("PUT")
	router.HandleFunc("/recipes", handlers.GetRecipesHandler).Methods("GET")
	router.HandleFunc("/recipes", handlers.AddRecipesHandler).Methods("POST")

	router.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")
	router.HandleFunc("/users/{user_id}", handlers.GetUserByIDHandler).Methods("GET")
	router.HandleFunc("/users", handlers.PostSignupHandler).Methods("POST")
	router.HandleFunc("/users/{user_id}", handlers.DeleteUserHandler).Methods("DELETE")


	router.HandleFunc("/houses/{house_id}", handlers.GetHouseHandler).Methods("GET")
	router.HandleFunc("/houses", handlers.PostHouseHandler).Methods("POST")
	router.HandleFunc("/houses/{house_id}", handlers.DeleteHouseHandler).Methods("DELETE")
	router.HandleFunc("/houses/{house_id}", handlers.UpdateHouseHandler).Methods("PUT")

	router.HandleFunc("/storages/{house_id}", handlers.GetStoragesHandler).Methods("GET")
	router.HandleFunc("/storages/{house_id}", handlers.PostStoragesHandler).Methods("POST")

	router.HandleFunc("/invitations/users/{user_id}", handlers.GetUserInvitationsHandler).Methods("GET")
	router.HandleFunc("/invitations/houses/{house_id}", handlers.GetHouseInvitationsHandler).Methods("GET")
	router.HandleFunc("/invitations/join", handlers.InviteUserHandler).Methods("POST")
	router.HandleFunc("/invitations/respond", handlers.InviteResponseHandler).Methods("POST")
	router.HandleFunc("/houses/{house_id}/users/{user_id}", handlers.DeleteMemberHandler).Methods("DELETE")

	router.HandleFunc("/requests/houses/{house_id}", handlers.GetHouseJoinsHandler).Methods("GET")
	router.HandleFunc("/requests/users/{user_id}", handlers.GetUserJoinsHandler).Methods("GET")
	router.HandleFunc("/requests/join", handlers.RequestJoinHandler).Methods("POST")
	router.HandleFunc("/requests/respond", handlers.RespondRequestJoinHandler).Methods("POST")

	// WORKS FOR BOTH INVITE AND REQUEST
	router.HandleFunc("/invitations/{invite_id}", handlers.DeleteRequestHandler).Methods("DELETE")
	router.HandleFunc("/units", handlers.GetAllUnitsHandler).Methods("GET")

	router.HandleFunc("/schedules/{house_id}", handlers.GetScheduleHandler).Methods("GET")
	router.HandleFunc("/schedules/{house_id}", handlers.DeleteScheduleHandler).Methods("DELETE")
	router.HandleFunc("/schedules/{house_id}", handlers.ModifyScheduleHandler).Methods("POST")
	router.HandleFunc("/schedules/create/{house_id}", handlers.CreateScheduleHandler).Methods("POST")
	//router.HandleFunc("/schedules/new/{house_id}", handlers.NewFullScheduleHandler).Methods("GET")

	router.HandleFunc("/meals", handlers.GetMealTypesHandler).Methods("GET")

	router.HandleFunc("/days", handlers.GetDaysHandler).Methods("GET")

	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return router
}
