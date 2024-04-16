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

// MockCatchpointClient is a mock implementation of CatchpointClientInterface for testing.
type MockCatchpointClient struct {
	MockFetchNodeStatus            func(nodeId string) (*NodeStatusResponse, error)
	MockFetchSLAPurgeItems         func() (*SLAPurgeItemsResponse, error)
	MockFetchTestErrorsRaw         func() (*TestErrorsRawResponse, error)
	MockFetchAlerts                func() (*AlertsResponse, error)
	MockFetchNodeTestRuns          func(nodeId int) (*NodeTestRunResponse, error)
	MockFetchNodeRunRate           func(nodeId int) (*NodeRunRateResponse, error)
	MockFetchFetchNodeTestRunCount func(nodeId int) (*TestRunCountResponse, error)
}

func (m *MockCatchpointClient) FetchNodeStatus(nodeId string) (*NodeStatusResponse, error) {
	if m.MockFetchNodeStatus != nil {
		return m.MockFetchNodeStatus(nodeId)
	}
	// Default behavior or panic if not expected to be called
	return nil, nil
}

func (m *MockCatchpointClient) FetchSLAPurgeItems() (*SLAPurgeItemsResponse, error) {
	if m.MockFetchSLAPurgeItems != nil {
		return m.MockFetchSLAPurgeItems()
	}
	// Default behavior or panic if not expected to be called
	return nil, nil
}

func (m *MockCatchpointClient) FetchTestErrorsRaw() (*TestErrorsRawResponse, error) {
	if m.MockFetchTestErrorsRaw != nil {
		return m.MockFetchTestErrorsRaw()
	}
	// Default behavior or panic if not expected to be called
	return nil, nil
}

func (m *MockCatchpointClient) FetchAlerts() (*AlertsResponse, error) {
	if m.MockFetchAlerts != nil {
		return m.MockFetchAlerts()
	}
	// Default behavior or panic if not expected to be called
	return nil, nil
}

func (m *MockCatchpointClient) FetchNodeTestRuns(nodeId int) (*NodeTestRunResponse, error) {
	if m.MockFetchNodeTestRuns != nil {
		return m.MockFetchNodeTestRuns(nodeId)
	}
	// Default behavior or panic if not expected to be called
	return nil, nil
}

func (m *MockCatchpointClient) FetchNodeRunRate(nodeId int) (*NodeRunRateResponse, error) {
	if m.MockFetchNodeRunRate != nil {
		return m.MockFetchNodeRunRate(nodeId)
	}
	// Default behavior or panic if not expected to be called
	return nil, nil
}

func (m *MockCatchpointClient) FetchNodeTestRunCount(nodeId int) (*TestRunCountResponse, error) {
	if m.MockFetchFetchNodeTestRunCount != nil {
		return m.MockFetchFetchNodeTestRunCount(nodeId)
	}
	// Default behavior or panic if not expected to be called
	return nil, nil
}
