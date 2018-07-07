package main

import "fmt"

import "path/filepath"
import "runtime"
import "encoding/json"
import "io/ioutil"
import "strconv"
import "os"
import "strings"

func main2() {
	KnobsPath := "meal_planner/finch_knobs.json"
	_, filename, _, _ := runtime.Caller(0)
	dir, _ := filepath.Split(filepath.Dir(filename))
	file := filepath.Join(dir, KnobsPath)

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var conf map[string]interface{}
	json.Unmarshal(raw, &conf)
	for idx, _ := range conf {
		s := strconv.FormatFloat(conf[idx].(float64), 'f', 0, 64)
		if len(s) == 4 {
			newValue, _ := strconv.ParseFloat(strings.TrimSuffix(s, "000"), 64)
			conf[idx] = newValue
		} else if len(s) == 1 {
			extra := "000"
			newValue, _ := strconv.ParseFloat(s+string(extra), 64)
			conf[idx] = newValue
		}
	}
	newJson, _ := json.Marshal(conf)

	err = ioutil.WriteFile(file, newJson, 0644)
}
