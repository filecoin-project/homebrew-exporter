package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    listenAddress = flag.String("web.listen-address", ":9888", "Address to listen on for web interface.")
    metricPath    = flag.String("web.metrics-path", "/metrics", "Path under which to expose metrics.")
)

type homebrewCollector struct {
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

func newHomebrewCollector() *homebrewCollector {
  return &homebrewCollector{
    install30d: prometheus.NewDesc("homebrew.install.30d",
      "Results from https://formulae.brew.sh/api/analytics/install/30d.json",
      nil, nil,
    ),
    install90d: prometheus.NewDesc("homebrew.install.90d",
      "Results from https://formulae.brew.sh/api/analytics/install/90d.json",
      nil, nil,
    ),
    install365d: prometheus.NewDesc("homebrew.install.365d",
      "Results from https://formulae.brew.sh/api/analytics/install/365d.json",
      nil, nil,
    ),
    installOnRequest30d: prometheus.NewDesc("homebrew.install_on_request.30d",
      "Results from https://formulae.brew.sh/api/analytics/install-on-request/30d.json",
      nil, nil,
    ),
    installOnRequest90d: prometheus.NewDesc("homebrew.install_on_request.90d",
      "Results from https://formulae.brew.sh/api/analytics/install-on-request/90d.json",
      nil, nil,
    ),
    installOnRequest365d: prometheus.NewDesc("homebrew.install_on_request.365d",
      "Results from https://formulae.brew.sh/api/analytics/install-on-request/365d.json",
      nil, nil,
    ),
    buildError30d: prometheus.NewDesc("homebrew.build_error.30d",
      "Results from https://formulae.brew.sh/api/analytics/build-error/30d.json",
      nil, nil,
    ),
    buildError90d: prometheus.NewDesc("homebrew.build_error.90d",
      "Results from https://formulae.brew.sh/api/analytics/build-error/90d.json",
      nil, nil,
    ),
    buildError365d: prometheus.NewDesc("homebrew.build_error.365d",
      "Results from https://formulae.brew.sh/api/analytics/build-error/365d.json",
      nil, nil,
    ),
  }
}

func (collector *homebrewCollector) Describe(ch chan<- *prometheus.Desc) {

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

func (collector *homebrewCollector) Collect(ch chan<- prometheus.Metric) {

  // get homebrew metrics via http request

  //Write latest value for each metric in the prometheus metric channel.
  //Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
  m1 := prometheus.MustNewConstMetric(collector.install30d, prometheus.GaugeValue, metricValue)
  m2 := prometheus.MustNewConstMetric(collector.installOnRequest90d, prometheus.GaugeValue, metricValue)
  m2 := prometheus.MustNewConstMetric(collector.installOnRequest90d, prometheus.GaugeValue, metricValue)
  m1 = prometheus.NewMetricWithTimestamp(time.Now().Add(-time.Hour), m1)
  m2 = prometheus.NewMetricWithTimestamp(time.Now(), m2)
  ch <- m1
  ch <- m2
}

func main() {
  http.Handle("/console/metrics", promhttp.Handler())
  log.Fatal(http.ListenAndServe(":9101", nil))


    log.Fatal(homebrewExporter(*listenAddress, *metricPath))
}

func homebrewExporter(listenAddress, metricsPath string) error {
  homebrew := newHomebrewCollector()
  prometheus.MustRegister(homebrew)

    http.Handle(metricsPath, promhttp.Handler())
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`
            <html>
            <head><title>Homebrew Exporter Metrics</title></head>
            <body>
            <h1>Homebrew Exporter</h1>
            <p><a href='` + metricsPath + `'>Metrics</a></p>
            </body>
            </html>
        `))
    })

    return http.ListenAndServe(listenAddress, nil)
}