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
	"strings"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func NewMockCatchpointClient() *MockCatchpointClient {
	return &MockCatchpointClient{
		MockFetchNodeStatus: func(nodeId string) (*NodeStatusResponse, error) {
			// Correct structure with JSON tags as defined in your actual API client
			data := &NodeStatusResponse{
				Data: &struct {
					Nodes *[]struct {
						Id     int    `json:"id"`
						Name   string `json:"name"`
						Status *struct {
							Id   int    `json:"id"`
							Name string `json:"name"`
						} `json:"status"`
					} `json:"nodes"`
					HasMore  bool   `json:"hasMore"`
					Next     string `json:"next"`
					Previous string `json:"previous"`
				}{
					Nodes: &[]struct {
						Id     int    `json:"id"`
						Name   string `json:"name"`
						Status *struct {
							Id   int    `json:"id"`
							Name string `json:"name"`
						} `json:"status"`
					}{
						{
							Id:   1,
							Name: "Node 1",
							Status: &struct {
								Id   int    `json:"id"`
								Name string `json:"name"`
							}{
								Id:   1,
								Name: "active",
							},
						},
					},
					HasMore: false,
				},
			}
			return data, nil
		},
		MockFetchSLAPurgeItems: func() (*SLAPurgeItemsResponse, error) {
			// Correct structure with JSON tags as defined in your actual API client
			data := &SLAPurgeItemsResponse{
				Data: &struct {
					SLAItems      *[]SLAItem `json:"slaItems"`
					SLAPurgeItem  *SLAItem   `json:"slaPurgeItem"`
					ID            int        `json:"id"`
					SLAPurgeItems *[]struct {
						ID         int        `json:"id"`
						StatusType IDNamePair `json:"statusType"`
					} `json:"slaPurgeItems"`
					HasMore bool `json:"hasMore"`
				}{
					SLAItems: &[]SLAItem{
						{
							ID:         1,
							StatusType: IDNamePair{ID: 1, Name: "Active"},
						},
					},
					HasMore: false,
				},
			}
			return data, nil
		},
		MockFetchTestErrorsRaw: func() (*TestErrorsRawResponse, error) {
			return &TestErrorsRawResponse{
				Data: TestErrorsData{
					ResponseItems: []TestErrorsResponseItem{
						{
							StartTimeUtc: "2024-04-10T23:13:49.293Z",
							EndTimeUtc:   "2024-04-11T02:13:49.293Z",
							SummaryItems: []TestErrorSummaryItem{
								{
									Values: []float64{6},
									Dimensions: []TestErrorDimension{
										{ID: 1, Name: "ErrorType:DNS"},
										{ID: 101, Name: "HostIP:192.0.2.1"},
									},
								},
								{
									Values: []float64{4},
									Dimensions: []TestErrorDimension{
										{ID: 2, Name: "ErrorType:Connection"},
										{ID: 102, Name: "HostIP:198.51.100.1"},
									},
								},
								{
									Values: []float64{3},
									Dimensions: []TestErrorDimension{
										{ID: 3, Name: "ErrorType:SSL"},
										{ID: 103, Name: "HostIP:203.0.113.1"},
									},
								},
								{
									Values: []float64{2},
									Dimensions: []TestErrorDimension{
										{ID: 4, Name: "ErrorType:NoResponse"},
										{ID: 104, Name: "HostIP:192.0.2.2"},
									},
								},
							},
						},
					},
				},
				Completed: true,
			}, nil
		},
		MockFetchAlerts: func() (*AlertsResponse, error) {
			data := &struct {
				Alerts   []Alert `json:"alerts"`
				HasMore  bool    `json:"hasMore"`
				Next     string  `json:"next"`
				Previous string  `json:"previous"`
			}{
				Alerts: []Alert{
					{
						ID: "1",
						Level: IDNamePair{
							ID:   1,
							Name: "Critical",
						},
						Test: IDNamePair{
							ID:   1,
							Name: "Test 1",
						},
						Node: IDNamePair{
							ID:   1,
							Name: "Node 1",
						},
					},
					{
						ID: "2",
						Level: IDNamePair{
							ID:   2,
							Name: "Warning",
						},
						Test: IDNamePair{
							ID:   2,
							Name: "Test 2",
						},
						Node: IDNamePair{
							ID:   1,
							Name: "Node 1",
						},
					},
				},
			}
			return &AlertsResponse{Data: data, Completed: true}, nil
		},
		MockFetchNodeTestRuns: func(nodeId int) (*NodeTestRunResponse, error) {
			return &NodeTestRunResponse{
				Data: &NodeTestRunData{
					Node: NodeDetail{},
					TestRuns: []TestRun{
						{
							TestId:             1,
							TestName:           "Test 1",
							DivisionName:       "Division 1",
							Runs:               10,
							UsagePercentage:    75.0,
							DownTimePercentage: 5.0,
							MonitorGroup:       IDNamePair{ID: 1, Name: "Browser"},
						},
					},
					TotalTests: 10,
					HasMore:    false,
					Next:       "",
					Previous:   "",
				},
				Messages:    []Message{},
				Errors:      []Error{},
				Completed:   true,
				TraceId:     "some-trace-id",
				UsageLimits: UsageLimits{},
			}, nil
		},
		MockFetchNodeRunRate: func(nodeId int) (*NodeRunRateResponse, error) {
			data := &struct {
				Node             NodeDetails `json:"node"`
				RequestSlippages []TimeValue `json:"requestSlippages"`
				RunRates         []TimeValue `json:"runRates"`
				HasMore          bool        `json:"hasMore"`
				Next             string      `json:"next"`
				Previous         string      `json:"previous"`
			}{
				Node: NodeDetails{
					Node: NodeInfo{
						ID:   1,
						Name: "Node 1",
					},
				},
				RequestSlippages: []TimeValue{
					{
						ReportTime: time.Now(),
						Value:      100,
					},
				},
				RunRates: []TimeValue{
					{
						ReportTime: time.Now(),
						Value:      95,
					},
				},
				HasMore:  false,
				Next:     "",
				Previous: "",
			}

			return &NodeRunRateResponse{
				Data: data,
				Messages: []Message{
					{
						Information: "Sample fetch node run rate response generated.",
					},
				},
				Errors: []Error{
					{
						ID:      "0",
						Message: "No error",
					},
				},
				Completed: true,
				TraceId:   "sample-trace-id",
				UsageLimits: UsageLimits{
					ClientId:                1,
					LastRequestTimestamp:    time.Now(),
					Limits:                  map[string]int{"daily": 1000, "hourly": 100},
					Runs:                    map[string]int{"daily": 900, "hourly": 80},
					DivisionUsageStatistics: []interface{}{},
				},
			}, nil
		},
		MockFetchFetchNodeTestRunCount: func(nodeId int) (*TestRunCountResponse, error) {
			return &TestRunCountResponse{
				Data: TestRunCountData{
					Node: NodeInfo{
						ID:   1,
						Name: "Node 1",
					},
					AllTestRuns: []MonitorData{
						{
							MonitorSetType: IDNamePair{ID: 1, Name: "Browser"},
							Data:           []TimeValue{{Value: 25}},
						},
						{
							MonitorSetType: IDNamePair{ID: 2, Name: "Browser"},
							Data:           []TimeValue{{Value: 10}},
						},
					},
					UniqueTestRuns: []MonitorData{
						{
							MonitorSetType: IDNamePair{ID: 1, Name: "Browser"},
							Data:           []TimeValue{{Value: 25}},
						},
					},
				},
				Completed: true,
			}, nil
		},
	}
}
func TestCollector_Collect_Success(t *testing.T) {
	cfg := &Config{
		BearerToken: "testToken",
		NodeIds:     []int{1},
	}
	logger := log.NewNopLogger()
	col := NewCollector(logger, cfg)
	col.client = NewMockCatchpointClient()

	registry := prometheus.NewPedanticRegistry()
	if err := registry.Register(col); err != nil {
		t.Error("failed to register collector:", err)
		return
	}

	// Gather all metrics
	metrics, err := registry.Gather()
	if err != nil {
		t.Fatalf("gathering metrics failed: %v", err)
	}

	// Check if we collected 14 metrics
	if len(metrics) != 14 {
		t.Errorf("expected 14 metrics, got %d", len(metrics))
	}

	// Use the GatherAndCompare function to check all metrics against expected values
	expected := `
        # HELP catchpoint_node_status The operational status of a Catchpoint node (1 for active, 0 for inactive).
        # TYPE catchpoint_node_status gauge
        catchpoint_node_status{node_id="1", node_name="Node 1"} 1
        # HELP catchpoint_sla_purge_items_count Count of SLA purge items by status.
        # TYPE catchpoint_sla_purge_items_count gauge
        catchpoint_sla_purge_items_count{status_id="Active"} 1
        # HELP catchpoint_test_error_total_count The total count of all test errors.
        # TYPE catchpoint_test_error_total_count gauge
        catchpoint_test_error_total_count 15
        # HELP catchpoint_test_alerts_critical_count Total number of critical alerts by test.
        # TYPE catchpoint_test_alerts_critical_count gauge
        catchpoint_test_alerts_critical_count{node_id="1", node_name="Node 1", test_id="1", test_name="Test 1"} 1
		# HELP catchpoint_test_alerts_warning_count Total number of warning alerts by test.
        # TYPE catchpoint_test_alerts_warning_count gauge
        catchpoint_test_alerts_warning_count{node_id="1",node_name="Node 1",test_id="2",test_name="Test 2"} 1
		# HELP catchpoint_test_error_by_ip_count Error counts traced back to specific IP addresses.
        # TYPE catchpoint_test_error_by_ip_count gauge
        catchpoint_test_error_by_ip_count{ip="192.0.2.1"} 6
        catchpoint_test_error_by_ip_count{ip="192.0.2.2"} 2
        catchpoint_test_error_by_ip_count{ip="198.51.100.1"} 4
        catchpoint_test_error_by_ip_count{ip="203.0.113.1"} 3
        # HELP catchpoint_test_error_by_type_count Count of errors segmented by error type.
        # TYPE catchpoint_test_error_by_type_count gauge
        catchpoint_test_error_by_type_count{error_type="Connection"} 4
        catchpoint_test_error_by_type_count{error_type="DNS"} 6
        catchpoint_test_error_by_type_count{error_type="NoResponse"} 2
        catchpoint_test_error_by_type_count{error_type="SSL"} 3
        # HELP catchpoint_usage_percentage Usage percentage of test runs on a node
        # TYPE catchpoint_usage_percentage gauge
		catchpoint_usage_percentage{monitor_group="Browser",node_id="1",test_name="Test 1"} 75
        # HELP catchpoint_downtime_percentage Downtime percentage of test runs on a node
        # TYPE catchpoint_downtime_percentage gauge
		catchpoint_downtime_percentage{monitor_group="Browser",node_id="1",test_name="Test 1"} 5
        # HELP catchpoint_node_request_slippage The slippage in test requests timings on a node, showing delays in scheduled test executions
        # TYPE catchpoint_node_request_slippage gauge
        catchpoint_node_request_slippage{node_id="1", node_name="Node 1"} 100
        # HELP catchpoint_node_run_rate The rate of which test runs are successfully completed on a node
        # TYPE catchpoint_node_run_rate gauge
        catchpoint_node_run_rate{node_id="1", node_name="Node 1"} 95
        # HELP catchpoint_total_test_runs_count The total number of test runs on a node
        # TYPE catchpoint_total_test_runs_count gauge
        catchpoint_total_test_runs_count{node_id="1",node_name="Node 1"} 25
		# HELP catchpoint_unique_test_runs_count The number of unique test runs on a node
        # TYPE catchpoint_unique_test_runs_count gauge
        catchpoint_unique_test_runs_count{node_id="1",node_name="Node 1"} 25
		# HELP catchpoint_up Indicates whether the last scrape of metrics from Catchpoint was successful.
        # TYPE catchpoint_up gauge
        catchpoint_up 1
    `

	if err := testutil.GatherAndCompare(registry, strings.NewReader(expected)); err != nil {
		t.Errorf("collected metrics did not match expected: %v", err)
	}
}

func TestCollector_Collect_Failure(t *testing.T) {
	// Setup mock client to return an error on FetchNodeStatus
	mockClient := &MockCatchpointClient{
		MockFetchNodeStatus: func(nodeId string) (*NodeStatusResponse, error) {
			return nil, fmt.Errorf("simulated API error")
		},
		MockFetchSLAPurgeItems: func() (*SLAPurgeItemsResponse, error) {
			return nil, fmt.Errorf("simulated API error")
		},
		MockFetchTestErrorsRaw: func() (*TestErrorsRawResponse, error) {
			return nil, fmt.Errorf("simulated API error")
		},
		MockFetchAlerts: func() (*AlertsResponse, error) {
			return nil, fmt.Errorf("simulated API error")
		},
		MockFetchNodeTestRuns: func(nodeId int) (*NodeTestRunResponse, error) {
			return nil, fmt.Errorf("simulated API error")
		},
		MockFetchNodeRunRate: func(nodeId int) (*NodeRunRateResponse, error) {
			return nil, fmt.Errorf("simulated API error")
		},
		MockFetchFetchNodeTestRunCount: func(nodeId int) (*TestRunCountResponse, error) {
			return nil, fmt.Errorf("simulated API error")
		},
	}

	cfg := &Config{
		BearerToken: "testToken",
		NodeIds:     []int{1},
	}

	col := NewCollector(log.NewNopLogger(), cfg)
	col.client = mockClient // Inject the mock client

	registry := prometheus.NewPedanticRegistry()
	if err := registry.Register(col); err != nil {
		t.Error("failed to register collector:", err)
		return
	}

	// Collect all metrics into a Metrics slice.
	metrics, err := testutil.GatherAndCount(registry, "catchpoint_up")
	if err != nil {
		t.Errorf("unexpected collecting error: %v", err)
		return
	}

	// Verify we have exactly one 'up' metric and that it's set to 0.
	if metrics != 1 {
		t.Errorf("Expected exactly one 'up' metric, got %d", metrics)
	}

	// Check the value of the 'up' metric.
	if err := testutil.CollectAndCompare(col, strings.NewReader(`
        # HELP catchpoint_up Indicates whether the last scrape of metrics from Catchpoint was successful.
        # TYPE catchpoint_up gauge
        catchpoint_up 1
    `), "catchpoint_up"); err != nil {
		t.Error("up metric value is not 0 as expected on failure:", err)
	}
}
