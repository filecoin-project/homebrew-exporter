# homebrew-exporter
A Prometheus Exporter for parsing public homebrew metrics at https://formulae.brew.sh/analytics/

### To configure

| ENV variable | Default value | Description |
|--------------|---------------|-------------|
| `METRICS_PATH` | `"/metrics"`| The path to publish the metrics to. |
| `LISTEN_PORT`  | `"9888"`    | The port the metrics exporter listens on. |
| `HOMEBREW_FORMULAE` | REQUIRED | The list of formulae to grab metrics for. If blank, the exporter will exit immediately. |
