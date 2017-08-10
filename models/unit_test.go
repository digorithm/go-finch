package models

import (
	"testing"

	"encoding/json"

	_ "github.com/lib/pq"
)

func newUnitForTest(t *testing.T) *Unit {
	return NewUnit(newDbForTest(t))
}

func allUnitsJSON() map[string]int {

	units := make(map[string]int)

	units["kg"] = 1
	units["pound"] = 2
	units["lb"] = 3
	units["tbsp"] = 4
	units["tsp"] = 5
	units["ounce"] = 6
	units["quantity"] = 7
	units["litre"] = 8
	units["cup"] = 9
	units["grams"] = 10

	return units
}

func TestGetAllUnits(t *testing.T) {

	u := newUnitForTest(t)
	tearDown := TestSetup(t, u.db)
	defer tearDown(t, u.db)

	res, err := u.GetAllUnits(nil)

	if err != nil {
		t.Errorf("Get units failed: %v", err)
	}

	var resJSON map[string]int

	err = json.Unmarshal(res, &resJSON)

	if err != nil {
		t.Errorf("Unmarshal failed: %v", err)
	}

	expected := allUnitsJSON()

	for val := range resJSON {

		if resJSON[val] != expected[val] {
			t.Errorf("GetAllUnits failed, got: %d, want: %d", resJSON[val], expected[val])
		}
	}

}
