package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// NewHouse creates a new base that points to house table
func NewHouse(db *sqlx.DB) *House {
	house := &House{}
	house.db = db
	house.table = "house"
	house.hasID = true

	return house
}

// House is a base
type House struct {
	Base
}

//HouseRowFromSQLResult returns the house that was last inserted to house table
func (h *House) HouseRowFromSQLResult(tx *sqlx.Tx, sqlResult sql.Result) (*HouseRow, error) {

	houseID, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return h.GetByID(tx, houseID)
}

// AllHouses returns every house in the housetable
func (h *House) AllHouses(tx *sqlx.Tx) ([]*HouseRow, error) {
	houses := []*HouseRow{}
	query := fmt.Sprintf("SELECT * FROM %v", h.table)
	err := h.db.Select(&houses, query)

	return houses, err
}

// GetByID returns a houseRow with the given id
func (h *House) GetByID(tx *sqlx.Tx, id int64) (*HouseRow, error) {
	house := &HouseRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=$1", h.table)
	err := h.db.Get(house, query, id)

	return house, err
}

// CreateHouse creates a house and instantiates the schedule for it
func (h *House) CreateHouse(tx *sqlx.Tx, name string, groceryDay string, household int64) (*HouseRow, error) {

	if name == "" {
		return nil, errors.New("House name cannot be blank")
	}

	data := make(map[string]interface{})
	data["name"] = name
	data["grocery_day"] = groceryDay
	data["household_number"] = household

	sqlResult, err := h.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return h.HouseRowFromSQLResult(tx, sqlResult)
}

// GetHouseUsers retrieves all users of a given house id and returns user id, email, password, name, ownership id and description
func (h *House) GetHouseUsers(tx *sqlx.Tx, houseID int64) ([]UserOwnTypeRow, error) {

	query := "SELECT U.ID, U.EMAIL, U.PASSWORD, U.USERNAME, O.OWN_TYPE, O.DESCRIPTION FROM USER_INFO U INNER JOIN MEMBER_OF M ON M.USER_ID = U.ID INNER JOIN OWNERSHIP O ON O.OWN_TYPE = M.OWN_TYPE WHERE M.HOUSE_ID = $1"

	data, err := h.GetCompoundModel(tx, query, houseID)

	users := createUserOwnTypeRows(data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return users, err
}

// GetHouseRecipes retrieves all recipes of a given house id and returns
// recipe id, name and serving size
func (h *House) GetHouseRecipes(tx *sqlx.Tx, houseID int64) ([]RecipeRow, error) {

	query := "SELECT R.ID, R.NAME, R.SERVES_FOR FROM RECIPE R INNER JOIN HOUSE_RECIPE H ON R.ID = H.RECIPE_ID WHERE H.HOUSE_ID = $1"

	return h.GetRecipeForStruct(tx, query, houseID)
}

// UpdateHouseHold updates the number of residents in a given house id
// and returns the updated house
func (h *House) UpdateHouseHold(tx *sqlx.Tx, houseHold int64, houseID int64) (*HouseRow, error) {

	data := make(map[string]interface{})
	data["household_number"] = houseHold
	where := fmt.Sprintf("ID = %v", houseID)

	_, err := h.UpdateFromTable(tx, data, where)

	if err != nil {
		fmt.Println(err)
	}

	return h.GetByID(tx, houseID)

}
