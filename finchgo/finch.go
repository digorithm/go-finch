package finchgo

import (
	"encoding/json"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NewFinch(KnobsPath, SLAsPath string) *Finch {
	Finch := &Finch{}
	Finch.KnobsPath = KnobsPath
	Finch.HTTPMonitorHandler = promhttp.Handler()

	Finch.SetupKnobsConfiguration(KnobsPath)
	Finch.SetupSLAs(SLAsPath)
	go Finch.Observe()
	// go Finch.DatasetBuilder
	// DatasetBuilder will build a dataset every x min and concatenate to the previously created datasets. It will also call ML component to train models

	return Finch
}

type Finch struct {
	KnobsPath          string
	Knobs              *viper.Viper
	HTTPMonitorHandler http.Handler
	RequestHistogram   *prometheus.HistogramVec
	HTTPRequestCount   *prometheus.CounterVec
	HTTPRequestLatency *prometheus.SummaryVec
	KnobsGauge         *prometheus.GaugeVec
	SLAs               []SLA
}

type SLA struct {
	SLA       string  `json:"sla"`
	Endpoint  string  `json:"endpoint"`
	Method    string  `json:"method"`
	Metric    string  `json:"metric"`
	Threshold float64 `json:"threshold"`
	Agreement float64 `json:"agreement"`
}

func (f *Finch) SetupSLAs(SLAsPath string) {

	// Setup correct file path
	_, filename, _, _ := runtime.Caller(0)
	dir, _ := filepath.Split(filepath.Dir(filename))
	file := filepath.Join(dir, SLAsPath)

	jsonFile, err := ioutil.ReadFile(file)

	if err != nil {
		logrus.Fatalf("Error reading sla file, %s", err)
	}
	var SLAs []SLA
	json.Unmarshal(jsonFile, &SLAs)

	f.SLAs = SLAs

	logrus.Info("SLAs loaded")

}

func (f *Finch) SetupKnobsConfiguration(KnobsPath string) {

	// Setup correct file path
	_, filename, _, _ := runtime.Caller(0)
	dir, _ := filepath.Split(filepath.Dir(filename))
	file := filepath.Join(dir, KnobsPath)

	// Setup viper
	viper := viper.New()
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file, %s", err)
	}

	logrus.Info("Knobs file loaded")

	// Watch for file changes
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	f.Knobs = viper
}

func (f *Finch) InitMonitoring() {

	var buckets []float64

	for _, sla := range f.SLAs {
		buckets = append(buckets, sla.Threshold)
	}
	sort.Float64s(buckets)

	f.RequestHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "app",
			Name:      "request_duration_mseconds",
			Help:      "A histogram of the API HTTP request durations in mseconds.",
			Buckets:   buckets,
		},
		[]string{"method", "endpoint"},
	)

	// Metrics to be monitored
	f.HTTPRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Name:      "http_request_count",
			Help:      "The number of HTTP requests.",
		},
		[]string{"method", "endpoint"},
	)

	f.HTTPRequestLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "app",
			Name:      "http_request_latency",
			Help:      "The latency of HTTP requests.",
		},
		[]string{"method", "endpoint"},
	)

	f.KnobsGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "app",
		Name:      "knobs",
		Help:      "Knobs metrics",
	},
		[]string{"knob"},
	)

	prometheus.MustRegister(f.HTTPRequestCount)
	prometheus.MustRegister(f.HTTPRequestLatency)
	prometheus.MustRegister(f.RequestHistogram)
	prometheus.MustRegister(f.KnobsGauge)
}

func (f *Finch) MonitorWorkload(Method, basePath string) {
	f.HTTPRequestCount.WithLabelValues(Method, basePath).Inc()
}

func (f *Finch) MonitorLatency(Method, basePath string, duration float64) {
	f.HTTPRequestLatency.WithLabelValues(Method, basePath).Observe(float64(duration))
	f.RequestHistogram.WithLabelValues(Method, basePath).Observe(float64(duration))
}

func (f *Finch) MonitorKnobs() {
	knobs := f.getCurrentKnobs()
	for knobName, value := range knobs {
		f.KnobsGauge.WithLabelValues(knobName).Set(float64(value))
	}
}

func (f *Finch) Observe() {
	ticker := time.NewTicker(5 * time.Second)
	quitWatcher := make(chan struct{})

	// Watch metrics
	go func() {
		for {
			select {
			case <-ticker.C:
				f.getSLAMetrics()
				fmt.Printf("Knobs are: %v\n", f.getCurrentKnobs())
			case <-quitWatcher:
				ticker.Stop()
				return
			}
		}
	}()

	// After developing/debugging, change to ~15min or so.
	tickerBuilder := time.NewTicker(10 * time.Minute)
	quitBuilder := make(chan struct{})

	// Build dataset periodically
	go func() {
		for {
			select {
			case <-tickerBuilder.C:
				f.DatasetBuilder(true, "-8m")
			case <-quitBuilder:
				ticker.Stop()
				return
			}
		}
	}()

	trainingMode := true
	// If in training mode, run this Goroutine to tweak knobs periodically

	if trainingMode {
		tickerTweakKnobs := time.NewTicker(5 * time.Minute)
		quitTweakKnobs := make(chan struct{})

		go func() {
			for {
				select {
				case <-tickerTweakKnobs.C:
					f.mutateKnobs()
				case <-quitTweakKnobs:
					ticker.Stop()
				}
			}
		}()
	}

}

func (f *Finch) mutateKnobs() {
	// Read json file containing the knobs
	KnobsPath := "finch_knobs.json"
	_, filename, _, _ := runtime.Caller(0)
	dir, _ := filepath.Split(filepath.Dir(filename))
	file := filepath.Join(dir, KnobsPath)

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
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

func (f *Finch) DatasetBuilder(isTrainingDataset bool, NegativeStartTime string) {

	/*
		NegativeStartTime should be a string determining how many minutes/seconds ago you want to start building the dataset. For example: "-5m" will create a dataset from 5min ago to now.

		This method will build the dataset based on a given time series range. If it is training dataset, it will save it as a csv file, and train the models on this dataset. Otherwise, it will just a single row that will be used for predictions.
	*/

	EndTime := time.Now()
	StartTimeDuration, _ := time.ParseDuration(NegativeStartTime)
	StartTime := EndTime.Add(StartTimeDuration)

	EndTimeString := strconv.FormatInt(EndTime.Unix(), 10)
	StartTimeString := strconv.FormatInt(StartTime.Unix(), 10)

	Metrics := []string{"HTTPRequestCount", "HTTPRequestLatency", "IOWait", "MemoryUsage", "WriteTime", "CPUUsage", "ReadTime", "CPUIdle", "Knobs"}

	MetricQueries := buildRangeQueries(Metrics, StartTimeString, EndTimeString)

	MetricData := extractFromPrometheus(MetricQueries)

	MetricStructs := buildStructs(MetricData)

	buildDataset(MetricStructs, isTrainingDataset)

	if isTrainingDataset {
		f.trainModels()
	}

}

func (f *Finch) trainModels() {

	src := "src/github.com/digorithm/meal_planner/finchgo/dataset/dataset.csv"
	dst := "src/github.com/digorithm/meal_planner/finchgo/machine_learning/dataset.csv"
	mlComponentDirectory := "src/github.com/digorithm/meal_planner/finchgo/machine_learning/"

	err := copyDatasetFile(src, dst)

	if err != nil {
		fmt.Printf("Copy failed: %v\n", err)
		logrus.Fatal(err)
	}

	cmd := exec.Command("python3", "train_models.py")
	cmd.Dir = mlComponentDirectory

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out))
	// TODO: Find the best way to expose the training report to the user
}

func (f *Finch) predictOptimalKnobs() (map[string]float64, bool) {

	var successfullyPredicted bool

	src := "src/github.com/digorithm/meal_planner/finchgo/dataset/single.csv"
	dst := "src/github.com/digorithm/meal_planner/finchgo/machine_learning/single.csv"
	mlComponentDirectory := "src/github.com/digorithm/meal_planner/finchgo/machine_learning/"

	err := copyDatasetFile(src, dst)

	if err != nil {
		fmt.Printf("Copy failed: %v\n", err)
		logrus.Fatal(err)
	}

	cmd := exec.Command("python3", "-u", "predict_optimal_knobs.py")
	cmd.Dir = mlComponentDirectory

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}

	var predictedKnobs interface{}
	err = json.Unmarshal(out, &predictedKnobs)

	if err != nil {
		fmt.Println(err)
		successfullyPredicted = false
		return nil, successfullyPredicted
	}

	successfullyPredicted = true
	// Cast type to map[string]float64

	fmt.Printf("Predicted knobs:: %v", predictedKnobs)
}

func (f *Finch) getCurrentKnobs() map[string]int {
	Knobs := make(map[string]int)

	Keys := f.Knobs.AllKeys()

	for _, key := range Keys {
		Knobs[key] = f.Knobs.GetInt(key)
	}

	return Knobs
}

func (f *Finch) getInstantResult(s []byte) float64 {
	var p SinglePrometheusJSON
	err := json.Unmarshal(s, &p)

	if err != nil {
		logrus.Fatal(err)
	}
	if len(p.Data.Result) == 1 {
		Value, err := getFloat(p.Data.Result[0].Value[1])
		if err != nil {
			logrus.Fatal(err)
		}
		return Value
	}

	// If no requests were made, then it's 1
	return 1.0

}

func (f *Finch) getSLAMetrics() {

	SLAQueries := f.buildInstantQueries()
	adapt := false

	for sla, query := range SLAQueries {
		body, err := getBodyFromURL(query)

		if err != nil {
			logrus.Fatalf("Couldn't get SLA metrics:: %v", err)
		}

		Value := f.getInstantResult(body)

		fmt.Printf("Metrics for SLA %s-%s:: %v \n \n", sla.Method, sla.Endpoint, (Value * 100.0))

		if (Value * 100.0) < float64(sla.Agreement) {
			fmt.Printf("### SLA %s-%s violated. Current value: %v \n\n", sla.Method, sla.Endpoint, Value)
			fmt.Printf("Initiating adaptation process\n")
			adapt = true
		}

		// Improve this part. If we already triggered adaptation process and we are seeing current improvements, show that to the user

		// If value is below agreement, call ML component, predict best set of knobs

	}

	if adapt {
		f.DatasetBuilder(false, "-30s")
		f.predictOptimalKnobs()
	}

}

func (f *Finch) buildInstantQueries() map[SLA]string {

	SLAQueries := make(map[SLA]string)

	genericPrometheusQuery := "http://prometheus:9090/api/v1/query?query=sum(rate(app_request_duration_mseconds_bucket{le=%%22%v%%22,%%20endpoint=%%22%v%%22,%%20method=%%22%v%%22}[5m]))%%20/%%20sum(rate(app_request_duration_mseconds_count{endpoint=%%22%v%%22,%%20method=%%22%v%%22}[5m]))"

	prometheusQueryAllEndpoints := "http://prometheus:9090/api/v1/query?query=sum(rate(app_request_duration_mseconds_bucket{le=%%22%v%%22,%%20method=%%22%v%%22}[5m]))%%20/%%20sum(rate(app_request_duration_mseconds_count{%%20method=%%22%v%%22}[5m]))"

	for _, sla := range f.SLAs {

		endpoint := strings.Replace(sla.Endpoint, "/", "", -1)

		if sla.Endpoint == "*" {

			prometheusSLAQuery := fmt.Sprintf(prometheusQueryAllEndpoints, sla.Threshold, sla.Method, sla.Method)

			SLAQueries[sla] = prometheusSLAQuery
		} else {

			prometheusSLAQuery := fmt.Sprintf(genericPrometheusQuery, sla.Threshold, endpoint, sla.Method, endpoint, sla.Method)

			SLAQueries[sla] = prometheusSLAQuery
		}
	}

	return SLAQueries
}
