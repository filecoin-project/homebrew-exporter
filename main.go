package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HomebrewMetricsItem struct {
  Number int `json:"number"`
  Formula string `json:"formula"`
  Count string `json:"count"`
  Percent string `json:"percent"`
}

type HomebrewMetrics struct {
  Category string `json:"category"`
  TotalItems int `json:"total_items"`
  StartDate string `json:"start_date"`
  EndDate string `json:"end_date"`
  TotalCount int `json:"total_count"`
  Items	[]HomebrewMetricsItem `json:"items"`
}

type HomebrewCollector struct {
  Formulae []string
  install30d *prometheus.Desc
  install90d *prometheus.Desc
  install365d *prometheus.Desc
  installOnRequest30d *prometheus.Desc
  installOnRequest90d *prometheus.Desc
  installOnRequest365d *prometheus.Desc
  buildError30d *prometheus.Desc
  buildError90d *prometheus.Desc
  buildError365d *prometheus.Desc
}

func newHomebrewCollector(formulae []string) *HomebrewCollector {
  return &HomebrewCollector{
    Formulae: formulae,
    install30d: prometheus.NewDesc("homebrew_install_30d",
      "Results from https://formulae.brew.sh/api/analytics/install/30d.json",
      []string{"formula"}, nil,
    ),
    install90d: prometheus.NewDesc("homebrew_install_90d",
      "Results from https://formulae.brew.sh/api/analytics/install/90d.json",
      []string{"formula"}, nil,
    ),
    install365d: prometheus.NewDesc("homebrew_install_365d",
      "Results from https://formulae.brew.sh/api/analytics/install/365d.json",
      []string{"formula"}, nil,
    ),
    installOnRequest30d: prometheus.NewDesc("homebrew_install_on_request_30d",
      "Results from https://formulae.brew.sh/api/analytics/install-on-request/30d.json",
      []string{"formula"}, nil,
    ),
    installOnRequest90d: prometheus.NewDesc("homebrew_install_on_request_90d",
      "Results from https://formulae.brew.sh/api/analytics/install-on-request/90d.json",
      []string{"formula"}, nil,
    ),
    installOnRequest365d: prometheus.NewDesc("homebrew_install_on_request_365d",
      "Results from https://formulae.brew.sh/api/analytics/install-on-request/365d.json",
      []string{"formula"}, nil,
    ),
    buildError30d: prometheus.NewDesc("homebrew_build_error_30d",
      "Results from https://formulae.brew.sh/api/analytics/build-error/30d.json",
      []string{"formula"}, nil,
    ),
    buildError90d: prometheus.NewDesc("homebrew_build_error_90d",
      "Results from https://formulae.brew.sh/api/analytics/build-error/90d.json",
      []string{"formula"}, nil,
    ),
    buildError365d: prometheus.NewDesc("homebrew_build_error_365d",
      "Results from https://formulae.brew.sh/api/analytics/build-error/365d.json",
      []string{"formula"}, nil,
    ),
  }
}

func (collector *HomebrewCollector) Describe(ch chan<- *prometheus.Desc) {

  //Update this section with the each metric you create for a given collector
  ch <- collector.install30d
  ch <- collector.install90d
  ch <- collector.install365d
  ch <- collector.installOnRequest30d
  ch <- collector.installOnRequest90d
  ch <- collector.installOnRequest365d
  ch <- collector.buildError30d
  ch <- collector.buildError90d
  ch <- collector.buildError365d
}

func (collector *HomebrewCollector) collectMetric(url string, metric *prometheus.Desc, ch chan<- prometheus.Metric) {
  client := http.Client{}
  homebrewMetrics := HomebrewMetrics{}
  req, err := http.NewRequest(http.MethodGet, url, nil)
  if err != nil { panic(err) }
  res, err := client .Do(req)
  if err != nil { panic(err) }
  body, err := ioutil.ReadAll(res.Body)
  if err != nil { panic(err) }
  json.Unmarshal(body, &homebrewMetrics)
  layout := "2006-01-02"
  endDate, err := time.Parse(layout, homebrewMetrics.EndDate)
  if err != nil { panic(err) }

  for _, formula := range collector.Formulae {
    for _, item := range homebrewMetrics.Items {
      if (item.Formula == formula) {
        countString := strings.Replace(item.Count, ",", "", -1)
        count, err := strconv.ParseFloat(countString, 32)
        if err != nil { panic(err) }
        m := prometheus.MustNewConstMetric(metric, prometheus.GaugeValue, count, item.Formula)
        m = prometheus.NewMetricWithTimestamp(endDate, m)
        ch <- m
      }
    }
  }
}

func (collector *HomebrewCollector) Collect(ch chan<- prometheus.Metric) {
  collector.collectMetric("https://formulae.brew.sh/api/analytics/install/30d.json", collector.install30d, ch)
  collector.collectMetric("https://formulae.brew.sh/api/analytics/install/90d.json", collector.install90d, ch)
  collector.collectMetric("https://formulae.brew.sh/api/analytics/install/365d.json", collector.install365d, ch)
  collector.collectMetric("https://formulae.brew.sh/api/analytics/install-on-request/30d.json", collector.installOnRequest30d, ch)
  collector.collectMetric("https://formulae.brew.sh/api/analytics/install-on-request/90d.json", collector.installOnRequest90d, ch)
  collector.collectMetric("https://formulae.brew.sh/api/analytics/install-on-request/365d.json", collector.installOnRequest365d, ch)
  collector.collectMetric("https://formulae.brew.sh/api/analytics/build-error/30d.json", collector.buildError30d, ch)
  collector.collectMetric("https://formulae.brew.sh/api/analytics/build-error/90d.json", collector.buildError90d, ch)
  collector.collectMetric("https://formulae.brew.sh/api/analytics/build-error/365d.json", collector.buildError365d, ch)
  }

func main() {
  listenPort := os.Getenv("LISTEN_PORT")
  if listenPort == "" { listenPort = "9888" }
  metricsPath := os.Getenv("METRICS_PATH")
  if metricsPath == "" { metricsPath= "/metrics" }
  formulaeString := os.Getenv("HOMEBREW_FORMULAE")
  if formulaeString != "" {
    formulae := strings.Split(formulaeString, ", ")
    log.Fatal(homebrewExporter(listenPort, metricsPath, formulae))
  }
}

func homebrewExporter(listenPort string, metricsPath string, formulae []string) error {
  homebrew := newHomebrewCollector(formulae)
  prometheus.MustRegister(homebrew)

    http.Handle(metricsPath, promhttp.Handler())
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`
            <html>
            <head><title>Homebrew Metrics Exporter</title></head>
            <body>
            <h1>Homebrew Exporter</h1>
            <p><a href='` + metricsPath + `'>Metrics</a></p>
            </body>
            </html>
        `))
    })

  return http.ListenAndServe(fmt.Sprintf(":%s", listenPort), nil)
}
