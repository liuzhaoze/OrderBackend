package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
)

type PrometheusClient struct {
	registry     *prometheus.Registry
	counterVec   *prometheus.CounterVec
	histogramVec *prometheus.HistogramVec
}

func (p *PrometheusClient) CountCall(name string, status string) {
	p.counterVec.With(prometheus.Labels{"name": name, "status": status}).Inc()
}

func (p *PrometheusClient) RecordTime(name string, value float64) {
	p.histogramVec.WithLabelValues(name).Observe(value)
}

var prometheusClient *PrometheusClient

func NewPrometheusClient(metricsExportHost, metricsExportPort string, counterName, histogramName string) {
	prometheusClient = &PrometheusClient{
		registry: prometheus.NewRegistry(),
		counterVec: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: counterName,
			Help: "",
		}, []string{"name", "status"}),
		histogramVec: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    histogramName,
			Help:    "",
			Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024},
		}, []string{"name"}),
	}

	prometheusClient.registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	prometheusClient.registry.MustRegister(
		prometheusClient.counterVec,
		prometheusClient.histogramVec,
	)

	http.Handle("/metrics", promhttp.HandlerFor(prometheusClient.registry, promhttp.HandlerOpts{Registry: prometheusClient.registry}))

	go func() {
		logrus.Fatalf(
			"failed to start prometheus metrics endpoint: %v",
			http.ListenAndServe(fmt.Sprintf("%s:%s", metricsExportHost, metricsExportPort), nil),
		)
	}()
}

func GetPrometheusClient() *PrometheusClient {
	if prometheusClient == nil {
		logrus.Panicln("prometheus client is not initialized")
	}
	return prometheusClient
}
