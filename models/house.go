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

	users = createUserOwnTypeRows(users, data)

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

	recipes = createRecipeRows(recipes, data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return recipes, err
}

func (h *House) GetHouseStorage(tx *sqlx.Tx, house_id int64) ([]HouseStorageRow, error) {

	var storage []HouseStorageRow

	query := "SELECT I.NAME, S.AMOUNT, U.NAME FROM INGREDIENT I INNER JOIN ITEM_IN_STORAGE S ON I.ID = S.INGREDIENT_ID INNER JOIN UNIT U ON U.ID = S.UNIT_ID WHERE S.HOUSE_ID = $1"

	data, err := h.GetCompoundModel(nil, query, house_id)

	storage = createHouseStorageRows(storage, data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return storage, err
}

func (h *House) GetHouseSchedule(tx *sqlx.Tx, house_id int64) ([]HouseScheduleRow, error) {

	var schedule []HouseScheduleRow

	query := "SELECT W.DAY, T.TYPE, R.NAME FROM RECIPE R, WEEKDAY W, MEAL_TYPE T, SCHEDULE S WHERE S.HOUSE_ID = $1 AND S.WEEK_ID = W.ID AND S.TYPE_ID = T.ID AND S.RECIPE_ID = R.ID"

	data, err := h.GetCompoundModel(nil, query, house_id)

	schedule = createHouseScheduleRows(schedule, data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return schedule, err
}
