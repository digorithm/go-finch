package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
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

func downloadData() (HTTPRequestCount, HTTPRequestLatency, IOWait, MemoryUsage, WriteTime, IOTime, CPUUsage, ReadTime, CPUIdle []byte) {

	// TODO: document this piece of shit

	UNIXTimeStart := "1515702937"
	UNIXTimeEnd := "1515703237"

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
	WriteTimeURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=node_disk_write_time_ms&start=%v&end=%v&step=2", UNIXTimeStart, UNIXTimeEnd)

	WriteTimeData, err := getBodyFromURL(WriteTimeURL)
	checkURLErr(err)

	// Ms spent doing IO
	IOTimeURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=node_disk_io_time_ms&start=%v&end=%v&step=2", UNIXTimeStart, UNIXTimeEnd)

	IOTimeData, err := getBodyFromURL(IOTimeURL)
	checkURLErr(err)

	// CPU usage
	CPUUsageURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=avg(irate(node_cpu{job='node-exporter',mode='user'}[1m]))*100&start=%v&end=%v&step=2", UNIXTimeStart, UNIXTimeEnd)

	CPUUsageData, err := getBodyFromURL(CPUUsageURL)
	checkURLErr(err)

	// Read time
	ReadTimeURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=node_disk_read_time_ms&start=%v&end=%v&step=2", UNIXTimeStart, UNIXTimeEnd)

	ReadTimeData, err := getBodyFromURL(ReadTimeURL)
	checkURLErr(err)

	// CPU idle
	CPUIdleURL := fmt.Sprintf("http://localhost:9090/api/v1/query_range?query=avg(irate(node_cpu{job='node-exporter',mode='idle'}[1m]))*100&start=%v&end=%v&step=2", UNIXTimeStart, UNIXTimeEnd)

	CPUIdleData, err := getBodyFromURL(CPUIdleURL)
	checkURLErr(err)

	return HTTPRequestCountData, HTTPRequestLatencyData, IOWaitData, MemoryUsageData, WriteTimeData, IOTimeData, CPUUsageData, ReadTimeData, CPUIdleData
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

func buildDataset(HTTPRequestCount, HTTPRequestLatency, IOWait, MemoryUsage, WriteTime, IOTime, CPUUsage, ReadTime, CPUIdle []byte) ([]string, map[int]interface{}) {

	AnalyzeData := true
	HTTPRequestLatencyStruct := parseToPrometheusStruct(HTTPRequestLatency)

	HTTPRequestCountStruct := parseToPrometheusStruct(HTTPRequestCount)

	IOWaitStruct := parseToPrometheusStruct(IOWait)

	MemoryUsageStruct := parseToPrometheusStruct(MemoryUsage)

	WriteTimeStruct := parseToPrometheusStruct(WriteTime)

	IOTimeStruct := parseToPrometheusStruct(IOTime)

	CPUUsageStruct := parseToPrometheusStruct(CPUUsage)

	ReadTimeStruct := parseToPrometheusStruct(ReadTime)

	CPUIdleStruct := parseToPrometheusStruct(CPUIdle)

	if AnalyzeData {
		fmt.Println(HTTPRequestLatencyStruct, "\n\n")
		fmt.Println(HTTPRequestCountStruct, "\n\n")
		fmt.Println(IOWaitStruct, "\n\n")
		fmt.Println(MemoryUsageStruct, "\n\n")
		fmt.Println(WriteTimeStruct, "\n\n")
		fmt.Println(IOTimeStruct, "\n\n")
		fmt.Println(CPUUsageStruct, "\n\n")
		fmt.Println(ReadTimeStruct, "\n\n")
		fmt.Println(CPUIdleStruct, "\n\n")
	}

	PrometheusStructs := []PrometheusJSON{HTTPRequestCountStruct, HTTPRequestLatencyStruct, IOWaitStruct, MemoryUsageStruct, WriteTimeStruct, IOTimeStruct, CPUUsageStruct, ReadTimeStruct, CPUIdleStruct}

	// Check if they have the same number of samples
	NumberOfSamples := make([]int, 0)
	for _, s := range PrometheusStructs {
		for _, res := range s.Data.Result {
			NumberOfSamples = append(NumberOfSamples, len(res.Values))
		}
	}
	if len(uniques(NumberOfSamples)) != 1 {
		log.Fatal("Number of Samples isn't the same, debug time!")
	}

	KeysToFeatureName := make(map[int]string)
	KeysToFeatureName[0] = "workload"
	KeysToFeatureName[2] = "io_wait"
	KeysToFeatureName[3] = "memory_usage"
	KeysToFeatureName[6] = "cpu_usage"
	KeysToFeatureName[8] = "cpu_idle"

	// Create a 1D slice that will hold feature names of the dataset
	var featureNames []string
	featureNames = append(featureNames, "timestamp")

	// Note that it follows the order we built PrometheusStructs
	// Here we are just creating a slice that contains the feature names
	for k, s := range PrometheusStructs {
		for _, res := range s.Data.Result {
			if res.Metric.Name == "" {
				featureNames = append(featureNames, KeysToFeatureName[k])
			} else {
				if res.Metric.Name == "app_http_request_latency" {
					name := fmt.Sprintf(res.Metric.Name + "_" + res.Metric.Endpoint + "_" + res.Metric.Method + "_" + res.Metric.Quantile)
					featureNames = append(featureNames, name)
				} else {
					featureNames = append(featureNames, res.Metric.Name)
				}
			}
		}
	}

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
					datasetStruct[k+3] = featureValues
				}
			}
		}
	}

	return featureNames, datasetStruct
}

func createAndSaveCSV(dataset map[int]interface{}, featureNames []string) {
	file, err := os.Create("result.csv")
	checkURLErr(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for k := range featureNames {

		featureString := make([]string, 0)
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
}

func main() {

	featureNames, dataset := buildDataset(downloadData())

	for k, name := range featureNames {
		fmt.Println(k)
		fmt.Printf("For feature %v, dataset is:: %v \n\n\n\n", name, dataset[k])
	}

	createAndSaveCSV(dataset, featureNames)

	csvFile, _ := os.Open("result.csv")

	file, err := os.Create("final_result.csv")
	checkURLErr(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	transposeCsv(csvFile, file)
}
