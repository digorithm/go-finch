package models

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
)

func NewItemInStorage(db *sqlx.DB) *ItemInStorage {
	storage := &ItemInStorage{}
	storage.db = db
	storage.table = "ITEM_IN_STORAGE"
	storage.hasID = true

	return storage
}

type ItemInStorage struct {
	Base
}

func (i *ItemInStorage) GetHouseStorage(tx *sqlx.Tx, houseID int64) ([]HouseStorageRow, error) {

	query := "SELECT S.HOUSE_ID, I.ID, S.AMOUNT, S.UNIT_ID, I.NAME, U.NAME FROM INGREDIENT I INNER JOIN ITEM_IN_STORAGE S ON I.ID = S.INGREDIENT_ID INNER JOIN UNIT U ON U.ID = S.UNIT_ID WHERE S.HOUSE_ID = $1"

	data, err := i.GetCompoundModel(tx, query, houseID)

	storage := createHouseStorageRows(data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return storage, err
}

func (i *ItemInStorage) AddIngToStorage(tx *sqlx.Tx, houseID int64, ingID int64, amount float64, unitID int64) (ItemInStorageRow, error) {

	data := make(map[string]interface{})
	data["house_id"] = houseID
	data["ingredient_id"] = ingID
	data["amount"] = amount
	data["unit_id"] = unitID

	_, err := i.InsertIntoMultiKeyTable(tx, data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return i.GetStorageIngredient(tx, houseID, ingID)
}

func (i *ItemInStorage) UpdateStorage(tx *sqlx.Tx, houseID int64, ingID int64, newAmt float64, newUnt int64) (ItemInStorageRow, error) {

	emptyStorage := ItemInStorageRow{}

	res, _ := i.GetStorageIngredient(tx, houseID, ingID)

	if res == emptyStorage {

		return i.AddIngToStorage(tx, houseID, ingID, newAmt, newUnt)

	}

	return i.UpdateIngredient(tx, houseID, ingID, newAmt, newUnt)

}

func (i *ItemInStorage) UpdateIngredient(tx *sqlx.Tx, houseID int64, ingID int64, newAmt float64, newUnt int64) (ItemInStorageRow, error) {

	where := fmt.Sprintf("HOUSE_ID = %v AND INGREDIENT_ID = %v", houseID, ingID)

	data := make(map[string]interface{})
	data["amount"] = newAmt
	data["unit_id"] = newUnt

	_, err := i.UpdateFromTable(tx, data, where)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return i.GetStorageIngredient(tx, houseID, ingID)
}

func (i *ItemInStorage) GetStorageIngredient(tx *sqlx.Tx, houseID, ingID int64) (ItemInStorageRow, error) {

	var storage ItemInStorageRow

	query := fmt.Sprintf("SELECT * FROM ITEM_IN_STORAGE WHERE HOUSE_ID = %v AND INGREDIENT_ID = $1", houseID)

	res, err := i.GetCompoundModel(tx, query, ingID)

	if reflect.ValueOf(res).IsNil() {

		return storage, err

	}

	storage = createItemInStorage(res)

	return storage, err

}

func (i *ItemInStorage) AddIngredientList(JSONRequest []byte, HouseID int64) error {

	Ingredients := make([]map[string]interface{}, 0, 0)

	IngredientObj := NewIngredient(i.db)

	_ = json.Unmarshal(JSONRequest, &Ingredients)

	for _, Ingredient := range Ingredients {

		// Check if ingredient exists in the database
		IRow, _ := IngredientObj.GetByName(nil, Ingredient["name"].(string))

		// If not, add it to the database
		if IRow == nil {
			IRow, _ = IngredientObj.AddIngredient(nil, Ingredient["name"].(string))
		}

		_, err := i.UpdateStorage(nil, int64(HouseID), IRow.ID, Ingredient["amount"].(float64), int64(Ingredient["unit"].(float64)))

		if err != nil {
			return err
		}
	}

	return nil
}
