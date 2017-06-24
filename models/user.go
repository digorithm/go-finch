package models

import (
	"database/sql"
	"errors"
	"fmt"

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

	var houses []UserHouseRow

	query := "SELECT H.ID, H.NAME, O.OWN_TYPE, O.DESCRIPTION FROM HOUSE H INNER JOIN MEMBER_OF M ON M.HOUSE_ID = H.ID INNER JOIN OWNERSHIP O ON O.OWN_TYPE = M.OWN_TYPE WHERE M.USER_ID = $1"

	data, err := u.GetCompoundModel(nil, query, userID)

	houses = createUserHouseRows(houses, data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return houses, err
}

func (u *User) GetUserRecipes(tx *sqlx.Tx, userID int64) ([]RecipeRow, error) {

	query := "SELECT R.ID, R.NAME, R.TYPE, R.SERVES_FOR FROM RECIPE R INNER JOIN USER_RECIPE U ON R.ID = U.RECIPE_ID WHERE U.USER_ID = $1"

	return u.GetRecipeForStruct(tx, query, userID)
}

// AddRecipe adds a recipe and binds it to a user.
// It will come here as a JSON and then the respective Handler will break
// the json into 3 maps: recipe, ingredients and steps. This will be used
// to add the necessary information to add the recipe to the database.
func (u *User) AddRecipe(tx *sqlx.Tx, recipe map[string]interface{}, steps []map[string]interface{}) (FullRecipeRow, error) {

	var FullRecipe FullRecipeRow
	recipeObj := NewRecipe(u.db)
	ingredientObj := NewIngredient(u.db)

	// Add recipe metadata to DB
	recipe_result, err := recipeObj.InsertIntoTable(tx, recipe)
	if err != nil {
		return FullRecipe, err
	}

	// Add steps to the DB
	var step_db Base
	step_db.db = u.db
	step_db.table = "step"
	step_db.hasID = true

	var step_ingredient_db Base
	step_ingredient_db.db = u.db
	step_ingredient_db.table = "step_ingredient"
	step_ingredient_db.hasID = false

	for _, step := range steps {

		// Extract ingredients and save them to DB
		ingredientStrings := u.ExtractInterfaceSliceOfStrings(step["step_ingredients"])
		ingredientIDs, err := ingredientObj.AddIngredients(tx, ingredientStrings)

		if err != nil {
			return FullRecipe, err
		}

		// Add info about steps to the DB
		step_table_data := make(map[string]interface{})
		step_table_data["recipe_id"], _ = recipe_result.LastInsertId()
		step_table_data["text"] = step["text"]
		step_table_data["id"] = step["step_id"]
		step_table_result, err := step_db.InsertIntoTable(tx, step_table_data)

		if err != nil {
			fmt.Println(step_table_result)
			return FullRecipe, err
		}

		for _, ingId := range ingredientIDs {
			fmt.Println("AYOOO")
			// Add info about step_ingredient to the DB
			step_ingredient_data := make(map[string]interface{})
			step_ingredient_data["recipe_id"] = step_table_data["recipe_id"]
			step_ingredient_data["step_id"] = step_table_data["id"]
			step_ingredient_data["ingredient_id"] = ingId
			step_ingredient_data["unit_id"] = step["unit"]
			step_ingredient_data["amount"] = step["amount"]

			step_ingredient_result, err := step_ingredient_db.InsertIntoTable(tx, step_ingredient_data)
			if err != nil {
				return FullRecipe, err
			}
			fmt.Println(step_ingredient_result)
		}
	}

	return FullRecipe, err
}
