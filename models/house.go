package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewHouse(db *sqlx.DB) *House {
	house := &House{}
	house.db = db
	house.table = "house"
	house.hasID = true

	return house
}

type House struct {
	Base
}

func (h *House) houseRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*HouseRow, error) {
	houseId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return h.GetById(tx, houseId)
}

func (h *House) AllHouses(tx *sqlx.Tx) ([]*HouseRow, error) {
	houses := []*HouseRow{}
	query := fmt.Sprintf("SELECT * FROM %v", h.table)
	err := h.db.Select(&houses, query)

	return houses, err
}

func (h *House) GetById(tx *sqlx.Tx, id int64) (*HouseRow, error) {
	house := &HouseRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=$1", h.table)
	err := h.db.Get(house, query, id)

	return house, err
}

func (h *House) CreateHouse(tx *sqlx.Tx, name string) (*HouseRow, error) {

	if name == "" {
		return nil, errors.New("House name cannot be blank")
	}

	data := make(map[string]interface{})
	data["name"] = name

	sqlResult, err := h.InsertIntoTable(tx, data)

	if err != nil {
		return nil, err
	}

	return h.houseRowFromSqlResult(tx, sqlResult)
}

func (h *House) GetHouseUsers(tx *sqlx.Tx, house_id int64) ([]UserOwnTypeRow, error) {

	var users []UserOwnTypeRow

	query := "SELECT U.ID, U.EMAIL, U.PASSWORD, U.USERNAME, O.OWN_TYPE, O.DESCRIPTION FROM USER_INFO U INNER JOIN MEMBER_OF M ON M.USER_ID = U.ID INNER JOIN OWNERSHIP O ON O.OWN_TYPE = M.OWN_TYPE WHERE M.HOUSE_ID = $1"

	data, err := h.GetCompoundModel(nil, query, house_id)
	// Move data to users

	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Println(data)

	return users, err
}

func (h *House) GetHouseRecipes(tx *sqlx.Tx, house_id int64) ([]RecipeRow, error) {
	var recipes []RecipeRow

	query := "SELECT R.ID, R.NAME FROM RECIPE R INNER JOIN HOUSE_RECIPE H ON R.ID = H.RECIPE_ID WHERE H.HOUSE_ID = $1"

	data, err := h.GetCompoundModel(nil, query, house_id)

	// Move data to recipes

	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Println(data)

	return recipes, err
}
