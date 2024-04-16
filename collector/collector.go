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

package collector

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure Collector implements prometheus.Collector
var _ prometheus.Collector = (*Collector)(nil)

// Collector represents a collector for Catchpoint metrics.
type Collector struct {
	logger                  log.Logger
	Config                  *Config
	client                  CatchpointClientInterface
	nodeStatus              *prometheus.Desc
	slaPurgeItemsCountDesc  *prometheus.Desc
	testErrorTotalCount     *prometheus.Desc
	testErrorByType         *prometheus.Desc
	testErrorByIPCount      *prometheus.Desc
	testAlertsCriticalTotal *prometheus.Desc
	testAlertsWarningTotal  *prometheus.Desc
	testUsagePercentage     *prometheus.Desc
	testDownTimePercentage  *prometheus.Desc
	testRequestSlippage     *prometheus.Desc
	testRunRate             *prometheus.Desc
	totalTestRunsCount      *prometheus.Desc
	uniqueTestRunsCount     *prometheus.Desc
	up                      *prometheus.Desc
}

// Metric names
const (
	NodeStatusMetric              = "catchpoint_node_status"
	SLAPurgeItemsCountMetric      = "catchpoint_sla_purge_items_count"
	TestErrorTotalCountMetric     = "catchpoint_test_error_total_count"
	TestErrorByTypeMetric         = "catchpoint_test_error_by_type_count"
	TestErrorByIPCountMetric      = "catchpoint_test_error_by_ip_count"
	TestAlertsCriticalTotalMetric = "catchpoint_test_alerts_critical_count"
	TestAlertsWarningTotalMetric  = "catchpoint_test_alerts_warning_count"
	TestUsagePercentageMetric     = "catchpoint_usage_percentage"
	TestDownTimePercentageMetric  = "catchpoint_downtime_percentage"
	TestRequestSlippageMetric     = "catchpoint_node_request_slippage"
	TestRunRateMetric             = "catchpoint_node_run_rate"
	TotalTestRunsCountMetric      = "catchpoint_total_test_runs_count"
	UniqueTestRunsCountMetric     = "catchpoint_unique_test_runs_count"
	UpMetric                      = "catchpoint_up"
)

// Metric descriptions
const (
	NodeStatusDesc              = "The operational status of a Catchpoint node (1 for active, 0 for inactive)."
	SLAPurgeItemsCountDesc      = "Count of SLA purge items by status."
	TestErrorTotalCountDesc     = "The total count of all test errors."
	TestErrorByTypeDesc         = "Count of errors segmented by error type."
	TestErrorByIPCountDesc      = "Error counts traced back to specific IP addresses."
	TestAlertsCriticalTotalDesc = "Total number of critical alerts by test."
	TestAlertsWarningTotalDesc  = "Total number of warning alerts by test."
	TestUsagePercentageDesc     = "Usage percentage of test runs on a node"
	TestDownTimePercentageDesc  = "Downtime percentage of test runs on a node"
	TestRequestSlippageDesc     = "The slippage in test requests timings on a node, showing delays in scheduled test executions"
	TestRunRateDesc             = "The rate of which test runs are successfully completed on a node"
	TotalTestRunsCountDesc      = "The total number of test runs on a node"
	UniqueTestRunsCountDesc     = "The number of unique test runs on a node"
	UpDesc                      = "Indicates whether the last scrape of metrics from Catchpoint was successful."
)

// Label names
var (
	NodeLabels      = []string{"node_id", "node_name"}
	TestLabels      = []string{"test_id", "test_name", "node_id", "node_name"}
	TestRunLabels   = []string{"node_id", "test_name", "monitor_group"}
	StatusLabels    = []string{"status_id"}
	ErrorTypeLabels = []string{"error_type"}
	IPAddressLabels = []string{"ip"}
)

// NewCollector creates a new Collector.
func NewCollector(logger log.Logger, cfg *Config) *Collector {
	client := NewCatchpointClient(cfg.BearerToken)
	return &Collector{
		logger:                  logger,
		client:                  client,
		Config:                  cfg,
		nodeStatus:              prometheus.NewDesc(NodeStatusMetric, NodeStatusDesc, NodeLabels, nil),
		slaPurgeItemsCountDesc:  prometheus.NewDesc(SLAPurgeItemsCountMetric, SLAPurgeItemsCountDesc, StatusLabels, nil),
		testErrorTotalCount:     prometheus.NewDesc(TestErrorTotalCountMetric, TestErrorTotalCountDesc, nil, nil),
		testErrorByType:         prometheus.NewDesc(TestErrorByTypeMetric, TestErrorByTypeDesc, ErrorTypeLabels, nil),
		testErrorByIPCount:      prometheus.NewDesc(TestErrorByIPCountMetric, TestErrorByIPCountDesc, IPAddressLabels, nil),
		testAlertsCriticalTotal: prometheus.NewDesc(TestAlertsCriticalTotalMetric, TestAlertsCriticalTotalDesc, TestLabels, nil),
		testAlertsWarningTotal:  prometheus.NewDesc(TestAlertsWarningTotalMetric, TestAlertsWarningTotalDesc, TestLabels, nil),
		testUsagePercentage:     prometheus.NewDesc(TestUsagePercentageMetric, TestUsagePercentageDesc, TestRunLabels, nil),
		testDownTimePercentage:  prometheus.NewDesc(TestDownTimePercentageMetric, TestDownTimePercentageDesc, TestRunLabels, nil),
		testRequestSlippage:     prometheus.NewDesc(TestRequestSlippageMetric, TestRequestSlippageDesc, NodeLabels, nil),
		testRunRate:             prometheus.NewDesc(TestRunRateMetric, TestRunRateDesc, NodeLabels, nil),
		totalTestRunsCount:      prometheus.NewDesc(TotalTestRunsCountMetric, TotalTestRunsCountDesc, NodeLabels, nil),
		uniqueTestRunsCount:     prometheus.NewDesc(UniqueTestRunsCountMetric, UniqueTestRunsCountDesc, NodeLabels, nil),
		up:                      prometheus.NewDesc(UpMetric, UpDesc, nil, nil),
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.nodeStatus
	ch <- c.slaPurgeItemsCountDesc
	ch <- c.testErrorTotalCount
	ch <- c.testErrorByType
	ch <- c.testErrorByIPCount
	ch <- c.testAlertsCriticalTotal
	ch <- c.testAlertsWarningTotal
	ch <- c.testUsagePercentage
	ch <- c.testDownTimePercentage
	ch <- c.testRequestSlippage
	ch <- c.testRunRate
	ch <- c.totalTestRunsCount
	ch <- c.uniqueTestRunsCount
	ch <- c.up
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	// Check for nil configuration or client
	if c.Config == nil || c.client == nil {
		c.logger.Log("msg", "Collector config or client is nil")
		return
	}

	c.logger.Log("msg", "Starting collection")

	up := 1.0

	for _, nodeId := range c.Config.NodeIds {
		nodeStatusResponse, err := c.client.FetchNodeStatus(strconv.Itoa(nodeId))
		if err != nil {
			c.logger.Log("msg", "Failed to fetch node status", "node_id", nodeId, "error", err.Error())
			continue
		}

		if nodeStatusResponse.Data != nil && nodeStatusResponse.Data.Nodes != nil {
			for _, node := range *nodeStatusResponse.Data.Nodes { // Dereference the pointer to range over the slice
				statusValue := 0
				if node.Status != nil && node.Status.Name == "active" { // Check if Status is not nil before accessing Name
					statusValue = 1
				}

				// Use the node ID and name safely since we're now sure the necessary data exists
				ch <- prometheus.MustNewConstMetric(
					c.nodeStatus,
					prometheus.GaugeValue,
					float64(statusValue),
					strconv.Itoa(node.Id),
					node.Name,
				)
			}
		} else {
			// Emit 0 status for the node if no data is found
			ch <- prometheus.MustNewConstMetric(
				c.nodeStatus,
				prometheus.GaugeValue,
				0,
				strconv.Itoa(nodeId),
				"no_data",
			)
		}
	}

	response, err := c.client.FetchSLAPurgeItems()
	if err != nil {
		c.logger.Log("msg", "Failed to fetch SLA purge items", "err", err)
	} else {
		if response.Data != nil && response.Data.SLAItems != nil { // Check if Data and SLAItems are not nil
			// Assuming you want to count SLA purge items by status type
			countByStatus := make(map[string]int)
			for _, item := range *response.Data.SLAItems { // Dereference the pointer to range over the slice
				if item.StatusType.Name != "" { // Check if the Name field is not empty
					countByStatus[item.StatusType.Name]++
				}
			}

			// Emit metrics for each status type
			for status, count := range countByStatus {
				ch <- prometheus.MustNewConstMetric(
					c.slaPurgeItemsCountDesc,
					prometheus.GaugeValue,
					float64(count),
					status,
				)
			}
		} else {
			// Emit 0 count for sla purge items if no data is found
			ch <- prometheus.MustNewConstMetric(
				c.slaPurgeItemsCountDesc,
				prometheus.GaugeValue,
				0,
				"no_data",
			)
		}
	}

	errorsResponse, err := c.client.FetchTestErrorsRaw()
	if err != nil {
		c.logger.Log("msg", "Failed to fetch test errors", "error", err.Error())
	} else {
		totalErrors := 0.0
		errorsByType := make(map[string]float64)
		errorsByIP := make(map[string]float64)

		for _, item := range errorsResponse.Data.ResponseItems {
			for _, summary := range item.SummaryItems {
				for _, val := range summary.Values {
					totalErrors += val
				}

				for _, dimension := range summary.Dimensions {
					// Assume dimension.Name is in format "ErrorType:DNS"
					parts := strings.SplitN(dimension.Name, ":", 2)
					if len(parts) == 2 {
						dimensionType := parts[0]  // e.g., "ErrorType"
						dimensionValue := parts[1] // e.g., "DNS"

						if dimensionType == "ErrorType" && len(summary.Values) > 0 {
							errorsByType[dimensionValue] += summary.Values[0]
						} else if dimensionType == "HostIP" && len(summary.Values) > 0 {
							errorsByIP[dimensionValue] += summary.Values[0]
						}
					}
				}
			}
		}

		// Emit total errors metric
		ch <- prometheus.MustNewConstMetric(c.testErrorTotalCount, prometheus.GaugeValue, float64(totalErrors))

		// Emit error by type metrics
		for errorType, count := range errorsByType {
			ch <- prometheus.MustNewConstMetric(c.testErrorByType, prometheus.GaugeValue, float64(count), errorType)
		}

		// Emit error by IP metrics
		for ip, count := range errorsByIP {
			ch <- prometheus.MustNewConstMetric(c.testErrorByIPCount, prometheus.GaugeValue, float64(count), ip)
		}
	}

	alertsResponse, err := c.client.FetchAlerts()
	if err != nil {
		c.logger.Log("msg", "Failed to fetch test alerts", "error", err.Error())
	} else {

		for _, alert := range alertsResponse.Data.Alerts {

			if alert.Level.Name == "Critical" {
				ch <- prometheus.MustNewConstMetric(
					c.testAlertsCriticalTotal,
					prometheus.GaugeValue,
					1,
					strconv.Itoa(alert.Test.ID),
					alert.Test.Name,
					strconv.Itoa(alert.Node.ID),
					alert.Node.Name,
				)
			} else if alert.Level.Name == "Warning" {
				ch <- prometheus.MustNewConstMetric(
					c.testAlertsWarningTotal,
					prometheus.GaugeValue,
					1,
					strconv.Itoa(alert.Test.ID),
					alert.Test.Name,
					strconv.Itoa(alert.Node.ID),
					alert.Node.Name,
				)
			}
		}
	}

	for _, nodeId := range c.Config.NodeIds {
		response, err := c.client.FetchNodeTestRuns(nodeId)
		if err != nil {
			c.logger.Log("msg", "Failed to fetch node test runs", "error", err.Error())
			continue
		}

		for _, testRun := range response.Data.TestRuns {
			ch <- prometheus.MustNewConstMetric(
				c.testUsagePercentage,
				prometheus.GaugeValue,
				testRun.UsagePercentage,
				fmt.Sprintf("%d", nodeId),
				testRun.TestName,
				testRun.MonitorGroup.Name,
			)
			ch <- prometheus.MustNewConstMetric(c.testDownTimePercentage,
				prometheus.GaugeValue,
				testRun.DownTimePercentage,
				fmt.Sprintf("%d", nodeId),
				testRun.TestName,
				testRun.MonitorGroup.Name,
			)
		}
	}

	for _, nodeID := range c.Config.NodeIds {
		nodeRuneRateResponse, err := c.client.FetchNodeRunRate(nodeID)
		if err != nil {
			c.logger.Log("msg", "Failed to fetch node run rate", "error", err.Error())
			continue
		}

		nodeName := nodeRuneRateResponse.Data.Node.Node.Name
		nodeID := nodeRuneRateResponse.Data.Node.Node.ID

		ch <- prometheus.MustNewConstMetric(
			c.testRequestSlippage,
			prometheus.GaugeValue,
			float64(nodeRuneRateResponse.Data.RequestSlippages[len(nodeRuneRateResponse.Data.RequestSlippages)-1].Value),
			fmt.Sprintf("%d", nodeID),
			nodeName,
		)

		ch <- prometheus.MustNewConstMetric(
			c.testRunRate,
			prometheus.GaugeValue,
			float64(nodeRuneRateResponse.Data.RunRates[len(nodeRuneRateResponse.Data.RunRates)-1].Value),
			fmt.Sprintf("%d", nodeID),
			nodeName,
		)
	}

	for _, nodeId := range c.Config.NodeIds {
		response, err := c.client.FetchNodeTestRunCount(nodeId)
		if err != nil {
			c.logger.Log("msg", "Failed to fetch test run count", "error", err.Error())
			continue
		}

		nodeInfo := response.Data.Node
		if len(response.Data.AllTestRuns) == 0 {
			c.logger.Log("msg", "No 'All Test Runs' data found for node", "node_id", nodeInfo.ID)
			continue
		}

		lastTotalRunsData := response.Data.AllTestRuns[0].Data[len(response.Data.AllTestRuns[0].Data)-1]

		ch <- prometheus.MustNewConstMetric(
			c.totalTestRunsCount,
			prometheus.GaugeValue,
			float64(lastTotalRunsData.Value),
			strconv.Itoa(nodeInfo.ID), nodeInfo.Name,
		)

		// Handle Unique Test Runs if necessary, ensuring data is available.
		if len(response.Data.UniqueTestRuns) > 0 && len(response.Data.UniqueTestRuns[0].Data) > 0 {
			lastUniqueRunsData := response.Data.UniqueTestRuns[0].Data[len(response.Data.UniqueTestRuns[0].Data)-1]
			ch <- prometheus.MustNewConstMetric(
				c.uniqueTestRunsCount,
				prometheus.GaugeValue,
				float64(lastUniqueRunsData.Value),
				strconv.Itoa(nodeInfo.ID), nodeInfo.Name,
			)
		} else {
			c.logger.Log("msg", "No 'Unique Test Runs' data found for node", "node_id", nodeInfo.ID)
		}
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, up)
}
