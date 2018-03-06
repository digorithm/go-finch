package finchgo

import "math/rand"
import "encoding/json"

func RandomIntRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func JSONToMap(JSONResponse []byte) map[string]interface{} {
	Struct := make(map[string]interface{})
	_ = json.Unmarshal(JSONResponse, &Struct)

	return Struct
}

func JSONListToMap(JSONResponse []byte) []map[string]interface{} {
	Struct := make([]map[string]interface{}, 0, 0)
	_ = json.Unmarshal(JSONResponse, &Struct)

	return Struct
}
