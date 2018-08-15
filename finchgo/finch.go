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
	"math/rand"
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
	go Finch.StartMAPELoop()

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
	ViolatedSLAs       []SLA
	trainingMode       bool
	// SLAHasBeenOptimizedFor will tell us if we've predicted the knobs in order to optmize a given SLA
	SLAHasBeenOptimizedFor map[SLA]bool
	modelsHaveBeenTrained  bool
	knobHasBeenMutated     map[string]bool
	// ArtificialBlockingPoints holds knob -> true/false. This bool value determines if the configuration will affect the performance of this system proportially or inversely proportionally to its value
	artificialBlockingPoints       map[string]bool
	experimentOptimalConfiguration map[string]float64
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

	f.SLAHasBeenOptimizedFor = make(map[SLA]bool)

	for _, SLA := range SLAs {
		// Initialize SLAs optmization
		f.SLAHasBeenOptimizedFor[SLA] = false
	}

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
	allKnobs := viper.AllSettings()

	f.knobHasBeenMutated = make(map[string]bool)
	for knob := range allKnobs {
		f.knobHasBeenMutated[knob] = false
	}
	// Randomize initial configuration
	for knob := range allKnobs {
		var knobRandomValueString string

		knobNumberString := strings.Split(knob, "")[1]

		rand.Seed(time.Now().UnixNano())
		shouldAddZeroes := rand.Float32() < 0.5

		if shouldAddZeroes {
			knobRandomValueString = knobNumberString + string("000")
		} else {
			knobRandomValueString = knobNumberString
		}
		knobRandomValue, _ := strconv.ParseFloat(knobRandomValueString, 64)
		viper.Set(knob, knobRandomValue)
	}

	// Setup artificial blocking points
	randomBlockingPoints := make(map[string]bool)

	for knob := range allKnobs {
		// Randomly generate bool value
		rand.Seed(time.Now().UnixNano())
		randomBlockingPoints[knob] = rand.Float32() < 0.5
	}
	f.artificialBlockingPoints = randomBlockingPoints

	logrus.Info("Knobs file loaded")

	fmt.Printf("### Random blocking points:: %v ###\n### Current configuration after randomization:: %v ###\n", f.artificialBlockingPoints, viper.AllSettings())

	// Checking optimal configuration given randomization
	optimalConfiguration := make(map[string]float64)
	for knob, value := range allKnobs {
		valueString := strconv.FormatFloat(value.(float64), 'f', 0, 64)
		if f.artificialBlockingPoints[knob] {
			// it is proportional
			if len(valueString) == 1 {
				// Correct value, leave as is
				correctValue, _ := strconv.ParseFloat(valueString, 64)
				optimalConfiguration[knob] = correctValue
			} else {
				// Incorrect value, remove 3 zeroes
				correctValue, _ := strconv.ParseFloat(strings.TrimSuffix(valueString, "000"), 64)
				optimalConfiguration[knob] = correctValue
			}
		} else {
			// it is inversely proportional
			if len(valueString) == 1 {
				// Incorrect value, add 3 zeroes
				extra := "000"
				correctValue, _ := strconv.ParseFloat(valueString+string(extra), 64)
				optimalConfiguration[knob] = correctValue
			} else {
				// Correct value, leave as is
				correctValue, _ := strconv.ParseFloat(valueString, 64)
				optimalConfiguration[knob] = correctValue
			}
		}
	}
	fmt.Printf("### Optimal configuration:: %v ###\n", optimalConfiguration)
	f.experimentOptimalConfiguration = optimalConfiguration

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

func (f *Finch) StartMAPELoop() {
	f.trainingMode = true

	trainingAccuracyChannel := make(chan float64, 1)

	ticker := time.NewTicker(5 * time.Second)
	quitWatcher := make(chan struct{})
	// Watch metrics
	go f.MAPELoop(ticker, quitWatcher, trainingAccuracyChannel)

	// Build dataset every 15 minutes
	tickerBuilder := time.NewTicker(60 * time.Minute)
	quitBuilder := make(chan struct{})
	go f.contextBuilderLoop(tickerBuilder, quitBuilder, trainingAccuracyChannel)

	// If in training mode, run this Goroutine to tweak knobs periodically
	if f.trainingMode {
		tickerTweakKnobs := time.NewTicker(10 * time.Minute)
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

func (f *Finch) MAPELoop(ticker *time.Ticker, quitWatcher chan struct{}, trainingAccuracyChannel chan float64) {
	SLAMetricsHistory := make(map[SLA][]float64)

	// System states
	adaptationWasCarried := false
	isImproving := false
	f.modelsHaveBeenTrained = false

	currentTrainingAccuracy := 0.0

	for {
		select {
		case <-ticker.C:
			currentSLAMetrics := f.getSLAMetrics()
			SLAMetricsHistory = f.appendSLAMetricsHistory(SLAMetricsHistory, currentSLAMetrics)

			f.logSLAMetrics(currentSLAMetrics)

			if f.checkForViolation(currentSLAMetrics) {
				f.displayViolatedSLAs()

				if f.modelsHaveBeenTrained {
					// Only create adaptation plans if we have trained models
					fmt.Printf("### Models have been trained. Current accuracy:: %v ### \n", currentTrainingAccuracy)

					if adaptationWasCarried {
						// Analyze SLAMetricsHistory. If there's an improvement
						// set adaptationWasCarried to true, this will prevent the next if statement to execute the adaptation again
						// If there's no improvement, set adaptationWasCarried to false, this will make it try to adapt again
						isImproving = f.checkForImprovement(SLAMetricsHistory)
						fmt.Printf("### Improvement after deploying adaptation:: %v ###\n", isImproving)
						if !isImproving {
							adaptationWasCarried = false
						}
					} else {

						if !f.trainingMode {
							adaptationPlan := f.optimizeConfiguration()
							fmt.Printf("### Predicted optimal configuration:: %v ###\n", adaptationPlan)

							f.checkExperimentPrecision(adaptationPlan)

							f.carryAdaptationPlan(adaptationPlan)

							logrus.Infof("### Adaptation has been carried out ###")

							adaptationWasCarried = true
						}
						fmt.Println("### Training mode. Adaptations will happen only when models are trained ###")
					}
				} else {
					fmt.Println("### Waiting for models to be trained ### ")
				}
			} else {
				fmt.Println("### SLAs are in agreement ###")
			}

		case accuracy := <-trainingAccuracyChannel:
			// Received training result from contextBuilderLoop goroutine
			f.modelsHaveBeenTrained = true
			currentTrainingAccuracy = accuracy

		case <-quitWatcher:
			ticker.Stop()
			return
		}
	}
}

func (f *Finch) contextBuilderLoop(tickerBuilder *time.Ticker, quitBuilder chan struct{}, trainingAccuracyChannel chan float64) {
	for {
		select {
		case <-tickerBuilder.C:
			trainingAccuracy, success := f.DatasetBuilder(true, "-60m")
			if success {
				trainingAccuracyChannel <- trainingAccuracy
			}
		case <-quitBuilder:
			tickerBuilder.Stop()
			return
		}
	}
}

func (f *Finch) checkForImprovement(History map[SLA][]float64) bool {
	pastDataPoints := 5

	for _, values := range History {
		if len(values) < pastDataPoints+2 {
			// Not enough history
			return true
		}
		pastDataPointsValues := values[len(values)-pastDataPoints : len(values)-1]

		successiveDiffs := make([]float64, 0)

		for idx := range pastDataPointsValues {
			if idx != len(pastDataPointsValues)-1 {
				currentDiff := pastDataPointsValues[idx+1] - pastDataPointsValues[idx]
				successiveDiffs = append(successiveDiffs, currentDiff)
			}
		}
		// check if the sum successive diffs is positive
		sum := 0.0
		for _, value := range successiveDiffs {
			sum += value
		}

		// IF it's positive, that means there's an improvement going on
		if sum > 0.0 {
			return true
		}
	}

	return false
}

func (f *Finch) appendSLAMetricsHistory(History map[SLA][]float64, currentMetrics map[SLA]float64) map[SLA][]float64 {

	for sla, currentValue := range currentMetrics {
		History[sla] = append(History[sla], currentValue)
	}

	return History
}

func (f *Finch) logSLAMetrics(SLAMetrics map[SLA]float64) {
	for sla, value := range SLAMetrics {
		logrus.WithFields(logrus.Fields{
			"SLA":       sla.SLA,
			"Endpoint":  sla.Endpoint,
			"Method":    sla.Method,
			"Threshold": sla.Threshold,
			"Agreement": sla.Agreement,
		}).Infof("Current SLA value:: %v", value)
	}
}

func (f *Finch) mutateKnobs() {

	fmt.Printf("### %v ###\n", f.knobHasBeenMutated)
	for knob, wasMutated := range f.knobHasBeenMutated {
		if !wasMutated {
			fmt.Printf("### Mutating knob %v ###\n", knob)

			value := f.Knobs.GetFloat64(knob)
			s := strconv.FormatFloat(value, 'f', 0, 64)
			f.Knobs.Set(knob, value)
			if len(s) == 4 {
				newValue, _ := strconv.ParseFloat(strings.TrimSuffix(s, "000"), 64)
				f.Knobs.Set(knob, newValue)
			} else if len(s) == 1 {
				extra := "000"
				newValue, _ := strconv.ParseFloat(s+string(extra), 64)
				f.Knobs.Set(knob, newValue)
			}
			f.knobHasBeenMutated[knob] = true
			break
		}
	}
	if f.allKnobsHaveBeenMutated() {
		fmt.Println("### All knobs have been mutated ###\n")
		for knob, _ := range f.knobHasBeenMutated {
			f.knobHasBeenMutated[knob] = false
		}
	}
}

func (f *Finch) allKnobsHaveBeenMutated() bool {
	for _, wasMutated := range f.knobHasBeenMutated {
		if !wasMutated {
			return false
		}
	}
	return true
}

func (f *Finch) DatasetBuilder(isTrainingDataset bool, NegativeStartTime string) (float64, bool) {

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
		fmt.Println("### Initiating training process ###")
		trainingAverage, success := f.trainModels()
		if success {
			fmt.Printf("### Models training average accuracy:: %v ### \n", trainingAverage)

			if f.trainingMode {
				// If in training mode, adaptations will only happen when we create the dataset, in order to let parameter mutation do its job in peace
				adaptationPlan := f.optimizeConfiguration()
				fmt.Printf("### Predicted optimal configuration:: %v ###\n", adaptationPlan)

				f.checkExperimentPrecision(adaptationPlan)

				f.carryAdaptationPlan(adaptationPlan)

				logrus.Infof("### Adaptation has been carried out ###")
			}

			return trainingAverage, success
		}
	}
	return 0.0, false
}

func (f *Finch) trainModels() (float64, bool) {

	src := "src/github.com/digorithm/meal_planner/finchgo/dataset/dataset.csv"
	dst := "src/github.com/digorithm/meal_planner/finchgo/machine_learning/dataset.csv"
	mlComponentDirectory := "src/github.com/digorithm/meal_planner/finchgo/machine_learning/"

	err := copyDatasetFile(src, dst)

	if err != nil {
		fmt.Printf("Train model copy failed: %v\n", err)
		logrus.Fatal(err)
	}

	var trainingScore []float64
	trainingSuccessful := false

	cmd := exec.Command("python3", "train_models.py")
	cmd.Dir = mlComponentDirectory

	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(out))
	err = json.Unmarshal(out, &trainingScore)

	finalScore := getTrainingAverage(trainingScore)

	if err != nil {
		fmt.Println(err)
	}
	// If all good, set as successful
	trainingSuccessful = true

	return finalScore, trainingSuccessful

}

func (f *Finch) optimizeConfiguration() map[string]float64 {

	var optimalConfiguration map[string]float64

	optimalConfiguration = f.optimizeForAllSLAs()

	return optimalConfiguration
}

func (f *Finch) optimizeForSingleSLA(sla SLA) map[string]float64 {
	// Currently deprecated
	predictionSuccessful := false

	var predictedKnobs map[string]float64

	for !predictionSuccessful {
		fmt.Printf("### Optimizing for %v %v ###\n", sla.Method, sla.Endpoint)
		f.prepareSingleRow()

		SLIString := fmt.Sprintf("app_http_request_latency_%v_%v_0.99", sla.Method, strings.Replace(sla.Endpoint, "/", "", -1))

		cmd := exec.Command("python3", "-u", "predict_optimal_knobs.py", SLIString)

		mlComponentDirectory := "src/github.com/digorithm/meal_planner/finchgo/machine_learning/"
		cmd.Dir = mlComponentDirectory

		out, _ := cmd.CombinedOutput()

		fmt.Println(string(out))

		err := json.Unmarshal(out, &predictedKnobs)

		if err != nil || len(predictedKnobs) == 0 {
			fmt.Println(err)
			predictionSuccessful = false
		} else {
			predictionSuccessful = true
		}

		if !predictionSuccessful {
			fmt.Println("### Not enough data. Will try again in a few seconds ###")
		}

		time.Sleep(5 * time.Second)
	}
	f.SLAHasBeenOptimizedFor[sla] = true

	return predictedKnobs
}

func (f *Finch) optimizeForAllSLAs() map[string]float64 {

	predictionSuccessful := false

	var predictedKnobs map[string]float64

	for !predictionSuccessful {
		// This should happen only when we already optimized for each single SLA
		fmt.Printf("### Optimizing for all SLAs ###\n")

		f.prepareSingleRow()

		cmd := exec.Command("python3", "-u", "predict_optimal_knobs.py")

		mlComponentDirectory := "src/github.com/digorithm/meal_planner/finchgo/machine_learning/"
		cmd.Dir = mlComponentDirectory

		out, _ := cmd.CombinedOutput()

		err := json.Unmarshal(out, &predictedKnobs)

		if err != nil || len(predictedKnobs) == 0 {
			fmt.Println(err)
			predictionSuccessful = false
		} else {
			predictionSuccessful = true
		}

		if !predictionSuccessful {
			fmt.Println("### Not enough data. Will try again in a few seconds ###")
		}

		time.Sleep(5 * time.Second)
	}

	return predictedKnobs

}

func (f *Finch) prepareSingleRow() {
	f.DatasetBuilder(false, "-30s")
	src := "src/github.com/digorithm/meal_planner/finchgo/dataset/single.csv"
	dst := "src/github.com/digorithm/meal_planner/finchgo/machine_learning/single.csv"

	err := copyDatasetFile(src, dst)

	if err != nil {
		fmt.Printf("Predict knobs copy failed: %v\n", err)
		logrus.Fatal(err)
	}
}

func (f *Finch) shouldOptimizeSingleSLA() (bool, SLA) {

	// Currently deprecated

	// Optimizes for a single SLA if there is an SLA that hasn't been optimized for, otherwise, we optimize for all SLAs
	var emptySLA SLA

	for sla, hasBeenOptimized := range f.SLAHasBeenOptimizedFor {
		if !hasBeenOptimized {
			return true, sla
		}
	}
	return false, emptySLA
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

func (f *Finch) checkForViolation(SLAMetrics map[SLA]float64) bool {
	for sla, currentValue := range SLAMetrics {
		if currentValue < float64(sla.Agreement) {
			// Add to violated SLA if not there
			if len(f.ViolatedSLAs) == 0 {
				f.ViolatedSLAs = append(f.ViolatedSLAs, sla)
				return true
			} else {
				for _, violatedSLA := range f.ViolatedSLAs {
					if violatedSLA == sla {
						return true
					}
				}
				f.ViolatedSLAs = append(f.ViolatedSLAs, sla)
				return true
			}
		} else {
			// SLA isn't violated, if it's in the list, remove from it
			for idx, violatedSLA := range f.ViolatedSLAs {
				if sla == violatedSLA {
					f.ViolatedSLAs = removeSLAFromSlice(f.ViolatedSLAs, idx)
					break
				}
			}
		}
	}
	return false
}

func (f *Finch) displayViolatedSLAs() {
	for _, sla := range f.ViolatedSLAs {
		fmt.Printf("### Violated SLA:: %v %v ###\n", sla.Method, sla.Endpoint)
	}
}

func removeSLAFromSlice(s []SLA, i int) []SLA {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func (f *Finch) getSLAMetrics() map[SLA]float64 {

	SLAQueries := f.buildInstantQueries()
	SLAMetrics := make(map[SLA]float64)

	for sla, query := range SLAQueries {
		body, err := getBodyFromURL(query)

		if err != nil {
			logrus.Fatalf("Couldn't get SLA metrics:: %v", err)
		}

		currentValue := f.getInstantResult(body)
		SLAMetrics[sla] = currentValue * 100.0
	}
	return SLAMetrics
}

func (f *Finch) carryAdaptationPlan(predictedOptimalKnobs map[string]float64) {
	for knob, predictedValue := range predictedOptimalKnobs {
		f.Knobs.Set(knob, predictedValue)
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

func (f *Finch) ArtificialBlockingPoint(knob string) {
	if f.artificialBlockingPoints[knob] {
		// True means the blocking time is proportial to the knob value. The higher the value, the more it blocks
		sleepFor := time.Duration(f.Knobs.GetInt(knob)) * time.Millisecond
		time.Sleep(sleepFor)
	} else {
		// And the inverse is true
		sleepFor := time.Duration((1.0/f.Knobs.GetFloat64(knob))*10000) * time.Millisecond
		time.Sleep(sleepFor)
	}
}

func (f *Finch) checkExperimentPrecision(predictedKnobs map[string]float64) {
	for knob, value := range predictedKnobs {
		fmt.Printf("### Knob %v :: predicted = %v | correct = %v ###\n", knob, value, f.experimentOptimalConfiguration[knob])
	}
}
