package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHouseStorageEndpoint(t *testing.T) {

	endpoint := "/storages/1"
	method := "GET"

	request, _ := http.NewRequest(method, endpoint, nil)

	request = SetTestDBEnv(request)

	response := httptest.NewRecorder()

	RouterForTest().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}

	var res []map[string]interface{}
	_ = json.Unmarshal(body, &res)

	assert.Equal(t, res[0]["ingredient"], "potato")
	assert.Equal(t, res[0]["amount"].(float64), 5.0)
}
