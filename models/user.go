package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(db *sqlx.DB) *User {
	user := &User{}
	user.db = db
	user.table = "user_info"
	user.hasID = true

	return user
}

type User struct {
	Base
}

func (u *User) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*UserRow, error) {
	userId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return u.GetById(tx, userId)
}

// AllUsers returns all user rows.
func (u *User) AllUsers(tx *sqlx.Tx) ([]*UserRow, error) {
	users := []*UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v", u.table)
	err := u.db.Select(&users, query)

	return users, err
}

// GetById returns record by id.
func (u *User) GetById(tx *sqlx.Tx, id int64) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=$1", u.table)
	err := u.db.Get(user, query, id)

	return user, err
}

// GetByEmail returns record by email.
func (u *User) GetByEmail(tx *sqlx.Tx, email string) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE email=$1", u.table)
	err := u.db.Get(user, query, email)

	return user, err
}

// GetByUsername returns record by username.
func (u *User) GetByUsername(tx *sqlx.Tx, username string) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE username=$1", u.table)
	err := u.db.Get(user, query, username)

	return user, err
}

// GetByEmail returns record by email but checks password first.
func (u *User) GetUserByEmailAndPassword(tx *sqlx.Tx, email, password string) (*UserRow, error) {
	user, err := u.GetByEmail(tx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, err
}

// Signup create a new record of user.
func (u *User) Signup(tx *sqlx.Tx, email, username, password, passwordAgain string) (*UserRow, error) {
	if email == "" {
		return nil, errors.New("Email cannot be blank.")
	}
	if password == "" {
		return nil, errors.New("Password cannot be blank.")
	}
	if password != passwordAgain {
		return nil, errors.New("Password is invalid.")
	}
	if username == "" {
		return nil, errors.New("Username cannot be blank")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["email"] = email
	data["username"] = username
	data["password"] = hashedPassword

	sqlResult, err := u.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return u.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (u *User) UpdateEmailAndPasswordById(tx *sqlx.Tx, userID int64, email, password, passwordAgain string) (*UserRow, error) {
	data := make(map[string]interface{})

	if email != "" {
		data["email"] = email
	}

	if password != "" && passwordAgain != "" && password == passwordAgain {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
		if err != nil {
			return nil, err
		}

		data["password"] = hashedPassword
	}

	if len(data) > 0 {
		_, err := u.UpdateByID(tx, data, userID)
		if err != nil {
			return nil, err
		}
	}

	return u.GetById(tx, userID)
}

func (u *User) GetUserHouses(tx *sqlx.Tx, userID int64) ([]UserHouseRow, error) {

	query := "SELECT H.ID, H.NAME, O.OWN_TYPE, O.DESCRIPTION FROM HOUSE H INNER JOIN MEMBER_OF M ON M.HOUSE_ID = H.ID INNER JOIN OWNERSHIP O ON O.OWN_TYPE = M.OWN_TYPE WHERE M.USER_ID = $1"

	data, err := u.GetCompoundModel(tx, query, userID)

	houses := createUserHouseRows(data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return houses, err
}

// GetUserRecipes used by house and user
func (u *User) GetUserRecipes(tx *sqlx.Tx, userID int64) ([]RecipeRow, error) {

	query := "SELECT R.ID, R.NAME, R.SERVES_FOR FROM RECIPE R INNER JOIN USER_RECIPE U ON R.ID = U.RECIPE_ID WHERE U.USER_ID = $1"

	return u.GetRecipeForStruct(tx, query, userID)
}

// getMembership returns the membership of the user for the given houseID
func (u *User) getMembership(tx *sqlx.Tx, userID, houseID int64) (*OwnerRow, error) {

	query := fmt.Sprintf("SELECT O.OWN_TYPE, O.DESCRIPTION FROM OWNERSHIP O INNER JOIN MEMBER_OF M ON M.OWN_TYPE = O.OWN_TYPE WHERE M.HOUSE_ID = %v AND M.USER_ID = $1", houseID)

	data, err := u.GetCompoundModel(tx, query, userID)

	if err != nil {
		fmt.Printf("%v", err)
	}

	ownership := createOwnerRow(data)
	return ownership, err

}

// AddRecipe adds a recipe and binds it to a user.
// It will come as a JSON and then the respective Handler will break
// the json into 2 maps: recipe and steps. This will be used
// to add the necessary information to add the recipe to the database.
func (u *User) AddRecipe(tx *sqlx.Tx, jsonRecipe []byte, userID int64) ([]FullRecipeRow, error) {

	var FullRecipe []FullRecipeRow
	recipeObj := NewRecipe(u.db)

	// Add recipe metadata to DB

	recipeResult, err := u.addJsonRecipeTable(tx, jsonRecipe)
	if err != nil {
		return FullRecipe, err
	}

	recipeId, _ := recipeResult.LastInsertId()

	_, err = u.addJsonRecipeTypeTable(tx, recipeId, jsonRecipe)

	if err != nil {
		return FullRecipe, err
	}

	// Add both step and step_ingredient tables to the DB
	_, _, err = u.addJsonIngredientStepTable(tx, recipeId, jsonRecipe)

	if err != nil {
		return FullRecipe, err
	}

	var user_recipe_db Base
	user_recipe_db.db = u.db
	user_recipe_db.table = "user_recipe"
	user_recipe_db.hasID = false

	userRecipeData := make(map[string]interface{})
	userRecipeData["user_id"] = userID
	userRecipeData["recipe_id"] = recipeId
	user_recipe_db.InsertIntoTable(tx, userRecipeData)

	FullRecipe, err = recipeObj.GetFullRecipe(tx, recipeId)

	return FullRecipe, err
}

func (u *User) addJsonRecipeTable(tx *sqlx.Tx, jsonData []byte) (sql.Result, error) {
	recipeObj := NewRecipe(u.db)

	var err error

	recipeData := make(map[string]interface{})
	recipeData["name"], err = jsonparser.GetString(jsonData, "recipe_name")
	recipeData["serves_for"], err = jsonparser.GetString(jsonData, "serves_for")

	if err != nil {
		return nil, err
	}

	recipeResult, err := recipeObj.InsertIntoTable(tx, recipeData)

	return recipeResult, err
}

func (u *User) addJsonIngredientStepTable(tx *sqlx.Tx, recipeId int64, jsonData []byte) (sql.Result, sql.Result, error) {

	ingredientObj := NewIngredient(u.db)

	var err error
	var stepResult sql.Result
	var stepIngredientResult sql.Result

	var step_db Base
	step_db.db = u.db
	step_db.table = "step"
	step_db.hasID = true

	var step_ingredient_db Base
	step_ingredient_db.db = u.db
	step_ingredient_db.table = "step_ingredient"
	step_ingredient_db.hasID = false

	jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		// Insert data into step table
		stepData := make(map[string]interface{})
		stepData["recipe_id"] = recipeId
		stepData["id"], _ = jsonparser.GetInt(value, "step_id")
		stepData["text"], _ = jsonparser.GetString(value, "text")

		recipeType := make(map[string]interface{})
		recipeType["recipe_id"] = recipeId

		stepResult, err = step_db.InsertIntoTable(tx, stepData)

		if err != nil {
			fmt.Printf("Error while adding step into DB. Error: %v", err)
		}

		jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

			stepIngredientData := make(map[string]interface{})
			stepIngredientData["recipe_id"] = stepData["recipe_id"]
			stepIngredientData["unit_id"], _ = jsonparser.GetInt(value, "unit")
			stepIngredientData["amount"], _ = jsonparser.GetFloat(value, "amount")
			stepIngredientData["step_id"] = stepData["id"]

			// Check if ingredient ID exists in the DB
			stepIngredientDataName, _ := jsonparser.GetString(value, "name")
			iRow, _ := ingredientObj.GetByName(tx, stepIngredientDataName)
			if iRow == nil {
				addedIRow, _ := ingredientObj.AddIngredient(tx, stepIngredientDataName)
				stepIngredientData["ingredient_id"] = addedIRow.ID
			} else {
				stepIngredientData["ingredient_id"] = iRow.ID
			}
			stepIngredientResult, err = step_ingredient_db.InsertIntoTable(tx, stepIngredientData)
			if err != nil {
				fmt.Printf("Error while adding step ingredient into DB. Error: %v", err)
			}
		}, "step_ingredients")
	}, "steps")

	return stepResult, stepIngredientResult, err
}

func (u *User) addJsonRecipeTypeTable(tx *sqlx.Tx, recipeID int64, jsonData []byte) (sql.Result, error) {

	recipeTypeObj := NewRecipeType(u.db)

	var recipeTypeRes sql.Result
	var err error

	jsonparser.ArrayEach(jsonData, func(jsonData []byte, dataType jsonparser.ValueType, offset int, err error) {

		data := make(map[string]interface{})
		data["recipe_id"] = recipeID
		typeS := string(jsonData)
		data["type_id"] = getTypeID(string(typeS))

		recipeTypeRes, err = recipeTypeObj.InsertIntoTable(tx, data)

	}, "type")

	return recipeTypeRes, err

}

func getTypeID(typeR string) int {

	var ID int

	switch typeR {
	case "Breakfast":
		ID = 1
	case "Snack":
		ID = 2
	case "Lunch":
		ID = 3
	case "Dinner":
		ID = 4
	}

	return ID
}
