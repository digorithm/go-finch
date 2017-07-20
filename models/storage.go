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

//AddIngredientList adds all the ingredients received from the clientside to the house's storage
//if the ingredient received is not in the ingredient table, we add it to the table
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

// GetHouseStorage gets the house id, ingredient id, amount and unit of the ingredient,
// the name of the ingredient and the unit id of the house
func (i *ItemInStorage) GetHouseStorage(tx *sqlx.Tx, houseID int64) ([]HouseStorageRow, error) {

	query := "SELECT S.HOUSE_ID, I.ID, I.NAME, S.AMOUNT, S.UNIT_ID, U.NAME FROM INGREDIENT I INNER JOIN ITEM_IN_STORAGE S ON I.ID = S.INGREDIENT_ID INNER JOIN UNIT U ON U.ID = S.UNIT_ID WHERE S.HOUSE_ID = $1"

	data, err := i.GetCompoundModel(tx, query, houseID)

	storage := createHouseStorageRows(data)

	if err != nil {
		fmt.Printf("%v", err)
	}

	return storage, err
}

// AddIngToStorage adds to the storage of the given house_id
// the ingredient and the new amount it should have
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

	query := fmt.Sprintf("SELECT S.HOUSE_ID, S.INGREDIENT_ID, I.NAME, S.AMOUNT, S.UNIT_ID FROM ITEM_IN_STORAGE S INNER JOIN INGREDIENT I ON I.ID = S.INGREDIENT_ID WHERE S.HOUSE_ID = %v AND S.INGREDIENT_ID = $1 ORDER BY S.INGREDIENT_ID", houseID)

	res, err := i.GetCompoundModel(tx, query, ingID)

	if reflect.ValueOf(res).IsNil() {

		return storage, err

	}

	storage = createItemInStorage(res)

	return storage, err

}

func (i *ItemInStorage) DeleteStorage(tx *sqlx.Tx, houseID int64) error {

	where := fmt.Sprintf("house_id = %v", houseID)
	_, err := i.DeleteFromTable(tx, where)

	if err != nil {
		return fmt.Errorf("Delete Storage Failed: %v", err)
	}

	res, err := i.GetHouseStorage(tx, houseID)

	if res != nil {
		return fmt.Errorf("Delete Storage Failed")
	}

	return nil
}

func (i *ItemInStorage) DeleteStorageItems(tx *sqlx.Tx, respJSON []byte, houseID int64) ([]byte, error) {

	items := make([]map[string]interface{}, 0, 0)

	ingredient := NewIngredient(i.db)

	err := json.Unmarshal(respJSON, &items)

	if err != nil {
		err = fmt.Errorf("Unmarshal Failed: %v", err)
	}

	for _, item := range items {

		res, err := ingredient.GetByName(tx, item["name"].(string))

		if err != nil {
			err = fmt.Errorf("Finding Ingredient Failed: %v", err)
		}

		where := fmt.Sprintf("house_id = %v and ingredient_id = %v", houseID, res.ID)

		val, err := i.DeleteFromTable(tx, where)

		aff, _ := val.RowsAffected()

		if aff != 1 {
			err = fmt.Errorf("Delete Storage Item Failed: %v", err)
		}
	}

	storage, er := i.GetHouseStorage(tx, houseID)

	if er != nil {
		err = fmt.Errorf("Getting House Storage Failed: %v", er)
	}

	StorageJSON, er := json.Marshal(storage)

	if er != nil {
		err = fmt.Errorf("Getting House Storage Failed: %v", er)
	}

	return StorageJSON, err
}
