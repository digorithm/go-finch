package finchgo

import "math/rand"

import "encoding/json"

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
)

type Metric struct {
	Name     string `json:"__name__"`
	Device   string `json:"device"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
	Quantile string `json:"quantile"`
	Knob     string `json:knob`
}

type SingleResult struct {
	Value []interface{} `json:"value"`
}

type SinglePrometheusData struct {
	ResultType string         `json:"resultType"`
	Result     []SingleResult `json:"result"`
}

type SinglePrometheusJSON struct {
	Status string               `json:"status"`
	Data   SinglePrometheusData `json:"data"`
}

type RangeResult struct {
	Metric Metric
	Values []interface{} `json:"values"`
}

type RangePrometheusData struct {
	ResultType string        `json:"resultType"`
	Result     []RangeResult `json:"result"`
}

type RangePrometheusJSON struct {
	Status string              `json:"status"`
	Data   RangePrometheusData `json:"data"`
}

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

func buildStructs(MetricData map[string][]byte) map[string]RangePrometheusJSON {

	MetricStructs := make(map[string]RangePrometheusJSON)

	for metric, data := range MetricData {
		JSONStruct := parseToPrometheusStruct(data)
		MetricStructs[metric] = JSONStruct
	}

	return MetricStructs
}

func extractFromPrometheus(MetricQueries map[string]string) map[string][]byte {
	MetricData := make(map[string][]byte)
	for metric, query := range MetricQueries {
		data, err := getBodyFromURL(query)

		if err != nil {
			logrus.Fatal("Couldn't download data from Prometheus:: %v", err)
		}
		//fmt.Printf("Key %v: %v\n\n", metric, string(data))

		MetricData[metric] = data
	}

	return MetricData
}

func buildRangeQueries(Metrics []string, StartTime, EndTime string) map[string]string {
	MetricQuery := make(map[string]string)
	domain := "prometheus"
	port := "9090"

	if contains(Metrics, "HTTPRequestCount") {
		MetricQuery["HTTPRequestCount"] = fmt.Sprintf("http://%v:%v/api/v1/query_range?query=sum(irate(app_http_request_count[1m]))&start=%v&end=%v&step=10", domain, port, StartTime, EndTime)
	}

	if contains(Metrics, "HTTPRequestLatency") {
		MetricQuery["HTTPRequestLatency"] = fmt.Sprintf("http://%v:%v/api/v1/query_range?query=app_http_request_latency&start=%s&end=%s&step=10", domain, port, StartTime, EndTime)
	}

	if contains(Metrics, "IOWait") {
		MetricQuery["IOWait"] = fmt.Sprintf("http://%v:%v/api/v1/query_range?query=avg(irate(node_cpu{job='node-exporter',mode='iowait'}[1m]))*100&start=%v&end=%v&step=10", domain, port, StartTime, EndTime)
	}

	if contains(Metrics, "MemoryUsage") {
		MetricQuery["MemoryUsage"] = fmt.Sprintf("http://%v:%v/api/v1/query_range?query=((node_memory_MemTotal)%%20-%%20((node_memory_MemFree%%2Bnode_memory_Buffers%%2Bnode_memory_Cached)))%%20%%2F%%20node_memory_MemTotal%%20*%%20100&start=%v&end=%v&step=10", domain, port, StartTime, EndTime)
	}

	if contains(Metrics, "WriteTime") {
		MetricQuery["WriteTime"] = fmt.Sprintf("http://%v:%v/api/v1/query_range?query=irate(node_disk_sectors_written[5m])*512&start=%v&end=%v&step=10", domain, port, StartTime, EndTime)
	}

	if contains(Metrics, "ReadTime") {
		MetricQuery["ReadTime"] = fmt.Sprintf("http://%v:%v/api/v1/query_range?query=irate(node_disk_sectors_read[5m])*512&start=%v&end=%v&step=10", domain, port, StartTime, EndTime)
	}

	if contains(Metrics, "CPUUsage") {
		MetricQuery["CPUUsage"] = fmt.Sprintf("http://%v:%v/api/v1/query_range?query=100-(avg%%20by%%20(instance)%%20(irate(node_cpu{job='node-exporter',mode='idle'}[1m]))*100)&start=%v&end=%v&step=10", domain, port, StartTime, EndTime)
	}

	if contains(Metrics, "CPUIdle") {
		MetricQuery["CPUIdle"] = fmt.Sprintf("http://%v:%v/api/v1/query_range?query=avg(irate(node_cpu{job='node-exporter',mode='idle'}[1m]))*100&start=%v&end=%v&step=10", domain, port, StartTime, EndTime)
	}
	if contains(Metrics, "Knobs") {
		MetricQuery["Knobs"] = fmt.Sprintf("http://%v:%v/api/v1/query_range?query=app_knobs&start=%v&end=%v&step=10", domain, port, StartTime, EndTime)
	}
	return MetricQuery
}

func parseToPrometheusStruct(s []byte) RangePrometheusJSON {
	var p RangePrometheusJSON
	err := json.Unmarshal(s, &p)

	if err != nil {
		return p
	}

	// Cast the right value
	for res := range p.Data.Result {
		for val := range p.Data.Result[res].Values {
			p.Data.Result[res].Values[val].([]interface{})[0] = int(p.Data.Result[res].Values[val].([]interface{})[0].(float64))
		}
	}
	return p
}

func getFloat(unk interface{}) (float64, error) {
	var floatType = reflect.TypeOf(float64(0))
	var stringType = reflect.TypeOf("")

	switch i := unk.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	default:
		v := reflect.ValueOf(unk)
		v = reflect.Indirect(v)
		if v.Type().ConvertibleTo(floatType) {
			fv := v.Convert(floatType)
			return fv.Float(), nil
		} else if v.Type().ConvertibleTo(stringType) {
			sv := v.Convert(stringType)
			s := sv.String()
			return strconv.ParseFloat(s, 64)
		} else {
			return math.NaN(), fmt.Errorf("Can't convert %v to float64", v.Type())
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getBodyFromURL(MetricURL string) ([]byte, error) {

	resp, err := http.Get(MetricURL)

	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func getPrometheusValues(val interface{}) (int, float64) {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(val)
		tsValueRightType := s.Index(0).Interface().(int)
		featureValueRightType, _ := strconv.ParseFloat(s.Index(1).Interface().(string), 64)
		return tsValueRightType, featureValueRightType
	}
	return 0, 0.0
}

func copyDatasetFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func buildDataset(MetricStructs map[string]RangePrometheusJSON, isTrainingDataset bool) ([]string, map[int]interface{}) {

	// This map is a map[featureName]map[timestamp]featureValue
	var dataset = map[string]map[int]float64{}

	for featureName, data := range MetricStructs {
		var finalFeatureName string
		if len(data.Data.Result) > 1 {
			// Multiple
			for _, result := range data.Data.Result {

				if result.Metric.Name == "app_http_request_latency" {

					finalFeatureName = fmt.Sprintf("%v_%v_%v_%v", result.Metric.Name, result.Metric.Method, result.Metric.Endpoint, result.Metric.Quantile)

				} else if result.Metric.Name == "app_knobs" {

					finalFeatureName = fmt.Sprintf("%v_%v", result.Metric.Name, result.Metric.Knob)

				}
				dataset[finalFeatureName] = make(map[int]float64)

				for _, value := range result.Values {
					timestamp, featureValue := getPrometheusValues(value)
					dataset[finalFeatureName][timestamp] = featureValue
				}
			}
		} else {
			// Single
			for _, result := range data.Data.Result {

				dataset[featureName] = make(map[int]float64)

				for _, value := range result.Values {
					timestamp, featureValue := getPrometheusValues(value)
					dataset[featureName][timestamp] = featureValue
				}
			}
		}
	}

	validTimestamps := getTimestampsIntersection(dataset)

	// Remove timestamps that are not shared by all features
	for _, innerMap := range dataset {
		for timestamp := range innerMap {
			if !containsInt(validTimestamps, timestamp) {
				delete(innerMap, timestamp)
			}
		}
	}

	fileName := saveDataset(validTimestamps, dataset)

	finalFileName := transposeCSV(fileName)

	if isTrainingDataset {
		mergeDatasets(fileName)
	} else {
		os.Rename("/go/src/github.com/digorithm/meal_planner/finchgo/dataset/"+finalFileName, "/go/src/github.com/digorithm/meal_planner/finchgo/dataset/"+"single.csv")

		err := os.Remove(fileName)

		if err != nil {
			fmt.Printf("### error:: %v ### \n", err)
		}
	}

	return []string{}, nil

}

func saveDataset(validTimestamps []int, dataset map[string]map[int]float64) string {

	_, filename, _, _ := runtime.Caller(0)
	dir, _ := filepath.Split(filepath.Dir(filename))

	datasetPath := fmt.Sprintf(dir + "finchgo/dataset/")

	if _, err := os.Stat(datasetPath); os.IsNotExist(err) {
		os.Mkdir(datasetPath, 0777)
		fmt.Println("### Directory has been created ### \n")
	}

	fileName := uuid.Must(uuid.NewV4())
	fileNameString := fmt.Sprintf("%v/%v.csv", datasetPath, fileName)

	file, err := os.Create(fileNameString)

	if err != nil {
		logrus.Fatal(err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	keys := make([]string, 0, len(dataset))
	for k := range dataset {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, featureName := range keys {
		row := make([]string, 0)
		row = append(row, featureName)
		for _, timestamp := range validTimestamps {
			row = append(row, strconv.FormatFloat(dataset[featureName][timestamp], 'E', -1, 64))
		}
		writer.Write(row)
	}

	return fileNameString

}

func getTimestampsIntersection(dataset map[string]map[int]float64) []int {
	// This function returns a list of timestamps that are present in every feature of the dataset

	allTimestamps := make([][]int, 0)
	for _, innerMap := range dataset {
		timestamps := make([]int, 0)
		for timestamp := range innerMap {
			timestamps = append(timestamps, timestamp)
		}
		allTimestamps = append(allTimestamps, timestamps)
	}

	all := make([]int, 0)
	for _, timestamps := range allTimestamps {
		for _, timestamp := range timestamps {
			all = append(all, timestamp)
		}
	}

	all = removeDuplicates(all)

	timestampFrequencies := make(map[int]int)
	numberOfSlices := len(allTimestamps)

	for _, slice := range allTimestamps {
		for _, timestamp := range all {
			contains := containsInt(slice, timestamp)
			if contains {
				timestampFrequencies[timestamp] += 1
			}
		}
	}
	result := make([]int, 0)
	for timestamp, frequency := range timestampFrequencies {
		if frequency == numberOfSlices {
			result = append(result, timestamp)
		}
	}
	return result
}

func removeDuplicates(elements []int) []int {
	// Use map to record duplicates as we find them.
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func containsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func buildDatasetDeprecated(MetricStructs map[string]RangePrometheusJSON) ([]string, map[int]interface{}) {

	fmt.Println("MetricStructs is:: ")
	spew.Dump(MetricStructs)

	// Check if MetricStructs is good enough so that we can actually rewrite this whole shitty method.

	// BUG: Every endpoint must be called between the timerange. Otherwise we will have different values. Fix this later, as this isnt an important problem

	PrometheusStructs := []RangePrometheusJSON{MetricStructs["HTTPRequestCount"], MetricStructs["IOWait"], MetricStructs["MemoryUsage"], MetricStructs["WriteTime"], MetricStructs["CPUUsage"], MetricStructs["ReadTime"], MetricStructs["CPUIdle"], MetricStructs["Knobs"], MetricStructs["HTTPRequestLatency"]}

	// Check if they have the same number of samples
	NumberOfSamples := make([]int, 0)
	for _, s := range PrometheusStructs {
		for _, res := range s.Data.Result {
			NumberOfSamples = append(NumberOfSamples, len(res.Values))
			fmt.Println("Number of samples:: ", len(res.Values))
			fmt.Println("Feature:: ", res.Values)
			fmt.Println("Feature name:: ", res.Metric.Name, res.Metric.Method, res.Metric.Endpoint)
		}
	}

	if len(uniques(NumberOfSamples)) != 1 {
		logrus.Fatal("Number of Samples isn't the same, debug time!")
	}

	KeysToFeatureName := make(map[int]string)
	KeysToFeatureName[0] = "workload"
	KeysToFeatureName[1] = "io_wait"
	KeysToFeatureName[2] = "memory_usage"
	KeysToFeatureName[3] = "disk_write_bytes"
	KeysToFeatureName[4] = "cpu_usage"
	KeysToFeatureName[5] = "disk_read_bytes"
	KeysToFeatureName[6] = "cpu_idle"
	KeysToFeatureName[7] = "knobs"

	// Create a 1D slice that will hold feature names of the dataset
	var featureNames []string
	featureNames = append(featureNames, "timestamp")
	knobCount := 1
	// Note that it follows the order we built PrometheusStructs
	// Here we are just creating a slice that contains the feature names
	for k, s := range PrometheusStructs {
		for _, res := range s.Data.Result {
			if res.Metric.Name == "" {

				// fmt.Printf("Key %v Feature name is:: %v \n\n\n", k, KeysToFeatureName[k])
				featureNames = append(featureNames, KeysToFeatureName[k])
			} else {
				if res.Metric.Name == "app_http_request_latency" {

					// fmt.Printf("Key %v Feature name is:: %v \n\n\n", k, KeysToFeatureName[k])
					name := fmt.Sprintf(res.Metric.Name + "_" + res.Metric.Endpoint + "_" + res.Metric.Method + "_" + res.Metric.Quantile)
					featureNames = append(featureNames, name)
				} else if res.Metric.Name == "app_knobs" {
					name := fmt.Sprintf(res.Metric.Name + "_" + strconv.Itoa(knobCount))
					featureNames = append(featureNames, name)
					knobCount += 1
				} else {

					// fmt.Printf("Key %v Feature name is:: %v \n\n\n", k, KeysToFeatureName[k])
					// fmt.Printf("Feature name is:: %v \n\n\n", res.Metric.Name)
					featureNames = append(featureNames, res.Metric.Name)
				}
			}
		}
	}

	// Knobs is here. Seems correct.
	spew.Dump(featureNames)

	// This part is tricky, complex, and poorly written
	datasetStruct := make(map[int]interface{})

	firstIteration := true
	for k, s := range PrometheusStructs {
		for i, res := range s.Data.Result {
			if firstIteration {
				timestampsValues := make([]interface{}, 0)
				firstFeatureValues := make([]interface{}, 0)
				for _, val := range res.Values {

					switch reflect.TypeOf(val).Kind() {
					case reflect.Slice:
						s := reflect.ValueOf(val)
						tsValueRightType := s.Index(0).Interface().(int)
						firstFeatureValueRightType, _ := strconv.ParseFloat(s.Index(1).Interface().(string), 64)

						timestampsValues = append(timestampsValues, tsValueRightType)
						firstFeatureValues = append(firstFeatureValues, firstFeatureValueRightType)

					}
				}

				datasetStruct[0] = timestampsValues
				datasetStruct[1] = firstFeatureValues

				firstIteration = false
			} else {

				featureValues := make([]interface{}, 0)
				for _, val := range res.Values {

					switch reflect.TypeOf(val).Kind() {
					case reflect.Slice:
						s := reflect.ValueOf(val)
						featureValueRightType, _ := strconv.ParseFloat(s.Index(1).Interface().(string), 64)

						featureValues = append(featureValues, featureValueRightType)

					}
				}
				// Q: does featureValues contain all knob values? Is it supposed to? (check how it is with the other features)
				// fmt.Println("featureValues is:: ")
				// spew.Dump(featureValues)
				// Since we are adding indexes 0 and 1 first because of reasons,
				// here we either go for k + i (index for result) + 1 (because first index is 0) if the metric is latency, because it's the only metric that have multiple results... or we go for k+3, this 3 is because of reasons I don't quite understand. Changing the dataset might change the way we build it, unfortunately. But since this is an experiment, the features are already defined, so I'm not changing them for now.
				if PrometheusStructs[k].Data.Result[0].Metric.Name == "app_http_request_latency" {
					datasetStruct[k+i+1] = featureValues
				} else if PrometheusStructs[k].Data.Result[0].Metric.Name == "app_knobs" {
					datasetStruct[k+i+1] = featureValues
				} else {
					datasetStruct[k+1] = featureValues

				}
			}
		}
	}
	// if the data here is wrong, then the problem is in the previous block.
	// You should start by checking the slice featureValues
	fmt.Println("final datasetStruct is:: ")
	spew.Dump(datasetStruct)

	return featureNames, datasetStruct
}

func writeToCSV(dataset map[int]interface{}, featureNames []string, fileName string) {

	file, err := os.Create(fileName)

	if err != nil {
		logrus.Fatal(err)
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// fmt.Println("dataset: \n")
	// spew.Dump(dataset)
	// fmt.Println("featureNames: \n")
	// spew.Dump(featureNames)

	for k := range featureNames {
		// fmt.Printf("K is:: %v\n", k)
		// This will be written to the csv file
		featureString := make([]string, 0)
		// Append name of the feature
		featureString = append(featureString, featureNames[k])
		// fmt.Printf("featureString:: %v\n", featureString)
		// fmt.Printf("dataset in position k:: %v", dataset[k])
		switch reflect.TypeOf(dataset[k]).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(dataset[k])
			for i := 0; i < s.Len(); i++ {
				valueString := fmt.Sprintf("%v", s.Index(i))
				featureString = append(featureString, valueString)
			}
			writer.Write(featureString)
		}
	}
	/*
		// Sort the keys because of this (https://nathanleclaire.com/blog/2014/04/27/a-surprising-feature-of-golang-that-colored-me-impressed/)
		var keys []string
		for k := range knobs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Here we are manually setting the knobs as many time as the number of sample
		for _, k := range keys {
			featureString := make([]string, 0)
			featureString = append(featureString, k, fmt.Sprintf("%v", knobs[k]))
			switch reflect.TypeOf(dataset[2]).Kind() {
			case reflect.Slice:
				s := reflect.ValueOf(dataset[2])
				NumberOfSamples := s.Len()
				fmt.Println("Number of samples:: ", NumberOfSamples)
				for i := 1; i < NumberOfSamples; i++ {
					featureString = append(featureString, fmt.Sprintf("%v", knobs[k]))
				}
			}
			writer.Write(featureString)
		}
	*/
}

func transposeCSV(fileName string) string {
	csvFile, _ := os.Open(fileName)
	// Problem: fileName is the whole path...
	// We could split that into filename and the path

	path, f := filepath.Split(fileName)

	finalFileName := fmt.Sprintf("final_%v", f)

	file, err := os.Create(fmt.Sprintf("%v%v", path, finalFileName))

	if err != nil {
		logrus.Fatal(err)
	}

	defer file.Close()

	//writer := csv.NewWriter(file)
	//defer writer.Flush()

	err = transposeCsv(csvFile, file)
	if err != nil {
		logrus.Fatal(err)
	}

	return finalFileName
}

// This function will get previously generated datasets and merge them.
// We need this because in this experiment I will be running a simulation
// then stoping, generating the dataset, manually tweaking the knobs, and repeat. In the end I want a single dataset
func mergeDatasets(fileName string) {
	// First, copy previously merged dataset.csv being used to training to the dataset/ folder. Rename it to final_dataset.csv so that it can be merged with the new datasets being generated in runtime. All this IF such dataset exists

	dst := "src/github.com/digorithm/meal_planner/finchgo/dataset/final_dataset.csv"
	src := "src/github.com/digorithm/meal_planner/finchgo/machine_learning/dataset.csv"

	copyDatasetFile(src, dst)

	path, _ := filepath.Split(fileName)

	files, err := ioutil.ReadDir(path)
	fmt.Printf("All files:: %v\n", files)

	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Println("### Creating final dataset and merging all datasets to it ###")

	finalFile, err := os.Create(fmt.Sprintf("%vdataset.csv", path))

	if err != nil {
		logrus.Fatal(err)
	}

	writer := csv.NewWriter(finalFile)
	defer writer.Flush()

	firstIteration := true
	for _, f := range files {
		match, err := regexp.MatchString("^final_", f.Name())
		if err != nil {
			logrus.Fatal(err)
		}
		if match {
			fmt.Printf("File match:: %v\n", f.Name())
			csvFile, _ := os.Open(fmt.Sprintf("%v%v", path, f.Name()))
			reader := csv.NewReader(bufio.NewReader(csvFile))

			if firstIteration {
				fmt.Println("### Writing first dataset ###")
				fmt.Printf("Merging %v\n", f.Name())
				tempDataset, err := reader.ReadAll()
				if err != nil {
					logrus.Fatal(err)
				}
				writer.WriteAll(tempDataset)
				firstIteration = false
			} else {
				fmt.Println("### Merging another dataset ###")
				fmt.Printf("Merging %v\n", f.Name())
				tempDataset, err := reader.ReadAll()
				if err != nil {
					logrus.Fatal(err)
				}
				writer.WriteAll(tempDataset[1:])
			}
		}
	}
}

func saveDatasetDeprecated(dataset map[int]interface{}, featureNames []string) {

	_, filename, _, _ := runtime.Caller(0)
	dir, _ := filepath.Split(filepath.Dir(filename))

	datasetPath := fmt.Sprintf(dir + "finchgo/dataset/")

	fmt.Printf("### Dataset path is:: %v ### \n", datasetPath)

	if _, err := os.Stat(datasetPath); os.IsNotExist(err) {
		os.Mkdir(datasetPath, 0777)
		fmt.Println("### Directory has been created ### \n")
	}

	fileName := uuid.Must(uuid.NewV4())
	fileNameString := fmt.Sprintf("%v/%v.csv", datasetPath, fileName)

	fmt.Printf("### file name is %v (this might be wrong...) ### \n", fileNameString)

	writeToCSV(dataset, featureNames, fileNameString)
	transposeCSV(fileNameString)
	mergeDatasets(fileNameString)
}

func uniques(input []int) []int {
	u := make([]int, 0, len(input))
	m := make(map[int]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// Logic to transpose CSV
type readWriteSeekCloser interface {
	io.ReadWriteCloser
	io.Seeker
}

// FileBuffer saves the columns in temporary files.
type FileBuffer struct {
	names []string
	rwcs  []readWriteSeekCloser
	err   error

	size int
	sep  byte
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

// Remove removes all temporary files
func (b *FileBuffer) Remove() {
	for _, n := range b.names {
		if fileExists(n) {
			if err := os.Remove(n); err != nil {
				logrus.Fatalf("error removing '%s': %v", n, err)
			}
		}
	}
}

// WriteTo writes the content from the temporary files
// into the result file
func (b *FileBuffer) WriteTo(w io.Writer) (int64, error) {
	size := b.size
	if size == 0 {
		size = 32 * 1024
	}
	var sum int64
	buf := make([]byte, size)
	for _, r := range b.rwcs {
		_, err := r.Seek(0, io.SeekStart)
		if err != nil {
			return 0, err
		}
		n, err := io.CopyBuffer(w, r, buf)
		if err != nil {
			return n, err
		}
		sum += n

		i, err := w.Write([]byte("\n"))
		if err != nil {
			return int64(i), err
		}
		sum += int64(i)
	}
	return sum, nil
}

// Store stores the content of the csv.Reader in
// temporary files
func (b *FileBuffer) Store(r *csv.Reader) error {
	for {
		line, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		err = b.append(line)
		if err != nil {
			return err
		}
	}
}

func (b *FileBuffer) append(line []string) error {
	if len(b.rwcs) == 0 {
		for _ = range line {
			f, err := ioutil.TempFile("", "transposer")
			if err != nil {
				return err
			}
			b.names = append(b.names, f.Name())
			b.rwcs = append(b.rwcs, f)
		}
	} else if len(line) != len(b.rwcs) {
		return errors.New("")
	}

	for i, s := range line {
		_, err := b.rwcs[i].Write(append([]byte(s), b.sep))
		if err != nil {
			return err
		}
	}
	return nil
}

func transposeCsv(csvFile io.Reader, w io.Writer) error {
	r := csv.NewReader(csvFile)
	buf := &FileBuffer{
		size: 32 * 1024,
		sep:  byte(','),
	}
	defer buf.Remove()

	err := buf.Store(r)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	return err
}

func getTrainingAverage(scores []float64) float64 {
	var total float64 = 0
	for _, value := range scores {
		total += value
	}
	return total / float64(len(scores))
}
