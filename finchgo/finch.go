package finchgo

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
	"runtime"
)

func NewFinch(KnobsPath string) *Finch {
	Finch := &Finch{}
	Finch.KnobsPath = KnobsPath
	Finch.HTTPMonitorHandler = promhttp.Handler()

	Finch.SetupKnobsConfiguration(KnobsPath)

	return Finch
}

type Finch struct {
	KnobsPath          string
	Knobs              *viper.Viper
	HTTPMonitorHandler http.Handler
	RequestHistogram   *prometheus.HistogramVec
	HTTPRequestCount   *prometheus.CounterVec
	HTTPRequestLatency *prometheus.SummaryVec
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

	f.RequestHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "app",
			Name:      "request_duration_mseconds",
			Help:      "A histogram of the API HTTP request durations in mseconds.",
			Buckets:   []float64{100, 150, 200, 1000, 1500},
		},
		[]string{"method", "endpoint"},
	)
	fmt.Println("I'm here")

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

	prometheus.MustRegister(f.HTTPRequestCount)
	prometheus.MustRegister(f.HTTPRequestLatency)
	prometheus.MustRegister(f.RequestHistogram)
}

func (f *Finch) MonitorWorkload(Method, basePath string) {
	f.HTTPRequestCount.WithLabelValues(Method, basePath).Inc()
}

func (f *Finch) MonitorLatency(Method, basePath string, duration float64) {
	f.HTTPRequestLatency.WithLabelValues(Method, basePath).Observe(float64(duration))
	f.RequestHistogram.WithLabelValues(Method, basePath).Observe(float64(duration))
}
