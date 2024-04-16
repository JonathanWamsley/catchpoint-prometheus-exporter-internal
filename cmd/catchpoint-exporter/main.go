// Copyright 2024 Grafana Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"catchpoint-prometheus-exporter/collector"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

var (
	webConfig    = webflag.AddFlags(kingpin.CommandLine, ":8080") // Change the default port if needed
	metricPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").Envar("CATCHPOINT_EXPORTER_WEB_TELEMETRY_PATH").String()
	bearerToken  = kingpin.Flag("bearer-token", "The bearer token for authentication with Catchpoint.").Envar("CATCHPOINT_EXPORTER_BEARER_TOKEN").Required().String()
	port         = kingpin.Flag("port", "The port to bind the HTTP server.").Default("8080").Envar("CATCHPOINT_EXPORTER_PORT").String()
	requestDelay = kingpin.Flag("request-delay", "Delay between API requests in seconds to manage rate limiting.").Default("1").Int()
	nodeIds      = kingpin.Flag("node-ids", "Comma-separated list of node IDs to include in the metrics.").Envar("CATCHPOINT_EXPORTER_NODE_IDS").String()
)

const (
	exporterName    = "catchpoint_exporter"
	landingPageHtml = `<html>
<head><title>Catchpoint Exporter</title></head>
<body>
	<h1>Catchpoint Exporter</h1>
	<p><a href='%s'>Metrics</a></p>
</body>
</html>`
)

func main() {
	kingpin.Version(version.Print(exporterName))
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(promlogConfig)

	nodeIdsSlice := strings.Split(*nodeIds, ",")
	nodeIdsIntSlice := make([]int, 0, len(nodeIdsSlice))
	for _, idStr := range nodeIdsSlice {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			level.Error(logger).Log("msg", "Invalid node ID", "id", idStr, "err", err)
			os.Exit(1)
		}
		nodeIdsIntSlice = append(nodeIdsIntSlice, id)
	}

	c := &collector.Config{
		BearerToken: *bearerToken,
		RateLimiter: collector.NewRateLimiter(*requestDelay),
		NodeIds:     nodeIdsIntSlice,
	}

	if err := c.Validate(); err != nil {
		level.Error(logger).Log("msg", "Configuration is invalid.", "err", err)
		os.Exit(1)
	}

	col := collector.NewCollector(logger, c)
	prometheus.MustRegister(col)

	serveMetrics(logger, *metricPath)
}

func serveMetrics(logger log.Logger, metricPath string) {
	landingPage := []byte(fmt.Sprintf(landingPageHtml, metricPath))

	http.Handle(metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.Write(landingPage)
	})

	srv := &http.Server{Addr: ":" + *port}
	if err := web.ListenAndServe(srv, webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error running HTTP server", "err", err)
		os.Exit(1)
	}
}
