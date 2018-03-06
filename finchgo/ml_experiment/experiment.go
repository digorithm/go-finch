package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
)

type Value struct {
	v []interface{}
}

type Metric struct {
	Name     string `json:"__name__"`
	Device   string `json:"device"`
	Instance string `json:"instance"`
	Job      string `json:"job"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
	Quantile string `json:"quantile"`
}

type Result struct {
	Metric Metric        `json:"metric"`
	Values []interface{} `json:"values"`
}

type PrometheusData struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type PrometheusJSON struct {
	Status string         `json:"status"`
	Data   PrometheusData `json:"data"`
}

func checkURLErr(err error) {
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}

func getBodyFromURL(MetricURL string) ([]byte, error) {

	resp, err := http.Get(MetricURL)

	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func downloadData() (HTTPRequestCount, HTTPRequestLatency, IOWait, MemoryUsage, WriteTime, CPUUsage, ReadTime, CPUIdle []byte) {

	// TODO: document this piece of shit

	UNIXTimeStart := "1519758804"
	UNIXTimeEnd := "1519759620"

	// HTTP request count
	HTTPRequestCountURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=sum(irate(app_http_request_count[1m]))&start=%v&end=%v&step=2", UNIXTimeStart, UNIXTimeEnd)

	HTTPRequestCountData, err := getBodyFromURL(HTTPRequestCountURL)
	checkURLErr(err)

	// HTTP request latency
	HTTPRequestLatencyURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=app_http_request_latency&start=%s&end=%s&step=2", UNIXTimeStart, UNIXTimeEnd)

	HTTPRequestLatencyData, err := getBodyFromURL(HTTPRequestLatencyURL)
	checkURLErr(err)

	// CPU IO wait
	IOWaitURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=avg(irate(node_cpu{job='node-exporter',mode='iowait'}[1m]))*100&start=%v&end=%v&step=2", UNIXTimeStart, UNIXTimeEnd)

	IOWaitData, err := getBodyFromURL(IOWaitURL)
	checkURLErr(err)

	// Memory usage
	// Note that if URLs contain %something you have to double the % so it can parse/format as a literal value. https://stackoverflow.com/questions/35681595/escape-variables-with-printf-golang
	MemoryUsageURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=((node_memory_MemTotal)%%20-%%20((node_memory_MemFree%%2Bnode_memory_Buffers%%2Bnode_memory_Cached)))%%20%%2F%%20node_memory_MemTotal%%20*%%20100&start=%v&end=%v&step=2", UNIXTimeStart, UNIXTimeEnd)

	MemoryUsageData, err := getBodyFromURL(MemoryUsageURL)
	checkURLErr(err)

	// Write time
	WriteTimeURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=irate(node_disk_sectors_written[5m])*512&start=%v&end=%v&step=2", UNIXTimeStart, UNIXTimeEnd)

	WriteTimeData, err := getBodyFromURL(WriteTimeURL)
	checkURLErr(err)

	// CPU usage
	CPUUsageURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=100-(avg%%20by%%20(instance)%%20(irate(node_cpu{job='node-exporter',mode='idle'}[1m]))*100)&start=%v&end=%v&step=2", UNIXTimeStart, UNIXTimeEnd)

	CPUUsageData, err := getBodyFromURL(CPUUsageURL)
	checkURLErr(err)

	// Read time
	ReadTimeURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=irate(node_disk_sectors_read[5m])*512&start=%v&end=%v&step=2", UNIXTimeStart, UNIXTimeEnd)

	ReadTimeData, err := getBodyFromURL(ReadTimeURL)
	checkURLErr(err)

	// CPU idle
	CPUIdleURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=avg(irate(node_cpu{job='node-exporter',mode='idle'}[1m]))*100&start=%v&end=%v&step=2", UNIXTimeStart, UNIXTimeEnd)

	CPUIdleData, err := getBodyFromURL(CPUIdleURL)
	checkURLErr(err)

	return HTTPRequestCountData, HTTPRequestLatencyData, IOWaitData, MemoryUsageData, WriteTimeData, CPUUsageData, ReadTimeData, CPUIdleData
}

func parseToPrometheusStruct(s []byte) PrometheusJSON {
	var p PrometheusJSON
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

func buildDataset(HTTPRequestCount, HTTPRequestLatency, IOWait, MemoryUsage, WriteTime, CPUUsage, ReadTime, CPUIdle []byte) ([]string, map[int]interface{}) {

	// BUG: Every endpoint must be called between the timerange. Otherwise we will have different values. Fix this later, as this isnt an important problem

	AnalyzeData := true
	HTTPRequestLatencyStruct := parseToPrometheusStruct(HTTPRequestLatency)

	HTTPRequestCountStruct := parseToPrometheusStruct(HTTPRequestCount)

	IOWaitStruct := parseToPrometheusStruct(IOWait)

	MemoryUsageStruct := parseToPrometheusStruct(MemoryUsage)

	WriteTimeStruct := parseToPrometheusStruct(WriteTime)

	CPUUsageStruct := parseToPrometheusStruct(CPUUsage)

	ReadTimeStruct := parseToPrometheusStruct(ReadTime)

	CPUIdleStruct := parseToPrometheusStruct(CPUIdle)

	if AnalyzeData {
		fmt.Println(HTTPRequestLatencyStruct, "\n\n")
		fmt.Println(HTTPRequestCountStruct, "\n\n")
		fmt.Println(IOWaitStruct, "\n\n")
		fmt.Println(MemoryUsageStruct, "\n\n")
		fmt.Println(WriteTimeStruct, "\n\n")
		fmt.Println(CPUUsageStruct, "\n\n")
		fmt.Println(ReadTimeStruct, "\n\n")
		fmt.Println(CPUIdleStruct, "\n\n")
	}

	PrometheusStructs := []PrometheusJSON{HTTPRequestCountStruct, IOWaitStruct, MemoryUsageStruct, WriteTimeStruct, CPUUsageStruct, ReadTimeStruct, CPUIdleStruct, HTTPRequestLatencyStruct}

	// Check if they have the same number of samples
	NumberOfSamples := make([]int, 0)
	for _, s := range PrometheusStructs {
		for _, res := range s.Data.Result {
			NumberOfSamples = append(NumberOfSamples, len(res.Values))
			// fmt.Println("Number of samples:: ", len(res.Values))
			// fmt.Println("Feature:: ", res.Values)
			// fmt.Println("Feature name:: ", res.Metric.Name, res.Metric.Method, res.Metric.Endpoint)
		}
	}
	if len(uniques(NumberOfSamples)) != 1 {
		log.Fatal("Number of Samples isn't the same, debug time!")
	}

	KeysToFeatureName := make(map[int]string)
	KeysToFeatureName[0] = "workload"
	KeysToFeatureName[1] = "io_wait"
	KeysToFeatureName[2] = "memory_usage"
	KeysToFeatureName[3] = "disk_write_bytes"
	KeysToFeatureName[4] = "cpu_usage"
	KeysToFeatureName[5] = "disk_read_bytes"
	KeysToFeatureName[6] = "cpu_idle"

	// Create a 1D slice that will hold feature names of the dataset
	var featureNames []string
	featureNames = append(featureNames, "timestamp")

	// Note that it follows the order we built PrometheusStructs
	// Here we are just creating a slice that contains the feature names
	for k, s := range PrometheusStructs {
		for _, res := range s.Data.Result {
			if res.Metric.Name == "" {

				fmt.Printf("Key %v Feature name is:: %v \n\n\n", k, KeysToFeatureName[k])
				featureNames = append(featureNames, KeysToFeatureName[k])
			} else {
				if res.Metric.Name == "app_http_request_latency" {

					fmt.Printf("Key %v Feature name is:: %v \n\n\n", k, KeysToFeatureName[k])
					name := fmt.Sprintf(res.Metric.Name + "_" + res.Metric.Endpoint + "_" + res.Metric.Method + "_" + res.Metric.Quantile)
					featureNames = append(featureNames, name)
				} else {

					fmt.Printf("Key %v Feature name is:: %v \n\n\n", k, KeysToFeatureName[k])
					fmt.Printf("Feature name is:: %v \n\n\n", res.Metric.Name)
					featureNames = append(featureNames, res.Metric.Name)
				}
			}
		}
	}
	fmt.Println("Feature names: ", featureNames)
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
				// Since we are adding indexes 0 and 1 first because of reasons,
				// here we either go for k + i (index for result) + 1 (because first index is 0) if the metric is latency, because it's the only metric that have multiple results... or we go for k+3, this 3 is because of reasons I don't quite understand. Changing the dataset might change the way we build it, unfortunately. But since this is an experiment, the features are already defined, so I'm not changing them for now.
				if PrometheusStructs[k].Data.Result[0].Metric.Name == "app_http_request_latency" {
					datasetStruct[k+i+1] = featureValues
				} else {
					datasetStruct[k+1] = featureValues
				}
			}
		}
	}

	return featureNames, datasetStruct
}

// This function will get previously generated datasets and merge them.
// We need this because in this experiment I will be running a simulation
// then stoping, generating the dataset, manually tweaking the knobs, and repeat. In the end I want a single dataset
func mergeDatasets() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("### Creating final dataset and merging all datasets to it ###")
	finalFile, err := os.Create("dataset.csv")
	checkURLErr(err)

	writer := csv.NewWriter(finalFile)
	defer writer.Flush()

	firstIteration := true
	for _, f := range files {
		match, err := regexp.MatchString("^final_", f.Name())
		checkURLErr(err)
		if match {
			csvFile, _ := os.Open(f.Name())
			reader := csv.NewReader(bufio.NewReader(csvFile))
			if firstIteration {
				fmt.Println("### Writing first dataset ###")
				tempDataset, err := reader.ReadAll()
				checkURLErr(err)
				writer.WriteAll(tempDataset)
				firstIteration = false
			} else {
				fmt.Println("### Merging another dataset ###")
				tempDataset, err := reader.ReadAll()
				checkURLErr(err)
				writer.WriteAll(tempDataset[1:])
			}
		}
	}
}

func transposeCSV(fileName string) {
	fmt.Println("filename:: ", fileName)
	csvFile, _ := os.Open(fileName)

	file, err := os.Create(fmt.Sprintf("final_%v", fileName))
	checkURLErr(err)
	defer file.Close()

	//writer := csv.NewWriter(file)
	//defer writer.Flush()

	err = transposeCsv(csvFile, file)
	if err != nil {
		log.Fatal(err)
	}
}

func writeToCSV(knobs map[string]float64, dataset map[int]interface{}, featureNames []string, fileName string) {

	file, err := os.Create(fileName)
	checkURLErr(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for k := range featureNames {

		// This will be written to the csv file
		featureString := make([]string, 0)
		// Append name of the feature
		featureString = append(featureString, featureNames[k])

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
}

func createAndSaveCSV(knobs map[string]float64, dataset map[int]interface{}, featureNames []string) {

	fileName := uuid.Must(uuid.NewV4())
	fileNameString := fmt.Sprintf("%v.csv", fileName)

	writeToCSV(knobs, dataset, featureNames, fileNameString)
	transposeCSV(fileNameString)
}

func main() {

	featureNames, dataset := buildDataset(downloadData())

	knobs := make(map[string]float64)

	knobs["pg_shared_buffers_mb"] = 128
	knobs["pg_effective_cache_size_mb"] = 128
	knobs["pg_work_mem_kb"] = 1024
	knobs["pg_wal_buffers_mb"] = 128
	knobs["pg_checkpoint_completion_target"] = 0.7
	knobs["pg_maintenance_work_mem_mb"] = 16
	knobs["pg_default_statistics_target"] = 100
	knobs["pg_random_page_cost"] = 4
	knobs["pg_max_wal_size_gb"] = 2
	knobs["pg_min_wal_size_mb"] = 1000

	createAndSaveCSV(knobs, dataset, featureNames)
	mergeDatasets()

}
