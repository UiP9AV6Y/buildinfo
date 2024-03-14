# buildinfo metrics

Helper library for creating a [Prometheus](github.com/prometheus/client_golang)
metrics collector using the embedded build information as data source.

```golang
package main

import (
  "log"
  "net/http"

  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"

  bicol "github.com/UiP9AV6Y/buildinfo/prometheus/collector"

  "example.com/version"
)

func main() {
  // Create a non-global registry.
  reg := prometheus.NewRegistry()

  // Create a collector and register it with the custom registry.
  if err := reg.Register(bicol.New("example", version.BuildInfo())); err != nil {
    log.Fatal(err)
  }

  // Expose metrics and custom registry via an HTTP server
  // using the HandleFor function. "/metrics" is the usual endpoint for that.
  http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
  log.Fatal(http.ListenAndServe(":8080", nil))
}
```
