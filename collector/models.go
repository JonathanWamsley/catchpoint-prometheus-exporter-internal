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

import "time"

// NodeStatusResponse represents the response structure for node status.
type NodeStatusResponse struct {
	Data *struct {
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
	} `json:"data"`
	Messages *[]struct {
		Information string `json:"information"`
		IgnoredPath *[]struct {
			Ops  string `json:"ops"`
			Path string `json:"path"`
		} `json:"ignoredPath"`
	} `json:"messages"`
	Errors *[]struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	} `json:"errors"`
	Completed   bool   `json:"completed"`
	TraceId     string `json:"traceId"`
	UsageLimits *struct {
		ClientId             int       `json:"clientId"`
		LastRequestTimestamp time.Time `json:"lastRequestTimestamp"`
		Limits               *struct {
			Minute int `json:"minute"`
			Hour   int `json:"hour"`
			Day    int `json:"day"`
		} `json:"limits"`
		Runs *struct {
			Minute int `json:"minute"`
			Hour   int `json:"hour"`
			Day    int `json:"day"`
		} `json:"runs"`
		DivisionUsageStatistics *[]struct {
			DivisionId         int `json:"divisionId"`
			ConsumerStatistics *[]struct {
				ConsumerId   int `json:"consumerId"`
				RequestCount *struct {
					Minute int `json:"minute"`
					Hour   int `json:"hour"`
					Day    int `json:"day"`
				} `json:"requestCount"`
				MaxPerDay int `json:"maxPerDay"`
			} `json:"consumerStatistics"`
		} `json:"divisionUsageStatistics"`
	} `json:"usageLimits"`
}

// SLAPurgeItemsResponse represents the response structure for SLA purge items.
type SLAPurgeItemsResponse struct {
	Data *struct {
		SLAItems      *[]SLAItem `json:"slaItems"`
		SLAPurgeItem  *SLAItem   `json:"slaPurgeItem"`
		ID            int        `json:"id"`
		SLAPurgeItems *[]struct {
			ID         int        `json:"id"`
			StatusType IDNamePair `json:"statusType"`
		} `json:"slaPurgeItems"`
		HasMore bool `json:"hasMore"`
	} `json:"data"`
	Messages *[]struct {
		Information string `json:"information"`
		IgnoredPath *[]struct {
			Ops  string `json:"ops"`
			Path string `json:"path"`
		} `json:"ignoredPath"`
	} `json:"messages"`
	Errors *[]struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	} `json:"errors"`
	Completed   bool   `json:"completed"`
	TraceId     string `json:"traceId"`
	UsageLimits *struct {
		ClientId                              int       `json:"clientId"`
		LastRequestTimestamp                  time.Time `json:"lastRequestTimestamp"`
		Limits, Runs, DivisionUsageStatistics map[string]int
	} `json:"usageLimits"`
}

// SLAItem represents an SLA item.
type SLAItem struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	Reason        string       `json:"reason"`
	StatusType    IDNamePair   `json:"statusType"`
	IntervalStart time.Time    `json:"intervalStart"`
	IntervalEnd   time.Time    `json:"intervalEnd"`
	PurgeRuns     IDNamePair   `json:"purgeRuns"`
	Tests         []IDNamePair `json:"tests"`
}

// IDNamePair represents an identifier-name pair.
type IDNamePair struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TestErrorsRawResponse represents the raw test errors response structure.
type TestErrorsRawResponse struct {
	Data        TestErrorsData `json:"data"`
	Messages    []Message      `json:"messages"`
	Errors      []Error        `json:"errors"`
	Completed   bool           `json:"completed"`
	TraceId     string         `json:"traceId"`
	UsageLimits UsageLimits    `json:"usageLimits"`
}

// TestErrorsData encapsulates the data part of the test errors response.
type TestErrorsData struct {
	ResponseItems []TestErrorsResponseItem `json:"responseItems"`
}

// TestErrorsResponseItem represents a single item in the response.
type TestErrorsResponseItem struct {
	StartTimeUtc   string                 `json:"startTimeUtc"`
	EndTimeUtc     string                 `json:"endTimeUtc"`
	TimeZoneOffset string                 `json:"timeZoneOffSet"`
	HasMoreRecords bool                   `json:"hasMoreRecords"`
	Dimensions     []TestErrorDimension   `json:"dimensions"`
	Metrics        []TestErrorMetric      `json:"metrics"`
	Items          []interface{}          `json:"items"`
	SummaryItems   []TestErrorSummaryItem `json:"summaryItems"`
}

// TestErrorDimension represents a dimension of test errors.
type TestErrorDimension struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TestErrorMetric represents a metric within the test errors.
type TestErrorMetric struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

// TestErrorSummaryItem encapsulates a summary of errors.
type TestErrorSummaryItem struct {
	Values     []float64            `json:"values"`
	Dimensions []TestErrorDimension `json:"dimensions"`
}

// Message structure for any messages in the response.
type Message struct {
	Information string `json:"information"`
	IgnoredPath *[]struct {
		Ops  string `json:"ops"`
		Path string `json:"path"`
	} `json:"ignoredPath"`
}

// Error structure for detailing errors within the response.
type Error struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// UsageLimits details the usage limits of the client.
type UsageLimits struct {
	ClientId                int            `json:"clientId"`
	LastRequestTimestamp    time.Time      `json:"lastRequestTimestamp"`
	Limits                  map[string]int `json:"limits"`
	Runs                    map[string]int `json:"runs"`
	DivisionUsageStatistics []interface{}  `json:"divisionUsageStatistics"`
}

// AlertsResponse struct for alert responses.
type AlertsResponse struct {
	Data *struct {
		Alerts   []Alert `json:"alerts"`
		HasMore  bool    `json:"hasMore"`
		Next     string  `json:"next"`
		Previous string  `json:"previous"`
	} `json:"data"`
	Messages    []Message   `json:"messages"`
	Errors      []Error     `json:"errors"`
	Completed   bool        `json:"completed"`
	TraceId     string      `json:"traceId"`
	UsageLimits UsageLimits `json:"usageLimits"`
}

// Alert struct for individual alerts.
type Alert struct {
	ID                  string     `json:"id"`
	ReportTime          time.Time  `json:"reportTime"`
	ProcessingTimestamp time.Time  `json:"processingTimestamp"`
	StartTime           time.Time  `json:"startTime"`
	Level               IDNamePair `json:"level"`
	Test                IDNamePair `json:"test"`
	TestPath            []string   `json:"testPath"`
	TestPathStyled      string     `json:"testPathStyled"`
	TestURL             string     `json:"testUrl"`
	TestSmartboardURL   string     `json:"testSmartboardUrl"`
	TestPerformanceURL  string     `json:"testPerformanceUrl"`
	TestStatisticalURL  string     `json:"testStatisticalUrl"`
	TestScatterplotURL  string     `json:"testScatterplotUrl"`
	TestRecordURL       string     `json:"testRecordUrl"`
	TestType            IDNamePair `json:"testType"`
	AlertType           IDNamePair `json:"alertType"`
	AlertSubtype        IDNamePair `json:"alertSubtype"`
	TriggerType         IDNamePair `json:"triggerType"`
	OperationType       IDNamePair `json:"operationType"`
	WarningTrigger      float64    `json:"warningTrigger"`
	CriticalTrigger     float64    `json:"criticalTrigger"`
	DurationInSeconds   int        `json:"durationInSecond"`
	Node                IDNamePair `json:"node"`
	AcknowledgedBy      IDNamePair `json:"acknowledgedBy"`
	Labels              []Label    `json:"labels"`
}

// Label struct for labels associated with alerts.
type Label struct {
	Color  string   `json:"color"`
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// NodeTestRunResponse struct for node test run responses.
type NodeTestRunResponse struct {
	Data        *NodeTestRunData `json:"data"`
	Messages    []Message        `json:"messages"`
	Errors      []Error          `json:"errors"`
	Completed   bool             `json:"completed"`
	TraceId     string           `json:"traceId"`
	UsageLimits UsageLimits      `json:"usageLimits"`
}

// NodeTestRunData struct for data in node test run responses.
type NodeTestRunData struct {
	Node       NodeDetail `json:"node"`
	TestRuns   []TestRun  `json:"testRuns"`
	TotalTests int        `json:"totalTests"`
	HasMore    bool       `json:"hasMore"`
	Next       string     `json:"next"`
	Previous   string     `json:"previous"`
}

// NodeDetail struct for detail of a node.
type NodeDetail struct {
	// Fields would be added here as needed
}

// TestRun struct for individual test runs.
type TestRun struct {
	TestId             int        `json:"testId"`
	TestName           string     `json:"testName"`
	DivisionName       string     `json:"divisionName"`
	Runs               int        `json:"runs"`
	UsagePercentage    float64    `json:"usagePercentage"`
	DownTimePercentage float64    `json:"downTimePercentage"`
	MonitorGroup       IDNamePair `json:"monitorGroup"`
}

// NodeRunRateResponse struct for node run rate responses.
type NodeRunRateResponse struct {
	Data *struct {
		Node             NodeDetails `json:"node"`
		RequestSlippages []TimeValue `json:"requestSlippages"`
		RunRates         []TimeValue `json:"runRates"`
		HasMore          bool        `json:"hasMore"`
		Next             string      `json:"next"`
		Previous         string      `json:"previous"`
	} `json:"data"`
	Messages    []Message   `json:"messages"`
	Errors      []Error     `json:"errors"`
	Completed   bool        `json:"completed"`
	TraceId     string      `json:"traceId"`
	UsageLimits UsageLimits `json:"usageLimits"`
}

// NodeDetails struct for details of a node.
type NodeDetails struct {
	Node             NodeInfo   `json:"node"`
	RequestSlippages []Slippage `json:"requestSlippages"`
	RunRates         []RunRate  `json:"runRates"`
	HasMore          bool       `json:"hasMore"`
	Next             string     `json:"next"`
	Previous         string     `json:"previous"`
}

// Slippage struct for representing a single slippage event.
type Slippage struct {
	ReportTime time.Time `json:"reportTime"`
	Value      float64   `json:"value"`
}

// RunRate struct for representing a single run rate event.
type RunRate struct {
	ReportTime time.Time `json:"reportTime"`
	Value      float64   `json:"value"`
}

// NodeInfo struct for detailed information about a node.
type NodeInfo struct {
	IsPaused                bool       `json:"isPaused"`
	Country                 IDNamePair `json:"country"`
	State                   IDNamePair `json:"state"`
	Continent               IDNamePair `json:"continent"`
	Latitude                string     `json:"latitude"`
	Longitude               string     `json:"longitude"`
	RunRate                 float64    `json:"runRate"`
	InstanceCount           int        `json:"instanceCount"`
	ActiveInstanceCount     int        `json:"activeInstanceCount"`
	Capacity                int        `json:"capacity"`
	Instances               []Instance `json:"instances"`
	ID                      int        `json:"id"`
	Name                    string     `json:"name"`
	Status                  IDNamePair `json:"status"`
	NetworkType             IDNamePair `json:"networkType"`
	OsType                  IDNamePair `json:"osType"`
	Size                    int        `json:"size"`
	IsIPv6                  bool       `json:"isIPv6"`
	NodeToNodeAddress       string     `json:"nodeToNodeAddress"`
	InternetServiceProvider IDNamePair `json:"internetServiceProvider"`
	City                    IDNamePair `json:"city"`
	Package                 IDNamePair `json:"package"`
	UtilizedInstances       int        `json:"utilizedInstances"`
}

// Instance struct for details about an instance.
type Instance struct {
	ID                int        `json:"id"`
	HostName          string     `json:"hostName"`
	Status            IDNamePair `json:"status"`
	OperatingSystem   IDNamePair `json:"operatingSystem"`
	MacAddress        string     `json:"macAddress"`
	InternalIpAddress string     `json:"internalIpAddress"`
	ActivationKey     string     `json:"activationKey"`
	Core              int        `json:"core"`
	Memory            int        `json:"memory"`
	Node              string     `json:"node"`
}

// TimeValue struct used for various time-related data points.
type TimeValue struct {
	ReportTime time.Time `json:"reportTime"`
	Value      int       `json:"value"`
}

// TestRunCountResponse struct for test run count responses.
type TestRunCountResponse struct {
	Data        TestRunCountData `json:"data"`
	Messages    []Message        `json:"messages"`
	Errors      []Error          `json:"errors"`
	Completed   bool             `json:"completed"`
	TraceId     string           `json:"traceId"`
	UsageLimits UsageLimits      `json:"usageLimits"`
}

// TestRunCountData struct for data in test run count responses.
type TestRunCountData struct {
	Node           NodeInfo      `json:"node"`
	AllTestRuns    []MonitorData `json:"allTestRuns"`
	UniqueTestRuns []MonitorData `json:"uniqueTestRuns"`
	HasMore        bool          `json:"hasMore"`
	Next           string        `json:"next"`
	Previous       string        `json:"previous"`
}

// MonitorData struct for monitor data.
type MonitorData struct {
	MonitorSetType IDNamePair  `json:"monitorSetType"`
	Data           []TimeValue `json:"data"`
}
